package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/sascha-andres/reuse/flag"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var (
	verbose, printSecrets bool
)

// hintsVariableShouldBeSecret contains a list of substrings that suggest a variable should be kept confidential.
var hintsVariableShouldBeSecret = []string{
	"password",
	"token",
}

func init() {
	flag.SetEnvPrefix("printenv")
	flag.SetSeparated()
	flag.BoolVar(&verbose, "verbose", false, "verbose output")
	flag.BoolVar(&printSecrets, "print-secrets", false, "print secrets")
	log.SetPrefix("[PRINTENV] ")
	log.SetFlags(log.LstdFlags | log.LUTC | log.Lshortfile)
}

func main() {
	flag.Parse()
	separated := flag.GetSeparated()

	env := os.Environ()
	if len(separated) == 0 {
		kv := make(map[string]string)
		keyLength := 0
		for _, e := range env {
			b, a, found := strings.Cut(e, "=")
			if !found {
				continue
			}
			kv[b] = a
			if len(b) > keyLength {
				keyLength = len(b)
			}
		}
		keys := maps.Keys(kv)
		slices.SortFunc(keys, func(i, j string) int {
			if i < j {
				return -1
			}
			if i > j {
				return 1
			}
			return 0
		})
		for _, k := range keys {
			val := kv[k]
			if !printSecrets && isSecret(k) {
				val = "<<REDACTED>>"
			}
			fmt.Printf("%-*s = %s\n", keyLength, k, val)
		}
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

// isSecret checks if the provided key contains substrings that suggest it should be treated as confidential.
func isSecret(k string) bool {
	for i := range hintsVariableShouldBeSecret {
		if strings.Contains(strings.ToUpper(k), strings.ToUpper(hintsVariableShouldBeSecret[i])) {
			return true
		}
	}
	return false
}
