package images

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"io"
	"mensa-queue/internal/repository"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type PayloadLoginResponse struct {
	Message string `json:"message"`
	User    struct {
		Id        uint      `json:"id"`
		Email     string    `json:"email"`
		Verified  bool      `json:"_verified"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	} `json:"user"`
	Token string `json:"token"`
	Exp   int    `json:"exp"`
}

func getPayloadBearer() (string, error) {
	payloadEmail := os.Getenv("PAYLOAD_EMAIL")
	payloadPassword := os.Getenv("PAYLOAD_PASSWORD")
	jsonBody := []byte(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, payloadEmail, payloadPassword))
	bodyReader := bytes.NewReader(jsonBody)
	url := fmt.Sprintf("%s/api/users/login", os.Getenv("PAYLOAD_URL"))
	req, err := http.Post(url, "application/json", bodyReader)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Println("failed to read response body: %w\n", err)
		return "", err
	}

	var payloadLogin PayloadLoginResponse
	if err = json.Unmarshal(body, &payloadLogin); err != nil {
		fmt.Printf("failed to unmarshal response body: %v\n", err)
		return "", err
	}

	fmt.Printf("payload login response: %+v\n", payloadLogin)

	return payloadLogin.Token, nil
}

type PayloadUploadResponse struct {
	Doc struct {
		ID           uint        `json:"id"`
		Alt          string      `json:"alt"`
		UpdatedAt    time.Time   `json:"updatedAt"`
		CreatedAt    time.Time   `json:"createdAt"`
		Url          string      `json:"url"`
		ThumbnailURL interface{} `json:"thumbnailURL"`
		Filename     string      `json:"filename"`
		MimeType     string      `json:"mimeType"`
		Filesize     int         `json:"filesize"`
		Width        int         `json:"width"`
		Height       int         `json:"height"`
		FocalX       int         `json:"focalX"`
		FocalY       int         `json:"focalY"`
	} `json:"doc"`
	Message string `json:"message"`
}

func saveRecipeImage(data *RecipeData, ctx context.Context) {
	buf := new(bytes.Buffer)
	bw := multipart.NewWriter(buf) // body writer

	// Create the file field in the multipart form
	fileName := fmt.Sprintf("%d_%d.png", time.Now().Unix(), data.ID)
	fileWriter, err := bw.CreateFormFile("file", fileName) // Ensure "file" is the correct field name
	if err != nil {
		fmt.Printf("failed to create form file: %v\n", err)
		return
	}

	// Check if file data is empty
	if len(data.FileData) == 0 {
		fmt.Println("File data is empty, nothing to upload!")
		return
	}

	// Write the file data into the form
	if _, err := fileWriter.Write(data.FileData); err != nil {
		fmt.Printf("failed to write file data: %v\n", err)
		return
	}

	if err := bw.WriteField("_payload", `{"alt": "AI Image"}`); err != nil {
		fmt.Println("failed to write form payload: %w", err)
		return
	}

	if err := bw.Close(); err != nil {
		fmt.Println("failed to close writer: %w", err)
		return
	}

	url := fmt.Sprintf("%s/api/media", os.Getenv("PAYLOAD_URL"))
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		fmt.Println("failed to build request: %w", err)
		return
	}
	//req.Header.Add("Content-Type", "multipart/form-data")
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+bw.Boundary())
	bearer, err := getPayloadBearer()
	if err != nil {
		fmt.Println("failed to get payload bearer: %w", err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearer))

	fmt.Printf("Auth Header: %+v\n", req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("failed to make request: %w", err)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("failed to read response body: %w\n", err)
		return
	}

	payloadUpload := PayloadUploadResponse{}
	if err = json.Unmarshal(body, &payloadUpload); err != nil {
		fmt.Println("failed to unmarshal response body: %w\n", err)
		return
	}

	fmt.Printf("Payload Upload: %+v\n", payloadUpload)

	// TODO: Add Image ID to Recipe
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		panic(err)
	}
	repo := repository.New(conn)
	err = repo.SetRecipeAIImage(ctx, repository.SetRecipeAIImageParams{
		ID:            *data.ID,
		AiThumbnailID: int32(payloadUpload.Doc.ID),
	})
	if err != nil {
		fmt.Printf("failed to set recipe AI image: %v\n", err)
	}

	data.FileData = []byte{}
	_ = resp.Body.Close()
}
