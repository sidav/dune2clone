package map_generator

type tileCode int

const (
	SAND tileCode = iota
	BUILDABLE_TERRAIN
	RICH_RESOURCES
	MEDIUM_RESOURCES
	POOR_RESOURCES
	ROCKS
	RESOURCE_VEIN
)
