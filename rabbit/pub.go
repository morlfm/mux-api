package rabbit

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	application "github.com/morlfm/rest-api/application/employee"
	"github.com/morlfm/rest-api/application/model"
	"github.com/streadway/amqp"
)

var (
	emp model.Employee
)

func Publish(app *application.App) {
	// connect rabbit via amqp
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}

	emps, err := app.GetEmployees(emp)
	if err != nil {
		return
	}

	sendEmps, ersr := json.Marshal(emps)
	if ersr != nil {
		fmt.Println(ersr)
		return
	}

	// checkING  connection
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"TestQueueRain",
		false,
		false,
		false,
		false,
		nil,
	)

	fmt.Println(q)
	if err != nil {
		fmt.Println(err)
	}
	// publish a message
	err = ch.Publish(
		"",
		"TestQueueRain",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        sendEmps,
		},
	)
	if err != nil {
		fmt.Println(err)
	}

	json, err := json.Marshal(emps)
	if err != nil {
		fmt.Println(err)
	}
	file, err := os.Create("../rabbit/outputs/pub.json")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	file.Write(json)
	file.Close()
	fmt.Println("Message Published")
}