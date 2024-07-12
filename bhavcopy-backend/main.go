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
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func csvGenerator(w http.ResponseWriter, req *http.Request) {
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
		// err = dataProcessor.Downloadzip(obj)
		// if err != nil {
		// 	w.Header().Set("Content-Type", "text/csv")
		// 	w.Header().Set("Content-Disposition", "attachment;filename=TheCSVFileName.csv")
		// 	w.Write([]byte{})
		// 	return
		// }
		var dd [][]string = nil
		if strings.ToLower(obj.Exchange) == "nse" {
			ddDate := obj.Date[0:2] + config.MonthMapping[obj.Date[2:5]] + obj.Date[5:9]
			err = dataProcessor.DownloadDeliverableDataNSE(ddDate)
			if err != nil {
				fmt.Printf("err: %s", err)
			}
			dd = utils.GetDeliverableData()
		}

		// fmt.Println("Done downloading zip file nse")
		// csvData = dataProcessor.ReadZipfile()
		//fmt.Println("Done reading zip file nse")
		utils.SaveCSV(dd, nil)
		conn.UpdateToGithub(obj)
		fmt.Println("Done uploading to github")

		e := os.Remove(config.LocalFilePath)
		if e != nil {
			log.Print(e)
		}
		e = os.Remove(config.LocalDeliverablePath)
		if e != nil {
			log.Print(e)
		}
		if dd != nil {
			csvData = conn.ReadIfFileExistsFromGit(obj)
		}
	}

	bytesBuffer := dataProcessor.FilterCsvData(csvData, obj)
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=TheCSVFileName.csv")
	w.Write(bytesBuffer.Bytes())
}

func optionsGenerator(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	symbol := req.URL.Query().Get("symbol")
	var obj config.Symboles
	fmt.Println(symbol)
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

	jsonData := conn.ReadIfFileExistsFromGitOptions(obj)

	if jsonData == nil {
		fmt.Println(obj.Date + " File not in git")
		config.GetFilePath(obj)

		dat, _ := utils.ReadJSON("./config/nse.json")
		var exchangeConfig config.ExchangeConfig
		_ = json.Unmarshal(dat, &exchangeConfig)

		req, _ := http.NewRequest("GET", fmt.Sprintf("%s?symbol=%s", config.NSEOPTIONURL, symbol), nil)
		req.Header.Add("Host", "www.nseindia.com")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != 200 {
			fmt.Println(err.Error())
			w.Header().Set("Content-Type", "application/json")
			return
		}

		defer resp.Body.Close()
		jsonData, _ = ioutil.ReadAll(resp.Body)
		fmt.Println("Done downloading zip file nse")
		conn.UpdateToGithubOptions(jsonData, obj)
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonData)
}

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	router := mux.NewRouter()
	router.HandleFunc("/", welcome).Methods("GET")
	router.HandleFunc("/getbhavcopy", csvGenerator).Methods("POST")
	router.HandleFunc("/optionchain", optionsGenerator).Methods("POST")

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

func validateInput(obj config.Symboles) bool {
	re := regexp.MustCompile(config.RegexDate)
	if !re.MatchString(obj.Date) {
		return false
	}
	return true
}
