package goexmars

import "fmt"

// FightConfig holds the simulation parameters for a fight.
// It supports fluent configuration via SetXxx chainable methods.
type FightConfig struct {
	CoreSize      int `json:"core_size"`
	Cycles        int `json:"cycles"`
	MaxProcess    int `json:"max_process"`
	Rounds        int `json:"rounds"`
	MaxWarriorLen int `json:"max_warrior_len"`
	MinSep        int `json:"min_sep"`
	PSpaceSize    int `json:"p_space_size"`
	FixPos        int `json:"fix_pos"`
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

// Validate checks whether the config contains a sane set of values.
//
// PSpaceSize may be zero to use exmars' default behavior. FixPos may be zero to
// use exmars' default placement behavior.
func (c FightConfig) Validate() error {
	if c.CoreSize <= 0 {
		return fmt.Errorf("invalid CoreSize: %d", c.CoreSize)
	}
	if c.Cycles <= 0 {
		return fmt.Errorf("invalid Cycles: %d", c.Cycles)
	}
	if c.MaxProcess <= 0 {
		return fmt.Errorf("invalid MaxProcess: %d", c.MaxProcess)
	}
	if c.Rounds <= 0 {
		return fmt.Errorf("invalid Rounds: %d", c.Rounds)
	}
	if c.MaxWarriorLen <= 0 {
		return fmt.Errorf("invalid MaxWarriorLen: %d", c.MaxWarriorLen)
	}
	if c.MinSep <= 0 {
		return fmt.Errorf("invalid MinSep: %d", c.MinSep)
	}
	if c.PSpaceSize < 0 {
		return fmt.Errorf("invalid PSpaceSize: %d", c.PSpaceSize)
	}
	if c.MaxWarriorLen > c.CoreSize {
		return fmt.Errorf("MaxWarriorLen (%d) must be <= CoreSize (%d)", c.MaxWarriorLen, c.CoreSize)
	}
	if c.MinSep > c.CoreSize {
		return fmt.Errorf("MinSep (%d) must be <= CoreSize (%d)", c.MinSep, c.CoreSize)
	}
	if c.PSpaceSize > 0 && c.PSpaceSize > c.CoreSize {
		return fmt.Errorf("PSpaceSize (%d) must be <= CoreSize (%d)", c.PSpaceSize, c.CoreSize)
	}
	if c.FixPos < 0 {
		return fmt.Errorf("invalid FixPos: %d", c.FixPos)
	}
	if c.FixPos >= c.CoreSize {
		return fmt.Errorf("FixPos (%d) must be < CoreSize (%d)", c.FixPos, c.CoreSize)
	}
	return nil
}
