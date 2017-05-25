package entities

import (
	"github.com/LokiTheMango/jatdg/game/level"
)

type Camera struct {
	x       int
	y       int
	removed bool
	level   level.Level
}

func NewCamera() Mob {
	return &Camera{}
}

func (camera *Camera) Move(xa int, ya int) {
	camera.x += xa
	camera.y += ya
}
func (camera *Camera) Update() {

}
func (camera *Camera) Render() {

}
func (camera *Camera) GetX() int {
	return camera.x
}
func (camera *Camera) GetY() int {
	return camera.y
}
func (camera *Camera) Remove() {

}
func (camera *Camera) IsRemoved() bool {
	return false
}