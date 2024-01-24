package lib

import (
	"bufio"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"os/exec"
)

func NixpackVersion() string {

	prg := "nixpacks"

	arg1 := "-V"

	cmd := exec.Command(prg, arg1)
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	log.Debug(string(stdout))

	return string(stdout)
}

func NixpackBuild(path string) string {

	prg := "nixpacks"

	arg1 := "build"

	cmd := exec.Command(prg, arg1, path)
	log.Debug("executing: " + cmd.String())
	stderr, err := cmd.StderrPipe()
	cmd.Start()
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}

	return "hello"
}
