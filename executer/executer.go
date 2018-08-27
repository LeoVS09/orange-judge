package executer

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

func Run() {
	cmd := exec.Command("ls")
	//cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("FIles in dir: %s\n", out.String())
}
