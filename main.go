package main

import (
	"fmt"
	"image/color"
	"runtime"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type App struct{}

// Global variables
var (
	frameCount int32
	windowSize = rl.Vector2{
		X: 1280,
		Y: 720,
	}
	Player = Entity{
		transform: rl.Rectangle{
			X:      windowSize.X / 4,
			Y:      windowSize.Y - (33 * 5) - 50,
			Width:  33 * 5,
			Height: 33 * 5,
		},
		world_position: rl.Vector2{X: windowSize.X / 4, Y: windowSize.Y - 33*5 - 50},
		speed:          20,
		sprint_speed:   30,
		isRunnging:     false,
		isFacingRight:  true,
		isAttacking:    false,
		attackType:     0,
		scale:          5,
		sprite:         nil,
		sprite_set:     0,
		entity_type:    PLAYER,
		hitbox:         rl.Rectangle{X: windowSize.X/4 + (33 * 5 / 4), Y: windowSize.Y - (33 * 6), Width: (33 * 5) / 2, Height: 33 * 5},
	}
	deltaTime           float32  = 0
	betweenAttacksTimer float32  = TIME_FOR_ATTACK_2
	isHitboxDebug       bool     = false
	arr_entities        []Entity = make([]Entity, 0)
)

// Background layers
var (
	bg_layer_0          *rl.Texture2D
	bg_layer_1          *rl.Texture2D
	bg_layer_2          *rl.Texture2D
	bg_layer_3          *rl.Texture2D
	bg_layer_4          *rl.Texture2D
	bg_layer_5          *rl.Texture2D
	bg_layer_6          *rl.Texture2D
	bg_layer_7          *rl.Texture2D
	bg_layer_8          *rl.Texture2D
	scrolling_back      float32 = 0.0
	scrolling_clouds    float32 = 0.0
	scrolling_backmid   float32 = 0.0
	scrolling_mid       float32 = 0.0
	scrolling_midmid    float32 = 0.0
	scrolling_midmidmid float32 = 0.0
	scrolling_midfore   float32 = 0.0
	scrolling_fore      float32 = 0.0
	scrolling_forefore  float32 = 0.0
)

// Constants
const (
	TIME_FOR_ATTACK_2 float32 = 0.1
	ATTACK_NONE       int32   = 0
	ATTACK_LIGHT      int32   = 1
	ATTACK_LIGHT_BACK int32   = 2
	ATTACK_HEAVY      int32   = 3
	VELOCITY
)

func main() {
	app := App{}
	rl.InitWindow(int32(windowSize.X), int32(windowSize.Y), "Demon Slayer")
	defer rl.CloseWindow()

	app.Setup()

	// Main loop for window / game
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// Update deltaTime
		deltaTime = 1.0 / float32(rl.GetFPS())

		// Call respective methods for udpating and drawing
		app.Update()
		app.Draw()

		rl.EndDrawing()
		if frameCount%500 == 0 {
			runtime.GC()
		}
	}
}

func (a *App) Setup() {
	// Set target FPS
	rl.SetTargetFPS(60)

	rl.SetTraceLogLevel(rl.LogError)

	// Load sprites for background
	bg_layer_0 = get_texture(backgroundSprites[0])
	bg_layer_1 = get_texture(backgroundSprites[1])
	bg_layer_2 = get_texture(backgroundSprites[2])
	bg_layer_3 = get_texture(backgroundSprites[3])
	bg_layer_4 = get_texture(backgroundSprites[4])
	bg_layer_5 = get_texture(backgroundSprites[5])
	bg_layer_6 = get_texture(backgroundSprites[6])
	bg_layer_7 = get_texture(backgroundSprites[7])
	bg_layer_8 = get_texture(backgroundSprites[8])
}

