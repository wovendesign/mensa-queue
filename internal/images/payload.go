package images

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"mensa-queue/internal/payload"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type PayloadLoginResponse struct {
	Message string `json:"message"`
	User    struct {
		Id        string    `json:"id"`
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
	payloadPassword := os.Getenv("PAYLOAD_API_KEY")
	jsonBody := []byte(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, payloadEmail, payloadPassword))
	bodyReader := bytes.NewReader(jsonBody)
	req, err := http.Post("localhost:3001/api/users/login", "application/json", bodyReader)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Println("failed to read response body: %w\n", err)
		return "", err
	}

	var payloadLogin PayloadLoginResponse
	if err = json.Unmarshal(body, payloadLogin); err != nil {
		fmt.Println("failed to unmarshal response body: %w\n", err)
		return "", err
	}

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

func saveRecipeImage(data *RecipeData) {
	buf := new(bytes.Buffer)
	bw := multipart.NewWriter(buf) // body writer

	p1w, _ := bw.CreateFormField("file")
	if _, err := p1w.Write(data.FileData); err != nil {
		fmt.Println("failed to write form file: %w", err)
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

	req, err := http.NewRequest("POST", "URL", buf)
	if err != nil {
		fmt.Println("failed to build request: %w", err)
		return
	}
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+bw.Boundary())
	bearer, err := getPayloadBearer()
	if err != nil {
		fmt.Println("failed to get payload bearer: %w", err)
		return
	}
	req.Header.Add("Authentication", fmt.Sprintf("Bearer %s", bearer))

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

	// TODO: Add Image ID to Recipe
	dsn := "host=127.0.0.1 user=mensauser password=postgres dbname=mensahhub port=5432 sslmode=disable TimeZone=Europe/Berlin"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.Model(&payload.Recipe{}).Where("id = ?", data.ID).Update("ai_thumbnail_id", payloadUpload.Doc.ID)

	data.FileData = []byte{}
	_ = resp.Body.Close()
}
