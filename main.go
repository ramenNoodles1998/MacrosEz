package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func main() {
	fmt.Println("hello world")

	var temp *template.Template = template.Must(template.ParseGlob("./template/*.html"))
	var data = map[string]float64{
		"Protein": 0.0,
		"Carbs": 0.0,
		"Fat": 0.0,
		"Calories": 0.0,
		"ProteinValue": 0.0,
		"CarbsValue": 0.0,
		"FatValue": 0.0,
	}
		

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		temp.Execute(w, data)
	})


	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		protein, _ := strconv.ParseFloat(r.Form["protein"][0], 64)
		carbs, _ := strconv.ParseFloat(r.Form["carbs"][0], 64)
		fat, _ := strconv.ParseFloat(r.Form["fat"][0], 64)
		//TODO: save log here

		data = map[string]float64{
			"Protein": data["Protein"] + protein,
			"Carbs": data["Carbs"] + carbs,
			"Fat": data["Fat"] + fat,
			"Calories": data["Calories"] + (protein * 4) + (carbs * 4) + (fat * 9),
			"ProteinValue": 0.0,
			"CarbsValue": 0.0,
			"FatValue": 0.0,
		}
		temp.Execute(w, data)
	})

	http.ListenAndServe(":8080", nil)
}