func (a *App) Update() {
	// Update app frame count
	frameCount++

	// Check player input for character movement
	if Player.isAttacking {
		if Player.attackType == ATTACK_LIGHT && Player.animation_phase < len(sprite_table[Player.sprite_set][ATTACK_ANIM])-1 {
			a.play_attack_animation()
		} else if Player.attackType == ATTACK_LIGHT_BACK && Player.animation_phase < len(sprite_table[Player.sprite_set][ATTACK2_ANIM])-1 {
			a.play_attack_animation()
		} else if Player.attackType == ATTACK_HEAVY && Player.animation_phase < len(sprite_table[Player.sprite_set][ATTACK3_ANIM])-1 {
			a.play_attack_animation()
		} else if betweenAttacksTimer > 0 {
			if rl.IsMouseButtonDown(rl.MouseButtonLeft) && Player.attackType != ATTACK_LIGHT_BACK {
				Player.attackType = ATTACK_LIGHT_BACK
				reset_for_animation()
				a.play_attack_animation()
			}
			betweenAttacksTimer -= 1 * deltaTime
		} else {
			Player.isAttacking = false
			reset_for_animation()
			Player.attackType = ATTACK_NONE
			betweenAttacksTimer = TIME_FOR_ATTACK_2
		}
	} else if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		Player.attackType = ATTACK_LIGHT
		frameCount = 0
		Player.isAttacking = true
	} else if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
		Player.attackType = ATTACK_HEAVY
		reset_for_animation()
		Player.isAttacking = true
	} else if rl.IsKeyDown(rl.KeyA) {
		a.run_left()
	} else if rl.IsKeyDown(rl.KeyD) {
		a.run_right()
	} else {
		Player.isRunnging = false
		Player.isAttacking = false
	}
	if rl.IsKeyDown(rl.KeyLeftShift) && !Player.isAttacking {
		Player.current_speed = Player.sprint_speed
	} else {
		Player.current_speed = Player.speed
	}

	// Check player input for debug
	if rl.IsKeyPressed(rl.KeyF1) { // Show hitboxes
		isHitboxDebug = !isHitboxDebug
	}
	if rl.IsKeyPressed(rl.KeyO) { // Add new entity to entity list
		spawn_entity("Demon", rl.Rectangle{X: float32(rl.GetRandomValue(0, int32(windowSize.X))), Y: -50, Width: 44, Height: 44}, ENEMY, 1, 6)
	}
	if rl.IsKeyPressed(rl.KeyL) {
		arr_entities = []Entity{}
	}

	// Reset / idle player when not doing anything
	if !Player.isRunnging && !Player.isAttacking {
		a.idle_player()
	}

	// Constantly update hitboxes to keep up with entity transform
	udpate_hitbox(&Player)

	// Debug print output
	print_debug_info()
}

func (a *App) Draw() {
	draw_background()

	// Draw entities
	for entity := range arr_entities {
		arr_entities[entity].draw()
	}

	// Draw the player sprite
	if Player.isFacingRight {
		rl.DrawTextureRec(*Player.sprite, rl.Rectangle{X: 0, Y: 0, Width: Player.transform.Width, Height: Player.transform.Height}, rl.Vector2{X: Player.transform.X, Y: Player.transform.Y}, color.RGBA{255, 255, 255, 255})
	} else {
		rl.DrawTextureRec(*Player.sprite, rl.Rectangle{X: 0, Y: 0, Width: -Player.transform.Width, Height: Player.transform.Height}, rl.Vector2{X: Player.transform.X, Y: Player.transform.Y}, color.RGBA{255, 255, 255, 255})
	}

	if isHitboxDebug {
		rl.DrawRectangleRoundedLines(Player.hitbox, 0, 0, 1, rl.Red)
		for index := range arr_entities {
			rl.DrawRectangleLines(int32(arr_entities[index].world_position.X), int32(arr_entities[index].world_position.Y), int32(arr_entities[index].transform.Width), int32(arr_entities[index].transform.Height), rl.Red)
		}
	}

	// Smoothly move clouds constantly
	scrolling_clouds += 1

	// Draw foreground over player
	rl.DrawTextureRec(*bg_layer_8, rl.Rectangle{X: float32(scrolling_forefore), Y: 0, Width: float32(bg_layer_6.Width), Height: float32(bg_layer_6.Height)}, rl.Vector2{X: 0, Y: 50}, color.RGBA{255, 255, 255, 255})

	// Draw fps for debugging
	rl.DrawFPS(10, 10)
}

