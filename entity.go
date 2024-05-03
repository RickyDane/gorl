package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type EntityType int
type EntityState int

const (
	DEMON EntityType = iota
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
	is_colliding    bool
}

func (e *Entity) update() {
	// Update hitbox to keep up with entity transform
	e.update_hitbox()
	e.check_collisions()
	// e.print_debug_info()
}

func (e *Entity) draw() {

	// Change the position of the entity relative to the players world position
	e.transform.X = e.world_position.X - Player.world_position.X/3

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
		rl.DrawTextureRec(*e.sprite, rl.Rectangle{X: 0, Y: 0, Width: e.transform.Width, Height: e.transform.Height}, rl.Vector2{X: e.transform.X, Y: e.transform.Y}, rl.White)
	} else {
		rl.DrawTextureRec(*e.sprite, rl.Rectangle{X: 0, Y: 0, Width: -e.transform.Width, Height: e.transform.Height}, rl.Vector2{X: e.transform.X, Y: e.transform.Y}, rl.White)
	}

	if isHitboxDebug {
		// Draw entity name centered over head
		rl.DrawText(e.name, int32(e.transform.X+float32(rl.MeasureText(e.name, 20))+(e.hitbox.Width/2)), int32(e.hitbox.Y), 20, rl.Red)
	}
}

func (e *Entity) update_hitbox() {
	if e.entity_type == DEMON {
		e.hitbox.Width = 50
		e.hitbox.X = e.transform.X + e.hitbox.Width + 50
		e.hitbox.Height = float32(e.transform.Height)
		e.hitbox.Y = e.transform.Y
	} else if e.entity_type == PLAYER {
		e.hitbox.Width = 50
		e.hitbox.X = e.transform.X + e.hitbox.Width/2 + 50/1.5
		e.hitbox.Height = float32(e.transform.Height)
		e.hitbox.Y = e.transform.Y
	}
}

func (e *Entity) check_collisions() {
	// Checking for collisions with other entities
	// TODO: impl

	// Check for collisions with player
	if Player.hitbox.X+Player.hitbox.Width >= e.hitbox.X && Player.hitbox.X+Player.hitbox.Width < e.hitbox.X+e.hitbox.Width {
		e.is_colliding = true
	} else if e.hitbox.X+e.hitbox.Width >= Player.hitbox.X && e.hitbox.X+e.hitbox.Width < Player.hitbox.X+Player.hitbox.Width {
		e.is_colliding = true
	} else {
		e.is_colliding = false
	}
}

func (e *Entity) print_debug_info() {
	fmt.Print("\033[H\033[2J")
	println("##-- Entity --##")
	println("Name: ", e.name)
	println("is_colliding: ", e.is_colliding)
}
