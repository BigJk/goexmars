package goexmars

// Fight1Warrior lets one warrior fight against himself.
func Fight1Warrior(w1 string, cfg FightConfig) (resultWin1 int, resultWin2 int, resultEqual int) {
	requireLibrary()

	var win1 int32
	var win2 int32
	var equal int32

	fight1Warrior(
		w1,
		int32(cfg.CoreSize),
		int32(cfg.Cycles),
		int32(cfg.MaxProcess),
		int32(cfg.Rounds),
		int32(cfg.MaxWarriorLen),
		int32(cfg.MinSep),
		int32(cfg.PSpaceSize),
		int32(cfg.FixPos),
		&win1,
		&win2,
		&equal,
	)

	return int(win1), int(win2), int(equal)
}

// Fight2Warriors lets two warrior fight.
func Fight2Warriors(w1 string, w2 string, cfg FightConfig) (resultWin1 int, resultWin2 int, resultEqual int) {
	requireLibrary()

	var win1 int32
	var win2 int32
	var equal int32

	fight2Warriors(
		w1,
		w2,
		int32(cfg.CoreSize),
		int32(cfg.Cycles),
		int32(cfg.MaxProcess),
		int32(cfg.Rounds),
		int32(cfg.MaxWarriorLen),
		int32(cfg.MinSep),
		int32(cfg.PSpaceSize),
		int32(cfg.FixPos),
		&win1,
		&win2,
		&equal,
	)

	return int(win1), int(win2), int(equal)
}
