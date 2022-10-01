package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	defaultFont rl.Font
	// index of array is faction color.
	tilesAtlaces       = map[string]*spriteAtlas{}
	buildingsAtlaces   = map[string]*spriteAtlas{}
	unitChassisAtlaces = map[string]*spriteAtlas{}
	turretsAtlaces     = map[string]*spriteAtlas{}
	projectilesAtlaces = map[string]*spriteAtlas{}
	effectsAtlaces     = map[string]*spriteAtlas{}

	uiAtlaces = map[string]*spriteAtlas{}
)

func loadResources() {
	// a := int32(255)
	// defaultFont = rl.LoadFontEx("", 96, &a, 255)
	defaultFont = rl.LoadFont("resources/flexi.ttf")
	// rl.GenTextureMipmaps(&defaultFont.Texture)
	loadSprites()
}

func loadSprites() {
	tilesAtlaces = make(map[string]*spriteAtlas)
	buildingsAtlaces = make(map[string]*spriteAtlas)
	unitChassisAtlaces = make(map[string]*spriteAtlas)
	turretsAtlaces = make(map[string]*spriteAtlas)
	projectilesAtlaces = make(map[string]*spriteAtlas)
	uiAtlaces = make(map[string]*spriteAtlas)

	drawLoadingScreen("LOADING: TERRAIN")
	currPath := "resources/sprites/terrain/"
	tilesAtlaces["sand1"] = CreateAtlasFromFile(currPath+"sand1.png", 0, 0, 16, 16, 16, 16, 1, false, false)
	tilesAtlaces["sand2"] = CreateAtlasFromFile(currPath+"sand2.png", 0, 0, 16, 16, 16, 16, 1, false, false)
	tilesAtlaces["sand3"] = CreateAtlasFromFile(currPath+"sand3.png", 0, 0, 16, 16, 16, 16, 1, false, false)
	tilesAtlaces["buildable1"] = CreateAtlasFromFile(currPath+"buildable1.png", 0, 0, 16, 16, 16, 16, 1, false, false)
	tilesAtlaces["rock1"] = CreateAtlasFromFile(currPath+"rocks.png", 0, 0, 16, 16, 16, 16, 1, false, false)
	tilesAtlaces["buildabledamaged"] = CreateAtlasFromFile(currPath+"buildable_damaged.png", 0, 0, 16, 16, 16, 16, 1, false, false)
	tilesAtlaces["resourcevein"] = CreateAtlasFromFile(currPath+"resource_vein.png", 0, 0, 16, 16, 16, 16, 1, false, false)
	tilesAtlaces["melangerich"] = CreateAtlasFromFile(currPath+"melange_rich.png", 0, 0, 16, 16, 16, 16, 1, false, false)
	tilesAtlaces["melangemedium"] = CreateAtlasFromFile(currPath+"melange_medium.png", 0, 0, 16, 16, 16, 16, 1, false, false)
	tilesAtlaces["melangepoor"] = CreateAtlasFromFile(currPath+"melange_poor.png", 0, 0, 16, 16, 16, 16, 1, false, false)

	drawLoadingScreen("LOADING: BUILDINGS")
	currPath = "resources/sprites/buildings/"
	// WARNING: IT HAS FRAMES
	buildingsAtlaces["underconstruction"] = CreateAtlasFromFile(currPath+"under_construction.png", 0, 0, 32, 32, 16, 16, 7, false, false)

	buildingsAtlaces["base"] = CreateAtlasFromFile(currPath+"base.png", 0, 0, 64, 64, 32, 32, 1, false, true)
	buildingsAtlaces["powerplant1"] = CreateAtlasFromFile(currPath+"powerplant.png", 0, 0, 64, 64, 32, 32, 1, false, true)
	buildingsAtlaces["powerplant2"] = CreateAtlasFromFile(currPath+"anjaopterix/powerplant.png", 0, 0, 64, 64, 32, 32, 4, false, true)
	buildingsAtlaces["fusionreactor"] = CreateAtlasFromFile(currPath+"super_reactor.png", 0, 0, 64, 64, 48, 48, 23, false, true)
	buildingsAtlaces["barracks"] = CreateAtlasFromFile(currPath+"barracks.png", 0, 0, 64, 64, 32, 32, 1, false, true)
	buildingsAtlaces["factory"] = CreateAtlasFromFile(currPath+"factory.png", 0, 0, 96, 64, 48, 32, 1, false, true)
	buildingsAtlaces["airfactory"] = CreateAtlasFromFile(currPath+"airfactory.png", 0, 0, 64, 96, 32, 48, 1, false, true)
	buildingsAtlaces["refinery1"] = CreateAtlasFromFile(currPath+"refinery1.png", 0, 0, 96, 64, 48, 32, 1, false, true)
	buildingsAtlaces["refinery2"] = CreateAtlasFromFile(currPath+"refinery2.png", 0, 0, 96, 64, 48, 32, 1, false, true)
	buildingsAtlaces["depot"] = CreateAtlasFromFile(currPath+"depot.png", 0, 0, 48, 32, 48, 32, 1, false, true)
	buildingsAtlaces["silo"] = CreateAtlasFromFile(currPath+"silo.png", 0, 0, 32, 64, 16, 32, 1, false, true)
	buildingsAtlaces["turret_base"] = CreateAtlasFromFile(currPath+"turret_base.png", 0, 0, 32, 32, 16, 16, 1, false, true)
	buildingsAtlaces["bld_aaturret"] = CreateAtlasFromFile(currPath+"aa_turret.png", 0, 0, 32, 32, 16, 16, 1, false, true)
	turretsAtlaces["bld_turret_cannon"] = CreateDirectionalAtlasFromFile(currPath+"cannon_turret.png", 32, 16, 1, 2, true)
	turretsAtlaces["bld_turret_minigun"] = CreateDirectionalAtlasFromFile(currPath+"minigun_turret.png", 32, 16, 1, 2, true)
	buildingsAtlaces["fortress"] = CreateAtlasFromFile(currPath+"fortress.png", 0, 0, 32, 32, 32, 32, 1, false, true)
	turretsAtlaces["bld_fortress_cannon"] = CreateDirectionalAtlasFromFile(currPath+"fortress_turret.png", 32, 32, 1, 2, true)

	drawLoadingScreen("LOADING: UNITS")
	currPath = "resources/sprites/units/"
	unitChassisAtlaces["placeholder"] = CreateDirectionalAtlasFromFile(currPath+"placeholder.png", 32, 16, 1, 2, true)
	turretsAtlaces["placeholder"] = CreateDirectionalAtlasFromFile(currPath+"placeholder_turret.png", 32, 16, 1, 2, true)
	unitChassisAtlaces["infantry"] = CreateDirectionalAtlasFromFile(currPath+"infantry.png", 32, 16, 1, 2, true)
	unitChassisAtlaces["infantryrocket"] = CreateDirectionalAtlasFromFile(currPath+"infantry_rocket.png", 32, 16, 1, 2, true)
	unitChassisAtlaces["infantryrecon"] = CreateDirectionalAtlasFromFile(currPath+"infantry_recon.png", 32, 16, 1, 2, true)
	unitChassisAtlaces["infantryheavy"] = CreateDirectionalAtlasFromFile(currPath+"infantry_heavy.png", 32, 16, 1, 2, true)
	unitChassisAtlaces["tank"] = CreateDirectionalAtlasFromFile(currPath+"tank_chassis.png", 32, 16, 1, 2, true)
	turretsAtlaces["tank"] = CreateDirectionalAtlasFromFile(currPath+"tank_cannon.png", 32, 16, 1, 2, true)
	unitChassisAtlaces["tank2"] = CreateDirectionalAtlasFromFile(currPath+"anjaopterix/tank_chassis.png", 32, 16, 2, 2, true)
	turretsAtlaces["tank2"] = CreateDirectionalAtlasFromFile(currPath+"anjaopterix/tank_turret.png", 32, 16, 1, 2, true)
	unitChassisAtlaces["devastator"] = CreateDirectionalAtlasFromFile(currPath+"devastator.png", 26, 16, 1, 2, true)
	turretsAtlaces["devastator"] = CreateDirectionalAtlasFromFile(currPath+"devastator_turret.png", 26, 16, 1, 2, true)
	unitChassisAtlaces["juggernaut"] = CreateDirectionalAtlasFromFile(currPath+"juggernaut.png", 32, 16, 1, 2, true)
	turretsAtlaces["juggernautmain"] = CreateDirectionalAtlasFromFile(currPath+"juggernaut_mainturret.png", 32, 16, 1, 2, true)
	turretsAtlaces["juggernautsec"] = CreateDirectionalAtlasFromFile(currPath+"juggernaut_secturret.png", 32, 16, 1, 2, true)
	unitChassisAtlaces["quad"] = CreateDirectionalAtlasFromFile(currPath+"quad.png", 32, 16, 1, 2, true)
	turretsAtlaces["msltank"] = CreateDirectionalAtlasFromFile(currPath+"missiletank_turret.png", 32, 16, 1, 2, true)
	turretsAtlaces["aamsltank"] = CreateDirectionalAtlasFromFile(currPath+"aatank_turret.png", 32, 16, 1, 2, true)
	unitChassisAtlaces["combatharvester"] = CreateDirectionalAtlasFromFile(currPath+"combat_harvester.png", 32, 16, 1, 2, true)
	unitChassisAtlaces["fastharvester"] = CreateDirectionalAtlasFromFile(currPath+"fast_harvester.png", 32, 16, 1, 2, true)
	unitChassisAtlaces["mcv"] = CreateDirectionalAtlasFromFile(currPath+"mcv.png", 32, 16, 1, 2, true)

	currPath = "resources/sprites/units/aircrafts/"
	unitChassisAtlaces["air_gunship"] = CreateDirectionalAtlasFromFile(currPath+"combat_plane.png", 32, 16, 1, 2, true)
	unitChassisAtlaces["air_fighter"] = CreateDirectionalAtlasFromFile(currPath+"fighter_plane.png", 32, 16, 1, 2, true)
	unitChassisAtlaces["air_transport"] = CreateDirectionalAtlasFromFile(currPath+"transport_plane.png", 32, 16, 1, 2, true)

	drawLoadingScreen("LOADING: PROJECTILES")
	currPath = "resources/sprites/projectiles/"
	projectilesAtlaces["shell"] = CreateDirectionalAtlasFromFile(currPath+"shell.png", 32, 16, 1, 2, false)
	projectilesAtlaces["bullets"] = CreateDirectionalAtlasFromFile(currPath+"bullets.png", 32, 8, 1, 2, false)
	projectilesAtlaces["missile"] = CreateDirectionalAtlasFromFile(currPath+"missile.png", 32, 16, 1, 2, false)
	projectilesAtlaces["aamissile"] = CreateDirectionalAtlasFromFile(currPath+"aamissile.png", 32, 16, 1, 2, false)
	projectilesAtlaces["omni"] = CreateDirectionalAtlasFromFile(currPath+"omni.png", 32, 16, 1, 2, false)

	drawLoadingScreen("LOADING: UI")
	currPath = "resources/sprites/ui/"
	uiAtlaces["factionflag"] = CreateAtlasFromFile(currPath+"building_faction_flag.png", 0, 0, 4, 4, 4, 4, 4, false, true)
	uiAtlaces["energyicon"] = CreateDirectionalAtlasFromFile(currPath+"energy_icon.png", 16, 8, 1, 1, false)
	uiAtlaces["repairicon"] = CreateDirectionalAtlasFromFile(currPath+"repair_icon.png", 16, 8, 1, 1, false)
	uiAtlaces["readyicon"] = CreateDirectionalAtlasFromFile(currPath+"ready_icon.png", 16, 8, 1, 1, false)
	uiAtlaces["veterancy"] = CreateDirectionalAtlasFromFile(currPath+"veterancy.png", 10, 5, 4, 1, false)

	drawLoadingScreen("LOADING: EFFECTS")
	currPath = "resources/sprites/effects/"
	effectsAtlaces["smallexplosion"] = CreateAtlasFromFile(currPath+"explosion_small.png", 0, 0, 4, 4, 4, 4, 16, false, false)
	effectsAtlaces["regularexplosion"] = CreateAtlasFromFile(currPath+"explosion.png", 0, 0, 16, 16, 16, 16, 3, false, false)
	effectsAtlaces["biggerexplosion"] = CreateAtlasFromFile(currPath+"explosion_bigger.png", 0, 0, 40, 40, 20, 20, 3, false, false)
}
