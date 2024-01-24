package lib

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/spf13/viper"
)

func Tst(uuid string, path string) string {
	// Here code co uploaduje do registru

	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	// curl -X GET -u prorocketeersdev:dckr_pat_AUm2LsPbOFP-fIQxRZfToWji41o "https://registry-1.docker.io/v2/_catalog"

	// dckr_pat_AUm2LsPbOFP-fIQxRZfToWji41o    heslo: ProRocketeers.com_777
	url := "https://index.docker.io/"
	username := "prorocketeersdev"                     // viper.Get("DOCKERHUB_USERNAME").(string) // anonymous
	password := "dckr_pat_AUm2LsPbOFP-fIQxRZfToWji41o" // viper.Get("DOCKERHUB_TOKEN").(string) // anonymous
	hub, err := registry.New(url, username, password)

	log.Debug(hub, err)

	repositories, err := hub.Repositories()

	log.Debug(repositories, err)

	return "with this uuid: " + uuid + " from this path: " + path
}
