package main

import (
	"net/http"
	"os"

	"github.com/go-git/go-git/v5"
	gitHttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

func initializeRepo(localPath string, remotePath string) (*git.Repository, error) {
	//check if folder at localPath already exists
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		//clone repo if it doesn't exist
		repo, err := git.PlainClone(localPath, false, &git.CloneOptions{
			URL: remotePath,
			Auth: &gitHttp.BasicAuth{
				Username: os.Getenv("GIT_USERNAME"),
				Password: os.Getenv("GIT_PASSWORD"),
			},
		})
		if err != nil {
			return nil, err
		}
		return repo, nil
	} else {
		//repo already exists, open local repo
		repo, err := git.PlainOpen(localPath)
		if err != nil {
			return nil, err
		}
		return repo, nil
	}
}

func main() {
	repo, err := initializeRepo(os.Getenv("LOCAL_GIT_JOURNAL_REPO"), os.Getenv("REMOTE_GIT_JOURNAL_REPO"))
	if err != nil {
		panic(err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		panic(err)
	}

	handler := &Handler{
		Repo:      repo,
		Worktree:  worktree,
		LocalPath: os.Getenv("LOCAL_GIT_JOURNAL_REPO"),
		Password:  os.Getenv("PASSWORD"),
	}
	handler.registerHandlers()

	http.ListenAndServe(os.Getenv("HTTP_LISTEN_PATH"), nil)
}
