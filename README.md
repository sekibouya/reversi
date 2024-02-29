# reversi
A person-to-person reversi game that also allows you to analyze the game.<br>
There are two modes, the first is Game Mode and the second is Analyze Mode.

# DEMO
Game Mode<br>
https://github.com/sekibouya/reversi/assets/99582134/7ed30bae-da7f-4fef-835d-350f9f205f56

Analyze Mode<br>
https://github.com/sekibouya/reversi/assets/99582134/68d0a81c-e45e-4d32-a3f6-bfd21f76ab3c

# Features
It is a simple reversi game made by GO and can analyze the game.

# Requirement
* ebiten/v2
* ebiten/v2/inpututil
* ebiten/v2/text
* ebiten/v2/vector
* golang.org/x/image/font/basicfont"

# Installation
```
go get "github.com/hajimehoshi/ebiten/v2"
go get "github.com/hajimehoshi/ebiten/v2/inpututil"
go get "github.com/hajimehoshi/ebiten/v2/text"
go get "github.com/hajimehoshi/ebiten/v2/vector"
go get "golang.org/x/image/font/basicfont"
```

# Usage
start: Enter go run command in the directory where the program files are stored
```
go run .
```
Basically, a stone can be placed by left-clicking with the mouse.<br>
Thereafter, follow the instructions displayed in the terminal.<br>
Clicking on a cell where a stone can be placed in Analyze Mode will allow you to consider a new move, and returning to the original procedure will take you back to the phase before the split.
