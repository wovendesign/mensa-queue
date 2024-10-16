package images

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type RecipeData struct {
	ID         *int32
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

func GenerateImages(images Recipes, ctx context.Context) {
	fillQueue(images)
	checkEmptyQueue(images)
	fetchImageName(images, ctx)
}

func fetchImageName(images Recipes, ctx context.Context) {
	for _, image := range images {
		fmt.Printf("Image: %+v\n", image)
		url := fmt.Sprintf("%s/history/%s", os.Getenv("COMFYUI_URL"), image.PromptID)
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
			fmt.Printf("HTTP request (%s) failed with status: %s\n", url, resp.Status)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("failed to read response body: %w\n", err)
			continue
		}

		data := make(map[string]interface{})

		// Unmarshal the JSON into the map
		err = json.Unmarshal([]byte(body), &data)
		if err != nil {
			log.Fatalf("Error unmarshalling JSON: %v", err)
		}

		// Access the "outputs" block for the specific UUID
		if uuidData, ok := data[image.PromptID]; ok {
			// Type assert to map for the UUID's value
			if uuidMap, ok := uuidData.(map[string]interface{}); ok {
				// Navigate to the "outputs" field
				if outputs, ok := uuidMap["outputs"].(map[string]interface{}); ok {
					// Now, check the output type "27" for "images"
					if output27, ok := outputs["27"].(map[string]interface{}); ok {
						if images, ok := output27["images"].([]interface{}); ok {
							// Loop over images (in case there are multiple)
							for _, imageData := range images {
								if imageInLoop, ok := imageData.(map[string]interface{}); ok {
									if filename, ok := imageInLoop["filename"].(string); ok {
										// Print the extracted filename
										image.FileName = filename
										uploadToPayload(image, ctx)
										break
									}
								}
							}
						}
					}
				}
			}
		} else {
			fmt.Printf("UUID %s not found in the data.\n", image.PromptID)
		}
	}
}

func fillQueue(images Recipes) {
	for _, image := range images {
		image.PromptSeed = rand.Int()

		requestData := strings.Replace(promptRequest, "{{.Prompt}}", image.Prompt, 1)
		requestData = strings.Replace(requestData, "{{.Seed}}", strconv.Itoa(image.PromptSeed), 1)

		url := fmt.Sprintf("%s/prompt", os.Getenv("COMFYUI_URL"))
		req, err := http.NewRequest("POST", url, strings.NewReader(requestData))
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
			fmt.Printf("HTTP request (%s) failed with status: %s\n", url, resp.Status)
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
	time.Sleep(time.Second*time.Duration(len(images)) + time.Second*0)
	var done bool

	for !done {
		time.Sleep(time.Second)

		url := fmt.Sprintf("%s/prompt", os.Getenv("COMFYUI_URL"))
		req, err := http.NewRequest("GET", url, nil)
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
			fmt.Printf("HTTP request (%s) failed with status: %s\n", url, resp.Status)
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

		_ = resp.Body.Close()
	}
}

func uploadToPayload(image *RecipeData, ctx context.Context) {
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

	saveRecipeImage(image, ctx)
}
