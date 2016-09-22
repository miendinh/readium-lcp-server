package logging

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	LogFile        *log.Logger
	complianceMode bool
)

const (
	BASIC_FUNCTION       = "2.3.2.1"
	SUCCESS_REGISTRATION = "2.3.2.2.1"
	REJECT_REGISTRATION  = "2.3.2.2.2"
	SUCCESS_RETURN       = "2.3.2.3.1"
	REJECT_RETURN        = "2.3.2.3.2"
	SUCCESS_RENEW        = "2.3.2.4.1"
	REJECT_RENEW         = "2.3.2.4.2"
)

//Init inits log file and opens it
func Init(logPath string, cm bool) error {
	complianceMode = cm
	if complianceMode == true {
		file, err := os.OpenFile(logPath, os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			return err
		}

		LogFile = log.New(file, "", log.LUTC)

		return nil
	}

	return nil
}

//WriteToFile writes result of function execution in log file
func WriteToFile(result string, status int, testId string) {
	if complianceMode == true {
		currentTime := time.Now().UTC().Format(time.RFC3339)

		var parsedStatus string
		if status != 0 {
			parsedStatus = strconv.Itoa(status)
		}
		LogFile.Println(result + "   " + currentTime + "   " + parsedStatus + "   " + testId)
	}
}

//ReadLogs reads logs from file
func ReadLogs(logPath string) ([]string, error) {
	var lines []string
	file, err := os.OpenFile(logPath, os.O_RDONLY, 0666)
	if err == nil {
		reader := bufio.NewReader(file)
		contents, err := ioutil.ReadAll(reader)

		if err == nil {
			lines = strings.Split(string(contents), "\n")
		}
	}
	lines = lines[:len(lines)-1]
	return lines, err
}

//CountTotal summarize the data in log file
func CountTotal(lines []string) (string, error) {
	var total, success, negative int
	var result string

	for _, value := range lines {
		splitted := strings.Split(value, "   ")

		if splitted[0] == "error" {
			negative++
		} else {
			success++
		}

		total++
	}

	result = "Total count: " + strconv.Itoa(total) + "; ended successfully: " + strconv.Itoa(success) + "; ended with error: " + strconv.Itoa(negative)
	return result, nil
}
