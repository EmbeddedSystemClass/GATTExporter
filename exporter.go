package gattexporter

import (
	"os"
	"bufio"
	"strings"
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

		}else{
			Exporter.writer.WriteString("Usage : CONNECT XX:XX:XX:XX:XX:XX")
		}
		break
	case "DISCONNECT":
		break
	case "LIST":
		if len(parts) == 2{

		}else{
			Exporter.writer.WriteString("Usage : LIST TIMEOUT_IN_SEC")
		}
		break
	}
}

