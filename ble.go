package gattexporter

import (
	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
	"encoding/json"
	"github.com/iris-contrib/errors"
	"fmt"
)

type bledevice struct{
	device gatt.Device
	initialized bool
	scanning bool
}

var ble bledevice = bledevice{}

var periph map[string]*gatt.Peripheral = make(map[string]*gatt.Peripheral)

//Init method initialize a BLE device
func(device *bledevice) Init() error{
	d, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		return err
	}

	// Register handlers.
	d.Handle(
		gatt.PeripheralDiscovered(onPeriphDiscovered),
		gatt.PeripheralConnected(onPeriphConnected),
		gatt.PeripheralDisconnected(onPeriphDisconnected),
	)
	d.Init(onStateChanged)
	return nil;
}

func(device *bledevice) Connect(mac string) error{
	if device.initialized{
		if device.scanning{
			device.device.StopScanning()
			device.scanning = false
		}
		if p := periph[mac]; p != nil{
			(*p).Device().Connect(*p)
		}else{
			return errors.New("Peripheral not known")
		}
	}else{
		return errors.New("Device not initialized yet")
	}

	return nil;
}

func(device *bledevice) Read(uuid string) error{
	return nil;
}

func(device *bledevice) Write(uuid string, value []byte) error{
	return nil;
}

func(device *bledevice) Disconnect() error{
	return nil;
}

func(device *bledevice) Scan() error{
	if device.initialized{
		device.device.Scan([]gatt.UUID{}, false)
		device.scanning = true
	}else{
		return errors.New("Device not initialized yet")
	}
	return nil
}

func onStateChanged(d gatt.Device, s gatt.State) {
	switch s {
	case gatt.StatePoweredOn:
		ble.initialized = true
		ble.device = d
		return
	default:
		ble.scanning = false
		ble.device.StopScanning()
	}
}

func onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	jsonData, err := json.Marshal(Peripheral{
		ID:p.ID(),
		Name:p.Name(),
		LocalName:a.LocalName,
		RSSI:rssi,
		TXPower:a.TxPowerLevel,
		ManufacturerData:a.ManufacturerData,
		ServiceData:a.ServiceData,
	})
	if err == nil{
		Exporter.data.WriteString(string(jsonData) + "\n")
		periph[p.ID()] = &p
	}
}

func onPeriphConnected(p gatt.Peripheral, err error) {
	Exporter.writer.WriteString("Connected")

	if err := p.SetMTU(500); err != nil {
		Exporter.writer.WriteString(fmt.Sprintf("Failed to set MTU, err: %s\n", err))
	}

	// Discovery services
	ss, err := p.DiscoverServices(nil)

	if err != nil {
		Exporter.writer.WriteString(fmt.Sprintf("Failed to discover services, err: %s\n", err))
		return
	}

	for _, s := range ss {
		msg := "Service: " + s.UUID().String()
		if len(s.Name()) > 0 {
			msg += " (" + s.Name() + ")"
		}
		Exporter.writer.WriteString(msg)

		// Discovery characteristics
		cs, err := p.DiscoverCharacteristics(nil, s)
		if err != nil {
			Exporter.writer.WriteString(fmt.Sprintf("Failed to discover characteristics, err: %s\n", err))
			continue
		}

		for _, c := range cs {
			msg := "  Characteristic  " + c.UUID().String()
			if len(c.Name()) > 0 {
				msg += " (" + c.Name() + ")"
			}
			msg += "\n    properties    " + c.Properties().String()
			Exporter.writer.WriteString(msg)

			// Read the characteristic, if possible.
			if (c.Properties() & gatt.CharRead) != 0 {
				b, err := p.ReadCharacteristic(c)
				if err != nil {
					Exporter.writer.WriteString(fmt.Sprintf("Failed to read characteristic, err: %s\n", err))
					continue
				}
				Exporter.writer.WriteString(fmt.Sprintf("    value         %x | %q\n", b, b))
			}

			// Discovery descriptors
			ds, err := p.DiscoverDescriptors(nil, c)
			if err != nil {
				Exporter.writer.WriteString(fmt.Sprintf("Failed to discover descriptors, err: %s\n", err))
				continue
			}

			for _, d := range ds {
				msg := "  Descriptor      " + d.UUID().String()
				if len(d.Name()) > 0 {
					msg += " (" + d.Name() + ")"
				}
				Exporter.writer.WriteString(msg)

				// Read descriptor (could fail, if it's not readable)
				b, err := p.ReadDescriptor(d)
				if err != nil {
					Exporter.writer.WriteString(fmt.Sprintf("Failed to read descriptor, err: %s\n", err))
					continue
				}
				Exporter.writer.WriteString(fmt.Sprintf("    value         %x | %q\n", b, b))
			}

			// Subscribe the characteristic, if possible.
			if (c.Properties() & (gatt.CharNotify | gatt.CharIndicate)) != 0 {
				f := func(c *gatt.Characteristic, b []byte, err error) {
					jsonData, _ := json.Marshal(BLEData{
						UUID:c.UUID().String(),
						Value:b,
					})
					if jsonData != nil{
						Exporter.data.WriteString(string(jsonData))
					}
				}
				if err := p.SetNotifyValue(c, f); err != nil {
					Exporter.writer.WriteString(fmt.Sprintf("Failed to subscribe characteristic, err: %s\n", err))
					continue
				}
			}

		}
	}
}

func onPeriphDisconnected(p gatt.Peripheral, err error) {
	Exporter.writer.WriteString("Disconnected")
}
