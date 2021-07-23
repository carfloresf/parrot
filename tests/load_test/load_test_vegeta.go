package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/nouney/randomstring"
	log "github.com/sirupsen/logrus"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

var mailRand string
var i int

func NewTargeterOrder() vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		tgt.Method = "POST"
		tgt.URL = "http://localhost:8080/order"
		fn, _ := tgt.Request()
		fn.SetBasicAuth(mailRand, "uno")

		tgt.Header = fn.Header
		tgt.Body = []byte(
			fmt.Sprintf(`{
										"email":"%s",
										"clientName":"Alfredo Flores",
										"price":0,
										"products":[{
											"name":"uno",
											"price": 1,
											"description":"uno uno",
											"amount":5
										},{
											"name":"dos",
											"price": 2,
											"description":"uno uno",
											"amount":15
										},
										{
											"name":"tres",
											"price": 20,
											"description":"uno uno",
											"amount":15111
										}]
									}`, mailRand))

		i++

		return nil
	}
}

func NewTargeterUser() vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		mailRand = strings.ToLower(randomstring.Generate(8)) + "@gmail.com"

		tgt.Method = "POST"
		tgt.URL = "http://localhost:8080/user"
		tgt.Body = []byte(fmt.Sprintf(`{"email":"%s","fullName":"Carlos Flores","password":"uno"}`, mailRand))

		return nil
	}
}

func main() {
	rate := vegeta.Rate{
		Freq: 200,
		Per:  time.Second,
	}

	duration := 10 * time.Second

	attackerCreateUser := vegeta.NewAttacker()
	targeterUser := NewTargeterUser()

	var metrics, metricsOrder vegeta.Metrics

	for res := range attackerCreateUser.Attack(targeterUser, rate, duration, "") {
		metrics.Add(res)
	}

	attackerCreateOrder := vegeta.NewAttacker()
	targeterOrder := NewTargeterOrder()

	for resOrder := range attackerCreateOrder.Attack(targeterOrder, rate, duration, "") {
		metricsOrder.Add(resOrder)
	}

	metrics.Close()
	metricsOrder.Close()

	log.Println("------ creación de cuentas --------")
	log.Printf("99th percentile: %s\n", metrics.Latencies.P99)
	log.Printf("Rate: %f req/s \n", metrics.Rate)
	log.Printf("statusCodes: %+v \n", metrics.StatusCodes)

	log.Println("------ creación de órdenes --------")
	log.Printf("99th percentile: %s\n", metricsOrder.Latencies.P99)
	log.Printf("Rate: %f req/s \n", metricsOrder.Rate)
	log.Printf("statusCodes: %+v \n", metricsOrder.StatusCodes)
}
