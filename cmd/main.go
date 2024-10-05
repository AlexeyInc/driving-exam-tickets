package main

import (
	"fmt"
	"strings"
	"time"

	"driving-exam-tickets/app"
)

var availableDates = []string{"09", "10", "11", "12", "15", "16", "17", "18", "19", "22", "23", "24"}

func main() {
	StartFreeTicketsCheck()
}

type Office struct {
	Name string
	ID   string
}

var availableOffises = []Office{
	{Name: "DanilaApostola", ID: app.DanilaApostolaOffiseID},
	{Name: "Bogdanivska", ID: app.BogdanivskaOffiseID},
}

func StartFreeTicketsCheck() {
	for {
		fmt.Println("Checking...")
		for _, office := range availableOffises {
			for _, dayMonth := range availableDates {
				res := app.CheckFreeTalons(office.ID, dayMonth)
				if res == nil {
					return
				}
				if len(res.Rows) > 0 {
					fmt.Printf("%s has free ticket! Date: %s, %v\n", office.Name, dayMonth, time.Now().Format(time.TimeOnly))
					app.PlaySiren()
					return
				} else if len(res.FreeDatesForOffice) > 0 {
					date := strings.Split(res.FreeDatesForOffice[0].ChDate, "-")
					if date[len(date)-1] != dayMonth {
						fmt.Printf("Date for %s should be shifted\n", office.Name)
					}
					fmt.Printf("%s: %s %s. Total left: %d, no free tickets for now\n", office.Name, dayMonth, time.Now().Month().String(), res.FreeDatesForOffice[0].Cnt)
				}
				time.Sleep(time.Second * app.CheckBetweenDatesTimeSec)
			}
			fmt.Println("<------------------------>")
		}

		fmt.Println("No appointment tickets available ¯\\_(ツ)_/¯")

		time.Sleep(time.Second * app.CheckTimeSec)
	}
}
