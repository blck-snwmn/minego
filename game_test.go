package minego

import (
	"bytes"
	"testing"
)

func TestGame_setBomb(t *testing.T) {
	g, _ := NewGame(3, 3, 1, bytes.NewBuffer([]byte{}))

	for hi := 0; hi < g.height(); hi++ {
		for wi := 0; wi < g.weight(); wi++ {
			if g.cells[hi][wi].hasBomb {
				t.Fatal("unexpected bomb")
			}
		}
	}
	if _, err := g.openCell(0, 0); err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	hasBomb := false
out:
	for hi := 0; hi < g.height(); hi++ {
		for wi := 0; wi < g.weight(); wi++ {
			if g.cells[hi][wi].hasBomb {
				hasBomb = true
				break out
			}
		}
	}
	if !hasBomb {
		t.Fatal("expected the bomb to be set, but not bomb")
	}
}