// Player specific behaviours
func (a *App) idle_player() {
	Player.state = IDLE_ANIM
	sprite := get_anim_sprite(&Player, 8)
	sprite.Width = Player.transform.ToInt32().Width
	sprite.Height = Player.transform.ToInt32().Height
	Player.sprite = sprite
}
func (a *App) run_left() {
	Player.state = RUN_ANIM
	Player.world_position.X -= Player.current_speed
	Player.isRunnging = true
	Player.isFacingRight = false
	sprite := get_anim_sprite(&Player, 100/int(Player.current_speed))
	sprite.Width = Player.transform.ToInt32().Width
	sprite.Height = Player.transform.ToInt32().Height
	Player.sprite = sprite
	scroll_background(-1)
}
func (a *App) run_right() {
	Player.state = RUN_ANIM
	Player.world_position.X += Player.current_speed
	Player.isRunnging = true
	Player.isFacingRight = true
	sprite := get_anim_sprite(&Player, 100/int(Player.current_speed))
	sprite.Width = Player.transform.ToInt32().Width
	sprite.Height = Player.transform.ToInt32().Height
	Player.sprite = sprite
	scroll_background(1)
}
func udpate_hitbox(entity *Entity) {
	entity.hitbox.Width = float32(entity.sprite.Width)
	entity.hitbox.Height = float32(entity.sprite.Height)
	entity.hitbox.X = Player.transform.X
	entity.hitbox.Y = Player.transform.Y
}
func (a *App) play_attack_animation() {
	sprite := rl.Texture2D{}
	switch Player.attackType {
	case ATTACK_LIGHT:
		Player.state = ATTACK_ANIM
		sprite = *get_anim_sprite(&Player, 4)
	case ATTACK_LIGHT_BACK:
		Player.state = ATTACK2_ANIM
		sprite = *get_anim_sprite(&Player, 2)
	case ATTACK_HEAVY:
		Player.state = ATTACK3_ANIM
		sprite = *get_anim_sprite(&Player, 4)
	}
	sprite.Width = Player.transform.ToInt32().Width
	sprite.Height = Player.transform.ToInt32().Height
	Player.sprite = &sprite
}

func draw_background() {
	if int32(scrolling_back) <= -bg_layer_0.Width*2 {
		scrolling_back = 0
	}
	if int32(scrolling_backmid) <= -bg_layer_1.Width*2 {
		scrolling_backmid = 0
	}
	if int32(scrolling_mid) <= -bg_layer_2.Width*2 {
		scrolling_mid = 0
	}
	if int32(scrolling_midfore) <= -bg_layer_3.Width*2 {
		scrolling_midfore = 0
	}
	if int32(scrolling_fore) <= -bg_layer_4.Width*2 {
		scrolling_fore = 0
	}
	if int32(scrolling_forefore) <= -bg_layer_5.Width*2 {
		scrolling_forefore = 0
	}

	bg_layer_0.Width = int32(windowSize.X)
	bg_layer_0.Height = int32(windowSize.Y)
	bg_layer_1.Width = int32(windowSize.X)
	bg_layer_1.Height = int32(windowSize.Y)
	bg_layer_2.Width = int32(windowSize.X)
	bg_layer_2.Height = int32(windowSize.Y)
	bg_layer_3.Width = int32(windowSize.X)
	bg_layer_3.Height = int32(windowSize.Y)
	bg_layer_4.Width = int32(windowSize.X)
	bg_layer_4.Height = int32(windowSize.Y)
	bg_layer_5.Width = int32(windowSize.X)
	bg_layer_5.Height = int32(windowSize.Y)
	bg_layer_6.Width = int32(windowSize.X)
	bg_layer_6.Height = int32(windowSize.Y)
	bg_layer_7.Width = int32(windowSize.X)
	bg_layer_7.Height = int32(windowSize.Y)
	bg_layer_8.Width = int32(windowSize.X)
	bg_layer_8.Height = int32(windowSize.Y)

	rl.DrawTextureRec(*bg_layer_0, rl.Rectangle{X: 0, Y: 0, Width: windowSize.X, Height: windowSize.Y}, rl.Vector2{X: 0, Y: 0}, color.RGBA{255, 255, 255, 255})
	rl.DrawTextureRec(*bg_layer_1, rl.Rectangle{X: float32(scrolling_back), Y: 0, Width: windowSize.X, Height: windowSize.Y}, rl.Vector2{X: 0, Y: 0}, color.RGBA{255, 255, 255, 255})
	rl.DrawTextureRec(*bg_layer_2, rl.Rectangle{X: float32(scrolling_clouds), Y: 0, Width: windowSize.X, Height: windowSize.Y}, rl.Vector2{X: 0, Y: 0}, color.RGBA{255, 255, 255, 255})
	rl.DrawTextureRec(*bg_layer_3, rl.Rectangle{X: float32(scrolling_backmid), Y: 0, Width: windowSize.X, Height: windowSize.Y}, rl.Vector2{X: 0, Y: 0}, color.RGBA{255, 255, 255, 255})
	rl.DrawTextureRec(*bg_layer_4, rl.Rectangle{X: float32(scrolling_mid), Y: 0, Width: windowSize.X, Height: windowSize.Y}, rl.Vector2{X: 0, Y: 0}, color.RGBA{255, 255, 255, 255})
	rl.DrawTextureRec(*bg_layer_5, rl.Rectangle{X: float32(scrolling_midmid), Y: 0, Width: windowSize.X, Height: windowSize.Y}, rl.Vector2{X: 0, Y: 0}, color.RGBA{255, 255, 255, 255})
	rl.DrawTextureRec(*bg_layer_6, rl.Rectangle{X: float32(scrolling_midmidmid), Y: 0, Width: windowSize.X, Height: windowSize.Y}, rl.Vector2{X: 0, Y: 0}, color.RGBA{255, 255, 255, 255})
	rl.DrawTextureRec(*bg_layer_7, rl.Rectangle{X: float32(scrolling_midfore), Y: 0, Width: windowSize.X, Height: windowSize.Y}, rl.Vector2{X: 0, Y: 0}, color.RGBA{255, 255, 255, 255})
}
func scroll_background(scroll_factor float32) {
	scrolling_back += Player.current_speed / 20 * scroll_factor
	scrolling_backmid += (Player.current_speed / 12) * scroll_factor
	scrolling_mid += (Player.current_speed / 6) * scroll_factor
	scrolling_midmid += (Player.current_speed / 5) * scroll_factor
	scrolling_midmidmid += (Player.current_speed / 4) * scroll_factor
	scrolling_midfore += (Player.current_speed / 3) * scroll_factor
	scrolling_fore += (Player.current_speed / 2) * scroll_factor
	scrolling_forefore += (Player.current_speed) * scroll_factor
}

