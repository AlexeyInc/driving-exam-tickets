package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

const (
	cookieConst = "WEBCHSID2=9b4p1qlk3f0vki60qfoi28gd0m; _identity=f0d02a9c35f47e20d8ee0d1daad22ecc82ca40df1b5ceaf2c40eb71022104b90a%3A2%3A%7Bi%3A0%3Bs%3A9%3A%22_identity%22%3Bi%3A1%3Bs%3A20%3A%22%5B2094253%2Cnull%2C28800%5D%22%3B%7D; _csrf=4b9d9cb422a1018f27559944105ac1e6973d022a77ef4c323e59e8dde65040a9a%3A2%3A%7Bi%3A0%3Bs%3A5%3A%22_csrf%22%3Bi%3A1%3Bs%3A32%3A%224HyYYo3jBMmKI0oUucRCM_CClnbmxSFz%22%3B%7D"
	tokenConst  = "1Pc0p-tFCUjNStwCeGZmpckv5c7jCNSSOjA8Ox2RwzPgv03-sio6Io8HsUkxVgnwvEy3ja5Xl9FWXl5WZcKFSQ=="
)

type Office struct {
	Name string
	ID   string
}

const (
	DanilaApostolaOffiseID string = "61" // 115
	BogdanivskaOffiseID    string = "177"
)
const (
	checkTimeSec             = 30
	checkBetweenDatesTimeSec = 2
)

var availableOffises = []Office{
	{Name: "DanilaApostola", ID: DanilaApostolaOffiseID},
	{Name: "Bogdanivska", ID: BogdanivskaOffiseID},
}
var availableDates = []string{"09", "10", "11", "12", "15", "16", "17", "18", "19", "22", "23", "24"}

func main() {
	StartFreeTicketsCheck()
}

func StartFreeTicketsCheck() {
	for {
		fmt.Println("Checking...")
		for _, office := range availableOffises {
			for _, dayMonth := range availableDates {
				res := CheckFreeTalons(office.ID, dayMonth)
				if len(res.Rows) > 0 {
					fmt.Printf("%s has free talon! Date: %s, %v\n", office.Name, dayMonth, time.Now().Format(time.TimeOnly))
					PlaySiren()
					return
				} else {
					date := strings.Split(res.FreeDatesForOffice[0].ChDate, "-")
					if date[len(date)-1] != dayMonth {
						fmt.Printf("Date for %s should be shifted\n", office.Name)
					}
					fmt.Printf("%s: %s %s. Total: %d and no free tickets\n", office.Name, dayMonth, time.Now().Month().String(), res.FreeDatesForOffice[0].Cnt)
				}
				time.Sleep(time.Second * checkBetweenDatesTimeSec)
			}
			fmt.Println("<------------------------>")
		}

		fmt.Println("No appointment tickets available ¯\\_(ツ)_/¯")

		time.Sleep(time.Second * checkTimeSec)
	}
}

type FreeDatesForOffice struct {
	Cnt    int    `json:"cnt"`
	ChDate string `json:"chdate"`
}

type Row struct {
	ID     int    `json:"id"`
	ChTime string `json:"chtime"`
}

type FreetimesResponse struct {
	Rows               []Row                `json:"rows"`
	TRows              []interface{}        `json:"trows"`
	FreeDatesForOffice []FreeDatesForOffice `json:"freedatesforoffice"`
}

func CheckFreeTalons(officeID, monthDay string) *FreetimesResponse {
	url := "https://eq.hsc.gov.ua/site/freetimes"

	data := "office_id=" + officeID + "&date_of_admission=2024-10-" + monthDay + "&question_id=55&es_date=&es_time=" // 61
	payload := strings.NewReader(data)

	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	// Set headers
	req.Header.Set("authority", "eq.hsc.gov.ua")
	req.Header.Set("method", "POST")
	req.Header.Set("path", "/site/freetimes")
	req.Header.Set("scheme", "https")
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-encoding", "gzip, deflate, br, zstd")
	req.Header.Set("accept-language", "en-US,en;q=0.9,ru;q=0.8")
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

	// Execute the request
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

	result := new(FreetimesResponse)
	err = json.Unmarshal([]byte(responseBody), result)
	if err != nil {
		fmt.Println("Error unmarshalling:", err)
		return nil
	}

	return result
}

func PlaySiren() {
	playCount := 3
	for i := 0; i < playCount; i++ {
		// Play the sound using afplay (macOS)
		cmd := exec.Command("afplay", "warning.wav")
		err := cmd.Start()
		if err != nil {
			fmt.Println("Error playing sound:", err)
			return
		}
		cmd.Wait()

		time.Sleep(time.Millisecond * 5)
	}
	fmt.Println("Good luck :)")
}
