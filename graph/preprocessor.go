// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package graph

import (
	"errors"
	"strings"
	yaml "gopkg.in/yaml.v2"
)

var (
	errImproperAlias = errors.New("alias can only use alphanumeric characters")
	errMissingAlias  = errors.New("no alias was specified")
	errMissingMatch  = errors.New("match for Alias may not be empty")
	errUnknownAlias  = errors.New("unknown Alias")
	directive = '$'
)

// PreTask intermediate step for processing before complete unmarshall
type PreTask struct {
	AliasSrc  []*string         `yaml:"alias-src"`
	AliasMap map[string]string  `yaml:"alias"`
}

//func (PreTask preTask) resolveMapAndValidate() error {
	// for _, item := range preTask.AliasMap {
	// 	if  re := regexp.MustCompile


	// 	value := item.Value
	// 	var err error
	// 	if (value.Contains(directive)) {
	// 		value, err = PreprocessString(preTask, value)	/// Check for errors	
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// 	preTask.resolvedMap[item.Key] = value
	// }
	// return nil


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


//Handles preprocessing of a string values
func PreprocessString(preTask *PreTask, str string) (string, error) {
	// Need to somehow read in the files associated with the Task TODO
	var out strings.Builder
	var command strings.Builder
	ongoingCmd := false

	// Search and Replace
	for _, char := range str{
		if ongoingCmd {
			//Maybe just checking if non alphanumeric, only allow alpha numeric aliases?
			if strings.Contains(")}/ .,;]&|'~\n\t", string(char)) { // Delineates the end of an alias
				resolvedCommand, commandPresent := preTask.AliasMap[command.String()]
				if !commandPresent {
					return "", errUnknownAlias
				}

				out.WriteString(resolvedCommand)

				if char != directive {
					ongoingCmd = false
					out.WriteRune(char)
				}
				command.Reset()

			} else {
				command.WriteRune(char)
			}
		} else if char == directive {

			if ongoingCmd { // Escape character triggered
				out.WriteRune(directive)
				ongoingCmd = false
				continue
			}

			ongoingCmd = true
			continue
		} else {
			out.WriteRune(char)
		}
	}

	return out.String(), nil

}

// PreprocessBytes Handles files or byte encoded data that can be parsed through pre processing
func PreprocessBytes(data []byte) ([]byte, error) {
	preTask := &PreTask{}

	if err := yaml.Unmarshal(data, preTask); err != nil {
		return preTask, err
	}

	// Search and Replace
	str := string(data[:])
	parsedStr, err := PreprocessString(preTask, str)
	return []byte(parsedStr), err
}
