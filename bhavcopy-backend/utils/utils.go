package utils

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/girishg4t/bhavcopy-backend/pkg/config"
)

// Contains tells whether a contains x.
func Contains(x string, a []string) bool {
	for _, n := range a {
		if strings.ToLower(x) == strings.ToLower(n) {
			return true
		}
	}
	return false
}

func ReadJSON(file string) ([]byte, error) {
	jsonFile, err := os.Open(file)
	dat, err := ioutil.ReadAll(jsonFile)
	return dat, err
}

func SaveCSV(data [][]string, dd [][]string) error {
	file, err := os.Create(config.LocalFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if dd != nil {
		// Write header
		header := append(data[0], "DELIV_QTY")
		header = append(header, "DELIV_PER")
		writer.Write(header)

		for _, value := range data[1:] {
			delQnt, delPer := findDeliverableData(value, dd[1:])
			value = append(value, delQnt)
			value = append(value, delPer)
			err := writer.Write(value)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetDeliverableData() [][]string {
	content, err := ioutil.ReadFile(config.LocalDeliverablePath)
	if err != nil {
		return nil
	}
	reader := csv.NewReader(strings.NewReader(string(content)))
	csvDeliverableData, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error", err)
		return nil
	}
	return csvDeliverableData
}

func findDeliverableData(cs []string, deliverableData [][]string) (string, string) {
	for _, d := range deliverableData {
		if TrimAndToUpper(cs[0]) == TrimAndToUpper(d[0]) &&
			TrimAndToUpper(cs[1]) == TrimAndToUpper(d[1]) &&
			TrimAndToUpper(cs[10]) == TrimAndToUpper(d[2]) {
			return TrimAndToUpper(d[13]), TrimAndToUpper(d[14])
		}
	}
	return "", ""
}

func TrimAndToUpper(d string) string {
	return strings.ToUpper(strings.Trim(d, " "))
}
