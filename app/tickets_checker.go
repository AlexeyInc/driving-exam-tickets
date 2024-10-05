package app

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
	CookieConst = "WEBCHSID2=0upil1bp6ehf2u2hf3e08c23nt; _identity=f0d02a9c35f47e20d8ee0d1daad22ecc82ca40df1b5ceaf2c40eb71022104b90a%3A2%3A%7Bi%3A0%3Bs%3A9%3A%22_identity%22%3Bi%3A1%3Bs%3A20%3A%22%5B2094253%2Cnull%2C28800%5D%22%3B%7D; _csrf=6e2cf44c914afd9c386d25f05809e567f1a46ac3008011e6103664c68dd57dd5a%3A2%3A%7Bi%3A0%3Bs%3A5%3A%22_csrf%22%3Bi%3A1%3Bs%3A32%3A%22F3rPz1KHpm0qyfubp5Ez818WkH4m-KK_%22%3B%7D"
	TokenConst  = "PDILG--VtBvoxDRTmwo9I-BqHxpLpuwbePtuSyrkMJ56AXlLlaT_U5ipBCLibEhBkF9aYHOX1EwTs1omB697wQ=="
)
const (
	DanilaApostolaOffiseID string = "61" // 115
	BogdanivskaOffiseID    string = "177"
)
const (
	CheckTimeSec             = 30
	CheckBetweenDatesTimeSec = 1
	PauseAfterServerError    = 10
)

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
	req.Header.Set("cookie", CookieConst)
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
	req.Header.Set("x-csrf-token", TokenConst)
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
		fmt.Printf("Response. StatusCode: %d, Body: %s\n", resp.StatusCode, responseBody)
		time.Sleep(time.Second * PauseAfterServerError)
	}

	return result
}

func PlaySiren() {
	playCount := 3
	for i := 0; i < playCount; i++ {
		// Play the sound using afplay (macOS)
		cmd := exec.Command("afplay", "../media/warning.wav")
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
