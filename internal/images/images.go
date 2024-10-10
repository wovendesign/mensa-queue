package images

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type RecipeData struct {
	ID         uint
	Prompt     string
	PromptID   string
	PromptSeed int
	FileName   string
	FileData   []byte
}

type PromptResponse struct {
	PromptId string `json:"prompt_id"`
	Number   int    `json:"number"`
}

type QueueResponse struct {
	ExecInfo struct {
		QueueRemaining int `json:"queue_remaining"`
	} `json:"exec_info"`
}

type Recipes []*RecipeData

func GenerateImages(images Recipes) {
	fillQueue(images)
	checkEmptyQueue(images)
	fetchImageName(images)
}

func fetchImageName(images Recipes) {
	for _, image := range images {
		url := fmt.Sprintf("https://ai.ericwaetke.de/history/%s", image.PromptID)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("failed to create request: %w\n", err)
			continue
		}
		req.SetBasicAuth("test", "test")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("failed to fetch image: %s\n", err)
			continue
		}

		// Check the HTTP status code
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("HTTP request failed with status: %s\n", resp.Status)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("failed to read response body: %w\n", err)
			continue
		}

		var output HistoryResponse
		err = json.Unmarshal(body, &output)
		if err != nil {
			fmt.Println("failed to unmarshal response body: %w\n", err)
			continue
		}

		for _, imageHistory := range output.Outputs[image.PromptID].Images {
			if imageHistory.Type == "output" {
				image.FileName = imageHistory.FileName
				uploadToPayload(image)
				break
			}
		}
	}
}

func fillQueue(images Recipes) {
	for _, image := range images {
		image.PromptSeed = rand.Int()

		requestData := strings.Replace(promptRequest, "{{.Prompt}}", image.Prompt, 1)
		requestData = strings.Replace(requestData, "{{.Seed}}", strconv.Itoa(image.PromptSeed), 1)

		req, err := http.NewRequest("POST", "https://ai.ericwaetke.de/prompt", strings.NewReader(requestData))
		if err != nil {
			fmt.Println("failed to create request: %w\n", err)
			continue
		}
		req.SetBasicAuth("test", "test")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("failed to fetch image: %s\n", err)
			continue
		}

		// Check the HTTP status code
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("HTTP request failed with status: %s\n", resp.Status)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("failed to read response body: %w\n", err)
			continue
		}

		var promptResponse PromptResponse
		if err = json.Unmarshal(body, &promptResponse); err != nil {
			fmt.Println("failed to unmarshal response body: %w\n", err)
			continue
		}

		image.PromptID = promptResponse.PromptId

		_ = req.Body.Close()
	}
}

func checkEmptyQueue(images Recipes) {
	time.Sleep(time.Second*time.Duration(len(images)) + time.Second*20)
	var done bool

	for !done {
		time.Sleep(time.Second)

		req, err := http.NewRequest("GET", "https://ai.ericwaetke.de/prompt", nil)
		if err != nil {
			fmt.Println("failed to make request: %w\n", err)
			continue
		}

		req.SetBasicAuth("test", "test")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("failed to fetch image: %\n", err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("HTTP request failed with status: %s\n", resp.Status)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("failed to read response body: %w\n", err)
			continue
		}

		queueResponse := &QueueResponse{}
		if err = json.Unmarshal(body, queueResponse); err != nil {
			fmt.Println("failed to unmarshal response body: %w\n", err)
			continue
		}

		if queueResponse.ExecInfo.QueueRemaining == 0 {
			done = true
		}

		_ = req.Body.Close()
	}
}

func uploadToPayload(image *RecipeData) {
	req, err := http.NewRequest("GET", "http://localhost:8188/view?filename="+image.FileName, nil)
	if err != nil {
		fmt.Println("failed to make request: %w\n", err)
		return
	}

	req.SetBasicAuth("test", "test")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("failed to fetch image: %s\n", err)
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if image.FileData, err = io.ReadAll(resp.Body); err != nil {
		fmt.Println("failed to read response body: %w\n", err)
		return
	}

	saveRecipeImage(image)
}
