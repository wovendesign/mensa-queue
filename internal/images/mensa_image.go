package images

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mensa-queue/internal/repository"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

func SaveMensaImage(imageURL string, recipeID int32, conn *pgx.Conn, ctx context.Context) error {
	imageData, err := downloadImage(imageURL)
	if err != nil {
		return fmt.Errorf("failed to download mensa image: %w", err)
	}

	mediaID, err := uploadToPayloadCMS(imageData, imageURL)
	if err != nil {
		return fmt.Errorf("failed to upload mensa image to Payload: %w", err)
	}

	repo := repository.New(conn)
	err = repo.SetRecipeAIImage(ctx, repository.SetRecipeAIImageParams{
		ID:            recipeID,
		AiThumbnailID: mediaID,
	})
	if err != nil {
		return fmt.Errorf("failed to set recipe image: %w", err)
	}

	return nil
}

func downloadImage(imageURL string) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to GET image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("image download failed with status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image body: %w", err)
	}

	return data, nil
}

func uploadToPayloadCMS(imageData []byte, imageURL string) (int32, error) {
	buf := new(bytes.Buffer)
	bw := multipart.NewWriter(buf)

	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return 0, fmt.Errorf("failed to parse image URL: %w", err)
	}

	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), path.Base(parsedURL.Path))
	fileWriter, err := bw.CreateFormFile("file", fileName)
	if err != nil {
		return 0, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := fileWriter.Write(imageData); err != nil {
		return 0, fmt.Errorf("failed to write image data: %w", err)
	}

	alt := strings.TrimSuffix(path.Base(parsedURL.Path), path.Ext(parsedURL.Path))

	altPayloadBytes, err := json.Marshal(map[string]string{"alt": alt})
	if err != nil {
		return 0, fmt.Errorf("failed to marshal alt payload: %w", err)
	}
	altPayload := string(altPayloadBytes)
	if err := bw.WriteField("_payload", altPayload); err != nil {
		return 0, fmt.Errorf("failed to write form payload: %w", err)
	}

	if err := bw.Close(); err != nil {
		return 0, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	url := fmt.Sprintf("%s/api/media", os.Getenv("PAYLOAD_URL"))
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return 0, fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Add("Content-Type", "multipart/form-data; boundary="+bw.Boundary())

	bearer, err := getPayloadBearer()
	if err != nil {
		return 0, fmt.Errorf("failed to get payload bearer: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearer))

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to make upload request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("payload upload failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	var payloadUpload PayloadUploadResponse
	if err = json.Unmarshal(body, &payloadUpload); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return int32(payloadUpload.Doc.ID), nil
}
