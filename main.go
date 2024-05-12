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
		time.Sleep(5 * time.Second)
		err := attemptToSendPrompt(&queue)
		if err != nil {
			if err.Error() == "image generator is not available" {
				continue
			}
			fmt.Println("Error", err)
		}
	}
}

func attemptToSendPrompt(queue *[]MensaQueue) error {
	// Ping to Check if image generator is available at 10.20.0.20:8080
	_, err := http.Get("http://10.20.0.20:8654")
	if err != nil {
		// Return new error with the error message "image generator is not available"
		return fmt.Errorf("image generator is not available")
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

	// log "New Prompt added. Total Prompts: <number of prompts in the list>
	log.Println("New Prompt added. Total Prompts:", len(*queue))
	updateHelperInHomeassistant(len(*queue))

	// return 200 OK
	w.WriteHeader(http.StatusOK)
}

func updateHelperInHomeassistant(queueLength int) error {
	// Update the helper in Home Assistant
	// homeassistant.local:8123/api/states/input_number.numbers_in_mensa_image_queue

	// Create a JSON object with the new queue length
	t := struct {
		State string `json:"state"`
	}{
		State: fmt.Sprintf("%d", queueLength),
	}

	// Send a POST request to the Home Assistant API
	// with the JSON object as the body
	// and Bearer token as the Authorization header
	t1, err := json.Marshal(t)
	if err != nil {
		log.Println("Error marshalling JSON")
		return err
	}
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://homeassistant.local:8123/api/states/input_number.numbers_in_mensa_image_queue", bytes.NewBuffer(t1))
	req.Header.Set("Authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJlODY4NDAxYTI1ZTM0MDUzODk5YTg2M2JiMmM5Y2UzMiIsImlhdCI6MTY1MTk0NzY0MywiZXhwIjoxOTY3MzA3NjQzfQ.dSfb-BfyJDmpKDZc_pLYF_6bZbNdNbVtTtHglsCGJZw")
	_, err = client.Do(req)
	if err != nil {
		log.Println("Error sending POST request")
		return err
	}

	return nil
}
