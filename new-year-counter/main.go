package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Response struct {
	DaysUntilNewYear int    `json:"days_until_new_year"`
	Message          string `json:"message"`
}

func daysUntilNewYear() int {
	now := time.Now()
	year := now.Year()
	newYear := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.Local)

	duration := newYear.Sub(now)
	days := int(duration.Hours() / 24)

	return days
}

func newYearHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	days := daysUntilNewYear()

	response := Response{
		DaysUntilNewYear: days,
		Message:          "До нового года осталось дней",
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/", newYearHandler)

	port := "3000"
	log.Printf("Сервер запущен на порту %s\n", port)
	log.Printf("Перейдите на http://localhost:%s\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
