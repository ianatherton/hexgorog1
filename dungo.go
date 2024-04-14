package main

import (
    "math/rand"
    "time"

    "github.com/nsf/termbox-go"
)

const (
    mapWidth  = 20
    mapHeight = 20
)

type position struct {
    q, r int
}

var (
    playerPos    position
    gameMap      map[position]rune
    hexNeighbors = []position{
        {-1, 0}, {-1, 1}, {0, 1}, {1, 0}, {1, -1}, {0, -1},
    }
    cursorPos position
)

func main() {
    err := termbox.Init()
    if err != nil {
        panic(err)
    }
    defer termbox.Close()

    // Enable mouse tracking
    termbox.SetInputMode(termbox.InputMouse)

    rand.Seed(time.Now().UnixNano())

    initGame()

    for {
        draw()

        ev := termbox.PollEvent()
        if ev.Type == termbox.EventKey && ev.Key == termbox.KeyCtrlC {
            return
        } else if ev.Type == termbox.EventMouse {
            handleMouseInput(ev)
        }

        if playerPos != cursorPos {
            movePlayerTowardsCursor()
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func initGame() {
    gameMap = createMap()
    playerPos = position{mapWidth / 2, mapHeight / 2}
    cursorPos = playerPos
}

func createMap() map[position]rune {
    gameMap := make(map[position]rune)
    for q := 0; q < mapWidth; q++ {
        for r := 0; r < mapHeight; r++ {
            gameMap[position{q, r}] = '.'
        }
    }
    return gameMap
}

func handleMouseInput(ev termbox.Event) {
    if ev.Key == termbox.MouseLeft {
        q := ev.MouseX / 2
        r := ev.MouseY
        if ev.MouseX%2 == 1 && r > 0 {
            r--
        }
        cursorPos = position{q, r}
    } else {
        cursorPos = position{ev.MouseX / 2, ev.MouseY}
    }
}

func movePlayerTowardsCursor() {
    path := findPath(playerPos, cursorPos)
    if len(path) > 0 {
        playerPos = path[0]
    }
}

func findPath(start, end position) []position {
    // Implement your pathfinding algorithm here (e.g., A*)
    // For simplicity, this example returns the direct path
    path := []position{}
    dq := end.q - start.q
    dr := end.r - start.r
    if abs(dq) > abs(dr) {
        stepQ := sign(dq)
        path = append(path, position{start.q + stepQ, start.r})
    } else {
        stepR := sign(dr)
        path = append(path, position{start.q, start.r + stepR})
    }
    return path
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}

func sign(x int) int {
    if x < 0 {
        return -1
    } else if x > 0 {
        return 1
    }
    return 0
}

func draw() {
    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

    for pos := range gameMap {
        drawTile(pos.q, pos.r, '.', termbox.ColorDefault)
    }

    drawTile(playerPos.q, playerPos.r, '@', termbox.ColorGreen|termbox.AttrBold)
    drawTile(cursorPos.q, cursorPos.r, 'X', termbox.ColorRed)

    termbox.Flush()
}

func drawTile(q, r int, tile rune, color termbox.Attribute) {
    x := q*2 + r%2
    y := r
    termbox.SetCell(x, y, tile, color, termbox.ColorDefault)
}