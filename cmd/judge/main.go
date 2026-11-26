// Command judge — пока инструмент разработчика: запускает программу через песочницу
// и печатает, сколько ресурсов она потратила.
//
//	go run ./cmd/judge -input tmp/input.txt -- tmp/ref
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/algoduel/algoduel-judge/sandbox"
)

func main() {
	input := flag.String("input", "", "файл, который подать на stdin")
	timeout := flag.Duration("timeout", 5*time.Second, "таймаут по реальному времени")
	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatal("usage: judge [-input file] [-timeout 5s] -- <program> [args...]")
	}

	spec := sandbox.Spec{
		Path:   flag.Arg(0),
		Args:   flag.Args()[1:],
		Limits: sandbox.Limits{WallTime: *timeout, OutputBytes: 64 << 20},
	}
	if *input != "" {
		f, err := os.Open(*input)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		spec.Stdin = f
	}

	var runner sandbox.Runner = sandbox.Basic{}
	u, err := runner.Run(context.Background(), spec)
	if err != nil {
		log.Fatal(err)
	}
	out, _ := json.MarshalIndent(u, "", "  ")
	fmt.Println(string(out))
	fmt.Printf("stdout: %d байт\n", len(u.Stdout))
}
