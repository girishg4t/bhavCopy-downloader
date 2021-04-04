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
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func csvGenerator(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
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

func welcome(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "hello\n")
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	router := httprouter.New()
	router.GET("/", welcome)
	router.POST("/getbhavcopy", csvGenerator)

	handler := cors.Default().Handler(router)

	http.ListenAndServe(":8080", handler)
}

func validateInput(obj config.Symboles) bool {
	re := regexp.MustCompile(config.RegexDate)
	if !re.MatchString(obj.Date) {
		return false
	}
	return true
}
