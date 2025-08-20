package game

import (
	"math/rand"
	"time"
)

// Инициализация генератора случайных чисел один раз при старте
func init() {
	rand.Seed(time.Now().UnixNano())
}

// Ячейка: true — белая сторона сферы, false — чёрная
type Cell bool

// Игровое поле 3x3, развернутое в массив длиной 9
type Board [9]Cell

// Создаёт новое случайное поле
func NewRandomBoard() Board {
	var b Board
	for i := 0; i < 9; i++ {
		b[i] = rand.Intn(2) == 1
	}
	return b
}

// Создаёт однородное поле: все белые (true) или все чёрные (false)
func NewUniformBoard(white bool) Board {
	var b Board
	for i := 0; i < 9; i++ {
		b[i] = Cell(white)
	}
	return b
}

// Переключает нажатую ячейку и ортогональных соседей (без диагоналей)
func (b *Board) Toggle(index int) {
	if index < 0 || index >= len(b) {
		return
	}
	// flip clicked
	(*b)[index] = !(*b)[index]

	row := index / 3
	col := index % 3

	// up
	if row > 0 {
		up := index - 3
		(*b)[up] = !(*b)[up]
	}
	// down
	if row < 2 {
		down := index + 3
		(*b)[down] = !(*b)[down]
	}
	// left
	if col > 0 {
		left := index - 1
		(*b)[left] = !(*b)[left]
	}
	// right
	if col < 2 {
		right := index + 1
		(*b)[right] = !(*b)[right]
	}
}
