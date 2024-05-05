package main

import (
	"container/list"
	"fmt"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type EntityType int
type EnemyType int

const (
	ENEMY EntityType = iota
	PLAYER
	SHOP
	INTERACTIVE_WORLD_OJECT
)
const (
	ENEMY_NONE EnemyType = iota
	DEMON
	// <-- Add more enemy types here
)
const (
	ENTITY_HIT_COOLDOWN float32 = 0.1
)

type Entity struct {
	id                int64
	name              string
	max_health        float32
	health            float32
	attack_damage     float32
	entity_type       EntityType
	enemy_type        EnemyType
	position          rl.Vector2
	size              rl.Vector2
	world_position    rl.Vector2
	speed             float32
	sprint_speed      float32
	current_speed     float32
	isRunning         bool
	isAttacking       bool
	attackType        int32
	isFacingRight     bool
	animation_phase   int64
	frame_count       int64
	hitbox            rl.Rectangle
	was_hit           bool
	hit_cooldown      float32
	is_colliding      bool
	colliding_objects list.List
	current_sprite    Sprite
	xp                int32
	xp_to_reach       int32
	level             int32
	sprite_color      rl.Color
}

func (e *Entity) update() {
	e.frame_count++
	// Update hitbox to keep up with entity transform
	e.update_hitbox()
	e.check_collisions()
	// e.print_debug_info()
}

func (e *Entity) draw() {

	// Change the position of the entity relative to the players world position
	e.position.X = e.world_position.X - Player.world_position.X/2
	e.position.Y = windowSize.Y - e.size.Y - 30

	if e.was_hit {
		draw_sprite(e, 1, rl.Red, 8)
		e.hit_cooldown -= 1 * deltaTime
		if e.hit_cooldown <= 0 {
			e.was_hit = false
			e.hit_cooldown = ENTITY_HIT_COOLDOWN
		}
		if e.health <= 0 {
			kill_entity(e)
		}
	} else {
		draw_sprite(e, 1, e.sprite_color, 8)
	}

	if isHitboxDebug {
		// Draw entity name centered over head
		rl.DrawRectangleLines(int32(e.hitbox.X), int32(e.hitbox.Y), int32(e.hitbox.Width), int32(e.hitbox.Height), rl.Red)
		rl.DrawText(e.name, int32(e.hitbox.X), int32(e.hitbox.Y), 5, rl.Red)
	}
}

func (e *Entity) update_hitbox() {
	if e.entity_type == PLAYER && !Player.isAttacking {
		e.hitbox.Width = e.size.X / 2
		e.hitbox.X = e.position.X + e.hitbox.Width/2
		e.hitbox.Height = e.size.Y
		e.hitbox.Y = e.position.Y
	} else if e.entity_type == ENEMY {
		e.hitbox.Width = e.size.X / 2
		e.hitbox.X = e.position.X + e.hitbox.Width/2
		e.hitbox.Height = e.size.Y
		e.hitbox.Y = e.position.Y
	} else {
		e.hitbox.Width = e.size.X
		e.hitbox.X = e.position.X
		e.hitbox.Height = e.size.Y
		e.hitbox.Y = e.position.Y
	}
}

func (e *Entity) check_collisions() {
	// Checking for collisions with other entities
	// TODO: impl

	// Check for "collision" with mouse position
	if is_entity_colliding(*e, Entity{hitbox: mouse_pos}) && (e.entity_type == INTERACTIVE_WORLD_OJECT || e.entity_type == SHOP) {
		e.sprite_color = rl.LightGray
	} else if is_entity_colliding(*e, Entity{hitbox: mouse_pos}) {
		e.draw_healthbar()
	} else {
		e.sprite_color = rl.White
	}

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
	width := int32(e.hitbox.Width)
	height := int32(5)
	percentage := int32(e.max_health / 100 * e.health)
	rl.DrawRectangle(e.hitbox.ToInt32().X, e.hitbox.ToInt32().Y-10, width, height, rl.White)
	rl.DrawRectangle(e.hitbox.ToInt32().X, e.hitbox.ToInt32().Y-10, width*percentage/100, height, rl.Red)
	rl.DrawText(strconv.FormatInt(int64(percentage), 10)+"%", e.hitbox.ToInt32().X, e.hitbox.ToInt32().Y-20, 10, rl.White)
}

func (e *Entity) hit(damage float32) {
	e.was_hit = true
	e.health -= damage
	rl.PlaySound(HIT)
	if e.health <= 0 {
		e.health = 0
		player_add_xp(25)
	}
}

// :u_functions
// util functions
func (e *Entity) print_debug_info() {
	fmt.Print("\033[H\033[2J")
	println("##-- Entity --##")
	println("Id: ", e.id)
	println("Name: ", e.name)
	fmt.Printf("Health: %f\n", e.health)
	println("is_colliding: ", e.is_colliding)
}
