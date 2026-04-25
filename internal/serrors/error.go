// Package serrors define all the custom errors in summon
package serrors

import "errors"

var ErrNoConfig = errors.New("no config found")
var ErrInvalidKeyCombo = errors.New("invalid key combination")
var ErrCaptureFailed = errors.New("hotkey capture ended without a valid combo")
