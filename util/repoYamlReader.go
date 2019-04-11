package util

import (
	log "github.com/Sirupsen/logrus"

	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ImportStruct TODO
type ImportStruct struct {
	Repositories []RepositoryStruct `yaml:"repositories"`
}

// RepositoryStruct TODO
type RepositoryStruct struct {
	RepoNameRegex string          `yaml:"repo-name-rgx"`
	MaxRepoCount  int             `yaml:"max-repo-count"`
	KeepTagRgx    string          `yaml:"keep-tag-rgx"`
	KeepNewer     KeepNewerStruct `yaml:"keep-newer"`
}

// KeepNewerStruct TODO
type KeepNewerStruct struct {
	Year  int `yaml:"year"`
	Month int `yaml:"month"`
	Day   int `yaml:"day"`
}

// ImportFile TODO documentation
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
