package lib

import (
	"github.com/go-git/go-git/v5"
	"github.com/gofiber/fiber/v2/log"
	"os"
)

func CloneRepository(repo string) error {
	direktus := "foobaz"
	err := os.RemoveAll(direktus)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Removing ", direktus)

	err = os.Mkdir(direktus, 0755)
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Creating ", direktus)

	// Clones the repository into the worktree (fs) and stores all the .git
	// content into the storer
	r, err := git.PlainClone(
		direktus,
		false,
		&git.CloneOptions{
			URL:  "https://github.com/Fenny/fiber-hello-world",
			Tags: git.NoTags,
		},
	)
	log.Debug("Cloning ", direktus)

	if err != nil {
		log.Fatal(err)
	}
	NixpackBuild(direktus)
	branches, err := r.Branches()

	if err != nil {
		log.Fatal(err)
	}
	println(branches, "BRancehros")
	// Prints the content of the CHANGELOG file from the cloned repository
	return nil
}
