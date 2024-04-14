package main

import (
    "math/rand"
    "time"

    "github.com/nsf/termbox-go"
)

const (
    minRoomSize = 5
    maxRoomSize = 10
    minHallSize = 3
    maxHallSize = 7
    mapWidth    = 100
    mapHeight   = 40
)

type position struct {
    x, y int
}

func main() {
    err := termbox.Init()
    if err != nil {
        panic(err)
    }
    defer termbox.Close()

    rand.Seed(time.Now().UnixNano())

    for {
        // Create a new game map
        gameMap, playerPos, exitPos := createMap()

        // Game loop
        for {
            // Clear the screen
            termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

            // Draw the game map
            drawMap(gameMap)

            // Draw the player
            drawHexagon(playerPos.x, playerPos.y, '&', termbox.ColorWhite)

            // Draw the exit
            drawHexagon(exitPos.x, exitPos.y, 'X', termbox.ColorRed|termbox.AttrBold)

            // Flush the screen
            termbox.Flush()

            // Get input event
            ev := termbox.PollEvent()
            if ev.Type == termbox.EventKey {
                switch ev.Key {
                case termbox.KeyCtrlC:
                    return
                case termbox.KeyArrowUp:
                    playerPos = movePlayer(playerPos, 0, -1, gameMap)
                case termbox.KeyArrowDown:
                    playerPos = movePlayer(playerPos, 0, 1, gameMap)
                case termbox.KeyArrowLeft:
                    playerPos = movePlayer(playerPos, -1, 0, gameMap)
                case termbox.KeyArrowRight:
                    playerPos = movePlayer(playerPos, 1, 0, gameMap)
                default:
                    switch ev.Ch {
                    case 'w', 'W':
                        playerPos = movePlayer(playerPos, 0, -1, gameMap)
                    case 's', 'S':
                        playerPos = movePlayer(playerPos, 0, 1, gameMap)
                    case 'a', 'A':
                        playerPos = movePlayer(playerPos, -1, 0, gameMap)
                    case 'd', 'D':
                        playerPos = movePlayer(playerPos, 1, 0, gameMap)
                    }
                }
            }

            // Check if the player reached the exit
            if playerPos == exitPos {
                break
            }
        }
    }
}

func createMap() ([][]rune, position, position) {
    // Create the game map with random sized rooms and hallways
    gameMap := make([][]rune, mapHeight)
    for i := range gameMap {
        gameMap[i] = make([]rune, mapWidth)
        for j := range gameMap[i] {
            gameMap[i][j] = ' '
        }
    }

    // Generate random rooms
    numRooms := rand.Intn(10) + 5
    var rooms [][]position
    for i := 0; i < numRooms; i++ {
        roomSize := rand.Intn(maxRoomSize-minRoomSize+1) + minRoomSize
        roomPos := position{
            x: rand.Intn(mapWidth-roomSize) + 1,
            y: rand.Intn(mapHeight-roomSize) + 1,
        }
        room := make([]position, 0)
        for y := roomPos.y; y < roomPos.y+roomSize; y += 2 {
            for x := roomPos.x; x < roomPos.x+roomSize; x++ {
                if isValidHexCell(x, y, gameMap) {
                    gameMap[y][x] = '.'
                    room = append(room, position{x, y})
                }
            }
        }
        rooms = append(rooms, room)
    }

    // Generate hallways
    for i := 0; i < len(rooms)-1; i++ {
        room1 := rooms[i]
        room2 := rooms[i+1]
        hall := generateHallway(room1[rand.Intn(len(room1))], room2[rand.Intn(len(room2))], gameMap)
        for _, pos := range hall {
            gameMap[pos.y][pos.x] = '.'
        }
    }

    // Find a valid starting position for the player
    var playerStart position
    for _, room := range rooms {
        playerStart = room[rand.Intn(len(room))]
        if gameMap[playerStart.y][playerStart.x] == '.' {
            break
        }
    }

    // Find a valid position for the exit
    var exitPos position
    for {
        room := rooms[rand.Intn(len(rooms))]
        exitPos = room[rand.Intn(len(room))]
        if gameMap[exitPos.y][exitPos.x] == '.' && exitPos != playerStart {
            break
        }
    }

    return gameMap, playerStart, exitPos
}

func generateHallway(start, end position, gameMap [][]rune) []position {
    var hallway []position
    xDiff := end.x - start.x
    yDiff := end.y - start.y

    xDir := 1
    if xDiff < 0 {
        xDir = -1
    }
    yDir := 1
    if yDiff < 0 {
        yDir = -1
    }

    // Generate horizontal part of the hallway
    for x := start.x; x != end.x; x += xDir {
        if isValidHexCell(x, start.y, gameMap) {
            hallway = append(hallway, position{x, start.y})
        }
        if isValidHexCell(x, start.y+yDir, gameMap) {
            hallway = append(hallway, position{x, start.y + yDir})
        }
    }

    // Generate vertical part of the hallway
    for y := start.y; y != end.y; y += yDir {
        if isValidHexCell(end.x, y, gameMap) {
            hallway = append(hallway, position{end.x, y})
        }
        if isValidHexCell(end.x+xDir, y, gameMap) {
            hallway = append(hallway, position{end.x + xDir, y})
        }
    }

    return hallway
}

func drawMap(gameMap [][]rune) {
    for y := 0; y < len(gameMap); y += 2 {
        for x := 0; x < len(gameMap[y]); x++ {
            if isValidHexCell(x, y, gameMap) {
                drawHexagon(x, y, gameMap[y][x], termbox.ColorDefault)
            }
        }
    }
}

func movePlayer(playerPos position, dx, dy int, gameMap [][]rune) position {
    newPos := position{x: playerPos.x + dx, y: playerPos.y + dy}
    if isValidHexCell(newPos.x, newPos.y, gameMap) && gameMap[newPos.y][newPos.x] != ' ' {
        playerPos = newPos
    }
    return playerPos
}

func isValidHexCell(x, y int, gameMap [][]rune) bool {
    return y >= 0 && y < len(gameMap) && x >= 0 && x < len(gameMap[y]) &&
        ((y%2 == 0 && x%2 == 0) || (y%2 != 0 && x%2 != 0))
}

func drawHexagon(x, y int, ch rune, fg termbox.Attribute) {
    // Calculate the top-left corner of the hexagon
    startX := x * 3
    startY := y * 2

    // Draw the top and bottom lines of the hexagon
    for i := 0; i < 3; i++ {
        termbox.SetCell(startX+i, startY, '-', fg, termbox.ColorDefault)
        termbox.SetCell(startX+i, startY+3, '-', fg, termbox.ColorDefault)
    }

    // Draw the middle lines of the hexagon
    termbox.SetCell(startX, startY+1, '\\', fg, termbox.ColorDefault)
    termbox.SetCell(startX+1, startY+1, ch, fg, termbox.ColorDefault)
    termbox.SetCell(startX+2, startY+1, '/', fg, termbox.ColorDefault)
    termbox.SetCell(startX, startY+2, '/', fg, termbox.ColorDefault)
    termbox.SetCell(startX+1, startY+2, ch, fg, termbox.ColorDefault)
    termbox.SetCell(startX+2, startY+2, '\\', fg, termbox.ColorDefault)
}