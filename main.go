package main

import (
	"fmt"
	"os"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "\x1b[31m%s\x1b[0m\n", err)
		os.Exit(1)
	}
}

func run() error {
	g := NewGame()
	ebiten.SetWindowSize(SCREEN_SIZE, SCREEN_SIZE)
	ebiten.SetWindowTitle("Reversi")
	if err := ebiten.RunGame(g); err != nil {
		return err
	}

	return nil
}
