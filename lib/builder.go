package lib

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"goham-9000/database"
	model2 "goham-9000/model"
	"math/rand"
	"os/exec"
	"strings"
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

func NixpackBuildStep(projectId string) (model2.Repository, error) {

	var record, err = database.GetProjectById(projectId)
	if err != nil {
		return model2.Repository{}, err
	}
	prg := "nixpacks"
	arg1 := "build"
	arg2 := "--name"
	imageName := randStringBytes()

	cmd := exec.Command(prg, arg1, ResolveProjectPath(record), arg2, imageName)
	log.Debug("________________________")
	log.Debug("executing: " + cmd.String())
	var isBuildFailed, message = buildImage(cmd)
	if isBuildFailed == 1 {
		log.Error(message)
		return model2.Repository{}, errors.New(message)
	}
	if err != nil {
		log.Error(err)
		return model2.Repository{}, err
	}
	log.Debug("Done")
	log.Debug("________________________")
	record.Status = database.P_IMG_BUILT
	record.ImgUrl = imageName
	_, err = database.UpdateProjectStatus(projectId, database.P_IMG_BUILT)
	if err != nil {
		return model2.Repository{}, err
	}
	return record, nil
}

func buildImage(cmd *exec.Cmd) (int, string) {
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(stderr)
	var isBuildFailed, message = 0, ""
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
		if strings.Contains(m, "error") {
			isBuildFailed = 1
			message = m
		}
	}
	err := cmd.Wait()
	if err != nil {
		isBuildFailed = 1
		message = err.Error()
	}
	return isBuildFailed, message
}
