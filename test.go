package main

import (
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	m1 := map[string]struct{}{"ke1": {}, "key2": {}}
	m2 := map[string]struct{}{"ke1": {}, "key3": {}}
	keys := map[string]struct{}{}
	for k := range m1 {
		keys[k] = struct{}{}
	}
	for k := range m2 {
		keys[k] = struct{}{}
	}

	fmt.Println(keys)
}
