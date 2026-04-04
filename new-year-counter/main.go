package main

import (
	"encoding/json"
<<<<<<< Updated upstream
=======
	"fmt"
>>>>>>> Stashed changes
	"log"
	"net/http"
	"time"
)

<<<<<<< Updated upstream
type Response struct {
	DaysUntilNewYear int    `json:"days_until_new_year"`
	Message          string `json:"message"`
=======
type newYearResponse struct {
	Days int `json:"days"`
>>>>>>> Stashed changes
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

<<<<<<< Updated upstream
	response := Response{
		DaysUntilNewYear: days,
		Message:          "До нового года осталось дней",
	}

	json.NewEncoder(w).Encode(response)
=======
	response := newYearResponse {
		Days: days,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Ошибка при форматирование ответа", http.StatusInternalServerError)
		return
	}
>>>>>>> Stashed changes
}

func main() {
	http.HandleFunc("/", newYearHandler)

	port := "3000"
	log.Printf("Сервер запущен на порту %s\n", port)
	log.Printf("Перейдите на http://localhost:%s\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
