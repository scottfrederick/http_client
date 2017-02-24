package main

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	httpClient *http.Client
)

const (
	MaxIdleConnections int = 20
	RequestTimeout     int = 5
)

func init() {
	httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: time.Duration(RequestTimeout) * time.Second,
	}
}

func main() {
	var endPoint string = "https://api.cf.scottfrederick.io/v2/info"

	for {
		for i := 0; i < 3; i++ {
			log.Println("Sending request")
			response, err := httpClient.Get(endPoint)
			if err != nil && response == nil {
				log.Fatalf("Error sending request: %+v", err)
			} else {
				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					log.Fatalf("Error parsing response: %+v", err)
				}

				log.Println("Got response:", string(body))

				response.Body.Close()
			}
			time.Sleep(time.Duration(1) * time.Second)
		}
		time.Sleep(time.Duration(5) * time.Minute)
	}

}
