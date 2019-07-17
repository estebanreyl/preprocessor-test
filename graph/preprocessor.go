// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package graph

import (
	"errors"
	yaml "gopkg.in/yaml.v2"
	"strings"
	"net/http"
)

var (
	errImproperAlias = errors.New("alias can only use alphanumeric characters")
	errMissingAlias  = errors.New("no alias was specified")
	errMissingMatch  = errors.New("match for Alias may not be empty")
	errUnknownAlias  = errors.New("unknown Alias")
	directive        = '$'
)

// PreTask intermediate step for processing before complete unmarshall
type PreTask struct {
	AliasSrc []*string         `yaml:"alias-src"`
	AliasMap map[string]string `yaml:"alias"`
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

func (preTask *PreTask) loadExternalAlias() error {
	// Iterating in reverse to easily and efficiently handle hierarchy. The later
	// declared the higher in the hierarchy of alias definitions.
	for i := len(preTask.AliasSrc)-1; i >= 0; i-- {
		aliasUri := preTask.AliasSrc[i]
		
		// Need to determine if aliasFile is a local file or a network resource.
		// Include http?
		if (err := addAliasFromFile(PreTask, aliasUri); err != nil) {
			return err;
		}
	}
}

/* Fetches and Parses out remote alias files and adds their content
  to the passed in PreTask. Note alias definitions already in preTask
  will not be overwritten. */
func addAliasFromRemote (preTask *PreTask, url string) error {
	remoteClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	res, getErr := remoteClient.Do(req)
	if getErr != nil {
		return getErr
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return readErr
	}

	return readAliasFromBytes (data, preTask)
}

/* Parses out local alias files and adds their content to the passed in
   PreTask. Note alias definitions already in preTask will not be 
   overwritten. */
func addAliasFromFile (preTask *PreTask, fileUri string) error {
	
	data, fileReadingError := ioutil.ReadFile(fileUri)
	if (fileReadingError) {
		return fileReadingError
	}
	return readAliasFromBytes (data, preTask)
}

/* Parses out alias  definitions from a given bytes array and appends
   them to the PreTask. Note alias definitions already in preTask will
   not be overwritten even if present in the array. */
func readAliasFromBytes (data []byte, preTask *PreTask) error {

	aliasMap := &map[string]string{}

	if err := yaml.Unmarshal(data, fileAliasMap); err != nil {
		return err
	}

	for key, value := range aliasMap { 
        if (!preTask.aliasMap.Contains(key)) {
			preTask.aliasMap[key] = value;
		}
	}
	return nil
}

// Handles preprocessing of a string values
func PreprocessString(preTask *PreTask, str string) (string, error) {
	// Load Remote/Local alias definitions
	preTask.loadExternalAlias()
	//preTask.loadGlobalDefinitions TODO?
	var out strings.Builder
	var command strings.Builder
	ongoingCmd := false

	// Search and Replace
	for _, char := range str {
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
