package helper

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"ta-kasir/config"
	"ta-kasir/model"
	"ta-kasir/model/request"
	"time"

	"github.com/fogleman/gg"
	"github.com/jung-kurt/gofpdf"
)

func GenerateRandomNumber(length int) string {
	const charset = "0123456789"

	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}

func GenerateImage(width, height int, idPenjualan, namaKasir string, productsID []int, dataPesanan []request.Pesanan, subTotals []float64, pembayaran request.Bayar) (string, error) {

	dc := gg.NewContext(width, height)

	// set background color
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.SetColor(color.Black)
	
	// HEADER 1
	err := dc.LoadFontFace("./storage/fonts/Poppins-Regular.ttf", 18)
	if err != nil {
		log.Fatal(err)
	}
	dc.DrawString("Simple Cash", 100, 30,)
	
	// HEADER 2
	err = dc.LoadFontFace("./storage/fonts/Poppins-Regular.ttf", 16)
	if err != nil {
		log.Fatal(err)
	}
	dc.DrawString("Jl Tanimbar No.22 Malang", 50, 55)

	// GARIS PEMISAH 1
	err = dc.LoadFontFace("./storage/fonts/Poppins-Regular.ttf", 10)
	if err != nil {
		log.Fatal(err)
	}
	dc.DrawString("-----------------------------------------------", 20, 75)
	
	// TANGGAL DAN NAMA KASIR
	timeString := time.Now().Format("02.01.2006")
	err = dc.LoadFontFace("./storage/fonts/Poppins-Regular.ttf", 11)
	if err != nil {
		log.Fatal(err)
	}
	dc.DrawString(timeString + "/", 20, 95)
	dc.DrawString(idPenjualan, 84, 95)
	dc.DrawString(namaKasir, 220, 95)

	// GARIS PEMISAH 2
	err = dc.LoadFontFace("./storage/fonts/Poppins-Regular.ttf", 10)
	if err != nil {
		log.Fatal(err)
	}
	dc.DrawString("-----------------------------------------------", 20, 115)


	// PRODUK DAN JUMLAH YANG DIBELI
	err = dc.LoadFontFace("./storage/fonts/Poppins-Regular.ttf", 12)
	if err != nil {
		log.Fatal(err)
	}
	var productsName []string
	db := config.ConnectDatabase()
	err = db.Debug().Model(&model.Produk{}).Where("id_produk IN (?)", productsID).Pluck("nama_produk", &productsName).Error
		if err != nil {
			log.Printf("Error fetching product names: %v", err)
			return "", err
		}

	// looping for each product
	var yStart float64 = 135
	Subtotal := float64(width) - 75
	lineHeight := float64(20)  
	for i, pesanan := range dataPesanan {
        y := yStart + float64(i*20)

        // Tampilkan nama produk
        dc.DrawString(productsName[i], 20, y)

        // Tampilkan jumlah produk dengan jarak yang ditentukan
        dc.DrawString(fmt.Sprintf("x%d", pesanan.JumlahProduk), 180, y)

        // Tampilkan jumlah dengan jarak ke kanan
        dc.DrawString(fmt.Sprintf("Rp. %.0f", subTotals[i]), Subtotal, y)

		if i == len(productsName)-1 {
			// Gambar garis pembatas terakhir
			dc.DrawLine(160, y+lineHeight, float64(width)-20, y+lineHeight)
			dc.Stroke()
		
			// Tampilkan "Subtotal" dan jumlahnya setelah garis pembatas terakhir
			dc.DrawString("Subtotal ", 120, y+lineHeight+20)
			dc.DrawString(":", 200, y+lineHeight+20)
		
			// Tampilkan jumlah subtotal
			dc.DrawString(fmt.Sprintf("Rp. %.0f", pembayaran.Amount), Subtotal, y+lineHeight+20)
		
			// Tampilkan kata "Admin Fees" dan jumlahnya setelah subtotal
			dc.DrawString("Admin Fees ", 120, y+lineHeight+40)
			dc.DrawString(":", 200, y+lineHeight+40)
		
			// Tampilkan jumlah admin fees
			dc.DrawString(fmt.Sprintf("Rp. %.0f", pembayaran.BiayaAdmin), Subtotal, y+lineHeight+40)
		
			// Gambar garis pembatas
			dc.DrawLine(160, y+lineHeight+60, float64(width)-20, y+lineHeight+60)
			dc.Stroke()
		
			// Tampilkan "Total" dan jumlahnya setelah garis pembatas terakhir
			dc.DrawString("Total ", 120, y+lineHeight+80)
			dc.DrawString(":", 200, y+lineHeight+80)
		
			// Tampilkan jumlah total
			dc.DrawString(fmt.Sprintf("Rp. %.0f", pembayaran.Grandtotal), Subtotal, y+lineHeight+80)

			// Tampilkan ucapan terima kasih
			dc.LoadFontFace("./storage/fonts/Poppins-Regular.ttf", 10)
			dc.DrawString("Thank you for shopping at Simple Cash", 60, y+lineHeight+150)
		}
		
    }

	receiptPath := fmt.Sprintf("./storage/receipt/Receipt Simple Cash (%s).png", idPenjualan) 
	err = dc.SavePNG(receiptPath)
	if err != nil {
		log.Printf("Error saving PNG: %v", err)
		return "", err
	}

	return receiptPath, nil
}

