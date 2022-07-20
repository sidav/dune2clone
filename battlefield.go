package main

type battlefield struct {
	tiles     [][]tile
	buildings []*building
	units     []*unit
}

func (b *battlefield) create(w, h int) {
	b.tiles = make([][]tile, w)
	for i := range b.tiles {
		b.tiles[i] = make([]tile, h)
		for j := range b.tiles[i] {
			b.tiles[i][j].code = TILE_SAND
		}
	}
	b.buildings = append(b.buildings, &building{
		topLeftX: 1,
		topLeftY: 1,
		code:     BLD_BASE,
	})

	b.units = append(b.units, &unit{
		code: UNT_TANK,
		x:    0,
		y:    0,
	})
}
