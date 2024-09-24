package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type ClientObj struct {
	Name        string
	Guards      int
	Hours       int
	Occurrences []int
	Rate        int
	Overtime    int
	LateHours   int
	Result      float64
}

func main() {
	http.HandleFunc("/", formHandler)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (c ClientObj) TotalPay() float64 {

	totalHours := 0
	overtimeHours := c.Overtime * c.Hours * c.Guards
	for _, val := range c.Occurrences {
		totalHours += c.Hours * val * c.Guards
	}

	//if total hours is not 0
	if totalHours > 0 && c.LateHours > 0 && overtimeHours > 0 {
		//subtract late hours
		totalHours -= c.LateHours
		//subtract over time
		totalHours -= overtimeHours
	}

	totalPay := float64(totalHours*c.Rate + overtimeHours*(c.Rate*2))
	totalPay *= 1.13
	return totalPay

}

func formHandler(w http.ResponseWriter, r *http.Request) {

	var obj ClientObj

	if r.Method == http.MethodPost {

		//parse input data
		obj.Name = r.FormValue("Client")

		guardsStr := r.FormValue("Guards")
		obj.Guards, _ = strconv.Atoi(guardsStr)

		hoursStr := r.FormValue("Hours")
		obj.Hours, _ = strconv.Atoi(hoursStr)

		occStr := r.FormValue("Occurrences")
		//parse occurrences string
		occStrArray := strings.Split(occStr, ",")
		occIntArray := make([]int, len(occStrArray))
		for i, val := range occStrArray {
			valInt, _ := strconv.Atoi(val)
			occIntArray[i] = valInt
		}
		obj.Occurrences = occIntArray

		rateStr := r.FormValue("Rate")
		obj.Rate, _ = strconv.Atoi(rateStr)

		overtimeStr := r.FormValue("Overtime")
		obj.Overtime, _ = strconv.Atoi(overtimeStr)

		lateStr := r.FormValue("LateHours")
		obj.LateHours, _ = strconv.Atoi(lateStr)

		result := obj.TotalPay()
		result = math.Round(result*100) / 100
		obj.Result = result

	}

	t, err := template.ParseFiles("homepage.html")
	if err != nil {
		fmt.Println("Error parsing html file: ", err)
	}

	err = t.Execute(w, obj)
	if err != nil {
		fmt.Println("Error executing sending html file: ", err)
	}

}
