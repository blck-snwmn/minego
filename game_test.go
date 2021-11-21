package minego

import (
	"testing"
)

func TestGame_open(t *testing.T) {
	g := NewGame(10, 10, 20)
	g.Show()
	g.OpenCell(2, 1)
	g.Show()
}
