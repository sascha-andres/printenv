package main

import (
	"fmt"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"os"
	"strings"
)

func main() {
	env := os.Environ()
	kv := make(map[string]string)
	keylength := 0
	for _, e := range env {
		b, a, found := strings.Cut(e, "=")
		if !found {
			continue
		}
		kv[b] = a
		if len(b) > keylength {
			keylength = len(b)
		}
	}
	keys := maps.Keys(kv)
	slices.SortFunc(keys, func(i, j string) bool {
		return i < j
	})
	for _, k := range keys {
		fmt.Printf("%-*s = %s\n", keylength, k, kv[k])
	}
}
