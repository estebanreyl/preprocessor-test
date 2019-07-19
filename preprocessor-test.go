package main

import (
	"fmt"

	"github.com/Azure/acr-builder-preprocessor/graph"
)

func main() {
	// Backwards Compatibility
	fmt.Println("Reading in acb.yaml - Backwards compatibility")
	graph.TestFile("input/simple_alias.yaml")
	fmt.Println("Test Completed")

	//-------- String Testing --------
	//String definitions
	//String with local src
	//String with global src
	//String Nested definitions
	//String all combined

	//-------- File Testing --------
	//File definitions
	//File external definitions
	//File with local/remote src
	//File with global src
	//File Nested definitions
	//File all combined
}
