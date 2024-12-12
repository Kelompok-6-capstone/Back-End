package controller

import (
	"calmind/helper"
	"calmind/usecase"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type CustServiceController struct {
	CustServiceUsecase usecase.CustServiceUsecase
}

func NewCustServiceController(custServiceUsecase usecase.CustServiceUsecase) *CustServiceController {
	return &CustServiceController{CustServiceUsecase: custServiceUsecase}
}

func (c *CustServiceController) GetResponse(ctx echo.Context) error {

	var request struct {
		Message string `json:"message"`
	}

	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Input tidak valid")
	}

	// Static FAQ responses
	standardResponses := map[string]string{
		"1": "Metode pembayaran yang bisa digunakan Kartu kredit/debit E-wallet (GoPay, OVO, dll.) Transfer bank Silahkan pilih metode yang paling nyaman bagi Anda 😊",
		"2": "Setelah memilih dokter dan metode pembayaran, Anda akan diarahkan ke halaman pembayaran. Ikuti petunjuk di layar untuk menyelesaikan pembayaran. Pilih dokter dan klik buat janji Isi keluhan yang sedang dialami Pilih metode pembayaran Ikuti petunjuk di layar untuk menyelesaikan pembayaran.",
		"3": "Anda harus melakukan konsultasi jika mengalami gejala-gejala berikut: Kesulitan mengelola emosi atau stres. Memiliki pikiran negatif yang mengganggu aktivitas sehari-hari. Mengalami perubahan perilaku yang signifikan. Ingin mengembangkan diri atau meningkatkan kualitas hidup.",
		"4": "Konsultasi bisa dilakukan dengan cara sebagai berikut: Pilih dokter yang ingin dihubungi. Isi form keluhan sesuai dengan keadaan yang dialami. Lakukan pembayaran Anda akan terhubung dengan dokter.",
		"5": "Mendaftar akun CalmMind sangat mudah. Klik tombol 'Daftar' pada halaman awal. Isi data pada form sesuai instruksi yang diberikan. Konfirmasi pada email yang terdaftar. Akun berhasil didaftarkan.",
		"6": "Jika tidak bisa masuk ke akun Anda, ada beberapa hal yang harus Anda lakukan. Periksa kembali koneksi internet. Pastikan nama pengguna dan kata sandi yang dimasukkan sudah benar. Jika masih mengalami masalah silahkan reset kata sandi dan hubungi Admin.",
		"7": "Jangan khawatir, Anda bisa mereset kata sandi dengan mudah. Klik tombol 'Lupa Kata Sandi' pada halaman masuk. Ikuti petunjuk yang diberikan. Anda akan menerima email berisi tautan untuk mengatur ulang kata sandi.",
	}

	lowerMessage := strings.ToLower(request.Message)
	if response, exists := standardResponses[lowerMessage]; exists {
		return helper.JSONSuccessResponse(ctx, map[string]interface{}{
			"message":  request.Message,
			"response": response,
		})
	}

	// Save question for admin
	// err := c.CustServiceUsecase.SaveCustService(claims.UserID, request.Message)
	// if err != nil {
	// 	return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan pertanyaan")
	// }

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"message":  request.Message,
		"response": "Jawaban atas pertanyaan anda tidak tersedia. Silahkan hubungi Admin melalui kontak berikut https://wa.me/+6283820440747",
	})
}

func (c *CustServiceController) GetUnansweredMessages(ctx echo.Context) error {
	messages, err := c.CustServiceUsecase.GetUnansweredMessages()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil pesan yang belum terjawab")
	}

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"unanswered_messages": messages,
	})
}

func (c *CustServiceController) AnswerMessage(ctx echo.Context) error {
	var request struct {
		ID     int    `json:"id"`
		Answer string `json:"answer"`
	}

	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Input tidak valid")
	}

	err := c.CustServiceUsecase.AnswerMessage(request.ID, request.Answer)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menjawab pesan")
	}

	return helper.JSONSuccessResponse(ctx, "Pesan berhasil dijawab.")
}

func (c *CustServiceController) GetQuestion(ctx echo.Context) error {
	// Static FAQ responses
	standardResponses := map[string]string{
		"1": "Metode pembayaran?",
		"2": "Proses pembayaran?",
		"3": "Kapan harus melakukan konsultasi?",
		"4": "Bagaimana cara melakukan konsultasi?",
		"5": "Bagaimana cara mendaftar akun?",
		"6": "Apa yang harus dilakukan jika tidak bisa masuk?",
		"7": "Bagaimana jika lupa kata sandi?",
	}

	return helper.JSONSuccessResponse(ctx, standardResponses)
}
