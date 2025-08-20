package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"perfectworld/internal/api"
	"perfectworld/internal/game"
)

func main() {
	// Путь к статике: по умолчанию ./public, можно переопределить через PUBLIC_DIR
	wd, _ := os.Getwd()
	publicDir := filepath.Join(wd, "public")
	if v := os.Getenv("PUBLIC_DIR"); v != "" {
		publicDir = filepath.Clean(v)
	}

	// Состояние приложения (игровое поле): по умолчанию все белые
	board := game.NewUniformBoard(true)

	// Роутер API и статики
	mux := api.Router(&board, publicDir)

	log.Println("listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
