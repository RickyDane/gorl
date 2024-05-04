package main

import (
	"container/list"
	"fmt"
	"strconv"

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

const (
	ENTITY_HIT_COOLDOWN float32 = 0.1
)

type Entity struct {
	id                int64
	name              string
	health            float32
	attack_damage     float32
	entity_type       EntityType
	transform         rl.Rectangle
	world_position    rl.Vector2
	speed             float32
	sprint_speed      float32
	current_speed     float32
	isRunnging        bool
	isAttacking       bool
	attackType        int32
	isFacingRight     bool
	scale             float32
	sprite            *rl.Texture2D
	sprite_set        int
	state             EntityState
	animation_phase   int
	hitbox            rl.Rectangle
	was_hit           bool
	hit_cooldown      float32
	is_colliding      bool
	colliding_objects list.List
	current_sprite    Sprite
}

func (e *Entity) update() {
	// Update hitbox to keep up with entity transform
	e.update_hitbox()
	e.check_collisions()
	e.print_debug_info()
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

	// Draw entity sprite and if the entity was hit, paint it slightly red
	if e.isFacingRight {
		if e.was_hit {
			rl.DrawTextureRec(*e.sprite, rl.Rectangle{X: 0, Y: 0, Width: e.transform.Width, Height: e.transform.Height}, rl.Vector2{X: e.transform.X, Y: e.transform.Y}, rl.Red)
			e.hit_cooldown -= 1 * deltaTime
			if e.hit_cooldown <= 0 {
				e.was_hit = false
				e.hit_cooldown = ENTITY_HIT_COOLDOWN
			}
		} else {
			rl.DrawTextureRec(*e.sprite, rl.Rectangle{X: 0, Y: 0, Width: e.transform.Width, Height: e.transform.Height}, rl.Vector2{X: e.transform.X, Y: e.transform.Y}, rl.White)
		}
	} else {
		if e.was_hit {
			rl.DrawTextureRec(*e.sprite, rl.Rectangle{X: 0, Y: 0, Width: -e.transform.Width, Height: e.transform.Height}, rl.Vector2{X: e.transform.X, Y: e.transform.Y}, rl.Red)
			e.hit_cooldown -= 1 * deltaTime
			if e.hit_cooldown <= 0 {
				e.was_hit = false
				e.hit_cooldown = ENTITY_HIT_COOLDOWN
			}
		} else {
			rl.DrawTextureRec(*e.sprite, rl.Rectangle{X: 0, Y: 0, Width: -e.transform.Width, Height: e.transform.Height}, rl.Vector2{X: e.transform.X, Y: e.transform.Y}, rl.White)
		}
	}

	if isHitboxDebug {
		// Draw entity name centered over head
		rl.DrawText(e.name, int32(e.transform.X+float32(rl.MeasureText(e.name, 20))+(e.hitbox.Width/2)), int32(e.hitbox.Y), 20, rl.Red)
	}

	e.draw_healthbar()
}

func (e *Entity) update_hitbox() {
	if e.entity_type == DEMON {
		e.hitbox.Width = 75
		e.hitbox.X = e.transform.X + e.hitbox.Width + 50
		e.hitbox.Height = float32(e.transform.Height) / 1.5
		e.hitbox.Y = windowSize.Y - e.hitbox.Height - 50
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
	if is_entity_colliding(*e, Player) && !contains(e.colliding_objects, Player) {
		e.colliding_objects.PushBack(&Player)
		e.is_colliding = true
	} else if !is_entity_colliding(*e, Player) {
		e.colliding_objects = *list.New()
		e.is_colliding = false
	}
}

func (e *Entity) draw_healthbar() {
	width := int32(100)
	height := int32(10)
	percentage := int32(e.health / float32(width) * e.health)
	rl.DrawRectangle(e.hitbox.ToInt32().X, e.hitbox.ToInt32().Y-25, width, height, rl.White)
	rl.DrawRectangle(e.hitbox.ToInt32().X, e.hitbox.ToInt32().Y-25, percentage, height, rl.Red)
	rl.DrawText(strconv.FormatInt(int64(percentage), 10)+"%", e.hitbox.ToInt32().X, e.hitbox.ToInt32().Y-25, 10, rl.White)
}

func (e *Entity) print_debug_info() {
	fmt.Print("\033[H\033[2J")
	println("##-- Entity --##")
	println("Id: ", e.id)
	println("Name: ", e.name)
	fmt.Printf("Health: %f\n", e.health)
	println("is_colliding: ", e.is_colliding)
}
