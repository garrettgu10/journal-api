package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Handler struct {
	Repo      *git.Repository
	Worktree  *git.Worktree
	LocalPath string
	Password  string
}

func (handler *Handler) hello(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "Hello World")
	return nil
}

type createNewNoteRequest struct {
	Contents string `json:"contents"`
	Year     int    `json:"year"`
	Month    int    `json:"month"`
	Day      int    `json:"day"`
	Password string `json:"password"`
}

func (handler *Handler) createNewNote(w http.ResponseWriter, r *http.Request) error {
	//set cors policy
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//parse json request body
	var newNote createNewNoteRequest
	err := json.NewDecoder(r.Body).Decode(&newNote)
	if err != nil {
		return err
	}

	//check password
	if newNote.Password != handler.Password {
		return fmt.Errorf("invalid password")
	}

	//pull latest changes from remote
	err = handler.Worktree.Pull(&git.PullOptions{
		RemoteName: "origin",
	})
	if err != nil && err.Error() != "already up-to-date" {
		return err
	}

	var folderPath = filepath.Join(handler.LocalPath,
		strconv.Itoa(newNote.Year),
		fmt.Sprintf("%02d", newNote.Month))
	var filePath = filepath.Join(folderPath, fmt.Sprintf("%02d", newNote.Day))
	prefixNewLine := "\n"

	//create folder if it doesn't exist
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err = os.MkdirAll(folderPath, 0755)
		if err != nil {
			return err
		}
	}

	//create file if it doesn't exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		file.Close()
		prefixNewLine = "" // no need to prefix new line since our file is new
	}

	//append contents to end of file
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(prefixNewLine + newNote.Contents)
	if err != nil {
		return err
	}

	//return "OK"
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Wrote journal entry for day %d-%02d-%02d", newNote.Year, newNote.Month, newNote.Day)))

	return nil
}

func (handler *Handler) commit(w http.ResponseWriter, r *http.Request) error {
	//set cors policy
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//add all changes to staging area
	err := handler.Worktree.AddWithOptions(&git.AddOptions{
		All: true,
	})
	if err != nil {
		return err
	}
	//commit all changes
	_, err = handler.Worktree.Commit("Add new changes to journal", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Git Journal API",
			Email: "donotemail@ggu.systems",
			When:  time.Now(),
		},
	})

	if err != nil {
		return err
	}

	//push changes to remote
	err = handler.Repo.Push(&git.PushOptions{
		RemoteName: "origin",
	})
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Committed changes to remote"))

	return nil
}

func wrapHandler(handler func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err)
		}
	}
}

func (handler *Handler) registerHandlers() {
	http.HandleFunc("/", wrapHandler(handler.hello))
	http.HandleFunc("/create", wrapHandler(handler.createNewNote))
	http.HandleFunc("/commit", wrapHandler(handler.commit))
}
