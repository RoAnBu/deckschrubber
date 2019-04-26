package main

import (
	"os"
	"testing"

	"github.com/RoAnBu/deckschrubber/util"
)

func TestRepoFileValid(t *testing.T) {
	envRepoFile, envRepoFileExists := os.LookupEnv("DS_REPO_FILE_PATH")
	if envRepoFileExists == false {
		t.Skipf("No Repository file environment variable found, skipping validating repo file")
	}
	actualResult, err := util.ImportFile(envRepoFile)

	if err != nil {
		t.Fatalf("Got error while parsing yaml file: %v", err)
	}

	repoList := make([]repoEntry, len(actualResult.Repositories))
	for i, repository := range actualResult.Repositories {
		repoList[i], err = createRepoEntry(repository.RepoNameRegex, repository.RemoveTagRgx, repository.KeepTagRgx,
			repository.KeepNewer.Day, repository.KeepNewer.Month, repository.KeepNewer.Year, repository.KeepLatest)
		if err != nil {
			t.Fatalf("Value import error while creating repository entry")
		}
	}
}