// Utility functions
func get_texture(sprite *rl.Image) *rl.Texture2D {
	result := rl.LoadTextureFromImage(sprite)
	return &result
}
func reset_for_animation() {
	Player.animation_phase = 0
	frameCount = 0
}
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %f MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %f MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %f MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}
func bToMb(b uint64) float32 {
	return float32(b) / 1024.0 / 1024.0
}
func print_debug_info() {
	fmt.Print("\033[H\033[2J")
	PrintMemUsage()
	fmt.Printf("Player position: x {%f} y{%f} \n", Player.transform.X, Player.transform.Y)
	fmt.Printf("Player world position: x: {%f} y: {%f}\n", Player.world_position.X, Player.world_position.Y)
	println("Player is using attack: ", Player.attackType)
	println("Player is attacking: ", Player.isAttacking)
	println("Player is running: ", Player.isRunnging)
	println("Player animation_phase: ", Player.animation_phase)
	fmt.Printf("Between attacks timer: %f \n", betweenAttacksTimer)
	println("Entities array length: ", len(arr_entities))
	// fmt.Printf("Entites array:  %+v", arr_entities)
}

// Engine logic
func spawn_entity(name string, transform rl.Rectangle, entity_type EntityType, sprite_set int, scale int) Entity {
	// Make the entity face the player
	is_facing_right := true
	if transform.X >= Player.world_position.X {
		is_facing_right = false
	}

	transform.Width *= float32(scale)
	transform.Height *= float32(scale)
	transform.Y += windowSize.Y - transform.Height

	new_entity := Entity{
		name:           name,
		transform:      transform,
		world_position: rl.Vector2{X: transform.X, Y: transform.Y},
		speed:          5,
		scale:          float32(scale),
		sprite:         nil,
		isFacingRight:  is_facing_right,
		sprite_set:     sprite_set,
		entity_type:    entity_type,
		state:          IDLE_ANIM,
	}
	arr_entities = append(arr_entities, new_entity)
	return new_entity
}
func kill_entity() {
	// Kill an entity
}
