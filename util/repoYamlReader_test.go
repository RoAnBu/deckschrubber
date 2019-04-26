package util

import "testing"

func TestImportFileEmpty(t *testing.T) {

	actualResult, err := ImportFile("test_files/repositories_empty.yml")

	if err != nil {
		t.Fatalf("Got error while parsing empty yaml file: %v", err)
	}
	if len(actualResult.Repositories) != 0 {
		t.Fatalf("Expected repositories array length of 0 but got %d with content %s", len(actualResult.Repositories), ImportStructToString(actualResult))
	}
	if actualResult.MaxRepos != 0 {
		t.Fatalf("Expected max repo value of 0 but got %d with content %s", actualResult.MaxRepos, ImportStructToString(actualResult))
	}
}

func TestImportFileInvalidValue(t *testing.T) {

	_, err := ImportFile("test_files/repositories_invalid_value.yml")

	if err == nil {
		t.Fatalf("Expected parsing to fail but got no error")
	}
}

func TestImportFileSimpleFilled(t *testing.T) {

	actualResult, err := ImportFile("test_files/repositories_simple_filled.yml")

	if err != nil {
		t.Fatalf("Got error while parsing filled yaml file: %v", err)
	}
	if len(actualResult.Repositories) != 1 {
		t.Errorf("Expected repositories array length of 1 but got %d with content %s", len(actualResult.Repositories), ImportStructToString(actualResult))
	}

	if actualResult.MaxRepos != 23 {
		t.Fatalf("Expected max repo value of 23 but got %d with content %s", actualResult.MaxRepos, ImportStructToString(actualResult))
	}

	// Test content

	repository := actualResult.Repositories[0]

	if repository.RepoNameRegex != "^test_repo$" {
		t.Errorf("Wrong field. Expected: %s got: %s", "^test_repo$", repository.RepoNameRegex)
	}

	if repository.RemoveTagRgx != ".*" {
		t.Errorf("Wrong field. Expected: %s got: %s", ".*", repository.RemoveTagRgx)
	}

	if repository.KeepTagRgx != "latest" {
		t.Errorf("Wrong field. Expected: %s got: %s", "latest", repository.KeepTagRgx)
	}

	if repository.KeepNewer.Day != 1 || repository.KeepNewer.Month != 2 || repository.KeepNewer.Year != 6 {
		t.Errorf("Wrong field. Expected: %d, %d, %d got: %d, %d, %d", 1, 2, 6, repository.KeepNewer.Day, repository.KeepNewer.Month, repository.KeepNewer.Year)
	}
	if repository.KeepLatest != 5 {
		t.Errorf("Wrong field. Expected: %d got: %d", 5, repository.KeepLatest)
	}
}

func TestImportFileSimpleGaps(t *testing.T) {

	actualResult, err := ImportFile("test_files/repositories_simple_gaps.yml")

	if err != nil {
		t.Fatalf("Got error while parsing filled yaml file: %v", err)
	}
	if len(actualResult.Repositories) != 1 {
		t.Fatalf("Expected repositories array length of 1 but got %d with content %s", len(actualResult.Repositories), ImportStructToString(actualResult))
	}

	if actualResult.MaxRepos != 25 {
		t.Fatalf("Expected max repo value of 23 but got %d with content %s", actualResult.MaxRepos, ImportStructToString(actualResult))
	}

	// Test content

	repository := actualResult.Repositories[0]

	if repository.RepoNameRegex != "^test_repo$" {
		t.Errorf("Wrong field. Expected: %s got: %s", "^test_repo$", repository.RepoNameRegex)
	}

	if repository.RemoveTagRgx != ".*" {
		t.Errorf("Wrong field. Expected: %s got: %s", ".*", repository.RemoveTagRgx)
	}

	if repository.KeepTagRgx != "" {
		t.Errorf("Wrong field. Expected: %s got: %s", "", repository.KeepTagRgx)
	}

	if repository.KeepNewer.Day != 0 || repository.KeepNewer.Month != 2 || repository.KeepNewer.Year != 6 {
		t.Errorf("Wrong field. Expected: %d, %d, %d got: %d, %d, %d", 0, 2, 6, repository.KeepNewer.Day, repository.KeepNewer.Month, repository.KeepNewer.Year)
	}

	if repository.KeepLatest != 0 {
		t.Errorf("Wrong field. Expected: %d got: %d", 0, repository.KeepLatest)
	}
}

