package main

import (
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func draw_ui() {
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
	rl.DrawText(strconv.FormatFloat(percentage, 'G', -1, 64), 50, 50, 20, rl.Beige)
}
