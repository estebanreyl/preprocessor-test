// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package graph

import (
	"errors"
	"strings"
	"fmt"
	yaml "gopkg.in/yaml.v2"
)

var (
	errImproperAlias = errors.New("Alias can only use alphanumeric characters")
	errMissingAlias  = errors.New("No alias was specified")
	errMissingMatch  = errors.New("Match for Alias may not be empty")
	errUnknownAlias  = errors.New("Unknown Alias")
)

// PreTask intermediate step for processing before complete unmarshall
type PreTask struct {
	AliasSrc []*string         `yaml:"alias-src"`
	AliasMap map[string]string `yaml:"alias"`
}

// // Validate validates the step and returns an error if the Step has problems.
// func (a *Alias) Validate() error {
// 	if a == nil {
// 		return nil
// 	}
// 	if len(a.val) > 1 {
// 		return improperFormat
// 	}
// 	return nil
// }

// // Equals determines whether or not two steps are equal.
// func (a *Alias) Equals(t *Alias) bool {
// 	if a == nil && t == nil {
// 		return true
// 	}
// 	if a == nil || t == nil {
// 		return false
// 	}
// 	return true
// 	// return a.val[1] == t.val[1] &&
// 	// 	a.val[0] == t.val[0]
// }

func PreprocessBytes(data []byte) ([]byte, error) {
	preTask := &PreTask{}

	if err := yaml.Unmarshal(data, preTask); err != nil {
		return nil, err
	}

	// Need to somehow read in the files associated with the Task TODO
	var out strings.Builder
	directive := '$'
	var command strings.Builder
	ongoing_cmd := false

	// Search and Replace
	str := string(data[:])
	for _, char := range str {
		if ongoing_cmd {
			//Maybe just checking if non alphanumeric, only allow alpha numeric aliases?
			if strings.Contains(")}/ .,;]&|'~\n\t", string(char)) { // Delineates the end of an alias
				out.WriteString(preTask.AliasMap[command.String()])
				if char != directive {
					ongoing_cmd = false
					out.WriteRune(char)
				}
				command.Reset()
			} else {
				command.WriteRune(char)
			}
		} else if char == directive {

			if ongoing_cmd { // Escape character triggered
				out.WriteRune(directive)
				ongoing_cmd = false
				continue
			}

			ongoing_cmd = true
			continue
		} else {
			out.WriteRune(char)
		}
	}
	fmt.Printf(out.String())
	return []byte(out.String()), nil

	// var aliases = make(map[string]int)a

	// for i := 0; i < len(preTask.AliasList); i++{
	// 	// Check if alias has
	// 	aliases[preTask.AliasList[i].alias] := preTask.AliasList[i].match
	// }
	// Second pass through these to redefine values in the header?
	// Should probably analyze them as they go, detect $
}
