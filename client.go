<<<<<<< HEAD
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// WeatherData is the structure sent to Kafka (simplified)
type WeatherData struct {
	Zip       string    `json:"zip"`
	Timestamp time.Time `json:"timestamp"`
	Forecast  any       `json:"forecast"`
}

func FetchForecast(zip string, days int, apiKey string) (*WeatherData, error) {
	// Using the 5 day/3 hour or One Call 3.0 is possible; here is a simple example using OpenWeatherMap 2.5 forecast
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?zip=%s,us&cnt=%d&appid=%s&units=metric", zip, days*8, apiKey)
	// Note: cnt param is count of data points; adjust for your preferred endpoint. This is an example.
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("api error: %s: %s", resp.Status, string(body))
	}
	var body any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	return &WeatherData{
		Zip:       zip,
		Timestamp: time.Now().UTC(),
		Forecast:  body,
	}, nil
}
=======
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// WeatherData is the structure sent to Kafka (simplified)
type WeatherData struct {
	Zip       string    `json:"zip"`
	Timestamp time.Time `json:"timestamp"`
	Forecast  any       `json:"forecast"`
}

func FetchForecast(zip string, days int, apiKey string) (*WeatherData, error) {
	// Using the 5 day/3 hour or One Call 3.0 is possible; here is a simple example using OpenWeatherMap 2.5 forecast
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?zip=%s,us&cnt=%d&appid=%s&units=metric", zip, days*8, apiKey)
	// Note: cnt param is count of data points; adjust for your preferred endpoint. This is an example.
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("api error: %s: %s", resp.Status, string(body))
	}
	var body any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	return &WeatherData{
		Zip:       zip,
		Timestamp: time.Now().UTC(),
		Forecast:  body,
	}, nil
}
>>>>>>> 428f5cd2f762435819e1f11314a19742522374ff
