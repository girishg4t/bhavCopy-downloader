package dataProcessor

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/girishg4t/bhavcopy-backend/pkg/config"
	"github.com/girishg4t/bhavcopy-backend/utils"
)

var exchangeConfig config.ExchangeConfig

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
			if !utils.Contains(each[0], obj.Stocks) {
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
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
	req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"119\", \"Chromium\";v=\"119\", \"Not?A_Brand\";v=\"24\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("Error in reading zip file " + string(err.Error()))
		fmt.Printf("err: %s", resp.Status)
		//os.Remove(config.LocalZipPath)
		return errors.New("not Allowed")
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	// Create the file
	out, err := os.Create(config.LocalZipPath)
	if err != nil {
		fmt.Println("Error in creating zip file")
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
		//api := fmt.Sprintf(config.NSEURLAPI, strings.ToUpper(t.Format("02Jan2006")))
		url := config.NSEURL + "sec_bhavdata_full_" + obj.Date[0:2] + config.MonthMapping[obj.Date[2:5]] + obj.Date[5:9] + ".csv"
		fmt.Println("NSE url " + url)
		return url
	}
	if obj.Exchange == "BSE" {
		dat, _ := utils.ReadJSON("./config/bse.json")
		json.Unmarshal(dat, &exchangeConfig)
		api := fmt.Sprintf(config.BSEURLAPI, obj.Date[0:2]+config.MonthMapping[obj.Date[2:5]]+obj.Date[7:9])
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

func DownloadDeliverableDataNSE(date string) error {
	url := fmt.Sprintf("https://nsearchives.nseindia.com/products/content/sec_bhavdata_full_%s.csv", date)
	// Create an HTTP client with HTTP/1.1 transport
	transport := &http.Transport{
		ForceAttemptHTTP2: false,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   15 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return err
	}

	// Set required headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", "https://www.nseindia.com/")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	var lastErr error
	for i := 0; i < 3; i++ {
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed to call API (attempt %d): %v", i+1, err)
			lastErr = err
			time.Sleep(2 * time.Second)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(resp.Body)
			log.Printf("API call to %s failed. Status: %d, Response: %s", url, resp.StatusCode, string(body))
			lastErr = fmt.Errorf("API call failed with status %d", resp.StatusCode)
			time.Sleep(2 * time.Second)
			continue
		}

		// Create the file
		out, err := os.Create(config.LocalDeliverablePath)
		if err != nil {
			fmt.Println("Error in creating zip file")
			fmt.Printf("err: %s", err)
		}
		defer out.Close()

		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			fmt.Printf("err: %s", err)
		}
	}
	log.Printf("Failed to create request: %v", lastErr)

	return nil
}
