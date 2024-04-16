package main

import (
	"fmt"
	"os"
	"strings"

	rediscache "github.com/encoder-run/operator/pkg/cache/redis"
	"github.com/go-git/go-git/v5" // with go modules enabled (GO111MODULE=on or outside GOPATH)
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/redis/go-redis/v9"
)

func main() {
	// action := "clone"
	url := "https://github.com/zachsmith1/gitinfo.git"

	// New redis storage
	s := rediscache.NewStorage(&redis.Options{
		Addr: "localhost:6379",
	}, url)

	// switch action {
	// case "clone":
	// 	clone(s, url)
	// default:
	// 	panic("unknown option")
	// }

	r, err := git.Open(s, nil)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			_, err := git.Clone(s, nil, &git.CloneOptions{
				URL: url,
			})
			CheckIfError(err)
			return
		} else {
			CheckIfError(err)
		}
	}
	// Fetch changes from the remote repository
	err = r.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		RefSpecs: []config.RefSpec{
			"+refs/heads/*:refs/remotes/origin/*",
		},
		Depth: 3,
	})
	if err != git.NoErrAlreadyUpToDate {
		CheckIfError(err)
	}

	err = r.Prune(git.PruneOptions{})
	CheckIfError(err)

	// Manually update local branch reference to match the remote tracking branch
	// Typically in a bare repo, this might be done in response to a push or a hook
	remoteRef, err := r.Reference(plumbing.NewRemoteReferenceName("origin", "main"), false)
	CheckIfError(err)

	// Update local main directly to point to the same commit
	localRef := plumbing.NewHashReference(plumbing.NewBranchReferenceName("main"), remoteRef.Hash())
	err = r.Storer.SetReference(localRef)
	CheckIfError(err)

	fmt.Printf("Updated local main to %s\n", remoteRef.Hash())

}

// CheckArgs should be used to ensure the right command line arguments are
// passed before executing an example.
func CheckArgs(arg ...string) {
	if len(os.Args) < len(arg)+1 {
		Warning("Usage: %s %s", os.Args[0], strings.Join(arg, " "))
		os.Exit(1)
	}
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}
