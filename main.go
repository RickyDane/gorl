package main

import (
	"container/list"
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"runtime"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type App struct{}

// Global variables
var (
	update_count int64
	windowSize   = rl.Vector2{
		X: 1280,
		Y: 720,
	}
	scaled_width  float32 = 0
	scaled_height float32 = 0
	window_scale  float32 = 0
	window_hdpi_scale rl.Vector2
	Player                = Entity{
		id:            int64(rl.GetRandomValue(0, 100000)),
		name:          "Player",
		health:        100,
		max_health:    100,
		attack_damage: 35,
		position: rl.Vector2{
			X: windowSize.X / 4,
			Y: windowSize.Y - 48*3 - 30,
		},
		size: rl.Vector2{
			X: 48 * 3,
			Y: 48 * 3,
		},
		world_position: rl.Vector2{X: windowSize.X / 4, Y: windowSize.Y*3 - 30},
		speed:          15,
		sprint_speed:   25,
		isRunning:      false,
		isFacingRight:  true,
		isAttacking:    false,
		attackType:     0,
		entity_type:    PLAYER,
		hitbox:         rl.Rectangle{X: windowSize.X / 4, Y: windowSize.Y - 48*3 - 30, Width: 0, Height: 0},
		xp_to_reach:    100,
		xp:             0,
		sprite_color:   rl.White,
	}
	mouse_pos rl.Rectangle = rl.Rectangle{
		X:      rl.GetMousePosition().X,
		Y:      rl.GetMousePosition().Y,
		Width:  5,
		Height: 5,
	}
	deltaTime             float32    = 0
	betweenAttacksTimer   float32    = TIME_FOR_ATTACK_2
	isHitboxDebug         bool       = false
	ls_entities           *list.List = list.New()
	arr_entities		  []*Entity = make([]*Entity, 0)
	is_fullscreen         bool       = false
	window_render_texture rl.RenderTexture2D
	entity_hovered        *Entity = &Entity{}
	is_ui_open            bool    = false
	ui_type               UiType
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
	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagWindowHighdpi)
	rl.InitWindow(int32(windowSize.X), int32(windowSize.Y), "Demon Slayer")
	defer rl.CloseWindow()

	// Runs one time by start of the game
	app.Setup()

	// Main loop for window / game
	for !rl.WindowShouldClose() {
		// Update deltaTime
		deltaTime = 1.0 / float32(rl.GetFPS())

		// Play background music with updating the stream buffer
		rl.UpdateMusicStream(BG_MUSIC)

		// Main engine update method implementation
		app.Update()

		// Draw game on texture
		rl.BeginTextureMode(window_render_texture)
		rl.ClearBackground(rl.White)

		// Main engine draw method implementation
		app.Draw()

		rl.EndTextureMode()
		// At the end draw the texture scaled to the window size
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawTexturePro(
			window_render_texture.Texture,
			rl.Rectangle{X: 0, Y: 0, Width: float32(window_render_texture.Texture.Width), Height: float32(-window_render_texture.Texture.Height)},
			rl.Rectangle{
				X:      scaled_width - scaled_width/2,
				Y:      scaled_height - scaled_height/2,
				Width:  windowSize.X * window_scale/window_hdpi_scale.X,
				Height: windowSize.Y * window_scale/window_hdpi_scale.Y,
			},
			rl.Vector2{X: 0, Y: 0},
			0,
			rl.White,
		)
		rl.EndDrawing()
	}
	rl.UnloadRenderTexture(window_render_texture)
}

