package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
	"strings"

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

	data["Logs"] = getLogs(svc)

	var temp *template.Template = template.Must(template.ParseGlob("./template/*.html"))
	//TODO: add add all button
	//h1 logs remove if no logs
	//add saving logs and deleting.
	//70s Macros title

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "index.html", data)
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		id := r.URL.Query().Get("id")

		macro, _ := strconv.ParseFloat(r.Form[id][0], 64)
		saveLog(svc, macro, id)
		var multiplier float64 = 4.0
		if(id == "fat") {
			multiplier = 9.0
		}

		var capitalizedId = strings.Title(id)
		var savedMacro = data[capitalizedId].(float64)

		data = map[string]interface{}{
			"Protein": 0.0,
			"Carbs": 0.0,
			"Fat": 0.0,
			"Calories": data["Calories"].(float64) + macro * multiplier,
			"Logs": getLogs(svc),
		}
		data[capitalizedId] = savedMacro + macro

		temp.ExecuteTemplate(w, id + ".html", data)
	})

	http.ListenAndServe(":8080", nil)
}

func saveLog(svc *dynamodb.DynamoDB, macro float64, Id string) {
	log := Log{
		PartitionKey: uuidRoman,
		SortKey: time.Now().String(),
	}
	if (Id == "protein") {
		log.Protein = macro
	} else if (Id == "carbs") {
		log.Carbs = macro
	} else {
		log.Fat = macro
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
