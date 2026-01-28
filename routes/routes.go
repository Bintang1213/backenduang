package routes

import (
	"gindev/controllers"
	"gindev/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS Configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://cedrick-unlunated-gwyn.ngrok-free.app"},
		AllowMethods:     []string{"POST", "GET", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With", "ngrok-skip-browser-warning"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api")
	{
		// Publik
		api.POST("/login", controllers.Login)

		// Terproteksi (Wajib Login)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// --- 1. ADMIN KEUANGAN (Master Data) ---
			keuangan := protected.Group("/keuangan").Use(middleware.RoleMiddleware("admin_keuangan"))
			{
				keuangan.GET("/dashboard", controllers.GetKeuanganDashboard)
				keuangan.GET("/tarif", controllers.GetAllTarif)
				keuangan.POST("/tarif", controllers.CreateTarif)
				keuangan.PUT("/tarif/:id", controllers.UpdateTarif)
				keuangan.DELETE("/tarif/:id", controllers.DeleteTarif)

				keuangan.GET("/riwayat", controllers.GetRiwayatTransaksi)
			}

			// --- 2. KASIR (Transaksi & View) ---
			kasir := protected.Group("/kasir").Use(middleware.RoleMiddleware("kasir"))
			{
				kasir.GET("/tarif", controllers.GetAllTarif)
				kasir.POST("/transaksi", controllers.CreateTransaksi)  // Modul Baru
				kasir.GET("/riwayat", controllers.GetRiwayatTransaksi) // Modul Baru
			}

			// --- 3. MANAJEMEN (Monitoring) ---
			manajemen := protected.Group("/manajemen").Use(middleware.RoleMiddleware("manajemen"))
			{
				manajemen.GET("/tarif", controllers.GetAllTarif)
				manajemen.GET("/riwayat", controllers.GetRiwayatTransaksi) // Bisa lihat semua transaksi
			}
		}
	}
	return r
}
