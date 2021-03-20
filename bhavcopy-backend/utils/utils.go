package utils

import (
	"encoding/csv"
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

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			return err
		}
	}
	return nil
}
