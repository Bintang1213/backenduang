package main

import (
    "gindev/config"
    "gindev/routes"
    "github.com/joho/godotenv"
)

func main() {
    godotenv.Load()
    config.ConnectDatabase()
    config.SeedUsers()

    r := routes.SetupRouter()
    r.Run(":8080")
}