package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "testing"
    "time"

    "github.com/cucumber/godog"
)

type featureState struct {
    statusCode int
    responseBody []byte
    lastCustomerID int
}

func TestCustomerAPIFeature(t *testing.T) {
    state := &featureState{}

    suite := godog.TestSuite{
        ScenarioInitializer: func(s *godog.ScenarioContext) {
            s.Given(`the customer service is running`, func() error {
                resp, err := http.Get("http://localhost:8081/health")
                if err != nil {
                    return err
                }
                defer resp.Body.Close()
                if resp.StatusCode != http.StatusOK {
                    return fmt.Errorf("customer service health returned %d", resp.StatusCode)
                }
                return nil
            })

            s.When(`I POST a customer payload to /customers`, func() error {
                uniqueToken := time.Now().UnixNano() % 100000000
                payload := []byte(fmt.Sprintf(`{
                    "first_name":"Jane",
                    "last_name":"Doe",
                    "email":"jane.godog.%d@example.com",
                    "phone":"555-3333",
                    "ssn":"555-55-%d",
                    "customer_status":"active"
                }`, uniqueToken, uniqueToken))

                resp, err := http.Post("http://localhost:8081/customers", "application/json", bytes.NewBuffer(payload))
                if err != nil {
                    return err
                }
                defer resp.Body.Close()

                body, err := io.ReadAll(resp.Body)
                if err != nil {
                    return err
                }

                state.statusCode = resp.StatusCode
                state.responseBody = body

                var c struct {
                    CustomerID int `json:"customer_id"`
                }
                if err := json.Unmarshal(body, &c); err == nil {
                    state.lastCustomerID = c.CustomerID
                }

                return nil
            })

            s.Then(`the response status should be 201`, func() error {
                if state.statusCode != http.StatusCreated {
                    return fmt.Errorf("expected 201 got %d", state.statusCode)
                }
                return nil
            })

            s.Then(`the customer should contain a customer_id`, func() error {
                var c struct {
                    CustomerID int `json:"customer_id"`
                }
                if err := json.Unmarshal(state.responseBody, &c); err != nil {
                    return err
                }
                if c.CustomerID == 0 {
                    return fmt.Errorf("customer_id was missing from create response")
                }
                return nil
            })

            s.Given(`a customer exists in the customers table`, func() error {
                // Keep the scenario-level state anchored to the stable seeded row.
                // This prevents a previous POST scenario from leaving the next HTTP
                // read/update steps pointed at the wrong customer_id.
                state.lastCustomerID = 1

                resp, err := http.Get(fmt.Sprintf("http://localhost:8081/customers/%d", state.lastCustomerID))
                if err != nil {
                    return err
                }
                defer resp.Body.Close()
                if resp.StatusCode != http.StatusOK {
                    return fmt.Errorf("seed customer lookup returned %d", resp.StatusCode)
                }
                return nil
            })

            s.When(`I GET /customers/{customer_id}`, func() error {
                if state.lastCustomerID == 0 {
                    state.lastCustomerID = 1
                }
                resp, err := http.Get(fmt.Sprintf("http://localhost:8081/customers/%d", state.lastCustomerID))
                if err != nil {
                    return err
                }
                defer resp.Body.Close()

                body, err := io.ReadAll(resp.Body)
                if err != nil {
                    return err
                }

                state.statusCode = resp.StatusCode
                state.responseBody = body
                return nil
            })

            s.Then(`the response status should be 200`, func() error {
                if state.statusCode != http.StatusOK {
                    return fmt.Errorf("expected 200 got %d", state.statusCode)
                }
                return nil
            })

            s.Then(`the JSON body should return the customer record`, func() error {
                var c struct {
                    CustomerID int `json:"customer_id"`
                    FirstName string `json:"first_name"`
                    LastName string `json:"last_name"`
                }
                if err := json.Unmarshal(state.responseBody, &c); err != nil {
                    return err
                }
                if c.CustomerID == 0 || c.FirstName == "" || c.LastName == "" {
                    return fmt.Errorf("customer payload is not readable")
                }
                return nil
            })

            s.When(`I PUT a modified customer payload to /customers/{customer_id}`, func() error {
                if state.lastCustomerID == 0 {
                    state.lastCustomerID = 1
                }

                uniqueToken := time.Now().UnixNano() % 100000000
                payload := []byte(fmt.Sprintf(`{
                    "first_name":"Jane",
                    "last_name":"Doe",
                    "email":"jane.godog.updated.%d@example.com",
                    "phone":"555-3334",
                    "ssn":"123-45-6789",
                    "customer_status":"active"
                }`, uniqueToken))

                req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:8081/customers/%d", state.lastCustomerID), bytes.NewBuffer(payload))
                if err != nil {
                    return err
                }
                req.Header.Set("Content-Type", "application/json")

                resp, err := http.DefaultClient.Do(req)
                if err != nil {
                    return err
                }
                defer resp.Body.Close()

                body, err := io.ReadAll(resp.Body)
                if err != nil {
                    return err
                }

                state.statusCode = resp.StatusCode
                state.responseBody = body
                return nil
            })

            s.Then(`the JSON body should contain the updated email field`, func() error {
                var c struct {
                    Email string `json:"email"`
                }
                if err := json.Unmarshal(state.responseBody, &c); err != nil {
                    return err
                }
                if c.Email == "" {
                    return fmt.Errorf("updated email was not returned in the response")
                }
                if !bytes.Contains(state.responseBody, []byte("jane.godog.updated")) {
                    return fmt.Errorf("updated email was not returned in the response")
                }
                return nil
            })
        },
        Options: &godog.Options{
            Format: "pretty",
            Paths:  []string{"customer_api.feature"},
        },
    }

    if suite.Run() != 0 {
        t.Fatal("godog feature suite failed")
    }
}
