package main

import (
    "math/rand"
    "time"

    "github.com/nsf/termbox-go"
)

const (
    mapWidth    = 20
    mapHeight   = 20
    targetFPS   = 8
    movementFPS = 1
)

type position struct {
    q, r int
}

var (
    playerPos    position
    gameMap      map[position]rune
    hexNeighbors = []position{{-1, 0}, {-1, 1}, {0, 1}, {1, 0}, {1, -1}, {0, -1}}
    targetPos    position
    hasTargetPos bool
    dashStates   = []rune{'|', '/', '-', '\\'}
    dashIndex    int
    frameCount   int
)

func main() {
    err := termbox.Init()
    if err != nil {
        panic(err)
    }
    defer termbox.Close()

    termbox.SetInputMode(termbox.InputMouse)

    rand.Seed(time.Now().UnixNano())

    initGame()

    go handleInput()

    targetFrameDuration := time.Second / targetFPS

    for {
        startTime := time.Now()

        updateGame()
        draw()

        elapsedTime := time.Since(startTime)
        if elapsedTime < targetFrameDuration {
            time.Sleep(targetFrameDuration - elapsedTime)
        }
    }
}

func initGame() {
    gameMap = createMap()
    playerPos = findEmptySpot()
    hasTargetPos = false
    frameCount = 0
}

func createMap() map[position]rune {
    gameMap := make(map[position]rune)

    // Fill the map with walls
    for q := 0; q < mapWidth; q++ {
        for r := 0; r < mapHeight; r++ {
            gameMap[position{q, r}] = '#'
        }
    }

    // Generate rooms
    numRooms := rand.Intn(5) + 3
    rooms := [][]position{}
    for i := 0; i < numRooms; i++ {
        roomWidth := rand.Intn(5) + 3
        roomHeight := rand.Intn(5) + 3
        roomQ := rand.Intn(mapWidth - roomWidth)
        roomR := rand.Intn(mapHeight - roomHeight)
        room := []position{}
        for q := roomQ; q < roomQ+roomWidth; q++ {
            for r := roomR; r < roomR+roomHeight; r++ {
                gameMap[position{q, r}] = '.'
                room = append(room, position{q, r})
            }
        }
        rooms = append(rooms, room)
    }

    // Connect rooms with hallways
    for i := 0; i < len(rooms)-1; i++ {
        room1Center := rooms[i][len(rooms[i])/2]
        room2Center := rooms[i+1][len(rooms[i+1])/2]
        connectRooms(room1Center, room2Center, gameMap)
    }

    return gameMap
}

func connectRooms(pos1, pos2 position, gameMap map[position]rune) {
    q1, r1 := pos1.q, pos1.r
    q2, r2 := pos2.q, pos2.r

    // Connect horizontally
    if q1 < q2 {
        for q := q1; q <= q2; q++ {
            if gameMap[position{q, r1}] == '#' {
                gameMap[position{q, r1}] = '.'
            }
            if gameMap[position{q, r1 + 1}] == '#' {
                gameMap[position{q, r1 + 1}] = '.'
            }
        }
    } else {
        for q := q2; q <= q1; q++ {
            if gameMap[position{q, r1}] == '#' {
                gameMap[position{q, r1}] = '.'
            }
            if gameMap[position{q, r1 + 1}] == '#' {
                gameMap[position{q, r1 + 1}] = '.'
            }
        }
    }

    // Connect vertically
    if r1 < r2 {
        for r := r1; r <= r2; r++ {
            if gameMap[position{q2, r}] == '#' {
                gameMap[position{q2, r}] = '.'
            }
            if gameMap[position{q2 + 1, r}] == '#' {
                gameMap[position{q2 + 1, r}] = '.'
            }
        }
    } else {
        for r := r2; r <= r1; r++ {
            if gameMap[position{q2, r}] == '#' {
                gameMap[position{q2, r}] = '.'
            }
            if gameMap[position{q2 + 1, r}] == '#' {
                gameMap[position{q2 + 1, r}] = '.'
            }
        }
    }
}

func findEmptySpot() position {
    for {
        q := rand.Intn(mapWidth)
        r := rand.Intn(mapHeight)
        if gameMap[position{q, r}] == '.' {
            return position{q, r}
        }
    }
}

func handleInput() {
    for {
        ev := termbox.PollEvent()
        if ev.Type == termbox.EventKey && ev.Key == termbox.KeyCtrlC {
            termbox.Close()
            return
        } else if ev.Type == termbox.EventMouse && ev.Key == termbox.MouseLeft {
            q := ev.MouseX / 2
            r := ev.MouseY
            if ev.MouseX%2 == 1 && r > 0 {
                r--
            }
            targetPos = position{q, r}
            hasTargetPos = true
        }
    }
}

func updateGame() {
    frameCount++
    if frameCount%movementFPS == 0 && hasTargetPos {
        movePlayerTowardsTarget()
    }
    updateDashIndicator()
}

func movePlayerTowardsTarget() {
    if playerPos == targetPos {
        hasTargetPos = false
        return
    }

    path := findPath(playerPos, targetPos)
    if len(path) > 0 {
        nextPos := path[0]
        if gameMap[nextPos] != '#' {
            playerPos = nextPos
        }
    }
}

func updateDashIndicator() {
    dashIndex = (dashIndex + 1) % len(dashStates)
}

func findPath(start, end position) []position {
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

    for pos, tile := range gameMap {
        color := termbox.ColorDefault
        if tile == '#' {
            color = termbox.ColorDarkGray
        }
        drawTile(pos.q, pos.r, tile, color)
    }

    drawTile(playerPos.q, playerPos.r, '@', termbox.ColorGreen|termbox.AttrBold)

    if hasTargetPos {
        drawTile(targetPos.q, targetPos.r, '*', termbox.ColorYellow)
    }

    drawDashIndicator()

    termbox.Flush()
}

func drawTile(q, r int, tile rune, color termbox.Attribute) {
    x := q*2 + r%2
    y := r
    termbox.SetCell(x, y, tile, color, termbox.ColorDefault)
}

func drawDashIndicator() {
    x := 0
    y := mapHeight
    termbox.SetCell(x, y, dashStates[dashIndex], termbox.ColorWhite, termbox.ColorDefault)
}