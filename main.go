package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	// "github.com/google/uuid"
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

const tableName string = "dev-macros"
const (
    YYYYMMDD = "2006-01-02"
)
//uuid.New().String()
const uuidRoman string = "123123"

var data = map[string]interface{}{
	"Protein": 0.0,
	"Carbs": 0.0,
	"Fat": 0.0,
	"Calories": 0.0,
	"Logs": []Log{},
}

func main() {
	fmt.Println("running main")

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	},)
	svc := dynamodb.New(sess)
	
	initData(svc)

	//do this for calories and macros.

	var temp *template.Template = template.Must(template.ParseGlob("./template/*.html"))
	//TODO: add add all button
	//h1 logs remove if no logs
	//add saving logs and deleting.
	//70s Macros title

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "index.html" ,data)
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		fmt.Printf("%s\n", r.Form["fat"])
		protein, _ := strconv.ParseFloat(r.Form["protein"][0], 64)
		carbs, _ := strconv.ParseFloat(r.Form["carbs"][0], 64)
		fat, _ := strconv.ParseFloat(r.Form["fat"][0], 64)
		saveLog(svc, protein, carbs, fat) 

		data = map[string]interface{}{
			"Protein": data["Protein"].(float64) + protein,
			"Carbs": data["Carbs"].(float64) + carbs,
			"Fat": data["Fat"].(float64) + fat,
			"Calories": data["Calories"].(float64) + (protein * 4) + (carbs * 4) + (fat * 9),
			"Logs": getLogs(svc),
		}

		temp.ExecuteTemplate(w, "index.html", data)
	})

	http.ListenAndServe(":8080", nil)
}

func initData(svc *dynamodb.DynamoDB) {
	logs := getLogs(svc)
	var protein float64
	var carbs float64
	var fat float64

	for i := range logs {
		protein += logs[i].Protein
		carbs += logs[i].Carbs
		fat += logs[i].Fat
	}

	data["Protein"] = protein
	data["Carbs"] = carbs 
	data["Fat"] = fat 
	data["Calories"] = (protein * 4) + (carbs * 4) + (fat * 9)
	data["Logs"] = logs 
}

func saveLog(svc *dynamodb.DynamoDB, protein float64, carbs float64, fat float64) {
	log := Log{
		Protein: protein,
		Carbs: carbs,
		Fat: fat,
		PartitionKey: uuidRoman,
		SortKey: time.Now().String(),
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
	result, err := svc.Query(&dynamodb.QueryInput{
		TableName: aws.String(tableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"PartitionKey": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue {
					{
						S: aws.String(uuidRoman),
					},
				},
			},
			"SortKey": {
				ComparisonOperator: aws.String("BEGINS_WITH"),
				AttributeValueList: []*dynamodb.AttributeValue {
					{
						S: aws.String(time.Now().UTC().Format(YYYYMMDD)),
					},
				},
			},
		},
	})

	if err != nil {
		fmt.Printf("Got error calling GetItem: %s", err)
		return []Log{}
	}

	if result.Items == nil {
		fmt.Printf("Could not find Logs")
		return []Log{}
	}
    
	logs := []Log{}

	for _, item := range result.Items {
		log := Log{}
		err = dynamodbattribute.UnmarshalMap(item, &log)
		logs = append(logs, log)
		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}
	}

	return logs 
}
