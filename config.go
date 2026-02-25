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

// SetCoreSize returns a copy of c with CoreSize set to v.
func (c FightConfig) SetCoreSize(v int) FightConfig {
	c.CoreSize = v
	return c
}

// SetCycles returns a copy of c with Cycles set to v.
func (c FightConfig) SetCycles(v int) FightConfig {
	c.Cycles = v
	return c
}

// SetMaxProcess returns a copy of c with MaxProcess set to v.
func (c FightConfig) SetMaxProcess(v int) FightConfig {
	c.MaxProcess = v
	return c
}

// SetRounds returns a copy of c with Rounds set to v.
func (c FightConfig) SetRounds(v int) FightConfig {
	c.Rounds = v
	return c
}

// SetMaxWarriorLen returns a copy of c with MaxWarriorLen set to v.
func (c FightConfig) SetMaxWarriorLen(v int) FightConfig {
	c.MaxWarriorLen = v
	return c
}

// SetMinSep returns a copy of c with MinSep set to v.
func (c FightConfig) SetMinSep(v int) FightConfig {
	c.MinSep = v
	return c
}

// SetPSpaceSize returns a copy of c with PSpaceSize set to v.
func (c FightConfig) SetPSpaceSize(v int) FightConfig {
	c.PSpaceSize = v
	return c
}

// SetFixPos returns a copy of c with FixPos set to v.
func (c FightConfig) SetFixPos(v int) FightConfig {
	c.FixPos = v
	return c
}
