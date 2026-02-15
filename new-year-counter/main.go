package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func daysUntilNewYear() int {
	now := time.Now()
	year := now.Year()
	newYear := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.Local)
	
	duration := newYear.Sub(now)
	days := int(duration.Hours() / 24)
	
	return days
}

func newYearHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	
	days := daysUntilNewYear()
	
	fmt.Fprintf(w, "До нового года: %d", days)
}

func main() {
	http.HandleFunc("/", newYearHandler)
	
	port := "3000"
	fmt.Printf("Сервер запущен на порту %s\n", port)
	fmt.Printf("Перейдите на http://localhost:%s\n", port)
	
	log.Fatal(http.ListenAndServe(":"+port, nil))
}