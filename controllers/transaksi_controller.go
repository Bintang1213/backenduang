package controllers

import (
	"fmt"
	"gindev/config"
	"gindev/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateTransaksi untuk menyimpan transaksi baru dengan nomor invoice urut
func CreateTransaksi(c *gin.Context) {
	var input struct {
		NamaPasien       string `json:"nama_pasien" binding:"required"`
		NoRekamMedis     string `json:"no_rekam_medis" binding:"required"`
		JenisKelamin     string `json:"jenis_kelamin"`
		MetodePembayaran string `json:"metode_pembayaran" binding:"required"`
		LayananIDs       []uint `json:"layanan_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak lengkap atau format salah"})
		return
	}

	// 1. Generate No Invoice Berurutan (INV/YYYYMMDD/0001)
	tgl := time.Now().Format("20060102")
	var count int64
	config.DB.Model(&models.Transaksi{}).
		Where("no_invoice LIKE ?", "INV/"+tgl+"/%").
		Count(&count)

	noInvoice := fmt.Sprintf("INV/%s/%04d", tgl, count+1)

	// 2. Mulai Transaction Database
	tx := config.DB.Begin()

	var total float64
	var details []models.TransaksiDetail

	// 3. Validasi & Ambil data dari Tabel Tarif (Denormalisasi)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan transaksi ke database"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, gin.H{"message": "Transaksi Berhasil", "data": transaksi})
}

// GetRiwayatTransaksi untuk melihat daftar semua transaksi
func GetRiwayatTransaksi(c *gin.Context) {
	var riwayat []models.Transaksi
	if err := config.DB.Preload("Details").Order("created_at desc").Find(&riwayat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data riwayat"})
		return
	}
	c.JSON(http.StatusOK, riwayat)
}

// GetRingkasanPemasukan untuk melihat total uang dan total transaksi
func GetRingkasanPemasukan(c *gin.Context) {
	var stats struct {
		TotalUang      float64 `json:"total_uang"`
		TotalTransaksi int64   `json:"total_transaksi"`
	}

	// Menggunakan COALESCE agar jika data kosong tetap mengembalikan 0 bukan null
	err := config.DB.Model(&models.Transaksi{}).
		Select("COALESCE(SUM(total_harga), 0) as total_uang, COUNT(id) as total_transaksi").
		Scan(&stats).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses data pemasukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   stats,
	})
}