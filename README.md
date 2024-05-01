# Reversi
A person-to-person reversi game that also allows you to analyze the game.<br>
There are two modes, the first is Game Mode and the second is Analyze Mode.

# DEMO
Game Mode<br>
https://github.com/sekibouya/reversi/assets/99582134/5a6e950e-a41a-44e7-bc2d-6206249ff801

Analyze Mode<br>
https://github.com/sekibouya/reversi/assets/99582134/68d0a81c-e45e-4d32-a3f6-bfd21f76ab3c

# Features
It is a simple reversi game made by GO and can analyze the game.

# Requirement
* image/color
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
Pressing the R key resets the board.<br>
After the game is finished, you can output the game record file by pressing the E key.<br>
After the game is finished, press the A key to switch to Analyze Mode to analyze the game.<br>
Pressing the I key loads the game file and switches to Analyze Mode.<br>
In Analyze Mode, you can advance or go back through the phases by pressing the arrow keys.<br>
Clicking on a cell where a stone can be placed in Analyze Mode will allow you to consider a new move, and returning to the original procedure will take you back to the phase before the split.

![スライド5](https://github.com/sekibouya/reversi/assets/99582134/dc536c48-eb58-4757-b3ae-d877d8fa1a85)
![スライド6](https://github.com/sekibouya/reversi/assets/99582134/608f516f-b4f7-4ecc-8371-794f72c98fb1)
![スライド7](https://github.com/sekibouya/reversi/assets/99582134/30fc509d-19fd-4afe-b13d-100829ff5a44)
![スライド8](https://github.com/sekibouya/reversi/assets/99582134/5f9c913e-e78d-46db-ab1f-08dceb6e21dc)
![スライド9](https://github.com/sekibouya/reversi/assets/99582134/d524be6b-f3e2-4e12-916b-baa46f3f1392)
![スライド10](https://github.com/sekibouya/reversi/assets/99582134/77ce95de-b35a-414a-af84-80aab7bb6181)
![スライド11](https://github.com/sekibouya/reversi/assets/99582134/e46f9da9-d5b7-4d09-b5c4-84649bc18650)
