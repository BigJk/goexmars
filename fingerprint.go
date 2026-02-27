package goexmars

import (
	"crypto/sha256"
	"encoding/hex"
)

// FingerprintOptions controls which normalized warrior parts are included in a fingerprint.
type FingerprintOptions struct {
	IncludeName   bool
	IncludeAuthor bool
	IncludeEnd    bool
}

// DefaultFingerprintOptions returns the default fingerprinting options.
//
// By default metadata is ignored so the fingerprint tracks normalized code and
// END position only.
func DefaultFingerprintOptions() FingerprintOptions {
	return FingerprintOptions{
		IncludeName:   false,
		IncludeAuthor: false,
		IncludeEnd:    true,
	}
}

// Fingerprint returns a SHA-256 fingerprint of the parsed warrior.
func (w ParsedWarrior) Fingerprint() string {
	return w.FingerprintWithOptions(DefaultFingerprintOptions())
}

// FingerprintWithOptions returns a SHA-256 fingerprint of the parsed warrior using opts.
func (w ParsedWarrior) FingerprintWithOptions(opts FingerprintOptions) string {
	text := w.Format(RedcodeFormatOptions{
		IncludeName:   opts.IncludeName,
		IncludeAuthor: opts.IncludeAuthor,
		IncludeEnd:    opts.IncludeEnd,
	})
	sum := sha256.Sum256([]byte(text))
	return hex.EncodeToString(sum[:])
}
