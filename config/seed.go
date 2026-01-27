package config

import (
	"fmt"
	"gindev/models"
	"golang.org/x/crypto/bcrypt"
)

func SeedUsers() {
	// 1. Daftar data akun yang akan dibuat
	users := []models.User{
		{Nama: "Admin Keuangan", Username: "admin", Role: "admin_keuangan"},
		{Nama: "Kasir Toko", Username: "kasir", Role: "kasir"},
		{Nama: "Manager", Username: "manager", Role: "manajemen"},
	}

	// 2. Hash password (sama untuk semua akun tes)
	passwordDefault := "12345678"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(passwordDefault), bcrypt.DefaultCost)

	fmt.Println("Memeriksa data seeder...")

	for _, user := range users {
		var existingUser models.User
		// Cek apakah username sudah terdaftar agar tidak terjadi error 'unique constraint'
		err := DB.Where("username = ?", user.Username).First(&existingUser).Error
		
		if err != nil { 
			user.Password = string(hashedPassword)
			if errCreate := DB.Create(&user).Error; errCreate != nil {
				fmt.Printf("Gagal membuat user %s: %v\n", user.Username, errCreate)
			} else {
				fmt.Printf("Berhasil membuat akun: %s | Role: %s\n", user.Username, user.Role)
			}
		}
	}
}