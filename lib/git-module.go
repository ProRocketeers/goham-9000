package lib

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/gofiber/fiber/v2/log"
	"os"
)

const RootGitDir = "GIT_WORK_DIR"

// CloneRepository clones a git repository into a given directory,
// always clean folder before clone
func CloneRepository(repo string, filename string) (string, error) {
	cloneDir := RootGitDir + "/" + filename
	log.Debug("Cloning ", repo, " into ", cloneDir, "...")

	log.Debug("Removing ", cloneDir)
	err := os.RemoveAll(cloneDir)

	if err != nil {
		log.Info("No file to remove", err)
		return "", err

	}

	log.Debug("Creating ", cloneDir)
	err = os.MkdirAll(cloneDir, 0755)
	if err != nil {
		log.Error(err)
		return "", err
	}

	log.Debug("Cloning ", cloneDir)

	//auth, err := ssh.DefaultAuthBuilder("git")
	//if err != nil {
	//	log.Fatalf("default auth builder: %v", err)
	//}
	//
	_, err = git.PlainClone(
		cloneDir,
		false,
		&git.CloneOptions{
			URL:  repo,
			Tags: git.NoTags,
			//Auth: auth,
		},
	)
	if err != nil {
		log.Error(err)
		return "", err

	}
	log.Debug("Repo cloned")

	return cloneDir, nil
}
func CloneRepoStep(projectId string) (string, error) {
	// get repo by key

	// resolve path

	// clone repo

	// update repo status in db

	return "", errors.New("not implemented")
}

// todo: maybe separate em
func CommitAndPush(filename string) error {
	fmt.Println("PlainOpen " + filename)
	r, err := git.PlainOpen(filename)
	if err != nil {
		log.Error(err)
		return err
	}

	auth, err := ssh.DefaultAuthBuilder("git")
	if err != nil {
		log.Fatalf("default auth builder: %v", err)
	}

	fmt.Println("set worktree")
	w, err := r.Worktree()
	if err != nil {
		fmt.Printf("%v", err)
		return err
	}

	fmt.Println("git add ", filename)
	w.Add(filename)

	fmt.Println("Commit our changes")
	w.Commit("Added my new file", &git.CommitOptions{})

	fmt.Println("Git push") // todo:
	err = r.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
	})
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func cleanup() {
	// todo
}
