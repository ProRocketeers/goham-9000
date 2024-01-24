package lib

import (
	"bufio"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"goham-9000/database"
	model2 "goham-9000/model"
	"math/rand"
	"os/exec"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz1234567890"

func randStringBytes() string {
	length := 5
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
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

func NixpackBuild(id string) model2.Repository {

	var record, err = database.GetProjectById(id)

	prg := "nixpacks"
	arg1 := "build"
	arg2 := "--name"
	imageName := randStringBytes()

	cmd := exec.Command(prg, arg1, record.URL, arg2, imageName)
	log.Debug("________________________")
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
	log.Debug("Done")
	log.Debug("________________________")
	record.Status = "BUILT"
	record.ImgUrl = imageName

	return record
}
