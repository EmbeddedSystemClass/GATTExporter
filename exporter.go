package gattexporter

import (
	"os"
	"bufio"
	"strings"
	"errors"
)

type exporter struct{
	reader *bufio.Reader
	writer *bufio.Writer
	data *bufio.Writer
}

var Exporter = exporter{
	reader:bufio.NewReader(os.Stdin),
	writer:bufio.NewWriter(os.Stderr),
	data:bufio.NewWriter(os.Stdout),
}

func(exporter *exporter) Start() {
	err := ble.Init()

	if err != nil {
		Exporter.writer.WriteString("Fail to initialize BLE device.")
	}

	text, err := Exporter.reader.ReadString('\n')
	for err == nil{
		Exporter.interpretCommand(text)
		text, err = Exporter.reader.ReadString('\n')
	}
}

func(exporter *exporter) interpretCommand(command string) {
	parts := strings.Split(command, " ")

	switch parts[0] {
	case "CONNECT":
		if len(parts) == 2{
			if err := ble.Connect(parts[1]); err != nil{
				Exporter.writer.WriteString("Error : " + err.Error())
			}
		}else{
			Exporter.writer.WriteString("Usage : CONNECT XX:XX:XX:XX:XX:XX")
		}
		break
	case "DISCONNECT":
		if err := ble.Disconnect(); err != nil{
			Exporter.writer.WriteString("Error : " + err.Error())
		}
		break
	case "SCAN":
		if err := ble.Scan(); err != nil{
			Exporter.writer.WriteString("Error : " + err.Error())
		}
		break
	case "READ":
		if len(parts) == 2{
			if err := ble.Read(parts[1]); err != nil{
				Exporter.writer.WriteString("Error : " + err.Error())
			}
		}else{
			Exporter.writer.WriteString("Usage : READ UUID")
		}
		break
	case "WRITE":
		if len(parts) == 3{
			b, err := convertStringToByteArray(parts[2])
			if err != nil{
				Exporter.writer.WriteString("Error : " + err.Error())
			}
			if err = ble.Write(parts[1], b); err != nil{
				Exporter.writer.WriteString("Error : " + err.Error())
			}
		}else{
			Exporter.writer.WriteString("Usage : WRITE UUID VALUE")
		}
		break
	}
}

func convertStringToByteArray(from string)(to []byte, err error){
	return nil,errors.New("Not implemented yet")
}
