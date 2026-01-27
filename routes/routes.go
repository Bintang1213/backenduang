package routes

import (
	"gindev/controllers"
	"gindev/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		// 1. Rute Publik
		api.POST("/login", controllers.Login)

		// 2. Rute Terproteksi (Harus Login)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Menu khusus Admin Keuangan
			keuangan := protected.Group("/keuangan").Use(middleware.RoleMiddleware("admin_keuangan"))
			{
				keuangan.GET("/dashboard", func(c *gin.Context) { 
                    c.JSON(200, gin.H{"msg": "Halo Admin Keuangan"}) 
                })
			}

			// Menu khusus Kasir
			kasir := protected.Group("/kasir").Use(middleware.RoleMiddleware("kasir"))
			{
				kasir.GET("/transaksi", func(c *gin.Context) { 
                    c.JSON(200, gin.H{"msg": "Halo Kasir"}) 
                })
			}

			// Menu khusus Manajemen
			manajemen := protected.Group("/manajemen").Use(middleware.RoleMiddleware("manajemen"))
			{
				manajemen.GET("/laporan", func(c *gin.Context) { 
                    c.JSON(200, gin.H{"msg": "Halo Manajemen"}) 
                })
			}
		}
	}
	return r
}