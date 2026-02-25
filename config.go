package goexmars

// FightConfig holds the simulation parameters for a fight.
// It supports fluent configuration via SetXxx chainable methods.
type FightConfig struct {
	CoreSize      int
	Cycles        int
	MaxProcess    int
	Rounds        int
	MaxWarriorLen int
	MinSep        int
	PSpaceSize    int
	FixPos        int
}

// NewFightConfig returns an empty config that can be configured fluently.
func NewFightConfig() FightConfig {
	return FightConfig{}
}

func (c FightConfig) SetCoreSize(v int) FightConfig {
	c.CoreSize = v
	return c
}

func (c FightConfig) SetCycles(v int) FightConfig {
	c.Cycles = v
	return c
}

func (c FightConfig) SetMaxProcess(v int) FightConfig {
	c.MaxProcess = v
	return c
}

func (c FightConfig) SetRounds(v int) FightConfig {
	c.Rounds = v
	return c
}

func (c FightConfig) SetMaxWarriorLen(v int) FightConfig {
	c.MaxWarriorLen = v
	return c
}

func (c FightConfig) SetMinSep(v int) FightConfig {
	c.MinSep = v
	return c
}

func (c FightConfig) SetPSpaceSize(v int) FightConfig {
	c.PSpaceSize = v
	return c
}

func (c FightConfig) SetFixPos(v int) FightConfig {
	c.FixPos = v
	return c
}
