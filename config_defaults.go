package goexmars

var DefaultConfig = FightConfig{
	Rounds:        250,
	CoreSize:      8000,
	Cycles:        80000,
	MaxWarriorLen: 100,
	MaxProcess:    8000,
	MinSep:        100,
	PSpaceSize:    500,
	FixPos:        0,
}

// Not fully supported: source preset uses Write limit = 4000 (CoreSize = 8000).
var FortressConfig = FightConfig{
	Rounds:        1000,
	CoreSize:      8000,
	Cycles:        80000,
	MaxWarriorLen: 400,
	MaxProcess:    80,
	MinSep:        4000,
	PSpaceSize:    500,
	FixPos:        0,
}

var TourneyConfig = FightConfig{
	Rounds:        250,
	CoreSize:      8192,
	Cycles:        100000,
	MaxWarriorLen: 300,
	MaxProcess:    8000,
	MinSep:        300,
	PSpaceSize:    512,
	FixPos:        0,
}

var LimitedProcessConfig = FightConfig{
	Rounds:        250,
	CoreSize:      8000,
	Cycles:        80000,
	MaxWarriorLen: 200,
	MaxProcess:    8,
	MinSep:        200,
	PSpaceSize:    500,
	FixPos:        0,
}

var MediumProcessConfig = FightConfig{
	Rounds:        250,
	CoreSize:      8000,
	Cycles:        80000,
	MaxWarriorLen: 100,
	MaxProcess:    64,
	MinSep:        100,
	PSpaceSize:    1,
	FixPos:        0,
}

var MetaswitchConfig = FightConfig{
	Rounds:        100,
	CoreSize:      8000,
	Cycles:        40000,
	MaxWarriorLen: 100,
	MaxProcess:    8000,
	MinSep:        100,
	PSpaceSize:    500,
	FixPos:        0,
}

var NanoConfig = FightConfig{
	Rounds:        250,
	CoreSize:      80,
	Cycles:        800,
	MaxWarriorLen: 5,
	MaxProcess:    80,
	MinSep:        5,
	PSpaceSize:    5,
	FixPos:        0,
}

var TinyConfig = FightConfig{
	Rounds:        250,
	CoreSize:      800,
	Cycles:        8000,
	MaxWarriorLen: 20,
	MaxProcess:    800,
	MinSep:        20,
	PSpaceSize:    50,
	FixPos:        0,
}

var TinyLimitedProcessConfig = FightConfig{
	Rounds:        250,
	CoreSize:      800,
	Cycles:        8000,
	MaxWarriorLen: 50,
	MaxProcess:    8,
	MinSep:        50,
	PSpaceSize:    50,
	FixPos:        0,
}