func (a *App) Setup() {
	// Set target FPS
	rl.SetTargetFPS(60)
	// Set log level to disable unneccessary output
	rl.SetTraceLogLevel(rl.LogError)

	setup_audio()

	// window_scale = float32(math.Min(float64(rl.GetScreenWidth())/float64(windowSize.X), float64(rl.GetScreenHeight())/float64(windowSize.Y)))
	window_scale = float32(calculate_window(rl.GetScreenWidth(), rl.GetScreenHeight()))
	window_render_texture = rl.LoadRenderTexture(int32(windowSize.X), int32(windowSize.Y))
	rl.SetTextureFilter(window_render_texture.Texture, rl.TextureFilterLinearMipNearest)
	window_hdpi_scale = rl.GetWindowScaleDPI();

	// Initial window size
	// rl.MaximizeWindow()

	// Load sprite atlas
	sprite_atlas = rl.LoadTexture("assets/sprite_atlas.png")

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

	// Spawn the first 3000 enemies across the world
	for i := 0; i < 100; i++ {
		spawn_entity("Demon", rl.Vector2{X: float32(rl.GetRandomValue(int32(windowSize.X), 100000)), Y: -25}, ENEMY, DEMON, 4)
	}

	// Spawn a shop
	spawn_entity("Shop", rl.Vector2{X: 600, Y: -25}, SHOP, ENEMY_NONE, 6)
}

func (a *App) Update() {
	// Update frame count | Important for animations!
	Player.frame_count++

	// Update mouse position
	mouse_pos = rl.Rectangle{
		X:      float32(rl.GetMouseX()) / window_scale * window_hdpi_scale.X,
		Y:      float32(rl.GetMouseY()) / window_scale * window_hdpi_scale.Y,
		Width:  mouse_pos.Width,
		Height: mouse_pos.Height,
	}

	// Update window information for scaling
	window_scale = float32(calculate_window(rl.GetRenderWidth(), rl.GetRenderHeight()))
	scaled_width = float32(rl.GetRenderWidth()) - (windowSize.X * window_scale)
	scaled_height = float32(rl.GetRenderHeight()) - (windowSize.Y * window_scale)

	is_entity_hovered := false
	for element := ls_entities.Front(); element != nil; element = element.Next() {
		entity := element.Value.(*Entity)
		entity.update()
		// Check if mouse is over entities and draw healthbar
		if is_entity_colliding(Entity{hitbox: mouse_pos}, *entity) {
			entity_hovered = entity
			is_entity_hovered = true
		}
	}
	/*for i := range arr_entities {
		entity := arr_entities[i]
		entity.update()
		// Check if mouse is over entities and draw healthbar
		if is_entity_colliding(Entity{hitbox: mouse_pos}, *entity) {
			entity_hovered = entity
			is_entity_hovered = true
		}
	}*/
	if !is_entity_hovered {
		entity_hovered = &Entity{}
	}

	// ----------- :start ## PLAYER SPECIFIC START ## -----------

	// :i_player Check player input for character movement
	if Player.isAttacking {
		if (Player.attackType == ATTACK_LIGHT || Player.attackType == ATTACK_LIGHT_BACK || Player.attackType == ATTACK_HEAVY) && Player.animation_phase < Player.current_sprite.frame_count-1 {
			a.play_attack_animation()
		} else if betweenAttacksTimer > 0 && Player.attackType != ATTACK_LIGHT_BACK {
			Player.attackType = ATTACK_NONE
			if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
				Player.attackType = ATTACK_LIGHT_BACK
				Player.current_sprite = attack2_anim
				reset_for_animation()
				a.play_attack_animation()
				rl.PlaySound(SWING)
				attack(Player.attack_damage)
			}
			betweenAttacksTimer -= 1 * deltaTime
		} else { // Combo for a second light attack
			reset_for_animation()
			betweenAttacksTimer = TIME_FOR_ATTACK_2
			Player.isAttacking = false
			Player.attackType = ATTACK_NONE
		}
	} else if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		if entity_hovered.entity_type == ENEMY {
			reset_for_animation()
			Player.attackType = ATTACK_LIGHT
			a.play_attack_animation()
			rl.PlaySound(SWING)
			Player.isAttacking = true
			attack(Player.attack_damage)
		}
	} else if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
		if entity_hovered.entity_type == ENEMY {
			reset_for_animation()
			Player.attackType = ATTACK_HEAVY
			a.play_attack_animation()
			rl.PlaySound(SLASH_STRONG)
			Player.isAttacking = true
			if pt_chance(0.5) {
				attack(Player.attack_damage * 2)
			} else {
				attack(Player.attack_damage)
			}
		}
	} else {
		Player.isAttacking = false
	}
	// Player movement left / right
	if rl.IsKeyDown(rl.KeyA) {
		a.run_left()
	} else if rl.IsKeyDown(rl.KeyD) {
		a.run_right()
	} else {
		Player.isRunning = false
	}
	// Input for sprinting
	if rl.IsKeyDown(rl.KeyLeftShift) && !Player.isAttacking {
		Player.current_speed = Player.sprint_speed
	} else {
		Player.current_speed = Player.speed
	}
	// Input for jumping
	// if rl.IsKeyPressed(rl.KeySpace) {
	// 	Player.position.Y -= 35
	// } else if Player.position.Y < windowSize.Y-Player.current_sprite.height-15 {
	// 	Player.position.Y += 150 * deltaTime
	// }

	// Input for interacting with entities
	if rl.IsKeyPressed(rl.KeyE) {
		if entity_hovered.entity_type == SHOP && !is_ui_open {
			ui_type = UI_SHOP
			is_ui_open = true
		} else if is_ui_open {
			is_ui_open = false
		}
	}

	// Check player input for debug
	if rl.IsKeyPressed(rl.KeyF1) { // Show hitboxes
		isHitboxDebug = !isHitboxDebug
	}
	if !Player.isAttacking && !Player.isRunning {
		a.idle_player()
	}
	// Constantly update player hitbox to keep up with player transform
	Player.update_hitbox()

	// ----------- :end ## PLAYER SPECIFIC END ## -----------

	if rl.IsKeyDown(rl.KeyO) { // Add new entity to entity list
		spawn_entity("Demon", rl.Vector2{X: float32(rl.GetRandomValue(0, int32(windowSize.X))), Y: -25}, ENEMY, DEMON, 4)
	}
	if rl.IsKeyPressed(rl.KeyL) {
		// Clear list of entities
		ls_entities = list.New()
	}

	// Check player input for some stuff
	if rl.IsKeyPressed(rl.KeyF11) {
		toggle_fullscreen()
	}

	// Debug print output
	print_debug_info()
}

