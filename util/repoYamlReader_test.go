package util

import "testing"

func TestImportFileEmpty(t *testing.T) {

	actualResult, err := ImportFile("test_files/repositories_empty.yml")

	if err != nil {
		t.Fatalf("Got error while parsing empty yaml file: %v", err)
		t.FailNow()
	}
	if len(actualResult.Repositories) != 0 {
		t.Fatalf("Expected repositories array length of 0 but got %d with content %s", len(actualResult.Repositories), ImportStructToString(actualResult))
		t.FailNow()
	}
}

func TestImportFileInvalidValue(t *testing.T) {

	_, err := ImportFile("test_files/repositories_invalid_value.yml")

	if err == nil {
		t.Fatalf("Expected parsing to fail but got no error")
		t.FailNow()
	}
}

func TestImportFileSimpleFilled(t *testing.T) {

	actualResult, err := ImportFile("test_files/repositories_simple_filled.yml")

	if err != nil {
		t.Fatalf("Got error while parsing filled yaml file: %v", err)
		t.FailNow()
	}
	if len(actualResult.Repositories) != 1 {
		t.Fatalf("Expected repositories array length of 1 but got %d with content %s", len(actualResult.Repositories), ImportStructToString(actualResult))
		t.FailNow()
	}

	// Test content

	repository := actualResult.Repositories[0]

	if repository.RepoNameRegex != "^test_repo$" {
		t.Fatalf("Wrong field. Expected: %s got: %s", "^test_repo$", repository.RepoNameRegex)
		t.Fail()
	}

	if repository.MaxRepoCount != 3 {
		t.Fatalf("Wrong field. Expected: %d got: %d", 3, repository.MaxRepoCount)
		t.Fail()
	}

	if repository.RemoveTagRgx != ".*" {
		t.Fatalf("Wrong field. Expected: %s got: %s", ".*", repository.RemoveTagRgx)
		t.Fail()
	}

	if repository.KeepTagRgx != "latest" {
		t.Fatalf("Wrong field. Expected: %s got: %s", "latest", repository.KeepTagRgx)
		t.Fail()
	}

	if repository.KeepNewer.Day != 1 || repository.KeepNewer.Month != 2 || repository.KeepNewer.Year != 6 {
		t.Fatalf("Wrong field. Expected: %d, %d, %d got: %d, %d, %d", 1, 2, 6, repository.KeepNewer.Day, repository.KeepNewer.Month, repository.KeepNewer.Year)
		t.Fail()
	}
}

func TestImportFileSimpleGaps(t *testing.T) {

	actualResult, err := ImportFile("test_files/repositories_simple_gaps.yml")

	if err != nil {
		t.Fatalf("Got error while parsing filled yaml file: %v", err)
		t.FailNow()
	}
	if len(actualResult.Repositories) != 1 {
		t.Fatalf("Expected repositories array length of 1 but got %d with content %s", len(actualResult.Repositories), ImportStructToString(actualResult))
		t.FailNow()
	}

	// Test content

	repository := actualResult.Repositories[0]

	if repository.RepoNameRegex != "^test_repo$" {
		t.Fatalf("Wrong field. Expected: %s got: %s", "^test_repo$", repository.RepoNameRegex)
		t.Fail()
	}

	if repository.MaxRepoCount != 3 {
		t.Fatalf("Wrong field. Expected: %d got: %d", 3, repository.MaxRepoCount)
		t.Fail()
	}

	if repository.RemoveTagRgx != ".*" {
		t.Fatalf("Wrong field. Expected: %s got: %s", ".*", repository.RemoveTagRgx)
		t.Fail()
	}

	if repository.KeepTagRgx != "" {
		t.Fatalf("Wrong field. Expected: %s got: %s", "", repository.KeepTagRgx)
		t.Fail()
	}

	if repository.KeepNewer.Day != 0 || repository.KeepNewer.Month != 2 || repository.KeepNewer.Year != 6 {
		t.Fatalf("Wrong field. Expected: %d, %d, %d got: %d, %d, %d", 0, 2, 6, repository.KeepNewer.Day, repository.KeepNewer.Month, repository.KeepNewer.Year)
		t.Fail()
	}
}

