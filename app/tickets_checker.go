package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	DanilaApostolaOffiseID string = "61" // 115
	BogdanivskaOffiseID    string = "177"
)
const (
	CheckTimeSec             = 330
	CheckBetweenDatesTimeSec = 3
	PauseAfterServerError    = 5
)

var (
	Cookie      string
	SecretToken string
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
	Rows               []Row                `json:"rows"` // actual free time slots for registration
	TRows              []any                `json:"trows"`
	FreeDatesForOffice []FreeDatesForOffice `json:"freedatesforoffice"`
}

func CheckFreeTalons(officeID, month, monthDay string) *FreetimesResponse {
	url := "https://eq.hsc.gov.ua/site/freetimes"
	data := fmt.Sprintf("office_id=%s&date_of_admission=%d-%s-%s&question_id=55&es_date=&es_time=", officeID, time.Now().Year(), month, monthDay)
	payload := strings.NewReader(data)

	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	setRequestHeaders(req)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil
	}
	defer resp.Body.Close()

	return parseResponseBody(resp)
}

func setRequestHeaders(req *http.Request) {
	req.Header.Set("authority", "eq.hsc.gov.ua")
	req.Header.Set("method", "POST")
	req.Header.Set("path", "/site/freetimes")
	req.Header.Set("accept-encoding", "gzip, deflate, br, zstd")
	req.Header.Set("accept-language", "en-US,en;q=0.9,ru;q=0.8")
	req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")

	req.Header.Set("cookie", Cookie)

	currentTime := time.Now()
	formattedDate := fmt.Sprintf("%d-%02d-%02d", currentTime.Year(), currentTime.Month(), currentTime.Day())
	req.Header.Set("referer", fmt.Sprintf("https://eq.hsc.gov.ua/site/step2?chdate=%s&question_id=55&id_es=", formattedDate))

	req.Header.Set("x-csrf-token", SecretToken)
	req.Header.Set("x-requested-with", "XMLHttpRequest")
}

func parseResponseBody(resp *http.Response) *FreetimesResponse {
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	responseBody := buf.String()

	result := new(FreetimesResponse)
	err := json.Unmarshal([]byte(responseBody), result)
	if err != nil {
		fmt.Println("Error unmarshalling:", err)
		fmt.Printf("Response. StatusCode: %d, Body: %s\n", resp.StatusCode, responseBody)
		time.Sleep(time.Second * PauseAfterServerError)
	}

	return result
}
