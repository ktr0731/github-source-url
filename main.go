package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/github/hub/github"
	"github.com/pkg/errors"
)

var (
	format = "%s/blob/%s/%s"
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
		flag.Usage()
		return 1, nil
	}

	if err := checkFlagCondition(); err != nil {
		return 1, errors.Wrap(err, "precondition failed")
	}

	absPath, err := filepath.Abs(flag.Arg(0))
	if err != nil {
		return 1, errors.Wrap(err, "failed to get absolute path")
	}

	// change dir and reset
	defer func(path string) func() {
		prev, _ := filepath.Abs(".")
		os.Chdir(path)
		return func() {
			os.Chdir(prev)
		}
	}(filepath.Dir(absPath))()

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

	path, err := regularizePath(absPath)
	if err != nil {
		return 1, err
	}
	fmt.Println(formatURL(proj.WebURL("", "", ""), br.ShortName(), path, *from, *to))
	return 0, nil
}

func formatURL(host, ref, path string, from, to uint) string {
	url := fmt.Sprintf(format, host, ref, path)
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

func regularizePath(path string) (string, error) {
	p, err := filepath.Abs(path)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse host as an URL")
	}

	out, err := exec.Command("git", "rev-parse", "-q", "--absolute-git-dir").CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, "failed to get project root")
	}

	return strings.Replace(p, filepath.Dir(string(out))+string(filepath.Separator), "", 1), nil
}
