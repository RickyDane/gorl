package main

import (
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type UiType int

const (
	UI_SHOP UiType = iota
)

func draw_ui() {
	if is_ui_open {
		switch ui_type {
		case UI_SHOP:
			draw_shop_ui()
		}
	}
	draw_xp_bar()
}

// Player xp bar
func draw_xp_bar() {
	bar_width := int32(200)
	bar_height := int32(10)
	percentage := float64(100) / float64(Player.xp_to_reach) * float64(Player.xp) / float64(100)
	rl.DrawRectangle(int32(windowSize.X)-202-19, 19, 202, 12, rl.White)
	rl.DrawRectangle(int32(windowSize.X)-bar_width-20, 20, int32(float64(bar_width)*percentage), bar_height, rl.Orange)
	rl.DrawText("Level: "+strconv.FormatInt(int64(Player.level), 10), int32(windowSize.X)-rl.MeasureText("Level: "+strconv.FormatInt(int64(Player.level), 10), 20)-bar_width-20-10, 15, 20, rl.Black)
}
func draw_shop_ui() {
	shop_ui_entity := Entity{
		current_sprite: panel,
		size: rl.Vector2{
			X: windowSize.X / 2,
			Y: windowSize.Y / 2,
		},
		position: rl.Vector2{
			X: windowSize.X / 4,
			Y: windowSize.Y / 4,
		},
	}
	draw_sprite(&shop_ui_entity, 1, rl.White, 1)
	rl.DrawText("Upgrades", int32(shop_ui_entity.position.X)+int32(shop_ui_entity.size.X)/2-rl.MeasureText("Upgrades", 30)/2+1, int32(shop_ui_entity.position.Y)+20+1, 30, rl.ColorAlpha(rl.Gray, 0.75))
	rl.DrawText("Upgrades", int32(shop_ui_entity.position.X)+int32(shop_ui_entity.size.X)/2-rl.MeasureText("Upgrades", 30)/2, int32(shop_ui_entity.position.Y)+20, 30, rl.ColorAlpha(rl.Black, 0.5))
}
