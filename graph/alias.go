// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package graph

import (
	"errors"
)

var (
	improperFormat  = errors.New("Alias is improperly defined")
	errMissingAlias = errors.New("No alias was specified")
	errMissingMatch = errors.New("Match for Alias may not be the empty string")
)

// Step is a step in the execution task.
type Alias struct {
	val map[string]string
}

// Validate validates the step and returns an error if the Step has problems.
func (a *Alias) Validate() error {
	if a == nil {
		return nil
	}
	if len(a.val) > 1 {
		return improperFormat
	}
	// if a.val[0] == "" {
	// 	return errMissingAlias
	// }
	// if a.val[0] == "" {
	// 	return errMissingMatch
	// }
	// if a.val[0] == "$" {
	// 	// Change default symbol
	// }
	return nil
}

// Equals determines whether or not two steps are equal.
func (a *Alias) Equals(t *Alias) bool {
	if a == nil && t == nil {
		return true
	}
	if a == nil || t == nil {
		return false
	}
	return true
	// return a.val[1] == t.val[1] &&
	// 	a.val[0] == t.val[0]
}