func TestImportFileLong(t *testing.T) {

	actualResult, err := ImportFile("test_files/repositories_long.yml")

	if err != nil {
		t.Fatalf("Got error while parsing filled yaml file: %v", err)
		t.FailNow()
	}
	if len(actualResult.Repositories) != 3 {
		t.Fatalf("Expected repositories array length of 3 but got %d with content %s", len(actualResult.Repositories), ImportStructToString(actualResult))
		t.FailNow()
	}

	// Test content

	repository := actualResult.Repositories[0]

	if repository.RepoNameRegex != "^test_repo$" {
		t.Fatalf("Wrong field. Expected: %s got: %s", "^test_repo$", repository.RepoNameRegex)
		t.Fail()
	}

	if repository.MaxRepoCount != 3 {
		t.Fatalf("Wrong field. Expected: %d got: %d", 3, repository.MaxRepoCount)
		t.Fail()
	}

	if repository.RemoveTagRgx != ".*" {
		t.Fatalf("Wrong field. Expected: %s got: %s", ".*", repository.RemoveTagRgx)
		t.Fail()
	}

	if repository.KeepTagRgx != "latest" {
		t.Fatalf("Wrong field. Expected: %s got: %s", "latest", repository.KeepTagRgx)
		t.Fail()
	}

	if repository.KeepNewer.Day != 1 || repository.KeepNewer.Month != 2 || repository.KeepNewer.Year != 6 {
		t.Fatalf("Wrong field. Expected: %d, %d, %d got: %d, %d, %d", 1, 2, 6, repository.KeepNewer.Day, repository.KeepNewer.Month, repository.KeepNewer.Year)
		t.Fail()
	}

	repository = actualResult.Repositories[1]

	if repository.RepoNameRegex != "^test_repo2$" {
		t.Fatalf("Wrong field. Expected: %s got: %s", "^test_repo2$", repository.RepoNameRegex)
		t.Fail()
	}

	if repository.MaxRepoCount != 5 {
		t.Fatalf("Wrong field. Expected: %d got: %d", 5, repository.MaxRepoCount)
		t.Fail()
	}

	if repository.RemoveTagRgx != "" {
		t.Fatalf("Wrong field. Expected: %s got: %s", "", repository.RemoveTagRgx)
		t.Fail()
	}

	if repository.KeepTagRgx != "latest" {
		t.Fatalf("Wrong field. Expected: %s got: %s", "latest", repository.KeepTagRgx)
		t.Fail()
	}

	if repository.KeepNewer.Day != 0 || repository.KeepNewer.Month != 2 || repository.KeepNewer.Year != 0 {
		t.Fatalf("Wrong field. Expected: %d, %d, %d got: %d, %d, %d", 0, 2, 0, repository.KeepNewer.Day, repository.KeepNewer.Month, repository.KeepNewer.Year)
		t.Fail()
	}

	repository = actualResult.Repositories[2]

	if repository.RepoNameRegex != "^test_repo3$" {
		t.Fatalf("Wrong field. Expected: %s got: %s", "^test_repo3$", repository.RepoNameRegex)
		t.Fail()
	}

	if repository.MaxRepoCount != 1 {
		t.Fatalf("Wrong field. Expected: %d got: %d", 1, repository.MaxRepoCount)
		t.Fail()
	}

	if repository.RemoveTagRgx != "^\\d{5}$" {
		t.Fatalf("Wrong field. Expected: %s got: %s", "^\\d{5}$", repository.RemoveTagRgx)
		t.Fail()
	}

	if repository.KeepTagRgx != "-latest" {
		t.Fatalf("Wrong field. Expected: %s got: %s", "latest", repository.KeepTagRgx)
		t.Fail()
	}

	if repository.KeepNewer.Day != 0 || repository.KeepNewer.Month != 0 || repository.KeepNewer.Year != 0 {
		t.Fatalf("Wrong field. Expected: %d, %d, %d got: %d, %d, %d", 0, 0, 0, repository.KeepNewer.Day, repository.KeepNewer.Month, repository.KeepNewer.Year)
		t.Fail()
	}
}
