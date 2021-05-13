package main

import (
	"bufio"
	"encoding/csv"
	"os"
)

type LogData struct {
	LogTime    string `bson:"log_time"`
	DeviceName string `bson:"device_name"`
	AttacName  string `bson:"attack_name"`
	RawPacket  string `bson:"raw_packet"`
}

// ReadCSV read csv to row
func ReadCSV(fileName string) []interface{} {

	// totalResult []
	var logDataResult []interface{}

	csvData, _ := os.Open("./bigdatasample/" + fileName)
	rdr := csv.NewReader(bufio.NewReader(csvData))

	rows, _ := rdr.ReadAll()

	for i, row := range rows {

		if i == 0 {
			continue
		}
		logData := &LogData{}

		for j := range row {
			switch j {
			case 1:
				logData.LogTime = rows[i][j]
			case 2:
				logData.DeviceName = rows[i][j]
			case 3:
				logData.AttacName = rows[i][j]
			case 4:
				logData.RawPacket = rows[i][j]
			}

		}

		logDataResult = append(logDataResult, logData)

	}

	return logDataResult
}
