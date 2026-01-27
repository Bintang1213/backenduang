package routes

import (
	"gindev/controllers"
	"gindev/middleware"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
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
			// --- 1. ADMIN KEUANGAN (FULL CRUD) ---
			keuangan := protected.Group("/keuangan").Use(middleware.RoleMiddleware("admin_keuangan"))
			{
				keuangan.GET("/dashboard", controllers.GetKeuanganDashboard)
				keuangan.GET("/tarif", controllers.GetAllTarif) 
				keuangan.POST("/tarif", controllers.CreateTarif)
				keuangan.PUT("/tarif/:id", controllers.UpdateTarif)
				keuangan.DELETE("/tarif/:id", controllers.DeleteTarif)
			}

			// --- 2. KASIR (HANYA VIEW TARIF) ---
			kasir := protected.Group("/kasir").Use(middleware.RoleMiddleware("kasir"))
			{
				kasir.GET("/tarif", controllers.GetAllTarif) // Kasir butuh ambil data tarif untuk input pembayaran
				kasir.GET("/transaksi", func(c *gin.Context) { 
					c.JSON(200, gin.H{"msg": "Menu Transaksi Kasir"}) 
				})
			}

			// --- 3. MANAJEMEN (VIEW LAPORAN & DASHBOARD) ---
			manajemen := protected.Group("/manajemen").Use(middleware.RoleMiddleware("manajemen"))
			{
				manajemen.GET("/tarif", controllers.GetAllTarif) // Manajemen bisa monitoring tarif
				manajemen.GET("/laporan", func(c *gin.Context) { 
					c.JSON(200, gin.H{"msg": "Menu Laporan Manajemen"}) 
				})
			}
		}
	}
	return r
}