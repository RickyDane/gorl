package main

import rl "github.com/gen2brain/raylib-go/raylib"

type EntityType int
type EntityState int

const (
	ENEMY EntityType = iota
	PLAYER
	INTERACTIVE_WORLD_OJECT
)

const (
	IDLE_ANIM EntityState = iota
	RUN_ANIM
	ATTACK_ANIM
	ATTACK2_ANIM
	ATTACK3_ANIM
	DEATH_ANIM
)

type Entity struct {
	name            string
	entity_type     EntityType
	transform       rl.Rectangle
	world_position  rl.Vector2
	speed           float32
	sprint_speed    float32
	current_speed   float32
	isRunnging      bool
	isAttacking     bool
	attackType      int32
	isFacingRight   bool
	scale           float32
	sprite          *rl.Texture2D
	sprite_set      int
	state           EntityState
	animation_phase int
	hitbox          rl.Rectangle
}

func (e *Entity) draw() {

	// Change the position of the entity relative to the players world position
	e.world_position.X = e.transform.X - Player.world_position.X/3

	// Unload previous texture
	if e.sprite != nil {
		rl.UnloadTexture(*e.sprite)
	}

	// Get sprite for entity and draw it
	sprite := get_anim_sprite(e, 8)
	sprite.Width = e.transform.ToInt32().Width
	sprite.Height = e.transform.ToInt32().Height
	e.sprite = sprite

	if e.isFacingRight {
		rl.DrawTextureRec(*e.sprite, rl.Rectangle{X: 0, Y: 0, Width: e.transform.Width, Height: e.transform.Height}, e.world_position, rl.White)
	} else {
		rl.DrawTextureRec(*e.sprite, rl.Rectangle{X: 0, Y: 0, Width: -e.transform.Width, Height: e.transform.Height}, e.world_position, rl.White)
	}

	// Draw entity name centered over head
	rl.DrawText(e.name, int32(e.world_position.X+float32(rl.MeasureText(e.name, 20))+(e.transform.Width/6)), int32(e.world_position.Y), 20, rl.Red)
}
