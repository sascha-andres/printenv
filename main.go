package main

import (
	"fmt"
	"github.com/sascha-andres/reuse/flag"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"os"
	"os/exec"
	"strings"
)

func init() {
	flag.SetEnvPrefix("reuse")
	flag.SetSeparated()
}
func main() {
	flag.Parse()
	separated := flag.GetSeparated()

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

	if len(separated) == 0 {
		return
	}

	fmt.Printf("executing separated: %s, %#v\n", separated[0], separated[1:])

	cmd := exec.Command(separated[0], separated[1:]...)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
