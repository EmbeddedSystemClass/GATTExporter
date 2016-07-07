package gattexporter

type bledevice struct{

}

var ble bledevice = bledevice{

}

func(device *bledevice) Connect(mac string) error{
	return nil;
}

func(device *bledevice) Disconnect() error{
	return nil;
}

func(device *bledevice) List(timeout int) error{
	return nil;
}