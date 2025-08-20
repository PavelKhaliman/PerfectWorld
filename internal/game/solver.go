package game

// Решатель кратчайшего пути через поиск в ширину (BFS) по пространству состояний 2^9

// Представление поля в виде битовой маски: 1 — белая, 0 — чёрная
// Порядок битов: индексы 0..8 соответствуют кнопкам 1..9

// Маски нажатий (какие клетки инвертируются при нажатии i)
var moveMasks [9]uint16

func init() {
	// Предрасчёт масок нажатий согласно правилам (клетка + ортогональные соседи)
	for i := 0; i < 9; i++ {
		mask := uint16(0)
		row := i / 3
		col := i % 3
		// self
		mask |= 1 << i
		// up
		if row > 0 {
			mask |= 1 << (i - 3)
		}
		// down
		if row < 2 {
			mask |= 1 << (i + 3)
		}
		// left
		if col > 0 {
			mask |= 1 << (i - 1)
		}
		// right
		if col < 2 {
			mask |= 1 << (i + 1)
		}
		moveMasks[i] = mask
	}
}

// Encode: Board -> битовая маска
func Encode(b Board) uint16 {
	var m uint16 = 0
	for i := 0; i < 9; i++ {
		if b[i] {
			m |= 1 << i
		}
	}
	return m
}

// Decode: битовая маска -> Board
func Decode(m uint16) Board {
	var b Board
	for i := 0; i < 9; i++ {
		b[i] = (m>>i)&1 == 1
	}
	return b
}

// ApplyMove применяет нажатие кнопки index (0..8) к маске
func ApplyMove(m uint16, index int) uint16 {
	return m ^ moveMasks[index]
}

// SolveShortest возвращает последовательность ходов (1..9) минимальной длины,
// чтобы привести начальное состояние к целевому (все белые, либо все чёрные).
// Возвращает пустой срез, если уже в целевом состоянии.
func SolveShortest(start Board, toAllWhite bool) []int {
	startMask := Encode(start)
	var target uint16
	if toAllWhite {
		target = 0x1FF // 9 единиц
	} else {
		target = 0x000
	}
	return solveBFS(startMask, target)
}

// Внутренний BFS по маскам
func solveBFS(start uint16, target uint16) []int {
	if start == target {
		return []int{}
	}
	// Очередь
	queue := make([]uint16, 0, 512)
	queue = append(queue, start)
	// Предки: состояние -> (предыдущее, ход)
	prev := make(map[uint16]struct {
		p    uint16
		move int
	}, 512)
	visited := make(map[uint16]bool, 512)
	visited[start] = true

	for qi := 0; qi < len(queue); qi++ {
		cur := queue[qi]
		for move := 0; move < 9; move++ {
			next := cur ^ moveMasks[move]
			if visited[next] {
				continue
			}
			visited[next] = true
			prev[next] = struct {
				p    uint16
				move int
			}{p: cur, move: move}
			if next == target {
				// восстановить путь
				return reconstruct(prev, start, target)
			}
			queue = append(queue, next)
		}
	}
	// Теоретически достижимо всегда, но на всякий случай
	return []int{}
}

func reconstruct(prev map[uint16]struct {
	p    uint16
	move int
}, start, target uint16) []int {
	path := make([]int, 0, 16)
	cur := target
	for cur != start {
		pr := prev[cur]
		// move в prev хранится 0..8; для вывода нужны 1..9
		path = append(path, pr.move+1)
		cur = pr.p
	}
	// развернуть
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}
