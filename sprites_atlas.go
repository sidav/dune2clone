package main

import (
	"dune2clone/geometry"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type spriteAtlas struct {
	// first index is color mask, second is sprite number (rotation is there), third is frame number (animation)
	atlas [][][]rl.Texture2D
	// spriteSize int // width of square sprite
}

func (sa *spriteAtlas) totalFrames() int {
	if sa != nil && len(sa.atlas) > 0 && len(sa.atlas[0]) > 0 {
		return len(sa.atlas[0][0])
	}
	return 0
}

func (sa *spriteAtlas) getSpriteByFrame(frameNum int) rl.Texture2D {
	return sa.atlas[0][0][frameNum]
}

func (sa *spriteAtlas) getSpriteByColorAndFrame(color, frameNum int) rl.Texture2D {
	if frameNum >= sa.totalFrames() {
		frameNum = frameNum % sa.totalFrames()
	}
	return sa.atlas[color][0][frameNum]
}

//func (sa *spriteAtlas) getSpriteByFrame(frameNum int) rl.Texture2D {
//	return sa.atlas[0][0][frameNum]
//}

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
	return sa.atlas[0][spriteGroup][num]
}

func (sa *spriteAtlas) getSpriteByColorDegreeAndFrameNumber(color, degree, num int) rl.Texture2D {
	rotFrame := geometry.DegreeToSectorNumber(degree, len(sa.atlas[color]))
	// +2 is because zero degree looks right, but first sprite in atlas looks up. +2 = +90degs.
	rotFrame = (2 + rotFrame) % len(sa.atlas[color])
	num = num % len(sa.atlas[color][rotFrame])
	return sa.atlas[color][rotFrame][num]
}
