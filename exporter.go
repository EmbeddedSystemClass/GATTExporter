package gattexporter

import (
	"os"
	"bufio"
	"strings"
	"strconv"
)

type exporter struct{
	reader *bufio.Reader
	writer *bufio.Writer
}

var Exporter = exporter{
	reader:bufio.NewReader(os.Stdin),
	writer:bufio.NewWriter(os.Stdout),
}

func(exporter *exporter) Start() {
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
			if err := ble.Connect(parts[1]) != nil; err != nil{
				Exporter.writer.WriteString("Error : " + err)
			}
		}else{
			Exporter.writer.WriteString("Usage : CONNECT XX:XX:XX:XX:XX:XX")
		}
		break
	case "DISCONNECT":
		if err := ble.Disconnect() != nil; err != nil{
			Exporter.writer.WriteString("Error : " + err)
		}
		break
	case "LIST":
		if len(parts) == 2{
			timeout, err := strconv.Atoi(parts[1])

			if err != nil{
				Exporter.writer.WriteString("Timeout must be an integer.")
			}

			if err := ble.List(timeout) != nil; err != nil{
				Exporter.writer.WriteString("Error : " + err)
			}
		}else{
			Exporter.writer.WriteString("Usage : LIST TIMEOUT_IN_SEC")
		}
		break
	}
}