func TestImportFileLong(t *testing.T) {

	actualResult, err := ImportFile("test_files/repositories_long.yml")

	if err != nil {
		t.Fatalf("Got error while parsing filled yaml file: %v", err)
	}
	if len(actualResult.Repositories) != 3 {
		t.Fatalf("Expected repositories array length of 3 but got %d with content %s", len(actualResult.Repositories), ImportStructToString(actualResult))
	}

	if actualResult.MaxRepos != 1500 {
		t.Fatalf("Expected max repo value of 23 but got %d with content %s", actualResult.MaxRepos, ImportStructToString(actualResult))
	}

	// Test content

	repository := actualResult.Repositories[0]

	if repository.RepoNameRegex != "^test_repo$" {
		t.Errorf("Wrong field. Expected: %s got: %s", "^test_repo$", repository.RepoNameRegex)
	}

	if repository.RemoveTagRgx != ".*" {
		t.Errorf("Wrong field. Expected: %s got: %s", ".*", repository.RemoveTagRgx)
	}

	if repository.KeepTagRgx != "latest" {
		t.Errorf("Wrong field. Expected: %s got: %s", "latest", repository.KeepTagRgx)
	}

	if repository.KeepNewer.Day != 1 || repository.KeepNewer.Month != 2 || repository.KeepNewer.Year != 6 {
		t.Errorf("Wrong field. Expected: %d, %d, %d got: %d, %d, %d", 1, 2, 6, repository.KeepNewer.Day, repository.KeepNewer.Month, repository.KeepNewer.Year)
	}

	if repository.KeepLatest != 0 {
		t.Errorf("Wrong field. Expected: %d got: %d", 0, repository.KeepLatest)
	}

	repository = actualResult.Repositories[1]

	if repository.RepoNameRegex != "^test_repo2$" {
		t.Errorf("Wrong field. Expected: %s got: %s", "^test_repo2$", repository.RepoNameRegex)
	}

	if repository.RemoveTagRgx != "" {
		t.Errorf("Wrong field. Expected: %s got: %s", "", repository.RemoveTagRgx)
	}

	if repository.KeepTagRgx != "latest" {
		t.Errorf("Wrong field. Expected: %s got: %s", "latest", repository.KeepTagRgx)
	}

	if repository.KeepNewer.Day != 0 || repository.KeepNewer.Month != 2 || repository.KeepNewer.Year != 0 {
		t.Errorf("Wrong field. Expected: %d, %d, %d got: %d, %d, %d", 0, 2, 0, repository.KeepNewer.Day, repository.KeepNewer.Month, repository.KeepNewer.Year)
	}

	if repository.KeepLatest != 0 {
		t.Errorf("Wrong field. Expected: %d got: %d", 0, repository.KeepLatest)
	}

	repository = actualResult.Repositories[2]

	if repository.RepoNameRegex != "^test_repo3$" {
		t.Errorf("Wrong field. Expected: %s got: %s", "^test_repo3$", repository.RepoNameRegex)
	}

	if repository.RemoveTagRgx != "^\\d{5}$" {
		t.Errorf("Wrong field. Expected: %s got: %s", "^\\d{5}$", repository.RemoveTagRgx)
	}

	if repository.KeepTagRgx != "-latest" {
		t.Errorf("Wrong field. Expected: %s got: %s", "latest", repository.KeepTagRgx)
	}

	if repository.KeepNewer.Day != 0 || repository.KeepNewer.Month != 0 || repository.KeepNewer.Year != 0 {
		t.Errorf("Wrong field. Expected: %d, %d, %d got: %d, %d, %d", 0, 0, 0, repository.KeepNewer.Day, repository.KeepNewer.Month, repository.KeepNewer.Year)
	}

	if repository.KeepLatest != 0 {
		t.Errorf("Wrong field. Expected: %d got: %d", 0, repository.KeepLatest)
	}
}
