package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Entity struct {
	transform     rl.Rectangle
	speed         float32
	isRunnging    bool
	isAttacking   bool
	attackType    int32
	isFacingRight bool
	scale         float32
	sprite        *rl.Texture2D
}
