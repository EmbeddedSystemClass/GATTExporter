package gattexporter

import "github.com/paypal/gatt"

type Peripheral struct {
	ID string
	Name string
	LocalName string
	TXPower int
	RSSI int
	ManufacturerData []byte
	ServiceData []gatt.ServiceData
}
