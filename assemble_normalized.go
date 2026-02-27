package goexmars

// AssembleNormalized assembles a warrior and returns normalized Redcode text
// formatted according to opts.
//
// Unlike Assemble, this output can omit metadata and/or the END line and is
// rendered via the ParsedWarrior formatter.
func AssembleNormalized(warrior string, cfg FightConfig, opts RedcodeFormatOptions) (string, error) {
	parsed, err := AssembleParsed(warrior, cfg)
	if err != nil {
		return "", err
	}
	return parsed.Format(opts), nil
}
