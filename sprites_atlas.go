package main

import (
	"dune2clone/geometry"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type spriteAtlas struct {
	// first index is sprite number (rotation is there), second is frame number (animation)
	atlas      [][]rl.Texture2D
	// spriteSize int // width of square sprite
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

func (sa *spriteAtlas) getSpriteByDegreeAndFrameNumber(degree, num int) rl.Texture2D {
	rotFrame := geometry.DegreeToRotationFrameNumber(degree, len(sa.atlas))
	num = num % len(sa.atlas[rotFrame])
	return sa.atlas[rotFrame][num]
}

