package main

import (
	"container/list"
	"fmt"
	"image/color"
	"math/rand"
	"runtime"
	"time"

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
		id:            int64(rl.GetRandomValue(0, 100000)),
		name:          "Player",
		health:        100,
		attack_damage: 10,
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
		hitbox:         rl.Rectangle{X: windowSize.X/4 + (33 * 5 / 4), Y: windowSize.Y - (33*5 - 50), Width: (33 * 2), Height: 33 * 5},
	}
	deltaTime           float32    = 0
	betweenAttacksTimer float32    = TIME_FOR_ATTACK_2
	isHitboxDebug       bool       = false
	arr_entities        []*Entity  = make([]*Entity, 0)
	ls_entities         *list.List = list.New()
	is_fullscreen       bool       = false
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
	// Make window resizable
	// rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(int32(windowSize.X), int32(windowSize.Y), "Demon Slayer")
	defer rl.CloseWindow()

	setup_audio()

	// Play music
	// rl.PlayMusicStream(BG_MUSIC)

	// Runs one time by start of the game
	app.Setup()

	// Main loop for window / game
	for !rl.WindowShouldClose() {
		// Update deltaTime
		deltaTime = 1.0 / float32(rl.GetFPS())

		app.Update()
		rl.UpdateMusicStream(BG_MUSIC)

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// Call respective methods for udpating and drawing
		app.Draw()

		rl.EndDrawing()
		if frameCount%1000 == 0 {
			runtime.GC()
		}
	}
}

func (a *App) Setup() {
	// Set target FPS
	rl.SetTargetFPS(60)
	// Set log level to disable unneccessary output
	rl.SetTraceLogLevel(rl.LogError)

	// Set initial window size
	windowSize.X = float32(rl.GetMonitorWidth(rl.GetCurrentMonitor())) / 2
	windowSize.Y = float32(rl.GetMonitorHeight(rl.GetCurrentMonitor())) / 2

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
	// Update app frame count: important for animations
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
				rl.PlaySound(SWING)
				attack(Player.attack_damage)
			}
			betweenAttacksTimer -= 1 * deltaTime
		} else {
			reset_for_animation()
			Player.isAttacking = false
			betweenAttacksTimer = TIME_FOR_ATTACK_2
			if pt_chance(0.75) && Player.attackType == ATTACK_HEAVY {
				attack(Player.attack_damage * 3)
			} else if Player.attackType == ATTACK_HEAVY {
				attack(Player.attack_damage * 1.5)
			} else {
				attack(Player.attack_damage)
			}
			Player.attackType = ATTACK_NONE
		}
	} else if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		reset_for_animation()
		Player.attackType = ATTACK_LIGHT
		rl.PlaySound(SWING)
		Player.isAttacking = true
	} else if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
		reset_for_animation()
		Player.attackType = ATTACK_HEAVY
		rl.PlaySound(SLASH_STRONG)
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
		spawn_entity("Demon", rl.Rectangle{X: float32(rl.GetRandomValue(0, int32(windowSize.X))), Y: -25, Width: 44, Height: 44}, DEMON, 1, 8)
	}
	if rl.IsKeyPressed(rl.KeyL) {
		// Clear list of entities
		arr_entities = []*Entity{}
	}

	// Check player input for some stuff
	// if rl.IsKeyPressed(rl.KeyF11) {
	// 	toggle_fullscreen()
	// }

	for element := ls_entities.Front(); element != nil; element = element.Next() {
		element.Value.(*Entity).update()
	}

	// Reset / idle player when not doing anything
	if !Player.isRunnging && !Player.isAttacking {
		// Constantly update player hitbox to keep up with player transform
		Player.update_hitbox()
		a.idle_player()
	}

	// Debug print output
	// print_debug_info()
}

