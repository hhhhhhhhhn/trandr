package main

type Grid [][]rune

func NewGrid(width, height int) (grid Grid) {
	for i := 0; i<width; i++ {
		column := []rune{}
		for j := 0; j<height; j++ {
			column = append(column, ' ')
		}
		grid = append(grid, column)
	}
	return grid
}

func (g Grid) DrawHorizontalLine(x1, x2, y int) {
	for i := x1; i <= x2; i++ {
		g[i][y] = '-'
	}
}

func (g Grid) DrawVerticalLine(x, y1, y2 int) {
	for j := y1; j <= y2; j++ {
		g[x][j] = '|'
	}
}

func (g Grid) DrawSquare(x1, y1, x2, y2 int) {
	g.DrawVerticalLine(x1, y1, y2)
	g.DrawVerticalLine(x2, y1, y2)
	g.DrawHorizontalLine(x1, x2, y1)
	g.DrawHorizontalLine(x1, x2, y2)
}

func (g Grid) DrawText(text string, x, y int) {
	for i := 0; i < len(text); i++ {
		g[i + x][y] = rune(text[i])
	}
}

func (g Grid) DrawCenteredText(text string, x, y int) {
	offset := len(text) / 2
	g.DrawText(text, x - offset, y)
}


func (g Grid) DrawTextBox(text string, x1, y1, x2, y2 int) {
	textX := (x1 + x2) / 2
	textY := (y1 + y2) / 2
	g.DrawSquare(x1, y1, x2, y2)
	g.DrawCenteredText(text, textX, textY)
}

func (g Grid) ToString() (str string) {
	width := len(g)
	height := len(g[0])

	for j := height-1; j >= 0; j-- {
		for i := 0; i < width; i++ {
			str += string(g[i][j])
		}
		str += "\n"
	}

	return str
}

//func main() {
//	g := NewGrid(30, 10)
//
//	g.DrawSquare(6, 2, 24, 8)
//	g.DrawCenteredText("Some cool text", 15, 5)
//	fmt.Print(g.ToString())
//}
