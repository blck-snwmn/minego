package main

import (
	"fmt"
	"os"

	"github.com/blck-snwmn/minego"
)

func main() {
	var h, w, bombNum int
	fmt.Scan(&h, &w, &bombNum)
	game, err := minego.NewGame(h, w, bombNum, os.Stdout)
	if err != nil {
		panic(err.Error())
	}
	var ctype string
	var ih, iw int

	game.Show()
	for {
		if _, err := fmt.Scan(&ctype); err != nil {
			fmt.Println(err)
			return
		}
		if _, err := fmt.Scan(&ih); err != nil {
			fmt.Println(err)
			return
		}
		if _, err := fmt.Scan(&iw); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("input (h, w) = (%d, %d)\n", ih, iw)

		cmd, err := minego.NewCommand(ctype, ih, iw)
		if err != nil {
			fmt.Println(err)
			continue
		}
		exploded, err := game.Do(cmd)
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
