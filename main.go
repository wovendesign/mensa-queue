package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type MensaQueue struct {
	FoodID    *int    `json:"food_id"`
	FoodTitle *string `json:"prompt"`
}

func main() {
	queue := []MensaQueue{}

	router := http.NewServeMux()
	router.HandleFunc("POST /prompt", func(w http.ResponseWriter, req *http.Request) {
		getRoot(w, req, &queue)
	})

	server := http.Server{
		Addr:    ":3333",
		Handler: router,
	}

	go func() {
		fmt.Println("Server is running on port 3333")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()
	// server.ListenAndServe()

	for {
		time.Sleep(5 * time.Minute)
		err := attemptToSendPrompt(&queue)
		if err != nil {
			fmt.Println("Error", err)
		}
	}
}

func attemptToSendPrompt(queue *[]MensaQueue) error {
	// Ping to Check if image generator is available at 10.20.0.20:8080
	_, err := http.Get("http://10.20.0.20:8654")
	if err != nil {
		// Return new error with the error message "image generator is not available"
		return fmt.Errorf("image generator is not available: %w", err)
	}

	for _, queue_item := range *queue {
		// Send prompt to image generator
		t, err := json.Marshal(queue_item)
		if err != nil {
			return err
		}
		log.Println(string(t))
		_, err = http.Post("http://10.20.0.20:8654/prompt", "application/json", bytes.NewBuffer(t))
		if err != nil {
			return err
		}
		// remove the prompt from the queue
		*queue = (*queue)[1:]
	}

	return nil
}

func getRoot(w http.ResponseWriter, req *http.Request, queue *[]MensaQueue) {
	bodyBytes, err1 := io.ReadAll(req.Body)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return
	}
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	decoder := json.NewDecoder(req.Body)
	var t MensaQueue
	err := decoder.Decode(&t)
	if err != nil {
		// bad JSON or unrecognized json field
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if t.FoodTitle == nil || *t.FoodTitle == "" || t.FoodID == nil {
		http.Error(w, "missing field 'prompt' or 'food_id' from JSON object", http.StatusBadRequest)
		return
	}

	// optional extra check
	if decoder.More() {
		http.Error(w, "extraneous data after JSON object", http.StatusBadRequest)
		return
	}

	// check if the if is already in the list
	for _, queue_item := range *queue {
		if *queue_item.FoodID == *t.FoodID {
			http.Error(w, "food ID already in the list", http.StatusBadRequest)
			return
		}
	}

	// add the new title to the list
	*queue = append(*queue, t)
	// return 200 OK
	w.WriteHeader(http.StatusOK)
}
