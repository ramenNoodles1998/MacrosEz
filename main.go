package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Log struct{
	//uuid
	PartitionKey string
	//Date
	SortKey string
	Protein float64 
	Carbs float64
	Fat float64
}

var tableName string = "dev-macros"

func main() {
	fmt.Println("hello world")

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	},)
	svc := dynamodb.New(sess)

	var temp *template.Template = template.Must(template.ParseGlob("./template/*.html"))
	var data = map[string]interface{}{
		"Protein": 0.0,
		"Carbs": 0.0,
		"Fat": 0.0,
		"Calories": 0.0,
		"ProteinValue": 0.0,
		"CarbsValue": 0.0,
		"FatValue": 0.0,
		"Logs": getLogs(svc),
	}
	//TODO: add add all button
	//h1 logs remove if no logs
	//add saving logs and deleting.
	//70s Macros title

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		temp.Execute(w, data)
	})


	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		protein, _ := strconv.ParseFloat(r.Form["protein"][0], 64)
		carbs, _ := strconv.ParseFloat(r.Form["carbs"][0], 64)
		fat, _ := strconv.ParseFloat(r.Form["fat"][0], 64)

		saveLog(svc, protein, carbs, fat)

		data = map[string]interface{}{
			"Protein": data["Protein"].(float64) + protein,
			"Carbs": data["Carbs"].(float64) + carbs,
			"Fat": data["Fat"].(float64) + fat,
			"Calories": data["Calories"].(float64) + (protein * 4) + (carbs * 4) + (fat * 9),
			"ProteinValue": 0.0,
			"CarbsValue": 0.0,
			"FatValue": 0.0,
			"Logs": getLogs(svc),
		}
		temp.Execute(w, data)
	})

	http.ListenAndServe(":8080", nil)
}

func saveLog(svc *dynamodb.DynamoDB, protein float64, carbs float64, fat float64) {
	log := Log{
		PartitionKey: uuid.New().String(),
		SortKey: time.Now().String(),
		Protein: protein,
		Carbs: carbs,
		Fat: fat,
	}

	av, err := dynamodbattribute.MarshalMap(log)
	if err != nil {
		fmt.Printf("Got error marshalling new movie item: %s", err)
		return
	}


	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)

	if err != nil {
		fmt.Printf("Got error calling PutItem: %s", err)
		return
	}
}

func getLogs(svc *dynamodb.DynamoDB) []Log {
	//get todays logs.
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"PartitionKey": {
				S: aws.String("90049e1c-54a2-4d4b-9dfb-dd0166c7bbbb"),
			},
			"SortKey": {
				S: aws.String("2024-02-29 11:11:00.8323678 -0700 MST m=+0.005003501"),
			},
		},
	})

	if err != nil {
		fmt.Printf("Got error calling GetItem: %s", err)
		return []Log{}
	}

	if result.Item == nil {
		fmt.Printf("Could not find Logs")
		return []Log{}
	}
    
	item := Log{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	fmt.Printf("sending log back %s", item.PartitionKey)

	return []Log { item }
}