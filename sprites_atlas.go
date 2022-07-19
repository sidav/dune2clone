package main

import rl "github.com/gen2brain/raylib-go/raylib"

type spriteAtlas struct {
	// first index is sprite number, second is frame number
	atlas      [][]rl.Texture2D
	spriteSize int // width of square sprite
}

func (sa *spriteAtlas) totalFrames() int {
	if sa != nil && len(sa.atlas) > 0 {
		return len(sa.atlas[0])
	}
	return 0
}

func (sa *spriteAtlas) getSpriteByDirectionAndFrameNumber(dx, dy, num int) rl.Texture2D {
	var spriteGroup uint8 = 0
	if dx == 1 {
		spriteGroup = 3
	}
	if dx == -1 {
		spriteGroup = 1
	}
	if dy == 1 {
		spriteGroup = 2
	}
	num = num % len(sa.atlas[spriteGroup])
	return sa.atlas[spriteGroup][num]
}
