package main

import (
	"fmt"
	"github.com/sascha-andres/reuse/flag"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"log"
	"os"
	"os/exec"
	"strings"
)

var verbose bool

func init() {
	flag.SetEnvPrefix("reuse")
	flag.SetSeparated()
	flag.BoolVar(&verbose, "verbose", false, "verbose output")
	log.SetPrefix("[REUSE] ")
	log.SetFlags(log.LstdFlags | log.LUTC | log.Lshortfile)
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

	if verbose {
		log.Printf("executing separated: %s, %#v", separated[0], separated[1:])
	}

	verbs := flag.GetVerbs()
	for _, verb := range verbs {
		if verbose {
			log.Printf("verb: %s", verb)
		}
		if strings.Contains(verb, "=") {
			env = append(env, verb)
		}
	}

	cmd := exec.Command(separated[0], separated[1:]...)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
