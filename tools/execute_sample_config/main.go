package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

func main() {
	dirs := []string{
		// "./plugins/inputs/...", "./plugins/outputs/...", "./plugins/processors/...", "./plugins/aggregators/...",
		"./plugins/inputs/...",
	}

	for _, d := range dirs {
		fmt.Println(d)
		fmt.Println("wut")
		// go func() {
		fmt.Println("wut")
		cmd := exec.Command("go", "generate", `-run="extract_sample_config/main.go`, d)
		var outb, errb bytes.Buffer
		cmd.Stdout = &outb
		cmd.Stderr = &errb
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("out:", outb.String(), "err:", errb.String())
		// }()
	}
}
