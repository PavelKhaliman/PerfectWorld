package api

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"

	"perfectworld/internal/game"
)

// Роутер API и статики. Принимает ссылку на игровое поле и путь к каталогу статики
func Router(board *game.Board, publicDir string) http.Handler {
	mux := http.NewServeMux()

	// Статика
	static := http.FileServer(http.Dir(publicDir))
	mux.Handle("/assets/", http.StripPrefix("/", static))
	mux.Handle("/css/", http.StripPrefix("/", static))
	mux.Handle("/js/", http.StripPrefix("/", static))

	// API: получить текущее состояние поля
	mux.HandleFunc("/api/state", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(board)
	})

	// API: переключить клетку по индексу (1..9)
	mux.HandleFunc("/api/toggle", func(w http.ResponseWriter, r *http.Request) {
		idxStr := r.FormValue("index")
		if idxStr != "" {
			if idx, err := strconv.Atoi(idxStr); err == nil {
				idx-- // convert 1-9 to 0-8
				board.Toggle(idx)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(board)
	})

	// API: сбросить поле — новая случайная конфигурация
	mux.HandleFunc("/api/reset", func(w http.ResponseWriter, r *http.Request) {
		*board = game.NewRandomBoard()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(board)
	})

	// API: решить кратчайшим путём
	// Параметры:
	// - target = white|black (по умолчанию white)
	// - state (опционально) — строка из 9 символов 0/1; если не задано — текущее поле
	mux.HandleFunc("/api/solve", func(w http.ResponseWriter, r *http.Request) {
		targetStr := r.FormValue("target")
		toWhite := targetStr != "black"

		var start game.Board
		stateStr := r.FormValue("state")
		if len(stateStr) == 9 {
			for i := 0; i < 9; i++ {
				start[i] = stateStr[i] == '1'
			}
		} else {
			start = *board
		}

		moves := game.SolveShortest(start, toWhite)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(struct {
			Moves  []int      `json:"moves"`
			Count  int        `json:"count"`
			From   game.Board `json:"from"`
			Target string     `json:"target"`
		}{
			Moves: moves,
			Count: len(moves),
			From:  start,
			Target: func() string {
				if toWhite {
					return "white"
				}
				return "black"
			}(),
		})
	})

	// Корень: index.html; прочие пути — статика
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			static.ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(publicDir, "index.html"))
	})

	return mux
}
