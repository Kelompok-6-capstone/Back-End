package helper

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/nedpals/supabase-go"
)

// GetSupabaseClient - Membuat client Supabase
func GetSupabaseClient() *supabase.Client {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_API_KEY")
	return supabase.CreateClient(supabaseURL, supabaseKey)
}

// UploadToSupabase - Mengunggah gambar ke Supabase Storage
func UploadToSupabase(file multipart.File, fileName string) (string, error) {
	client := GetSupabaseClient()
	bucketName := os.Getenv("SUPABASE_BUCKET_NAME")

	// Upload file ke Supabase Storage
	response := client.Storage.From(bucketName).Upload(fileName, file)
	if response.Message != "" { // Jika ada pesan error dari response
		return "", errors.New("gagal mengunggah file ke Supabase: " + response.Message)
	}

	// Buat URL gambar berdasarkan key
	fileURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", os.Getenv("SUPABASE_URL"), bucketName, response.Key)
	return fileURL, nil
}

// DeleteFromSupabase - Menghapus gambar dari Supabase Storage
func DeleteFromSupabase(fileName string) error {
	client := GetSupabaseClient()
	bucketName := os.Getenv("SUPABASE_BUCKET_NAME")

	// Hapus file dari Supabase Storage
	response := client.Storage.From(bucketName).Remove([]string{fileName})
	if response.Message != "" { // Jika ada pesan error dari response
		return errors.New("gagal menghapus file: " + response.Message)
	}

	return nil
}
