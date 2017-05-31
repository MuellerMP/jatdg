package game

import (
	"fmt"
	"math"

	"github.com/LokiTheMango/jatdg/game/entities"
	"github.com/LokiTheMango/jatdg/game/input"
	"github.com/LokiTheMango/jatdg/game/level"
	"github.com/LokiTheMango/jatdg/game/render"
	"github.com/LokiTheMango/jatdg/game/tiles"
)

// Game Object
type Game struct {
	screen        Screen
	level         level.Level
	input         input.Keyboard
	DrawRequested bool
	Towers        []entities.Entity
	Spawns        []entities.Entity
	Enemies       []entities.Mob
	camera        entities.Mob
	time          int
	numEnemies    int
}

//Constructor
func New() *Game {
	game := &Game{
		screen:        Screen{},
		input:         input.Keyboard{},
		DrawRequested: false,
		time:          0,
		numEnemies:    0,
	}
	return game
}

func (game *Game) Init(filePath string) {
	game.screen = NewScreen(filePath + "tiles.png")
	level, towerTiles, spawnTiles := level.NewLevel(game.screen.SpriteSheet, filePath+"level.png")
	game.level = level
	game.screen.SetLevel(&game.level)
	game.camera = entities.NewCamera(&game.level)
	game.createTowerEntities(towerTiles)
	game.createSpawnerEntities(spawnTiles)
	game.Enemies = make([]entities.Mob, 500)
}

func (game *Game) createTowerEntities(tiles []*tiles.Tile) {
	game.Towers = make([]entities.Entity, game.level.NumTowers)
	j := 0
	for i, _ := range game.Towers {
		for tiles[j] == nil {
			j++
		}
		game.Towers[i] = entities.NewTower(tiles[j])
	}
}

func (game *Game) createSpawnerEntities(tiles []*tiles.Tile) {
	game.Spawns = make([]entities.Entity, game.level.NumEnemySpawn)
	j := 0
	for i, _ := range game.Spawns {
		for tiles[j] == nil {
			j++
		}
		game.Spawns[i] = entities.NewEnemySpawn(tiles[j].X, tiles[j].Y)
	}
}

func (game *Game) Update() {
	game.spawnEnemies()
	for _, tower := range game.Towers {
		tower.Update()
	}
	game.moveObjects()
	game.clearScreen()
	game.render()
	for _, enemy := range game.Enemies {
		if enemy != nil {
			tile := *enemy.GetTile()
			game.screen.RenderMob(enemy.GetX(), enemy.GetY(), tile)
		}
	}
}

func (game *Game) spawnEnemies() {
	game.time++
	if game.time%1000 == 0 {
		for _, spawn := range game.Spawns {
			if spawn != nil {
				fmt.Println("created Enemy")
				x := spawn.GetX()
				y := spawn.GetY()
				tile := game.level.CreateEnemy(x, y)
				game.Enemies[game.numEnemies] = entities.NewEnemy(tile)
				game.numEnemies++
			}
		}
	}
}

func (game *Game) render() {
	x := game.camera.GetX()
	y := game.camera.GetY()
	game.screen.RenderLevel(x, y)
}

func (game *Game) clearScreen() {
	game.screen.ClearScreen()
}

func (game *Game) moveObjects() {
	xOffset := 0
	yOffset := 0
	if game.input.Up {
		yOffset--
	}
	if game.input.Down {
		yOffset++
	}
	if game.input.Left {
		xOffset--
	}
	if game.input.Right {
		xOffset++
	}
	if xOffset != 0 || yOffset != 0 {
		game.camera.Move(xOffset, yOffset)
		game.moveWithCollisionCheck(xOffset, 0)
		game.moveWithCollisionCheck(0, yOffset)
	}
}

func (game *Game) moveWithCollisionCheck(xa int, ya int) {
	for _, enemy := range game.Enemies {
		if enemy != nil {
			x := int((float64(enemy.GetX()+xa) / 128) + 0.5)
			y := (int((float64(enemy.GetY()+ya) / 32) + 0.5)) * game.level.Width
			if !game.level.Tiles[x+y].TileProperties.IsSolid {
				fmt.Println("moving enemy")
				enemy.Move(xa<<2, ya)
			} else {
				fmt.Println("solid tile")
				fmt.Println(game.level.Tiles[x+y].X)
				fmt.Println(game.level.Tiles[x+y].Y)
			}
		}
	}
}

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func (game *Game) UpdateInput(newInput input.Keyboard) {
	game.input = newInput
}

func (game *Game) GetSpriteSheet() render.SpriteSheet {
	return game.screen.SpriteSheet
}

func (game *Game) GetPixelArray() []byte {
	return game.screen.PixelArray
}
