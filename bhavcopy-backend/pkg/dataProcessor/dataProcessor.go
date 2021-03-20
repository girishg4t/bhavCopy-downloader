package dataProcessor

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/girishg4t/bhavcopy-backend/pkg/config"
	"github.com/girishg4t/bhavcopy-backend/utils"
)

var exchangeConfig config.ExchangeConfig
var monthMapping = map[string]string{
	"JAN": "01", "FEB": "02",
	"MAR": "03", "APR": "04", "MAY": "05",
	"JUN": "06", "JULY": "07", "AUG": "08",
	"SEP": "09", "OCT": "10", "NOV": "11",
	"DEC": "12",
}

func FilterCsvData(csvData [][]string, obj config.Symboles) *bytes.Buffer {
	b := &bytes.Buffer{} // creates IO Writer
	wr := csv.NewWriter(b)
	wr.Write(csvData[0])
	if len(obj.Stocks) == 0 {
		wr.WriteAll(csvData[1:])
		wr.Flush()
		return b
	}
	if obj.Exchange == "NSE" {
		for _, each := range csvData[1:] {
			if !utils.Contains(each[0], obj.Stocks) || obj.Fund[0:2] != each[1] {
				continue
			}
			wr.Write(each)
		}
	}
	if obj.Exchange == "BSE" {
		for _, each := range csvData[1:] {
			if !utils.Contains(each[0], obj.Stocks) {
				continue
			}
			wr.Write(each)
		}
	}
	wr.Flush()
	return b
}

func Downloadzip(obj config.Symboles) error {
	s := ReadIndicesConfig(obj)
	if s == "" {
		return errors.New("invalid input")
	}
	req, err := http.NewRequest("GET", s, nil)
	req.Header.Set("User-Agent", exchangeConfig.UserAgent)
	req.Header.Set("Referer", exchangeConfig.Referer)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("Error in reading zip file " + string(resp.StatusCode))
		fmt.Printf("err: %s", resp.Status)
		os.Remove(config.LocalZipPath)
		return errors.New("not Allowed")
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	// Create the file
	out, err := os.Create(config.LocalZipPath)
	if err != nil {
		fmt.Println("Error in reading zip file")
		fmt.Printf("err: %s", err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	return nil
}

func ReadIndicesConfig(obj config.Symboles) string {
	t, _ := time.Parse("02Jan2006", obj.Date)
	if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
		return ""
	}
	if obj.Exchange == "" || obj.Exchange == "NSE" {
		dat, _ := utils.ReadJSON("./config/nse.json")
		json.Unmarshal(dat, &exchangeConfig)
		api := fmt.Sprintf(config.NSEURLAPI, strings.ToUpper(t.Format("02Jan2006")))
		url := config.NSEURL + obj.Fund + "/" + obj.Date[5:9] + "/" + strings.ToUpper(obj.Date[2:5]) + "/" + api
		fmt.Println("NSE url " + url)
		return url
	}
	if obj.Exchange == "BSE" {
		dat, _ := utils.ReadJSON("./config/bse.json")
		json.Unmarshal(dat, &exchangeConfig)
		api := fmt.Sprintf(config.BSEURLAPI, obj.Date[0:2]+monthMapping[obj.Date[2:5]]+obj.Date[7:9])
		url := config.BSEURL + obj.Fund + "/" + api
		fmt.Println("BSE url " + url)
		return url
	}
	return ""
}

func ReadZipfile() [][]string {
	// Create a reader out of the zip archive
	zipReader, err := zip.OpenReader(config.LocalZipPath)
	if err != nil {
		log.Fatal(err)
	}
	defer zipReader.Close()

	// Iterate through each file found in
	for _, file := range zipReader.Reader.File {

		zippedFile, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer zippedFile.Close()

		reader := csv.NewReader(zippedFile)
		csvData, err := reader.ReadAll()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return csvData
	}
	return nil
}
