package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"driving-exam-tickets/app"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app.SecretToken = os.Getenv("SECRET_TOKEN")
	app.Cookie = os.Getenv("COOKIE")

	StartFreeTicketsCheck()
}

type Office struct {
	Name string
	ID   string
}

var availableOffises = []Office{
	// {Name: "DanilaApostola", ID: app.DanilaApostolaOffiseID},
	{Name: "Bogdanivska", ID: app.BogdanivskaOffiseID},
}

var availableDates []AvailableDate

func StartFreeTicketsCheck() {
	availableDates = getAvalableDates()

	if len(availableDates) == 0 {
		fmt.Println("Token renewal is required.")
		return
	}

	for {
		fmt.Println("Checking...")
		unavailableExamDaysCount := 0
		for _, office := range availableOffises {
			for _, date := range availableDates {
				if checkAndNotify(office, date, &unavailableExamDaysCount) {
					return
				}
				randCheckBetweenDatesTimeSec := time.Duration(rand.Intn(100) + app.CheckBetweenDatesTimeSec)
				time.Sleep(time.Millisecond * randCheckBetweenDatesTimeSec)

				if unavailableExamDaysCount > 3 {
					fmt.Println("unavailable exam days count > 3")
					unavailableExamDaysCount = 0
					break
				}
			}
			fmt.Println("<------------------------>")
		}

		fmt.Println("No appointment tickets available ¯\\_(ツ)_/¯")
		randCheckTimeSec := time.Duration(rand.Intn(30) + app.CheckTimeSec)
		time.Sleep(time.Second * randCheckTimeSec)
	}
}

func checkAndNotify(office Office, date AvailableDate, zeroFreeDatesDaysCount *int) bool {
	res := app.CheckFreeTalons(office.ID, date.Month, date.Day)
	if res == nil {
		return false
	}
	if len(res.Rows) > 0 {
		fmt.Printf("%s has free ticket! Date: %s, %v\n", office.Name, date.Day, time.Now().Format(time.TimeOnly))
		app.PlaySiren()
		return true
	} else if len(res.FreeDatesForOffice) > 0 {
		parsedDate := strings.Split(res.FreeDatesForOffice[0].ChDate, "-")
		if parsedDate[len(parsedDate)-1] != date.Day {
			fmt.Printf("Date for %s should be shifted\n", office.Name)
		} else {
			fmt.Printf("%s: %s %s. Total left: %d, no free tickets for now\n", office.Name, date.Day, date.Month, res.FreeDatesForOffice[0].Cnt)
		}
	} else {
		*zeroFreeDatesDaysCount++
	}
	return false
}

type AvailableDate struct {
	Month string
	Day   string
}

func getAvalableDates() []AvailableDate {
	url := "https://eq.hsc.gov.ua/site/step1?value=55"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Error creating getAvalableDates request:", err)
		return nil
	}

	// Set headers
	req.Header.Set("referer", "https://eq.hsc.gov.ua/site/step_cs?")
	req.Header.Set("cookie", app.Cookie)

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
	responseHTML := buf.String()

	return extractAvailableDates(responseHTML)
}

func extractAvailableDates(html string) []AvailableDate {
	re := regexp.MustCompile(`data-params='{"chdate":"(\d{4})-(\d{2})-(\d{2})"`)
	matches := re.FindAllStringSubmatch(html, -1)

	var availableDates []AvailableDate
	for _, match := range matches {
		if len(match) == 4 {
			availableDates = append(availableDates, AvailableDate{Month: match[2], Day: match[3]})
		}
	}
	if len(availableDates) == 0 {
		return []AvailableDate{}
	}
	return availableDates[3:] // from (today + 2 days)
}
