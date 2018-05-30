package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/github/hub/github"
	"github.com/pkg/errors"
)

// https://github.com/ktr0731/evans/blob/master/README.md
var (
	format        = "%s/blob/%s/%s"
	defaultRemote = "origin"
)

func main() {
	status, err := run(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	os.Exit(status)
}

var (
	from = flag.Uint("from", 0, "highlight from")
	to   = flag.Uint("to", 0, "highlight to (from flag required)")
)

func run(args []string) (int, error) {
	flag.Parse()

	if len(args) == 0 {
		fmt.Println("Usage: open-github-source <path>")
		return 1, nil
	}

	if err := checkFlagCondition(); err != nil {
		return 1, errors.Wrap(err, "precondition failed")
	}

	// change dir and reset
	defer func(path string) func() {
		prev, _ := filepath.Abs(".")
		os.Chdir(path)
		return func() {
			os.Chdir(prev)
		}
	}(filepath.Dir(flag.Arg(0)))()

	repo, err := github.LocalRepo()
	if err != nil {
		return 1, err
	}

	proj, err := repo.MainProject()
	if err != nil {
		return 1, err
	}

	br, err := repo.CurrentBranch()
	if err != nil {
		return 1, err
	}

	fmt.Println(formatURL(proj, br.ShortName(), flag.Arg(0), *from, *to))
	return 0, nil
}

func formatURL(proj *github.Project, ref, path string, from, to uint) string {
	url := fmt.Sprintf(format, proj.WebURL("", "", ""), ref, path)
	if from > 0 {
		url += fmt.Sprintf("#L%d", from)
	}
	if to > 0 {
		url += fmt.Sprintf("-L%d", to)
	}
	return url
}

func checkFlagCondition() error {
	if *to > 0 && *from == 0 {
		return errors.New("-from required")
	}

	if *from > 0 && *to > 0 && *to <= *from {
		return errors.Errorf("-to must be greater than -from: from=%d, to=%d", *from, *to)
	}
	return nil
}
