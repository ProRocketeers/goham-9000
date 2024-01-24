package lib

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	_ "os"
)

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.
type Deployment struct {
	Namespace      string
	Service        string
	Component      string
	Environment    string
	ReplicaCount   int
	Image          Image  `yaml:"image"`
	Ports          []Port `yaml:"ports"`
	Memory         string
	Envs           []Env `yaml:"envs,omitempty"`
	Ingress        Ingress
	ReadinessProbe ReadinessProbe `yaml:"readinessProbe"`
	LivenessProbe  LivenessProbe  `yaml:"livenessProbe"`
	Volumes        []Volumes      `yaml:"volumes,omitempty"`
	KubeSecrets    []KubeSecret   `yaml:"kubeSecrets"`
}

type Image struct {
	Repository  string   `yaml:"repository"`
	Tag         string   `yaml:"tag"`
	PullSecrets []string `yaml:"pullSecrets"`
}

type Port struct {
	Port string `yaml:"port"`
}

type Env struct {
}

type Ingress struct {
}

type ReadinessProbe struct {
}

type LivenessProbe struct {
}

type Volumes struct {
}

type KubeSecret struct {
}

func DeployEditor(folderPath string, id string) (string, error) {
	var fileToWrite = RootGitDir + "/" + folderPath + "/" + "test-deployment.yaml"

	log.Debug("Editing file: ", fileToWrite)
	// WRITE FILE
	fileData := Deployment{
		Namespace:    "default",
		Service:      id,
		Component:    "api",
		Environment:  "dev",
		ReplicaCount: 1,
		Image: Image{
			Repository: "harbor.cosik.com",
			Tag:        "abc12",
			//PullSecrets: [2] string {"gitlab-deploy-token-api", "dasf"}
		},
	}

	// todo: yaml editing magic
	generateDeployment(fileToWrite, fileData)
	log.Debug("File edited: ", fileToWrite)
	return "", nil
}

func generateDeployment(filename string, fileData Deployment) {
	// Writing the Person struct to a YAML file
	data, err := yaml.Marshal(fileData)
	if err != nil {
		fmt.Println("Error marshalling to YAML:", err)
		return
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Println("Error writing YAML file:", err)
	}
}
