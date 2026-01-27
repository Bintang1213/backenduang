package main

import (
	"gindev/config"
	"gindev/routes"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	
	// 1. Koneksi Database
	config.ConnectDatabase()

	// 2. Jalankan Seeder (mengisi akun otomatis)
	config.SeedUsers()

	// 3. Jalankan Server
	r := routes.SetupRouter()
	r.Run(":8080")
}