package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
)

// UploadToImgBB mengunggah gambar ke ImgBB dan mengembalikan URL gambar serta delete_url
func UploadToImgBB(apiKey string, filePath string, file multipart.File) (string, string, error) {
	url := "https://api.imgbb.com/1/upload?key=" + apiKey

	// Buat multipart writer
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Tambahkan file ke form data
	part, err := writer.CreateFormFile("image", filePath)
	if err != nil {
		return "", "", err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", "", err
	}

	writer.Close()

	// Kirim permintaan POST ke ImgBB
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	// Decode respons JSON
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	// Periksa apakah upload berhasil
	if result["status"].(float64) != 200 {
		return "", "", errors.New("Gagal mengunggah ke ImgBB")
	}

	// Ambil URL gambar dan delete_url
	data := result["data"].(map[string]interface{})
	imageURL := data["url"].(string)
	deleteURL := data["delete_url"].(string)

	return imageURL, deleteURL, nil
}

// DeleteFromImgBB menghapus gambar dari ImgBB menggunakan delete_url
func DeleteFromImgBB(deleteURL string) error {
	// Kirim permintaan DELETE ke ImgBB
	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode respons JSON
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	// Periksa apakah penghapusan berhasil
	if result["status"].(float64) != 200 {
		return errors.New("Gagal menghapus gambar dari ImgBB")
	}

	return nil
}
