package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
)

var (
	httpClient *http.Client
)

const (
	MaxIdleConnections int    = 2
	ApiEndpoint        string = "https://api.example.com/v2/info"
)

type cloudInfo struct {
	Name                     string  `json:"name"`
	Support                  string  `json:"support"`
	Build                    string  `json:"build"`
	Version                  float32 `json:"version"`
	Description              string  `json:"description"`
	AuthEndpoint             string  `json:"authorization_endpoint"`
	TokenEndpoint            string  `json:"token_endpoint"`
	ApiVersion               string  `json:"api_version"`
	LoggregatorEndpoint      string  `json:"loggregator_endpoint"`
	RoutingEndpoint          string  `json:"routing_endpoint"`
	LoggingEndpoint          string  `json:"logging_endpoint"`
	DopplerEndpoint          string  `json:"doppler_logging_endpoint"`
	MinCliVersion            string  `json:"min_cli_version"`
	MinRecommendedCliVersion string  `json:"min_recommended_cli_version"`
}

func init() {
	httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func main() {
	http.HandleFunc("/", handleInfo)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleInfo(w http.ResponseWriter, r *http.Request) {
	data, err := getInfo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	}
}

func getInfo() (cloudInfo, error) {
	log.Println("Sending request")

	response, err := httpClient.Get(ApiEndpoint)

	if err != nil && response == nil {
		log.Printf("Error sending request: %+v", err)
		return cloudInfo{}, err
	}

	log.Println("Got response")

	defer response.Body.Close()

	var info cloudInfo

	if err := json.NewDecoder(response.Body).Decode(&info); err != nil {
		return cloudInfo{}, err
	}

	return info, nil
}
