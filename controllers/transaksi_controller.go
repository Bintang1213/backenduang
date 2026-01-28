package controllers

import (
	"fmt"
	"gindev/config"
	"gindev/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateTransaksi(c *gin.Context) {
	var input struct {
		NamaPasien       string `json:"nama_pasien" binding:"required"`
		NoRekamMedis     string `json:"no_rekam_medis" binding:"required"`
		JenisKelamin     string `json:"jenis_kelamin"`
		MetodePembayaran string `json:"metode_pembayaran" binding:"required"`
		LayananIDs       []uint `json:"layanan_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak lengkap"})
		return
	}

	// 1. Logika Nomor Invoice Berurutan (INV/YYYYMMDD/0001)
	tgl := time.Now().Format("20060102")
	var count int64
	
	// Hitung transaksi yang sudah ada khusus hari ini saja
	config.DB.Model(&models.Transaksi{}).
		Where("no_invoice LIKE ?", "INV/"+tgl+"/%").
		Count(&count)

	noInvoice := fmt.Sprintf("INV/%s/%04d", tgl, count+1)

	// 2. Mulai Transaction Database
	tx := config.DB.Begin()

	var total float64
	var details []models.TransaksiDetail

	// 3. Validasi & Ambil data dari Tabel Tarif
	for _, id := range input.LayananIDs {
		var tarif models.Tarif
		if err := tx.First(&tarif, id).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Layanan ID %d tidak ditemukan", id)})
			return
		}

		total += tarif.Harga
		details = append(details, models.TransaksiDetail{
			TarifID:     tarif.ID,
			NamaLayanan: tarif.NamaTarif,
			Harga:       tarif.Harga,
		})
	}

	// 4. Buat Object Transaksi
	transaksi := models.Transaksi{
		NoInvoice:        noInvoice,
		NamaPasien:       input.NamaPasien,
		NoRekamMedis:     input.NoRekamMedis,
		JenisKelamin:     input.JenisKelamin,
		MetodePembayaran: input.MetodePembayaran,
		TotalHarga:       total,
		Details:          details,
	}

	// 5. Simpan ke Database
	if err := tx.Create(&transaksi).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan transaksi"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Transaksi Berhasil",
		"invoice": noInvoice,
		"data":    transaksi,
	})
}

// Fungsi untuk melihat semua riwayat transaksi
func GetRiwayatTransaksi(c *gin.Context) {
	var riwayat []models.Transaksi
	// Preload("Details") agar data item layanan ikut muncul
	if err := config.DB.Preload("Details").Order("created_at desc").Find(&riwayat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}
	c.JSON(http.StatusOK, riwayat)
}