func (a *App) Draw() {
	draw_background()

	// Draw entities
	for element := ls_entities.Front(); element != nil; element = element.Next() {
		element.Value.(*Entity).draw()
	}

	// Draw the player sprite
	if Player.isFacingRight {
		rl.DrawTextureRec(*Player.sprite, rl.Rectangle{X: 0, Y: 0, Width: Player.transform.Width, Height: Player.transform.Height}, rl.Vector2{X: Player.transform.X, Y: Player.transform.Y}, color.RGBA{255, 255, 255, 255})
	} else {
		rl.DrawTextureRec(*Player.sprite, rl.Rectangle{X: 0, Y: 0, Width: -Player.transform.Width, Height: Player.transform.Height}, rl.Vector2{X: Player.transform.X, Y: Player.transform.Y}, color.RGBA{255, 255, 255, 255})
	}

	// Draw lines to visualize the hitbox for debugging
	if isHitboxDebug {
		// Player hitbox
		rl.DrawRectangleRoundedLines(Player.hitbox, 0, 0, 1, rl.Red)
		// Entities hitboxes
		for element := ls_entities.Front(); element != nil; element = element.Next() {
			entity := element.Value.(*Entity)
			rl.DrawRectangleLines(int32(entity.hitbox.X), int32(entity.hitbox.Y), int32(entity.hitbox.Width), int32(entity.hitbox.Height), rl.Red)
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
	rl.StopSound(GRASS_RUNNING)
	Player.state = IDLE_ANIM
	sprite := get_anim_sprite(&Player, 8)
	sprite.Width = Player.transform.ToInt32().Width
	sprite.Height = Player.transform.ToInt32().Height
	if Player.sprite != nil {
		rl.UnloadTexture(*Player.sprite)
	}
	Player.sprite = sprite
	Player.hitbox.Width = 50
}
func (a *App) run_left() {
	if !rl.IsSoundPlaying(GRASS_RUNNING) {
		rl.PlaySound(GRASS_RUNNING)
	}
	Player.state = RUN_ANIM
	Player.world_position.X -= Player.current_speed
	Player.isRunnging = true
	Player.isFacingRight = false
	sprite := get_anim_sprite(&Player, 100/int(Player.current_speed))
	sprite.Width = Player.transform.ToInt32().Width
	sprite.Height = Player.transform.ToInt32().Height
	rl.UnloadTexture(*Player.sprite)
	Player.sprite = sprite
	scroll_background(-1)
}
func (a *App) run_right() {
	if !rl.IsSoundPlaying(GRASS_RUNNING) {
		rl.PlaySound(GRASS_RUNNING)
	}
	Player.state = RUN_ANIM
	Player.world_position.X += Player.current_speed
	Player.isRunnging = true
	Player.isFacingRight = true
	sprite := get_anim_sprite(&Player, 100/int(Player.current_speed))
	sprite.Width = Player.transform.ToInt32().Width
	sprite.Height = Player.transform.ToInt32().Height
	rl.UnloadTexture(*Player.sprite)
	Player.sprite = sprite
	scroll_background(1)
}
func (a *App) play_attack_animation() {
	sprite := rl.Texture2D{}
	switch Player.attackType {
	case ATTACK_LIGHT:
		Player.state = ATTACK_ANIM
		sprite = *get_anim_sprite(&Player, 4)
	case ATTACK_LIGHT_BACK:
		Player.state = ATTACK2_ANIM
		sprite = *get_anim_sprite(&Player, 3)
	case ATTACK_HEAVY:
		Player.state = ATTACK3_ANIM
		sprite = *get_anim_sprite(&Player, 6)
	}
	sprite.Width = Player.transform.ToInt32().Width
	sprite.Height = Player.transform.ToInt32().Height
	rl.UnloadTexture(*Player.sprite)
	Player.sprite = &sprite
	if Player.isFacingRight {
		Player.hitbox.Width = 100
		Player.hitbox.X = Player.transform.X + 100/2
	} else {
		Player.hitbox.X = Player.transform.X
		Player.hitbox.Width = 100
	}

}
func attack(attack_damage float32) {
	for element := ls_entities.Front(); element != nil; element = element.Next() {
		element := element.Value.(*Entity)
		if is_entity_colliding(Player, *element) {
			element.was_hit = true
			element.health -= attack_damage
			if element.health <= 0 {
				kill_entity(element)
				attack(attack_damage)
				break
			}
		}
	}
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

	// Unfortunately we have to do this so the background
	// is scaled properly to cover the window
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
	fmt.Printf("Player hitbox %v", Player.hitbox)
	println("Player is using attack: ", Player.attackType)
	println("Player is attacking: ", Player.isAttacking)
	println("Player is running: ", Player.isRunnging)
	println("Player animation_phase: ", Player.animation_phase)
	fmt.Printf("Between attacks timer: %f \n", betweenAttacksTimer)
	println("Is audio device ready: ", rl.IsAudioDeviceReady())
	println("Entities array length: ", len(arr_entities))
	println("Entities list length", ls_entities.Len())
	// fmt.Printf("Entites array:  %+v", arr_entities)
}
func contains(list list.List, target any) bool {
	for item := list.Front(); item != nil; item = list.Back().Next() {
		if item == target {
			return true
		}
	}
	return false
}
func pt_chance(percentage float32) bool {
	_rand := rand.New(rand.NewSource(time.Now().UnixMicro()))
	randomNumber := _rand.Intn(100) + 1
	return randomNumber <= int(percentage)
}

// Current not in use or worked on
// func toggle_fullscreen() {
// 	if !is_fullscreen {
// 		// rl.MaximizeWindow()
// 		windowSize.X = float32(rl.GetMonitorWidth(rl.GetCurrentMonitor()))
// 		windowSize.Y = float32(rl.GetMonitorHeight(rl.GetCurrentMonitor()))
// 		Player.scale = 9
// 		rl.MaximizeWindow()
// 		// rl.ToggleFullscreen()
// 		// rl.SetWindowSize(int(windowSize.X), int(windowSize.Y))
// 	} else {
// 		windowSize.X = float32(rl.GetMonitorWidth(rl.GetCurrentMonitor())) / 2
// 		windowSize.Y = float32(rl.GetMonitorHeight(rl.GetCurrentMonitor())) / 2
// 		Player.scale = 5
// 		rl.SetWindowSize(int(windowSize.X), int(windowSize.Y))
// 	}
// 	// Update player position
// 	Player.transform.Width = (33 * Player.scale)
// 	Player.transform.Height = (33 * Player.scale)
// 	Player.transform.Y = windowSize.Y - Player.transform.Height - 100
// 	is_fullscreen = !is_fullscreen
// }

// Engine logic
func setup_audio() {
	// Init audio
	rl.InitAudioDevice()
	// Load background music
	BG_MUSIC = rl.LoadMusicStream("assets/sounds/background-music.mp3")
	// Load player sounds
	SLASH = rl.LoadSound("assets/sounds/slash.wav")
	SWING = rl.LoadSound("assets/sounds/swing.wav")
	SLASH_STRONG = rl.LoadSound("assets/sounds/slash_strong.wav")
	GRASS_RUNNING = rl.LoadSound("assets/sounds/grass_running.wav")
}
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
		id:             int64(rl.GetRandomValue(0, 100000)),
		name:           name,
		health:         100,
		attack_damage:  10,
		transform:      transform,
		world_position: rl.Vector2{X: transform.X, Y: transform.Y},
		speed:          5,
		scale:          float32(scale),
		sprite:         nil,
		isFacingRight:  is_facing_right,
		sprite_set:     sprite_set,
		entity_type:    entity_type,
		state:          IDLE_ANIM,
		hit_cooldown:   ENTITY_HIT_COOLDOWN,
		is_colliding:   false,
	}
	arr_entities = append(arr_entities, &new_entity)
	ls_entities.PushBack(&new_entity)
	return new_entity
}
func kill_entity(entity *Entity) {
	for element := ls_entities.Front(); element != nil; element = element.Next() {
		if element.Value.(*Entity) == entity {
			ls_entities.Remove(element)
			fmt.Println("Killed entity")
			break
		}
	}
	// for index := range arr_entities {
	// 	e := arr_entities[index]
	// 	if e.id == entity.id {
	// 		arr_entities = remove_copy(arr_entities, index)
	// 		break
	// 	}
	// }
}
func is_entity_colliding(first_entity Entity, second_entity Entity) bool {
	source_hitbox := first_entity.hitbox
	target_hitbox := second_entity.hitbox
	if target_hitbox.X+target_hitbox.Width >= source_hitbox.X && target_hitbox.X+target_hitbox.Width < source_hitbox.X+source_hitbox.Width {
		return true
	} else if source_hitbox.X+source_hitbox.Width >= target_hitbox.X && source_hitbox.X+source_hitbox.Width < target_hitbox.X+target_hitbox.Width {
		return true
	} else {
		return false
	}
}
