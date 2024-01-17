package main__test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"ta-kasir/controller/produk"
	"testing"

	"github.com/gin-gonic/gin"
)

// Buat struktur tambahan untuk memudahkan analisis pesan kesalahan
type ErrorResponse struct {
	Message string `json:"message"`
}

func TestAddProduk(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Inisialisasi router Gin
	router := gin.Default()

	// Tentukan endpoint yang akan diuji
	router.POST("/addproduk", produk.AddProduk)

	// Persiapkan data JSON untuk produk
	jsonData := `{"nama": "TestProduk", "harga": 100, "stok": 50}`

	// Buat buffer untuk menyimpan data multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Tambahkan data JSON ke form
	writer.WriteField("json", jsonData)

	// Buat file palsu untuk diunggah
	fileContents := "editit"
	tempFile, err := ioutil.TempFile("", "tempfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Write([]byte(fileContents))
	tempFile.Close()

	fileWriter, err := writer.CreateFormFile("file", "file.txt")
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Open(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	// Salin konten file sementara ke form
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		t.Fatal(err)
	}

	// Tutup writer form
	writer.Close()

	// Buat permintaan palsu menggunakan httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/addproduk", &buf)
	if err != nil {
		t.Fatal(err)
	}

	// Set header Content-Type sesuai dengan form-data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Buat recorder untuk menyimpan respons
	rr := httptest.NewRecorder()

	// Kirimkan permintaan ke router Gin
	router.ServeHTTP(rr, req)

	// Periksa status code yang diharapkan
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v but got %v", http.StatusOK, status)

		// Cetak respons untuk membantu analisis
		fmt.Println(rr.Body.String())

		// Coba untuk mendapatkan pesan kesalahan sebagai string
		var errorResponse ErrorResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &errorResponse); err == nil {
			fmt.Println("Error message:", errorResponse.Message)
		}
	}

	// Periksa respons JSON yang dihasilkan
	// (Anda mungkin perlu membuat struktur untuk membaca respons JSON)
	var responseStruct struct {
		Status  int             `json:"status"`
		Error   json.RawMessage `json:"error"`
		Message string          `json:"message"`
		Data    string          `json:"data"`
	}
	err = json.NewDecoder(rr.Body).Decode(&responseStruct)
	if err != nil {
		t.Fatal(err)
	}
}
