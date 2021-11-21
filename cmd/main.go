package main

import (
	"fmt"
	"os"

	"github.com/blck-snwmn/minego"
)

func main() {
	var h, w, bombNum int
	fmt.Scan(&h, &w, &bombNum)
	game := minego.NewGame(h, w, bombNum, os.Stdout)

	var ih, iw int

	game.Show()
	for {
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
		game.Show()
		if exploded {
			fmt.Println("bomb is exploded. game over.")
			return
		}
		if game.Ends() {
			fmt.Println("you win")
			return
		}
	}
}