func (a *App) Draw() {
	draw_background()

	for element := ls_entities.Front(); element != nil; element = element.Next() {
		entity := element.Value.(*Entity)
		entity.draw()
	}
	/*for i := range arr_entities {
		entity := arr_entities[i]
		entity.draw()
	}*/

	// :dplayer Draw the player sprite
	draw_sprite(&Player, 1, Player.sprite_color, 6)

	// Draw lines to visualize the hitbox for debugging
	if isHitboxDebug {
		// Player hitbox
		rl.DrawRectangleRoundedLines(Player.hitbox, 0, 0, 1, rl.Red)
		rl.DrawRectangleRoundedLines(mouse_pos, 0, 0, 1, rl.Green)
	}

	// Smoothly move clouds constantly
	scrolling_clouds += 0.2

	// Draw foreground over player
	bg_layer_8.Width = int32(windowSize.X)
	bg_layer_8.Height = int32(windowSize.Y)
	rl.DrawTextureRec(*bg_layer_8, rl.Rectangle{X: float32(scrolling_forefore), Y: 0, Width: float32(bg_layer_8.Width), Height: float32(bg_layer_8.Height)}, rl.Vector2{X: 0, Y: 50}, color.RGBA{255, 255, 255, 255})

	// Draw available user interface
	draw_ui()

	// Draw fps for debugging
	rl.DrawFPS(10, 10)
}

