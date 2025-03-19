// cmd/main.go
package main

import (
	"domashka-backend/config"
	"domashka-backend/internal/app"
)

func main() {
	// Инициализация приложения

	app.Run(config.GetConfig())
}
