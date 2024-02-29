package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type Log struct{
	Protein float64 
	Carbs float64
	Fat float64
}

func main() {
	fmt.Println("hello world")
	var logs [0]Log

	var temp *template.Template = template.Must(template.ParseGlob("./template/*.html"))
	var data = map[string]interface{}{
		"Protein": 0.0,
		"Carbs": 0.0,
		"Fat": 0.0,
		"Calories": 0.0,
		"ProteinValue": 0.0,
		"CarbsValue": 0.0,
		"FatValue": 0.0,
		"Logs": logs,
	}
	//TODO: add add all button
	//h1 logs remove if no logs
	//add saving logs and deleting.
		

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		temp.Execute(w, data)
	})


	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		protein, _ := strconv.ParseFloat(r.Form["protein"][0], 64)
		carbs, _ := strconv.ParseFloat(r.Form["carbs"][0], 64)
		fat, _ := strconv.ParseFloat(r.Form["fat"][0], 64)
		//TODO: save log here
		data = map[string]interface{}{
			"Protein": data["Protein"].(float64) + protein,
			"Carbs": data["Carbs"].(float64) + carbs,
			"Fat": data["Fat"].(float64) + fat,
			"Calories": data["Calories"].(float64) + (protein * 4) + (carbs * 4) + (fat * 9),
			"ProteinValue": 0.0,
			"CarbsValue": 0.0,
			"FatValue": 0.0,
			"Logs": getLogs(),
		}
		temp.Execute(w, data)
	})

	http.ListenAndServe(":8080", nil)
}

func getLogs() []Log {
	return []Log {{ Protein: 0, Carbs: 0, Fat: 0,}}
}