// Player specific behaviours
func (a *App) idle_player() {
	rl.StopSound(GRASS_RUNNING)
	Player.current_sprite = idle_anim
}
func (a *App) run_left() {
	if !rl.IsSoundPlaying(GRASS_RUNNING) {
		rl.PlaySound(GRASS_RUNNING)
	}
	if !Player.isAttacking {
		Player.current_sprite = run_anim
	} else {
		Player.current_speed = Player.speed / 3
	}
	Player.world_position.X -= Player.current_speed
	Player.isRunning = true
	Player.isFacingRight = false
	scroll_background(-1)
}
func (a *App) run_right() {
	if !rl.IsSoundPlaying(GRASS_RUNNING) {
		rl.PlaySound(GRASS_RUNNING)
	}
	if !Player.isAttacking {
		Player.current_sprite = run_anim
	} else {
		Player.current_speed = Player.speed / 3
	}
	Player.world_position.X += Player.current_speed
	Player.isRunning = true
	Player.isFacingRight = true
	scroll_background(1)
}
func (a *App) play_attack_animation() {
	switch Player.attackType {
	case ATTACK_LIGHT:
		Player.current_sprite = attack1_anim
	case ATTACK_LIGHT_BACK:
		Player.current_sprite = attack2_anim
	case ATTACK_HEAVY:
		Player.current_sprite = attack3_anim
	}
}
func attack(attack_damage float32) {
	Player.hitbox.Width = Player.size.X / 1.5
	if !Player.isFacingRight {
		Player.hitbox.X = Player.position.X + 10
	}
	for element := ls_entities.Front(); element != nil; element = element.Next() {
		entity := element.Value.(*Entity)
		if entity.entity_type == ENEMY {
			if entity.was_hit {
				continue
			}
			if is_entity_colliding(Player, *entity) {
				entity.hit(attack_damage)
				attack(attack_damage)
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
	// TODO: Simplify to array?
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
	Player.frame_count = 0
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
	fmt.Printf("Player position: x {%f} y{%f} \n", Player.position.X, Player.position.Y)
	fmt.Printf("Player world position: x: {%f} y: {%f}\n", Player.world_position.X, Player.world_position.Y)
	fmt.Printf("Player hitbox %v\n", Player.hitbox)
	fmt.Printf("Player xp: %d\n", Player.xp)
	fmt.Printf("Player xp to reach: %d\n", Player.xp_to_reach)
	println("Player is using attack: ", Player.attackType)
	println("Player is attacking: ", Player.isAttacking)
	println("Player is running: ", Player.isRunning)
	println("Player animation_phase: ", Player.animation_phase)
	fmt.Printf("Between attacks timer: %f \n", betweenAttacksTimer)
	println("Is audio device ready: ", rl.IsAudioDeviceReady())
	println("Entities list length", ls_entities.Len())
	println("Entities array length", len(arr_entities))
	fmt.Printf("Window width: %f height: %f\n", windowSize.X, windowSize.Y)
	fmt.Printf("Window scale: %f\n", window_scale)
	fmt.Printf("%f %f\n", scaled_width, scaled_height)
	println("frameCount: ", Player.frame_count)
	println("update_count: ", update_count)
	fmt.Printf("Mouse pos: %v\n", rl.GetMousePosition())
	println("Entity hovered: ", entity_hovered.name)
	// fmt.Printf("Entites array:  %+v", arr_entities)
}
func contains(list list.List, target any) bool {
	for item := list.Front(); item != nil; item = item.Next() {
		if item == target {
			return true
		}
	}
	return false
}
func pt_chance(percentage float32) bool {
	randomNumber := rand.Float32() * 100
	return percentage <= randomNumber
}

// Current not in use or worked on
func toggle_fullscreen() {
	rl.ToggleFullscreen()
	if !is_fullscreen {
		rl.MaximizeWindow()
	} else {
		rl.RestoreWindow()
		rl.SetWindowSize(int(windowSize.X), int(windowSize.Y))
	}
	is_fullscreen = !is_fullscreen
}
func calculate_window(width, height int) float64 {
	// Seitenverhältnis des Referenzfensters
	refWindow := float64(windowSize.X) / float64(windowSize.Y)

	// Seitenverhältnis des aktuellen Fensters
	window := float64(width) / float64(height)

	// Berechne den Skalierungsfaktor basierend auf dem Seitenverhältnisunterschied
	scale := math.Min(float64(width)/float64(windowSize.X), float64(height)/float64(windowSize.Y))

	// Anpassen der Skala, um das Seitenverhältnis des Referenzfensters beizubehalten
	if window > refWindow {
		scale *= refWindow / window
	} else {
		scale *= window / refWindow
	}

	return scale
}
func player_add_xp(xp int32) {
	Player.xp += xp
	if Player.xp >= Player.xp_to_reach {
		Player.level += 1
		if Player.level == 1 {
			rl.PlaySound(FIRST_LEVEL_UP)
		} else {
			rl.PlaySound(LEVEL_UP)
		}
		Player.xp = Player.xp % Player.xp_to_reach
		Player.xp_to_reach = int32(float32(Player.xp_to_reach) * 1.25)
	}
}

// :e_logic
// # -- Engine logic -- #
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
	HIT = rl.LoadSound("assets/sounds/hit.wav")
	FIRST_LEVEL_UP = rl.LoadSound("assets/sounds/first_level_up.wav")
	LEVEL_UP = rl.LoadSound("assets/sounds/level_up.mp3")
}
func spawn_entity(name string, position rl.Vector2, entity_type EntityType, enemy_type EnemyType, scale float32) Entity {
	// Make the entity face the player
	is_facing_right := true
	position.Y = windowSize.Y - 48*4 - 30

	sprite := Sprite{}
	if entity_type == ENEMY {
		switch enemy_type {
		case DEMON:
			sprite = demon_idle
		}
		// <-- Add more cases here
		if position.X >= Player.world_position.X {
			is_facing_right = false
		}
	} else if entity_type == SHOP {
		sprite = shop
	}

	new_entity := Entity{
		id:            int64(rl.GetRandomValue(0, 1000000)),
		name:          name,
		health:        100,
		max_health:    100,
		attack_damage: 10,
		position:      position,
		size: rl.Vector2{
			X: 48 * scale,
			Y: 48 * scale,
		},
		world_position: rl.Vector2{X: position.X, Y: position.Y},
		speed:          5,
		isFacingRight:  is_facing_right,
		entity_type:    entity_type,
		enemy_type:     enemy_type,
		hit_cooldown:   ENTITY_HIT_COOLDOWN,
		current_sprite: sprite,
		sprite_color:   rl.White,
	}
	ls_entities.PushBack(&new_entity)
	// arr_entities = append(arr_entities, &new_entity)
	return new_entity
}
func kill_entity(entity *Entity) {
	for element := ls_entities.Front(); element != nil; element = element.Next() {
		if element.Value.(*Entity).id == entity.id {
			ls_entities.Remove(element)
			break
		}
	}
}
func is_entity_colliding(e1, e2 Entity) bool {
	// Check if the right side of the first entity is to the left of the left side of the second entity
	if e1.hitbox.X+e1.hitbox.Width <= e2.hitbox.X {
		return false // No overlap
	}
	// Check if the left side of the first entity is to the right of the right side of the second entity
	if e1.hitbox.X >= e2.hitbox.X+e2.hitbox.Width {
		return false // No overlap
	}
	// Check if the bottom side of the first entity is above the top side of the second entity
	if e1.hitbox.Y+e1.hitbox.Height <= e2.hitbox.Y {
		return false // No overlap
	}
	// Check if the top side of the first entity is below the bottom side of the second entity
	if e1.hitbox.Y >= e2.hitbox.Y+e2.hitbox.Height {
		return false // No overlap
	}
	// If none of the above conditions hold, there must be a collision
	return true
}
func draw_sprite(entity *Entity, scale float32, color rl.Color, speed int) {
	rl.DrawTexturePro(
		sprite_atlas,
		get_anim_transform(entity, speed),
		rl.Rectangle{X: entity.position.X, Y: entity.position.Y, Width: entity.size.X * scale, Height: entity.size.Y * scale},
		rl.Vector2{X: 0, Y: 0},
		0,
		color,
	)
}
