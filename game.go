package minego

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type direction int

func (d direction) offset() (h, w int) {
	switch d {
	case nop:
		return 0, 0
	case top:
		return -1, 0
	case topLeft:
		return -1, -1
	case left:
		return 0, -1
	case bottomLeft:
		return 1, -1
	case bottom:
		return 1, 0
	case bottomRight:
		return 1, 1
	case right:
		return 0, 1
	case topRight:
		return -1, 1
	default:
		return 0, 0
	}
}

const (
	nop direction = iota
	top
	left
	bottom
	right

	topLeft
	topRight
	bottomLeft
	bottomRight
)

type cell struct {
	hasBomb bool
	isOpen  bool
	bomb    int
}

func (c *cell) open() {
	c.isOpen = true
}

func (c cell) String() string {
	switch {
	case c.isOpen && c.bomb > 0:
		return strconv.Itoa(c.bomb)
	case c.isOpen && c.hasBomb:
		return "x"
	case c.isOpen:
		return "□"
	default:
		return "■"
	}
}

// NewGame は minesweeper の ゲームを生成します
// - Vertical    cell's size is h (index is 0~h-1)
// - Horizontal  cell's size is w (index is 0~w-1)
func NewGame(h, w, bobNum int) Game {
	// generate cells
	cells := make([][]cell, h)
	for i := 0; i < len(cells); i++ {
		cells[i] = make([]cell, w)
	}

	// set bomb
	rand.Seed(time.Now().UnixNano())
	// FIXME: loop never ends if  w*h < bomb
	setNum := 0
	for {
		if setNum == bobNum {
			break
		}

		wb := rand.Intn(w)
		hb := rand.Intn(h)
		if cells[hb][wb].hasBomb {
			continue
		}
		cells[hb][wb].hasBomb = true
		setNum++
	}
	return Game{
		maxHIndex: h - 1,
		maxWIndex: w - 1,
		cells:     cells,
	}
}

// Game は minesweeper のゲームを表します
type Game struct {
	maxHIndex int
	maxWIndex int
	cells     [][]cell
}

func (g *Game) Show() {
	fmt.Println("=========")
	for _, chs := range g.cells {
		for _, c := range chs {
			fmt.Print(c)
		}
		fmt.Println()
	}
	fmt.Println("=========")
}

// OpenCell open specified Game's cell
func (g *Game) OpenCell(h, w int) (bool, error) {
	if h < 0 || w < 0 || g.maxHIndex < h || g.maxWIndex < w {
		return false, errors.New("failed open: out of size")
	}
	c := g.cells[h][w]
	if c.isOpen {
		return false, nil
	}
	c.open()
	// TODO 爆発判定の仕方検討
	if c.hasBomb {
		return true, nil
	}
	g.openAdjacentCells(h, w)

	return false, nil
}

func (g *Game) hasBomb(h, w int) bool { return g.cells[h][w].hasBomb }

func (g *Game) open(h, w int) {
	c := g.cells[h][w]
	c.open()
	g.cells[h][w] = c
}

func (g *Game) hasAdjacentBomb(h, w int) bool {
	for _, d := range []direction{top, topLeft, left, bottomLeft, bottom, bottomRight, right, topRight} {
		ho, wo := d.offset()
		h := h + ho
		w := w + wo
		if h < 0 || w < 0 || g.maxHIndex < h || g.maxWIndex < w {
			continue
		}

		adjacentCell := g.cells[h][w]
		if adjacentCell.hasBomb {
			return true
		}
	}
	return false
}

// openAdjacentCells open adjacent cells starting from (h, w)
func (g *Game) openAdjacentCells(h, w int) {
	for _, d := range []direction{top, left, bottom, right} {
		ho, wo := d.offset()
		h := h + ho
		w := w + wo
		if h < 0 || w < 0 || g.maxHIndex < h || g.maxWIndex < w {
			continue
		}

		targetCell := g.cells[h][w]
		if targetCell.isOpen || targetCell.hasBomb || g.hasAdjacentBomb(h, w) {
			continue
		}
		g.open(h, w)
		g.openAdjacentCells(h, w)
	}
}