func GeneratePDF(idPenjualan, namaKasir string, productsID []int, dataPesanan []request.Pesanan, subTotals []float64, pembayaran request.Bayar) (string, error) {
    pdf := gofpdf.New("P", "mm", "A4", "") // Buat objek PDF baru dengan ukuran A4
    pdf.AddPage() // Tambahkan halaman pertama

    // Font dan ukuran teks
    pdf.SetFont("Arial", "", 12)

    // Header
    pdf.CellFormat(190, 10, "Simple Cash", "", 0, "C", false, 0, "")
    pdf.Ln(5)
    pdf.CellFormat(190, 10, "Jl Tanimbar No.22 Malang", "", 0, "C", false, 0, "")
    pdf.Ln(10)

    // Data pesanan
    var productsName []string
    db := config.ConnectDatabase()
    err := db.Debug().Model(&model.Produk{}).Where("id_produk IN (?)", productsID).Pluck("nama_produk", &productsName).Error
    if err != nil {
        log.Printf("Error fetching product names: %v", err)
        return "", err
    }

    for i, pesanan := range dataPesanan {
        // Tampilkan nama produk
        pdf.CellFormat(100, 10, fmt.Sprintf("%s", productsName[i]), "", 0, "", false, 0, "")

        // Tampilkan jumlah produk
        pdf.CellFormat(30, 10, fmt.Sprintf("x%d", pesanan.JumlahProduk), "", 0, "", false, 0, "")

        // Tampilkan subtotal
        pdf.CellFormat(30, 10, fmt.Sprintf("Rp. %.0f", subTotals[i]), "", 1, "R", false, 0, "")
    }

    // Total, Admin Fees, dan Grandtotal
    pdf.Ln(10)
    pdf.CellFormat(160, 10, "Total", "", 0, "", false, 0, "")
    pdf.CellFormat(30, 10, fmt.Sprintf("Rp. %.0f", pembayaran.Grandtotal), "", 1, "R", false, 0, "")
    pdf.CellFormat(160, 10, "Admin Fees", "", 0, "", false, 0, "")
    pdf.CellFormat(30, 10, fmt.Sprintf("Rp. %.0f", pembayaran.BiayaAdmin), "", 1, "R", false, 0, "")
    pdf.CellFormat(160, 10, "Subtotal", "", 0, "", false, 0, "")
    pdf.CellFormat(30, 10, fmt.Sprintf("Rp. %.0f", pembayaran.Amount), "", 1, "R", false, 0, "")

    // Ucapan terima kasih
    pdf.Ln(10)
    pdf.CellFormat(190, 10, "Thank you for shopping at Simple Cash", "", 0, "C", false, 0, "")

    receiptPath := fmt.Sprintf("./storage/receipt/Receipt Simple Cash (%s).pdf", idPenjualan)
    err = pdf.OutputFileAndClose(receiptPath)
    if err != nil {
        log.Printf("Error saving PDF: %v", err)
        return "", err
    }

    return receiptPath, nil
}



