package main

import (
	"encoding/json"
	"log"
	"net/http"

	"tech-debt-collector/internal/feedback"
)

func main() {
	port := ":8080"

	feedbackStorage := feedback.NewMemoryStorage()
	collector := feedback.NewCollector(feedbackStorage)

	// API endpoints
	http.HandleFunc("/api/feedback", feedbackHandler(collector, feedbackStorage))
	http.HandleFunc("/api/trends", trendsHandler(collector))
	http.HandleFunc("/health", healthHandler)

	// Static files
	http.Handle("/", http.FileServer(http.Dir("web/static")))

	log.Printf("Starting feedback dashboard on %s\n", port)
	log.Printf("Open http://localhost:8080 in your browser\n")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func feedbackHandler(c *feedback.Collector, storage feedback.FeedbackStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodPost {
			var fb feedback.UserFeedback
			if err := json.NewDecoder(r.Body).Decode(&fb); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if err := c.Record(r.Context(), &fb); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(fb)
			return
		}

		// GET feedbacks
		misclassified, err := c.GetMisclassified(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(misclassified)
	}
}

func trendsHandler(c *feedback.Collector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		trends, err := c.Analyze(r.Context(), 30)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(trends)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
