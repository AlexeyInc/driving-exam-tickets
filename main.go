package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	cookieConst = "WEBCHSID2=9b4p1qlk3f0vki60qfoi28gd0m; _identity=f0d02a9c35f47e20d8ee0d1daad22ecc82ca40df1b5ceaf2c40eb71022104b90a%3A2%3A%7Bi%3A0%3Bs%3A9%3A%22_identity%22%3Bi%3A1%3Bs%3A20%3A%22%5B2094253%2Cnull%2C28800%5D%22%3B%7D; _csrf=4b9d9cb422a1018f27559944105ac1e6973d022a77ef4c323e59e8dde65040a9a%3A2%3A%7Bi%3A0%3Bs%3A5%3A%22_csrf%22%3Bi%3A1%3Bs%3A32%3A%224HyYYo3jBMmKI0oUucRCM_CClnbmxSFz%22%3B%7D"
	tokenConst  = "1Pc0p-tFCUjNStwCeGZmpckv5c7jCNSSOjA8Ox2RwzPgv03-sio6Io8HsUkxVgnwvEy3ja5Xl9FWXl5WZcKFSQ=="
)

const (
	DanilaApostolaOffiseID string = "61"
	BogdanivskaOffiseID    string = "177"
)
const checkTimeSec = 10

var availableDates = []string{"09", "10", "11", "12", "15", "16", "17", "18", "19"}

func main() {

	for i := 0; i < 100; i++ {
		fmt.Println("Checking...")

		for _, dayMonth := range availableDates {
			res := CheckFreeTalons(
				DanilaApostolaOffiseID, dayMonth,
			)
			if res != nil {
				fmt.Printf("DanilaApostola has free talon! %v\n", time.Now())
			} else {
				fmt.Printf("DanilaApostola %v - empty\n", dayMonth)
			}
		}
		time.Sleep(time.Second)

		for _, dayMonth := range availableDates {
			res := CheckFreeTalons(
				DanilaApostolaOffiseID, dayMonth,
			)
			if res != nil {
				fmt.Printf("Bogdanivska has free talon! %v\n", time.Now())
			} else {
				fmt.Printf("Bogdanivska %v - empty\n", dayMonth)
			}
		}

		fmt.Println("No free talon")

		time.Sleep(time.Second * checkTimeSec)
	}
}

// Define the structure for "freedatesforoffice" array elements
type FreeDatesForOffice struct {
	Cnt    int    `json:"cnt"`
	ChDate string `json:"chdate"`
}

// Define the main structure
type Response struct {
	Rows               []interface{}        `json:"rows"`
	TRows              []interface{}        `json:"trows"`
	FreeDatesForOffice []FreeDatesForOffice `json:"freedatesforoffice"`
}

func CheckFreeTalons(officeID, monthDay string) []interface{} {
	// The URL for the POST request
	url := "https://eq.hsc.gov.ua/site/freetimes"

	// The payload (replace with your actual form data)
	data := "office_id=" + officeID + "&date_of_admission=2024-10-" + monthDay + "&question_id=55&es_date=&es_time=" // 61
	payload := strings.NewReader(data)

	// Create a new POST request
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	// Set headers as per your request
	req.Header.Set("authority", "eq.hsc.gov.ua")
	req.Header.Set("method", "POST")
	req.Header.Set("path", "/site/freetimes")
	req.Header.Set("scheme", "https")
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-encoding", "gzip, deflate, br, zstd")
	req.Header.Set("accept-language", "en-US,en;q=0.9,ru;q=0.8")
	// req.Header.Set("content-length", "75") // Content length should be dynamically set by http.Client, can remove this
	req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("cookie", cookieConst)
	req.Header.Set("origin", "https://eq.hsc.gov.ua")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://eq.hsc.gov.ua/site/step2?chdate=2024-10-04&question_id=55&id_es=")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="129", "Not=A?Brand";v="8", "Chromium";v="129"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")
	req.Header.Set("x-csrf-token", tokenConst)
	req.Header.Set("x-requested-with", "XMLHttpRequest")

	// Create a client and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil
	}
	defer resp.Body.Close()

	// Read and print the response
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	responseBody := buf.String()

	var responseObj Response
	err = json.Unmarshal([]byte(responseBody), &responseObj)
	if err != nil {
		fmt.Println("Error unmarshalling:", err)
		return nil
	}

	if len(responseObj.Rows) == 0 {
		return nil
	}
	return responseObj.Rows
}
