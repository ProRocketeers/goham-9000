package lib

import (
	"github.com/go-git/go-git/v5"
	"github.com/gofiber/fiber/v2/log"
	"os"
)

const ROOT_GIT_DIR = "GIT_WORK_DIR"

func CloneRepository(repo string, filename string) error {
	log.Debug("Cloning ", repo, " into ", filename, "...")
	cloneDir := ROOT_GIT_DIR + filename

	log.Debug("Removing ", cloneDir)
	err := os.RemoveAll(cloneDir)

	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Creating ", cloneDir)
	err = os.Mkdir(cloneDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Cloning ", cloneDir)

	// Clones the repository into the worktree (fs) and stores all the .git
	// content into the storer
	_, err = git.PlainClone(
		cloneDir,
		false,
		&git.CloneOptions{
			URL:  repo,
			Tags: git.NoTags,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Repo cloned")
	//
	//branches, err := r.Branches()
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//println(branches, "BRancehros")
	// Prints the content of the CHANGELOG file from the cloned repository
	return nil
}
