package main

// Must raises an error if it not nil
func Must(e error) {
	if e != nil {
		panic(e)
	}
}
