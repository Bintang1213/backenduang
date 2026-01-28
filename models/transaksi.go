package models

import "gorm.io/gorm"

type Transaksi struct {
	gorm.Model
	NoInvoice        string            `gorm:"unique;not null;index" json:"no_invoice"`
	NamaPasien       string            `gorm:"not null" json:"nama_pasien"`
	NoRekamMedis     string            `gorm:"not null" json:"no_rekam_medis"`
	JenisKelamin     string            `json:"jenis_kelamin"`
	MetodePembayaran string            `json:"metode_pembayaran"`
	TotalHarga       float64           `json:"total_harga"`
	Details          []TransaksiDetail `json:"details" gorm:"foreignKey:TransaksiID"`
}

type TransaksiDetail struct {
	gorm.Model
	TransaksiID uint    `json:"transaksi_id"`
	TarifID     uint    `json:"tarif_id"`
	NamaLayanan string  `json:"nama_layanan"`
	Harga       float64 `json:"harga"`
}