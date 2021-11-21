package main

import (
	"fmt"

	"github.com/blck-snwmn/minego"
)

func main() {
	w, h := 5, 6
	bobNum := 5
	game := minego.NewGame(h, w, bobNum)

	var ih, iw int

	for {
		game.Show()
		if _, err := fmt.Scan(&ih); err != nil {
			fmt.Println(err)
			return
		}
		if _, err := fmt.Scan(&iw); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("input (h, w) = (%d, %d)\n", ih, iw)
		exploded, err := game.OpenCell(ih, iw)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if exploded {
			fmt.Println("bomb is exploded. game over.")
			return
		}
	}
}
