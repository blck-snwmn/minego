package minego

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
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

var directions = []direction{top, topLeft, left, bottomLeft, bottom, bottomRight, right, topRight}

type cell struct {
	hasBomb bool
	hasFlag bool
	isOpen  bool
	bomb    int
}

func (c *cell) open() {
	c.isOpen = true
}

func (c *cell) setFlag() {
	c.hasFlag = true
}

func (c *cell) removeFlag() {
	c.hasFlag = false
}

func (c cell) String() string {
	if !c.isOpen {
		switch {
		case c.hasFlag:
			return "F"
		default:
			return "■"
		}
	}
	switch {
	case c.hasBomb:
		return "x"
	case c.bomb > 0:
		return strconv.Itoa(c.bomb)
	default:
		return "□"
	}
}

// CommandType is command type
type commandType string

const (
	open       commandType = "o"
	flagSet    commandType = "fs"
	flagRemove commandType = "fr"
)

func (ct commandType) isValid() bool {
	switch ct {
	case open, flagSet, flagRemove:
		return true
	default:
		return false
	}
}

// NewCommand generate Command
func NewCommand(t string, h, w int) (Command, error) {
	ct := commandType(t)
	if ct.isValid() {
		return Command{ct, h, w}, nil
	}
	return Command{}, errors.New("no exist command type")
}

// Command is command of minesweeper
type Command struct {
	typ  commandType
	h, w int
}

// NewGame は minesweeper の ゲームを生成します
// - Vertical    cell's size is h (index is 0~h-1)
// - Horizontal  cell's size is w (index is 0~w-1)
func NewGame(h, w, bobNum int, writter io.Writer) Game {
	// generate cells
	cells := make([][]cell, h)
	for i := 0; i < len(cells); i++ {
		cells[i] = make([]cell, w)
	}
	g := Game{
		cells:         cells,
		closedCellNum: h * w,
		bombNum:       bobNum,
		writer:        bufio.NewWriter(writter),
	}
	g.setBomb(h, w, bobNum)
	return g
}

// Game は minesweeper のゲームを表します
type Game struct {
	cells         [][]cell
	closedCellNum int
	setFlagNum    int
	bombNum       int

	writer *bufio.Writer
}

// Show show current game state
func (g *Game) Show() {
	// TODO buffering
	headerCellLen := len(g.cells[0])
	sep := strings.Repeat("=", (headerCellLen+1)*3)
	g.writer.WriteString(sep)
	g.writer.WriteString("\n")
	// fmt.Println(sep)

	// header
	g.writer.WriteString("   ")
	for i := 0; i < headerCellLen; i++ {
		g.writer.WriteString(fmt.Sprintf(" %02d", i))
	}

	// rows
	g.writer.WriteString("\n")
	for i, chs := range g.cells {
		g.writer.WriteString(fmt.Sprintf(" %02d", i))
		for _, c := range chs {
			g.writer.WriteString(fmt.Sprintf("%3s", c))
		}
		g.writer.WriteString("\n")
	}

	g.writer.WriteString(sep)
	g.writer.WriteString("\n")
	g.writer.Flush()
}

// Do は minesweeper を １サイクルすすめる
func (g *Game) Do(c Command) (bool, error) {
	if !g.isInGameArea(c.h, c.w) {
		return false, errors.New("failed open: out of size")
	}
	switch c.typ {
	case open:
		return g.openCell(c.h, c.w)
	case flagSet:
		g.setFlag(c.h, c.w)
	case flagRemove:
		g.removeFlag(c.h, c.w)
	}
	return false, nil
}

// openCell open specified Game's cell
func (g *Game) openCell(h, w int) (bool, error) {
	c := g.cells[h][w]
	if c.isOpen {
		return false, nil
	}
	if c.hasBomb {
		g.open(h, w)
		return true, nil
	}
	g.openAdjacentCells(h, w)

	return false, nil
}

// Ends returns whether the game is end
// if return true, game is end
func (g *Game) Ends() bool {
	fmt.Printf("closed cell=%d\n", g.closedCellNum)
	fmt.Printf("bomb=%d\n", g.bombNum)
	fmt.Printf("flag=%d\n", g.setFlagNum)
	return g.closedCellNum == g.bombNum
}

func (g *Game) setBomb(sizeH, sizeW, bobNum int) {
	// set bomb
	rand.Seed(time.Now().UnixNano())
	// FIXME: loop never ends if  w*h < bomb
	for setNum := 0; setNum != bobNum; {
		hb := rand.Intn(sizeH)
		wb := rand.Intn(sizeW)
		if g.cells[hb][wb].hasBomb {
			continue
		}
		g.cells[hb][wb].hasBomb = true
		g.incrementBomb(hb, wb)
		setNum++
	}
}

func (g *Game) incrementBomb(h, w int) {
	for _, d := range directions {
		ho, wo := d.offset()
		h := h + ho
		w := w + wo
		if !g.isInGameArea(h, w) {
			continue
		}

		g.cells[h][w].bomb++
	}
}

func (g *Game) hasBomb(h, w int) bool { return g.cells[h][w].hasBomb }

func (g *Game) open(h, w int) {
	c := g.cells[h][w]
	c.open()
	g.cells[h][w] = c
	g.closedCellNum--
}

func (g *Game) hasAdjacentBomb(h, w int) bool {
	return g.cells[h][w].bomb > 0
}

// openAdjacentCells open adjacent cells starting from (h, w)
func (g *Game) openAdjacentCells(h, w int) {
	g.open(h, w)
	if g.hasAdjacentBomb(h, w) {
		return
	}
	for _, d := range directions {
		ho, wo := d.offset()
		h := h + ho
		w := w + wo
		if !g.isInGameArea(h, w) {
			continue
		}

		targetCell := g.cells[h][w]
		if targetCell.isOpen || targetCell.hasBomb {
			continue
		}
		g.openAdjacentCells(h, w)
	}
}

func (g *Game) isInGameArea(h, w int) bool {
	return h >= 0 && w >= 0 && h < len(g.cells) && w < len(g.cells[0])
}

func (g *Game) setFlag(h, w int) {
	c := g.cells[h][w]
	c.setFlag()
	g.cells[h][w] = c
	if c.hasBomb {
		g.setFlagNum++
	}
}

func (g *Game) removeFlag(h, w int) {
	c := g.cells[h][w]
	c.removeFlag()
	g.cells[h][w] = c
	if c.hasBomb {
		g.setFlagNum--
	}
}
