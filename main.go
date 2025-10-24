package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/harik/weather_monitor/internal/api"
	"github.com/harik/weather_monitor/internal/kafka"
	"github.com/harik/weather_monitor/internal/metrics"
)

func main() {
	// Load API key
	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENWEATHERMAP_API_KEY not set")
	}

	// Load Kafka broker URL
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		// IMPORTANT: default for Docker networking
		brokers = "kafka:9092"
	}

	// Load Kafka topic
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		// consistent with .env file
		topic = "weather-data"
	}

	log.Printf("Using Kafka brokers: %s, topic: %s", brokers, topic)

	// init Prometheus metrics
	metrics.InitMetrics()

	// init Kafka producer
	producer, err := kafka.NewProducer([]string{brokers}, topic)
	if err != nil {
		log.Fatalf("failed to create kafka producer: %v", err)
	}
	log.Println("Kafka producer initialized successfully")
	defer producer.Close()

	// Start Kafka consumer in background
	go func() {
		log.Println("Starting Kafka consumer...")
		if err := kafka.StartConsumer([]string{brokers}, topic); err != nil {
			log.Printf("consumer stopped: %v", err)
		}
	}()

	// Setup HTTP routes
	r := mux.NewRouter()

	// Weather fetch endpoint
	r.HandleFunc("/fetch", func(w http.ResponseWriter, r *http.Request) {
		zip := r.URL.Query().Get("zip")
		daysStr := r.URL.Query().Get("days")
		if zip == "" || daysStr == "" {
			http.Error(w, "zip and days required", http.StatusBadRequest)
			return
		}

		days, err := strconv.Atoi(daysStr)
		if err != nil || days <= 0 || days > 16 {
			http.Error(w, "invalid days (1-16)", http.StatusBadRequest)
			return
		}

		// Run fetch concurrently (non-blocking)
		go func(zip string, days int) {
			start := time.Now()

			data, err := api.FetchForecast(zip, days, apiKey)
			metrics.IncAPICalls()

			if err != nil {
				log.Printf("api fetch error: %v", err)
				metrics.ObserveAPILatency(time.Since(start).Seconds())
				return
			}

			metrics.ObserveAPILatency(time.Since(start).Seconds())

			if err := producer.Publish(data); err != nil {
				log.Printf("kafka publish error: %v", err)
				metrics.IncKafkaPublishFailures()
			} else {
				metrics.IncKafkaPublishSuccess()
			}
		}(zip, days)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})
	}).Methods("GET")

	// Prometheus metrics endpoint
	r.Handle("/metrics", metrics.Handler())

	// HTTP server setup
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Weather Monitor service listening on :%s", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
