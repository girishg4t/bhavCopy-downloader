package config

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/moby/buildkit/session"
)

const RegexDate = `^(([0-9])|([0-2][0-9])|([3][0-1]))(JAN|FEB|MAR|APR|MAY|JUN|JULY|AUG|SEP|OCT|NOV|DEC)\d{4}$`

type Symboles struct {
	Stocks   []string `json:"Stocks"`
	Date     string   `json:"Date"`
	Exchange string   `json:"Exchange"`
	Fund     string   `json:"Fund"`
}

type ExchangeConfig struct {
	UserAgent      string `json:"userAgent"`
	Referer        string `json:"referer"`
	DeliverableUrl string `json:"deliverable_url"`
}

var GitRepo string
var GitUser string
var Email string
var filepath string
var NSEDrive string
var BSEDrive string
var GitAccessToken string
var BSEURL string
var NSEURL string
var NSEOPTIONURL string
var BSEURLAPI string
var NSEURLAPI string
var sess *session.Session
var OPTIONSDRIVE string

var (
	LocalFilePath        string
	LocalZipPath         string
	LocalDeliverablePath string
)

// GetEnvWithKey : get env value
func GetEnvWithKey(key string) string {
	return os.Getenv(key)
}

func LoadEnv() {
	if os.Getenv("Env") == "Development" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file")
			os.Exit(1)
		}
		fmt.Println("Env loaded for devlopment")
	}
	NSEDrive = GetEnvWithKey("NSE_DRIVE")
	BSEDrive = GetEnvWithKey("BSE_DRIVE")
	GitAccessToken = GetEnvWithKey("GIT_ACCESS_TOKEN")
	GitRepo = GetEnvWithKey("GIT_REPO")
	GitUser = GetEnvWithKey("GIT_USER")
	Email = GetEnvWithKey("EMAIL")
	BSEURL = GetEnvWithKey("BSE_URL")
	NSEURL = GetEnvWithKey("NSE_URL")
	BSEURLAPI = GetEnvWithKey("BSE_URL_API")
	NSEURLAPI = GetEnvWithKey("NSE_URL_API")
	OPTIONSDRIVE = GetEnvWithKey("OPTIONS_DRIVE")
	NSEOPTIONURL = GetEnvWithKey("NSE_OPTION_URL")

}

func GetSha() string {
	hasher := sha1.New()
	bv := []byte("mypassword")
	hasher.Write(bv)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

func GetFilePath(obj Symboles) {
	LocalZipPath = "./Data/" + obj.Date + ".zip"
	LocalFilePath = "./Data/" + obj.Date + ".csv"
	LocalDeliverablePath = "./Data/Deliverable" + obj.Date + ".csv"
}

var MonthMapping = map[string]string{
	"JAN": "01", "FEB": "02",
	"MAR": "03", "APR": "04", "MAY": "05",
	"JUN": "06", "JUL": "07", "AUG": "08",
	"SEP": "09", "OCT": "10", "NOV": "11",
	"DEC": "12",
}
