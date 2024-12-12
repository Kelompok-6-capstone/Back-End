package helper

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func ResponseAI(ctx context.Context, question string) (string, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("AI_API_KEY")))
	if err != nil {
		log.Printf("Error creating Gemini client: %v", err)
		return "", err
	}

	modelAI := client.GenerativeModel("gemini-pro")
	modelAI.SetTemperature(0) // Respons lebih deterministik

	resp, err := modelAI.GenerateContent(ctx, genai.Text(question))
	if err != nil {
		log.Printf("Error generating AI response: %v", err)
		return "", err
	}

	answer := resp.Candidates[0].Content.Parts[0]
	answerString := fmt.Sprintf("%v", answer)

	// Bersihkan simbol tambahan dari teks
	answerString = strings.ReplaceAll(answerString, "*", "")
	answerString = strings.ReplaceAll(answerString, "**", "")
	answerString = strings.ReplaceAll(answerString, "\n\n", " -")

	// Batasi panjang respons berdasarkan jumlah kalimat
	maxSentences := 5 // Jumlah maksimum kalimat
	sentences := strings.Split(answerString, ".") // Pisahkan berdasarkan titik
	if len(sentences) > maxSentences {
		answerString = strings.Join(sentences[:maxSentences], ".") + "."
	}

	return answerString, nil
}
