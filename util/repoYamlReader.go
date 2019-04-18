package util

import (
	"strconv"

	log "github.com/Sirupsen/logrus"

	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ImportStruct contains the repository entries imported from the yaml file
type ImportStruct struct {
	Repositories []RepositoryStruct `yaml:"repositories"`
}

// RepositoryStruct represents a repository entry of a yaml file
type RepositoryStruct struct {
	RepoNameRegex string          `yaml:"repo-name-rgx"`
	MaxRepoCount  int             `yaml:"max-repo-count"`
	RemoveTagRgx  string          `yaml:"remove-tag-rgx"`
	KeepTagRgx    string          `yaml:"keep-tag-rgx"`
	KeepNewer     KeepNewerStruct `yaml:"keep-newer"`
	KeepLatest    int             `yaml:"keep-latest"`
}

// KeepNewerStruct represents the age deadline for images in a repository
type KeepNewerStruct struct {
	Year  int `yaml:"year"`
	Month int `yaml:"month"`
	Day   int `yaml:"day"`
}

// ImportFile reads the given yaml file and parses the fields as provided in repositories_sample.yml
func ImportFile(filePath string) (ImportStruct, error) {
	log.Info("Starting yaml file import and parsing")
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorf("Error while trying to read a yaml file with path: %s", filePath)
		return ImportStruct{}, err
	}
	log.Debug("Successfully read file. Starting parsing...")
	var importContent ImportStruct

	err = yaml.Unmarshal(yamlFile, &importContent)
	if err != nil {
		log.Errorf("Error while parsing yaml file %s", filePath)
		return ImportStruct{}, err
	}
	log.Info("Successfully read and parsed yaml.")
	log.Debugf("Imported content:\n%v", importContent)
	return importContent, err
}

// ImportStructToString returns a formatted string with the content of the given import struct
func ImportStructToString(imp ImportStruct) string {
	const newLine string = "\n"
	var result string

	for _, repository := range imp.Repositories {
		result += "RepoNameRegex: " + repository.RepoNameRegex + newLine
		result += "MaxRepoCount: " + strconv.Itoa(repository.MaxRepoCount) + newLine
		result += "RemoveTagRgx: " + repository.RemoveTagRgx + newLine
		result += "KeepTagRgx: " + repository.KeepTagRgx + newLine
		result += "Day: " + strconv.Itoa(repository.KeepNewer.Day) + newLine
		result += "Month: " + strconv.Itoa(repository.KeepNewer.Month) + newLine
		result += "Year: " + strconv.Itoa(repository.KeepNewer.Year) + newLine + newLine
		result += "Latest: " + strconv.Itoa(repository.KeepLatest) + newLine + newLine
	}
	return result
}
