package models

import "gorm.io/gorm"

type Tarif struct {
	gorm.Model
	KodeTarif  string  `gorm:"unique;not null" json:"kode_tarif"`
	NamaTarif  string  `gorm:"not null" json:"nama_tarif"`
	Kategori   string  `gorm:"not null" json:"kategori"` 
	Harga      float64 `gorm:"not null" json:"harga"`
	Keterangan string  `json:"keterangan"`
	Status     string  `gorm:"default:'aktif'" json:"status"` // aktif / non-aktif
}