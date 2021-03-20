package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/girishg4t/bhavcopy-backend/pkg/config"
	"github.com/girishg4t/bhavcopy-backend/pkg/dataProcessor"
	"github.com/girishg4t/bhavcopy-backend/pkg/github"
	"github.com/girishg4t/bhavcopy-backend/utils"
)

func csvGenerator(w http.ResponseWriter, req *http.Request) {
	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	var obj config.Symboles

	if err := json.Unmarshal(body, &obj); err != nil {
		fmt.Printf("err: %s", err)
	}
	obj.Date = strings.ToUpper(obj.Date)
	if !validateInput(obj) {
		fmt.Fprintf(w, "Not a valid date format, it should be in ddMMMYYYY format, eg. 02Mar2020")
	}
	fmt.Println("api Parameter is correct")
	config.LoadEnv()
	conn := github.ConnectToGit(obj)
	fmt.Println("Read env variable")

	csvData := conn.ReadIfFileExistsFromGit(obj)

	if csvData == nil {
		fmt.Println(obj.Date + " File not in git")
		config.GetFilePath(obj)
		err = dataProcessor.Downloadzip(obj)
		if err != nil {
			w.Header().Set("Content-Type", "text/csv")
			w.Header().Set("Content-Disposition", "attachment;filename=TheCSVFileName.csv")
			w.Write([]byte{})
			return
		}
		fmt.Println("Done downloading zip file nse")
		csvData = dataProcessor.ReadZipfile()
		fmt.Println("Done reading zip file nse")
		utils.SaveCSV(csvData)
		conn.UpdateToGithub(obj)
		fmt.Println("Done uploading to github")
		e := os.Remove(config.LocalZipPath)
		if e != nil {
			log.Fatal(e)
		}
		e = os.Remove(config.LocalFilePath)
		if e != nil {
			log.Fatal(e)
		}
	}
	bytesBuffer := dataProcessor.FilterCsvData(csvData, obj)
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=TheCSVFileName.csv")
	w.Write(bytesBuffer.Bytes())
}

func welcome(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	http.HandleFunc("/getbhavcopy", csvGenerator)
	http.HandleFunc("/", welcome)
	http.ListenAndServe(":"+port, nil)
}

func validateInput(obj config.Symboles) bool {
	re := regexp.MustCompile(config.RegexDate)
	if !re.MatchString(obj.Date) {
		return false
	}
	return true
}
