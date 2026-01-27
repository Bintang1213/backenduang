package controllers

import (
	"fmt"
	"gindev/config"
	"gindev/models"
	"net/http"
	"github.com/gin-gonic/gin"
)

// Dashboard untuk Admin Keuangan
func GetKeuanganDashboard(c *gin.Context) {
	var totalTarif int64
	var tarifTerbaru []models.Tarif

	config.DB.Model(&models.Tarif{}).Where("status = ?", "aktif").Count(&totalTarif)
	config.DB.Order("created_at desc").Limit(5).Find(&tarifTerbaru)

	c.JSON(http.StatusOK, gin.H{
		"total_tarif_aktif": totalTarif,
		"tarif_terbaru":     tarifTerbaru,
	})
}

// Get All Tarif (Bisa diakses Admin Keuangan, Kasir, dan Manajemen)
func GetAllTarif(c *gin.Context) {
	var tarif []models.Tarif
	query := config.DB

	if nama := c.Query("nama"); nama != "" {
		query = query.Where("nama_tarif ILIKE ?", "%"+nama+"%")
	}
	if kat := c.Query("kategori"); kat != "" {
		query = query.Where("kategori = ?", kat)
	}

	query.Find(&tarif)
	c.JSON(http.StatusOK, tarif)
}

// Create Tarif (Hanya Admin Keuangan)
func CreateTarif(c *gin.Context) {
	var input models.Tarif
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var lastTarif models.Tarif
	config.DB.Unscoped().Order("id desc").First(&lastTarif)
	input.KodeTarif = fmt.Sprintf("TRF%03d", lastTarif.ID+1)

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan tarif"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Tarif berhasil dibuat", "data": input})
}

// Update Tarif (Hanya Admin Keuangan)
func UpdateTarif(c *gin.Context) {
	id := c.Param("id")
	var tarif models.Tarif

	if err := config.DB.First(&tarif, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tarif tidak ditemukan"})
		return
	}

	var input models.Tarif
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&tarif).Updates(input)
	c.JSON(http.StatusOK, gin.H{"message": "Tarif berhasil diupdate", "data": tarif})
}

// Delete Tarif (Hanya Admin Keuangan)
func DeleteTarif(c *gin.Context) {
	id := c.Param("id")
	var tarif models.Tarif

	if err := config.DB.First(&tarif, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tarif tidak ditemukan"})
		return
	}

	config.DB.Delete(&tarif)
	c.JSON(http.StatusOK, gin.H{"message": "Tarif berhasil dihapus"})
}