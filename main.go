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
	}
		

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		temp.Execute(w, data)
	})


	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		protein, _ := strconv.ParseFloat(r.Form["protein"][0], 64)
		carbs, _ := strconv.ParseFloat(r.Form["carbs"][0], 64)
		fat, _ := strconv.ParseFloat(r.Form["fat"][0], 64)
		data["Protein"] += protein * 4
		data["Carbs"] += carbs * 4
		data["Fat"] += fat * 9

		temp.Execute(w, data)
	})

	http.ListenAndServe(":8080", nil)
}