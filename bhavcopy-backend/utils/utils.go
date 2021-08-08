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

func SaveCSV(data [][]string) error {
	file, err := os.Create(config.LocalFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	dd := GetDeliverableData()
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

func GetDeliverableData() [][]string {
	content, _ := ioutil.ReadFile(config.LocalDeliverablePath)
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
		if trimAndToUpper(cs[0]) == trimAndToUpper(d[0]) &&
			trimAndToUpper(cs[1]) == trimAndToUpper(d[1]) &&
			trimAndToUpper(cs[10]) == trimAndToUpper(d[2]) {
			return trimAndToUpper(d[13]), trimAndToUpper(d[14])
		}
	}
	return "", ""
}

func trimAndToUpper(d string) string {
	return strings.ToUpper(strings.Trim(d, " "))
}
