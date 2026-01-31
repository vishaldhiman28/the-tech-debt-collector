package main

import "fmt"

func main() {
	// TODO: Add proper error handling
	result := doSomething()
	fmt.Println(result)

	// FIXME: This is a memory leak - need to free resources
	bigBuffer := make([]byte, 1000000)
	_ = bigBuffer

	// HACK: Quick fix for production, refactor later
	if true {
		fmt.Println("Temporary solution")
	}
}

func doSomething() string {
	// DEPRECATED: This function will be removed in v2.0
	return "old implementation"
}
