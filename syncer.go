package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/XiaoMengXinX/ProfileStatusSyncer/gh"
	log "github.com/sirupsen/logrus"
)

type logFormatter struct{}

type Quote struct {
	Emoji string `json:"emoji"`
	Text  string `json:"text"`
}

// Format is a formatter for logs
func (s *logFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006/01/02 15:04:05")
	var msg string
	msg = fmt.Sprintf("%s [%s] %s (%s:%d)\n", timestamp, strings.ToUpper(entry.Level.String()), entry.Message, path.Base(entry.Caller.File), entry.Caller.Line)
	return []byte(msg), nil
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors:          false,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})
	log.SetFormatter(new(logFormatter))
	log.SetReportCaller(true)
}

var githubToken = os.Getenv("GITHUB_TOKEN")

func main() {
	rand.Seed(time.Now().UnixNano())
	var err error
	defer func() {
		if err != nil {
			log.Errorln(err)
		}
	}()
	jsonFile, err := os.Open("quotes.json")
	if err != nil {
		log.Println(err)
	}
	log.Println("Successfully Opened quotes.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var quotes map[string]Quote
	json.Unmarshal(byteValue, &quotes)
	if githubToken == "" {
		err = fmt.Errorf("GITHUB_TOKEN is empty ")
		return
	}

	c, err := gh.NewClient(githubToken)
	if err != nil {
		return
	}
	quote := quotes[strconv.Itoa(rand.Intn(len(quotes)))]
	err = c.SetUserStatus(quote.Emoji, quote.Text)
	if err == nil {
		log.Println("Sync or update profile status successfully")
	}
}
