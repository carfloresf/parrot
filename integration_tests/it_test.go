// +build integration

package integration_tests

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/hellerox/parrot/api"
)

const databaseName = "parrot"
const databaseUser = "postgres"
const databasePassword = "secret"
const portAPI = "9999"

func TestMain(m *testing.M) {
	var db *sql.DB

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1).Intn(10000)

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	absDir, err := filepath.Abs(dir + "/../scripts/")
	if err != nil {
		log.Fatal(err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13",
		Name:       databaseName + fmt.Sprint(r1),
		Mounts:     []string{fmt.Sprintf("%s/init:/docker-entrypoint-initdb.d/", absDir)},
		Env:        []string{"POSTGRES_PASSWORD=" + databasePassword},
	})
	if err != nil {
		log.Fatalf("could not start resource: %s", err)
	}

	resource.Expire(480)

	databaseURL := fmt.Sprintf("postgres://"+databaseUser+":"+databasePassword+"@localhost:%s/%s?sslmode=disable", resource.GetPort("5432/tcp"), databaseName)
	log.Println(databaseURL)

	if err = pool.Retry(
		func() error {
			var err error
			db, err = sql.Open("postgres", databaseURL)
			if err != nil {
				return err
			}

			return db.Ping()
		},
	); err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	a := api.App{}

	a.Initialize(databaseURL, portAPI)

	code := m.Run()

	// Purge created Docker
	if err := pool.Purge(resource); err != nil {
		log.Errorf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestEndpoints_CreateUser(t *testing.T) {
	tests := []struct {
		name         string
		request      string
		expectedCode int
		expectedBody string
	}{
		{
			name: "Successful",
			request: `{
						"email":"carflores@gmail.com",
						"fullName":"Carlos Flores",
						"password":"uno"
						}`,
			expectedCode: http.StatusCreated,
			expectedBody: `{"status":"ok"}`,
		},
		{
			name: "Duplicate",
			request: `{
						"email":"carlos@mail.com",
						"fullName":"Carlos Flores",
						"password":"uno"
						}`,
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":"ERROR: duplicate key value violates unique constraint "user_pk" (SQLSTATE 23505)"}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var jsonStr = []byte(tt.request)
			req, err := http.NewRequest(http.MethodPost, "http://localhost:"+portAPI+"/user", bytes.NewBuffer(jsonStr))
			if err != nil {
				t.Fatal(err)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Error(err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, tt.expectedBody, string(body))
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}

func TestEndpoints_CreateOrder(t *testing.T) {
	tests := []struct {
		name         string
		request      string
		expectedCode int
		expectedBody string
		user         string
		password     string
	}{
		{
			name: "Successful",
			request: `{
						"email":"carlos@mail.com",
						"clientName":"Carlos Flores",
						"price":1234,
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
					}`,
			expectedCode: http.StatusCreated,
			expectedBody: `{"status":"ok","orderID":1}`,
			user:         "carflores@gmail.com",
			password:     "uno",
		},
		{
			name: "Unauthorized",
			request: `{
						"email":"carlos@mail.com",
						"fullName":"Carlos Flores",
						"password":"uno"
						}`,
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error": "Unauthorized"}`,
			user:         "carlos",
			password:     "carlos",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var jsonStr = []byte(tt.request)
			req, err := http.NewRequest(http.MethodPost, "http://localhost:"+portAPI+"/order", bytes.NewBuffer(jsonStr))
			if err != nil {
				t.Fatal(err)
			}

			req.SetBasicAuth(tt.user, tt.password)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Error(err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			assert.JSONEq(t, tt.expectedBody, string(body))
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}

func TestEndpoints_GenerateReport(t *testing.T) {
	tests := []struct {
		name         string
		request      string
		expectedCode int
		expectedBody string
		user         string
		password     string
	}{
		{
			name: "Successful",
			request: `{
						"from":"2015-01-28T17:41:52Z",
						"to": "2215-01-28T17:41:52Z"
					}`,
			expectedCode: http.StatusCreated,
			expectedBody: `{
							"data": [
								{
									"name": "tres",
									"totalAmount": 15111,
									"totalPrice": 302220
								},
								{
									"name": "dos",
									"totalAmount": 15,
									"totalPrice": 30
								},
								{
									"name": "uno",
									"totalAmount": 5,
									"totalPrice": 5
								}
							]
						}`,
			user:     "carflores@gmail.com",
			password: "uno",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var jsonStr = []byte(tt.request)
			req, err := http.NewRequest(http.MethodPost, "http://localhost:"+portAPI+"/report", bytes.NewBuffer(jsonStr))
			if err != nil {
				t.Fatal(err)
			}

			req.SetBasicAuth(tt.user, tt.password)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Error(err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			assert.JSONEq(t, tt.expectedBody, string(body))
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}
