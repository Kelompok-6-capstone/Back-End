package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// UploadFileToCloudinary mengunggah file ke Cloudinary
func UploadFileToCloudinary(file multipart.File, fileName string) (string, string, error) {
	cloudinaryURL := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/upload", os.Getenv("CLOUDINARY_CLOUD_NAME"))
	uploadPreset := os.Getenv("CLOUDINARY_UPLOAD_PRESET")

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("upload_preset", uploadPreset)

	// Tambahkan file ke form data
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return "", "", err
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", "", err
	}
	writer.Close()

	// Kirim request
	req, err := http.NewRequest("POST", cloudinaryURL, body)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("Gagal upload ke Cloudinary")
	}
	defer resp.Body.Close()

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	return result["secure_url"].(string), result["public_id"].(string), nil
}

// DeleteFileFromCloudinary menghapus file dari Cloudinary
func DeleteFileFromCloudinary(publicID string) error {
	cloudinaryURL := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/destroy", os.Getenv("CLOUDINARY_CLOUD_NAME"))
	payload := url.Values{"public_id": {publicID}}

	req, err := http.NewRequest("POST", cloudinaryURL, strings.NewReader(payload.Encode()))
	if err != nil {
		return err
	}
	req.SetBasicAuth(os.Getenv("CLOUDINARY_API_KEY"), os.Getenv("CLOUDINARY_API_SECRET"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Gagal menghapus file di Cloudinary")
	}
	defer resp.Body.Close()

	return nil
}
