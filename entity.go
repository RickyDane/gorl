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
	max_health        float32
	health            float32
	attack_damage     float32
	entity_type       EntityType
	position          rl.Vector2
	world_position    rl.Vector2
	speed             float32
	sprint_speed      float32
	current_speed     float32
	isRunning         bool
	isAttacking       bool
	attackType        int32
	isFacingRight     bool
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
	e.position.X = e.world_position.X - Player.world_position.X/3
	e.position.Y = windowSize.Y - e.current_sprite.height*4 - 25

	if e.was_hit {
		rl.DrawTexturePro(
			sprite_atlas,
			get_anim_transform(e, 6),
			rl.Rectangle{X: e.position.X, Y: e.position.Y, Width: e.current_sprite.width * 4, Height: e.current_sprite.height * 4},
			rl.Vector2{X: 0, Y: 0},
			0,
			rl.Red,
		)
		e.hit_cooldown -= 1 * deltaTime
		if e.hit_cooldown <= 0 {
			e.was_hit = false
			e.hit_cooldown = ENTITY_HIT_COOLDOWN
		}
	} else {
		rl.DrawTexturePro(
			sprite_atlas,
			get_anim_transform(e, 6),
			rl.Rectangle{X: e.position.X, Y: e.position.Y, Width: e.current_sprite.width * 4, Height: e.current_sprite.height * 4},
			rl.Vector2{X: 0, Y: 0},
			0,
			rl.White,
		)
	}

	if isHitboxDebug {
		// Draw entity name centered over head
		rl.DrawRectangleLines(int32(e.hitbox.X), int32(e.hitbox.Y), int32(e.hitbox.Width), int32(e.hitbox.Height), rl.Red)
		rl.DrawText(e.name, int32(e.position.X), int32(e.hitbox.Y), 5, rl.Red)
	}

	e.draw_healthbar()
}

func (e *Entity) update_hitbox() {
	if e.entity_type == PLAYER && !Player.isAttacking {
		e.hitbox.Width = e.current_sprite.width
		e.hitbox.X = e.position.X + e.current_sprite.width
		e.hitbox.Height = e.current_sprite.height * 3
		e.hitbox.Y = Player.position.Y
	} else if e.entity_type == DEMON {
		e.hitbox.Width = e.current_sprite.width * 2
		e.hitbox.X = e.position.X + e.current_sprite.width
		e.hitbox.Height = e.current_sprite.height * 4
		e.hitbox.Y = e.position.Y
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
	width := int32(e.hitbox.Width)
	height := int32(2)
	percentage := int32(e.max_health / 100 * e.health)
	rl.DrawRectangle(e.hitbox.ToInt32().X, e.hitbox.ToInt32().Y-5, width, height, rl.White)
	rl.DrawRectangle(e.hitbox.ToInt32().X, e.hitbox.ToInt32().Y-5, width*percentage/100, height, rl.Red)
	rl.DrawText(strconv.FormatInt(int64(percentage), 10)+"%", e.hitbox.ToInt32().X, e.hitbox.ToInt32().Y-15, 10, rl.White)
}

func (e *Entity) print_debug_info() {
	fmt.Print("\033[H\033[2J")
	println("##-- Entity --##")
	println("Id: ", e.id)
	println("Name: ", e.name)
	fmt.Printf("Health: %f\n", e.health)
	println("is_colliding: ", e.is_colliding)
}
