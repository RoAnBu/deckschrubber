package main

import (
	"errors"
	"sort"
	"strings"
	"syscall"
	"time"

	"encoding/json"

	"flag"

	"io"

	"fmt"
	"os"

	"regexp"

	"context"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/RoAnBu/deckschrubber/util"
	log "github.com/Sirupsen/logrus"
	distrContext "github.com/docker/distribution/context"
	schema2 "github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/client"
)

var (
	/** CLI flags */
	// Base URL of registry
	registryURL *string
	// Regexps for filtering repositories and tags
	repoRegexpStr, tagRegexpStr, negTagRegexpStr *string
	// File location for repository configuration file
	repoFile *string
	// Maximum age of image to consider for deletion
	day, month, year *int
	// Max number of repositories to be fetched from registry
	repoCount *int
	// Number of the latest n matching images of an repository that will be ignored
	latest *int
	// If true, application runs in debug mode
	debug *bool
	// If true, no actual deletion is done
	dry *bool
	// If true, version is shown and program quits
	ver *bool

	// Skip insecure TLS
	insecure *bool
	// Username and password
	uname, passwd *string

	basicAuthTransport *util.BasicAuthTransport
)

const (
	version string = "0.7.0pre"
)

type repoEntry struct {
	repoRegexString string
	repoRegex       *regexp.Regexp
	maxRepos        int
	tagRegex        *regexp.Regexp
	keepRegex       *regexp.Regexp
	keepNewer       time.Time
	keepLatest      int
}

func init() {
	/** CLI flags */
	// Max number of repositories to fetch from registry (default = 5)
	repoCount = flag.Int("repos", 5, "number of repositories to garbage collect")
	// Base URL of registry (default = http://localhost:5000)
	registryURL = flag.String("registry", "http://localhost:5000", "URL of registry")
	// Maximum age of iamges to consider for deletion in days (default = 0)
	day = flag.Int("day", 0, "max age in days")
	// Maximum age of months to consider for deletion in days (default = 0)
	month = flag.Int("month", 0, "max age in months")
	// Maximum age of iamges to consider for deletion in years (default = 0)
	year = flag.Int("year", 0, "max age in days")
	// Regexp for images (default = .*)
	repoRegexpStr = flag.String("repo", ".*", "matching repositories (allows regexp)")
	// Regexp for tags (default = .*)
	tagRegexpStr = flag.String("tag", ".*", "matching tags (allows regexp)")
	// Negative regexp for tags (default = empty)
	negTagRegexpStr = flag.String("ntag", "", "non matching tags (allows regexp)")
	// The number of the latest matching images of an repository that won't be deleted
	latest = flag.Int("latest", 1, "number of the latest matching images of an repository that won't be deleted")
	// File location of repo config file
	repoFile = flag.String("file", "", "A YAML file containing repositories which shall be cleaned")
	// Show debug log messages
	debug = flag.Bool("debug", false, "run in debug mode")
	// Dry run option (doesn't actually delete)
	dry = flag.Bool("dry", false, "does not actually deletes")
	// Shows version
	ver = flag.Bool("v", false, "shows version and quits")
	// Skip insecure TLS
	insecure = flag.Bool("insecure", false, "Skip insecure TLS verification")
	// Username and password
	uname = flag.String("user", "", "Username for basic authentication")
	passwd = flag.String("password", "", "Password for basic authentication")
}

func main() {
	flag.Parse()

	if *ver {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	// List with repositories to clean
	var repoList []repoEntry

	if *repoFile != "" {
		var err error
		repoList, err = importRepoFile(*repoFile)
		if err != nil {
			log.Fatalf("An error occurred while creating a repository struct. Error: %v", err)
		}
	} else {
		rEntry, err := createRepoEntry(*repoRegexpStr, *repoCount, *tagRegexpStr, *negTagRegexpStr, *day, *month, *year, *latest)
		if err != nil {
			log.Fatalf("An error occurred while creating a repository struct. Error: %v", err)
		}
		repoList = []repoEntry{rEntry}
	}

	// Add basic auth if user/pass is provided
	setupAuthentication(*uname, *passwd, *insecure)

	// Create registry object
	r, err := client.NewRegistry(*registryURL, basicAuthTransport)
	if err != nil {
		log.Fatalf("Could not create registry object! (err: %s", err)
	}

	var maxRepos int
	// calculate max repo value
	for _, repoStruct := range repoList {
		maxRepos += repoStruct.maxRepos
	}

	// List of all repositories fetched from the registry. The number
	// of fetched repositories depends on the number provided by the
	// user ('-repos' flag)
	entries := make([]string, maxRepos)

	// Empty context for all requests in sequel
	ctx := distrContext.Background()

	// Fetch all repositories from the registry
	numFilled, err := r.Repositories(ctx, entries, "")
	if err != nil && err != io.EOF {
		log.Fatalf("Error while fetching repositories! (err: %v)", err)
	}
	log.WithFields(log.Fields{"count": numFilled, "entries": entries}).Info("Successfully fetched repositories.")

	// loop through all repository expressions
	tagsDeletionRatioList := make([]map[string]tagsDeletionRatio, len(repoList))
	for i, repoStruct := range repoList {
		tagsDeletionRatioList[i] = schrubbRepositories(ctx, entries, numFilled, repoStruct.keepNewer, repoStruct.repoRegex, repoStruct.keepRegex, repoStruct.tagRegex, repoStruct.keepLatest)
	}

	// Print deletion stats
	var deleteMsg string
	if *dry {
		deleteMsg = "Would delete"
	} else {
		deleteMsg = "Deleted"
	}
	for _, repoMap := range tagsDeletionRatioList {
		for repo, stats := range repoMap {
			log.Infof("%s: \t%s %d out of %d, making it %.2f%%", deleteMsg, repo, stats.deletedTags, stats.totalTags,
				(float32(stats.deletedTags) / float32(stats.totalTags) * 100.0))
		}
	}
}

func createRepoEntry(repoRexeg string, maxRepos int, tagRegex string, keepRegex string, days int, months int, years int, latest int) (repoEntry, error) {
	defaultRgx := ".*"
	errorEntry := repoEntry{}
	repoEntry := repoEntry{}
	if repoRexeg == "" {
		return errorEntry, errors.New("Repo Regex can not be empty")
	}
	repoEntry.repoRegexString = repoRexeg

	repoEntry.repoRegex = regexp.MustCompile(repoRexeg)
	if maxRepos < 0 {
		return errorEntry, fmt.Errorf("Max Repos value for repository %s cannot be negative, was %d", repoRexeg, maxRepos)
	}
	if maxRepos == 0 {
		repoEntry.maxRepos = 1
	} else {
		repoEntry.maxRepos = maxRepos
	}

	// Regex for tags which images will be considered for deletion
	if tagRegex == "" {
		tagRegex = defaultRgx
	}
	repoEntry.tagRegex = regexp.MustCompile(tagRegex)
	// Regex for tags which images will never get deleted
	if keepRegex == "" {
		repoEntry.keepRegex = nil
	} else {
		repoEntry.keepRegex = regexp.MustCompile(keepRegex)
	}
	// keepNewer defines the youngest creation date for an image
	// to be considered for deletion
	repoEntry.keepNewer = time.Now().AddDate(years/-1, months/-1, days/-1)

	// keepLatest is the min number of newest images to keep
	repoEntry.keepLatest = latest

	log.Debug("Created Repository entry with content:")
	log.Debugf("repoRegex: %s, maxRepos: %d, tagRegex: %s, keepRegex: %s, keepNewer: %s, keepLatest: %d",
		repoRexeg, maxRepos, tagRegex, keepRegex, repoEntry.keepNewer.String(), latest)
	return repoEntry, nil
}

func importRepoFile(filePath string) ([]repoEntry, error) {
	importStruct, err := util.ImportFile(filePath)
	if err != nil {
		log.Fatalf("An error occurred while trying to import repository file. Error: \n %v", err)
	}
	repoList := make([]repoEntry, len(importStruct.Repositories))
	for i, repository := range importStruct.Repositories {
		repoList[i], err = createRepoEntry(repository.RepoNameRegex, repository.MaxRepoCount, repository.RemoveTagRgx, repository.KeepTagRgx,
			repository.KeepNewer.Day, repository.KeepNewer.Month, repository.KeepNewer.Year, repository.KeepLatest)
		if err != nil {
			return nil, err
		}
	}
	return repoList, nil
}

type tagsDeletionRatio struct {
	totalTags   int
	deletedTags int
}

func schrubbRepositories(ctx context.Context, entries []string, numFilled int, deadline time.Time, repoRegexp *regexp.Regexp,
	negTagRegexp *regexp.Regexp, tagRegexp *regexp.Regexp, latest int) map[string]tagsDeletionRatio {
	delStats := make(map[string]tagsDeletionRatio)
	for _, entry := range entries[:numFilled] {
		logger := log.WithField("repo", entry)

		matched := repoRegexp.MatchString(entry)

		if !matched {
			logger.WithFields(log.Fields{"entry": entry}).Debug("Ignore non matching repository (-repo=", (*repoRegexp).String(), ")")
			continue
		}

		// Establish repository object in registry
		repoName, err := reference.WithName(entry)

		if err != nil {
			logger.Fatalf("Could not parse repo from name! (err: %v)", err)
		}

		repo, err := client.NewRepository(repoName, *registryURL, basicAuthTransport)
		if err != nil {
			logger.WithFields(log.Fields{"entry": entry}).Fatalf("Could not create repo from name! (err: %v)", err)
		}
		logger.Debug("Successfully created repository object.")

		tagsService := repo.Tags(ctx)
		blobsService := repo.Blobs(ctx)
		manifestService, err := repo.Manifests(ctx)
		if err != nil {
			logger.Fatalf("Couldn't fetch manifest service! (err: %v)", err)
		}

		tagsData, err := tagsService.All(ctx)
		if err != nil {
			logger.Fatalf("Couldn't fetch tags! (err: %v)", err)
		}

		var tags []Image

		// Fetch information about each tag of a repository
		// This involves fetching the manifest, its details,
		// and the corresponding blob information
		tagFetchDataErrors := false
		for _, tag := range tagsData {
			tagLogger := logger.WithField("tag", tag)

			tagLogger.Debug("Fetching tag...")
			desc, err := tagsService.Get(ctx, tag)
			if err != nil {
				tagLogger.Error("Could not fetch tag!")
				tagFetchDataErrors = true
				break
			}

			tagLogger.Debug("Fetching manifest...")
			mnfst, err := manifestService.Get(ctx, desc.Digest)
			if err != nil {
				tagLogger.Error("Could not fetch manifest!")
				tagFetchDataErrors = true
				break
			}

			tagLogger.Debug("Parsing manifest details...")
			_, p, err := mnfst.Payload()
			if err != nil {
				tagLogger.Error("Could not parse manifest detail!")
				tagFetchDataErrors = true
				break
			}

			m := new(schema2.DeserializedManifest)
			m.UnmarshalJSON(p)

			tagLogger.Debug("Fetching blob")
			b, err := blobsService.Get(ctx, m.Manifest.Config.Digest)
			if err != nil {
				tagLogger.Error("Could not fetch blob!")
				tagFetchDataErrors = true
				break
			}

			var blobInfo BlobInfo
			json.Unmarshal(b, &blobInfo)

			tags = append(tags, Image{entry, tag, blobInfo.Created, desc})
		}

		if tagFetchDataErrors {
			// In case of error at any one tag, skip entire repo
			// (avoid acting on incomplete data, which migth lead to
			// deleting images that are actually in use)
			logger.Error("Error obtaining tag data - skipping this repo")
			continue
		}

		sort.Sort(ImageByDate(tags))

		logger.Debug("Analyzing tags...")

		tagCount := len(tags)

		if tagCount == 0 {
			logger.Debug("Ignore repository with no matching tags")
			continue
		}

		deletableTags := make(map[int]Image)
		nonDeletableTags := make(map[int]Image)

		ignoredTags := 0

		for tagIndex := len(tags) - 1; tagIndex >= 0; tagIndex-- {
			tag := tags[tagIndex]
			markForDeletion := false
			tagLogger := logger.WithField("tag", tag.Tag)

			// Provides a text which is followed by the tag and ntag flag values. The
			// latter iff defined.
			withTagParens := func(text string) string {
				xs := []string{fmt.Sprintf("-tag=%s", *tagRegexpStr)}
				if *negTagRegexpStr != "" {
					xs = append(xs, fmt.Sprintf("-ntag=%s", *negTagRegexpStr))
				}
				return fmt.Sprintf("%s (%s)", text, strings.Join(xs, ", "))
			}

			// Check whether the tag matches. If that's the case, don't stop there, and
			// check for the negative regexp as well.
			matched := tagRegexp.MatchString(tag.Tag)
			if matched && negTagRegexp != nil {
				negTagMatch := negTagRegexp.MatchString(tag.Tag)
				matched = !negTagMatch
			}

			if matched {
				tagLogger.Debug(withTagParens("Tag matches, considering for deletion"))
				if tag.Time.Before(deadline) {
					if ignoredTags < latest {
						tagLogger.WithField("time", tag.Time).Infof("Ignore %d latest matching tags (-latest=%d)", latest, latest)
						ignoredTags++
					} else {
						tagLogger.WithField("tag", tag.Tag).WithField("time", tag.Time).Infof("Marking tag as outdated")
						markForDeletion = true
					}
				} else {
					tagLogger.Info("Tag not outdated")
					ignoredTags++
				}
			} else {
				tagLogger.Info(withTagParens("Ignore non matching tag"))
			}

			if markForDeletion {
				deletableTags[tagIndex] = tag
			} else {
				nonDeletableTags[tagIndex] = tag
			}
		}

		// This approach is actually a workaround for the problem that Docker
		// Distribution doesn't implement TagService.Untag operation at the time of
		// this writing.
		// Actually we have to delete the underlying image (specified via its Digest
		// value), taking care not to delete images that are referenced by tags which
		// we want to preserve
		nonDeletableDigests := make(map[string]string)
		for _, tag := range nonDeletableTags {
			if nonDeletableDigests[tag.Descriptor.Digest.String()] == "" {
				nonDeletableDigests[tag.Descriptor.Digest.String()] = tag.Tag
			} else {
				nonDeletableDigests[tag.Descriptor.Digest.String()] = nonDeletableDigests[tag.Descriptor.Digest.String()] + ", " + tag.Tag
			}
		}
		var numDeletedTags int
		digestsDeleted := make(map[string]bool)
		for _, tag := range deletableTags {
			if !digestsDeleted[tag.Descriptor.Digest.String()] {
				if nonDeletableDigests[tag.Descriptor.Digest.String()] == "" {
					// TODO Add deletion stats
					logger.WithField("tag", tag.Tag).Info("All tags for this image digest marked for deletion")
					numDeletedTags++
					if !*dry {
						logger.WithField("tag", tag.Tag).WithField("time", tag.Time).WithField("digest", tag.Descriptor.Digest).Infof("Deleting image (-dry=%v)", *dry)
						err := manifestService.Delete(context.Background(), tag.Descriptor.Digest)
						if err == nil {
							digestsDeleted[tag.Descriptor.Digest.String()] = true
						} else {
							logger.WithField("tag", tag.Tag).WithField("err", err).Error("Could not delete image!")
						}
					} else {
						logger.WithField("tag", tag.Tag).WithField("time", tag.Time).Infof("Not actually deleting image (-dry=%v)", *dry)
					}
				} else {
					logger.WithField("tag", tag.Tag).WithField("alsoUsedByTags", nonDeletableDigests[tag.Descriptor.Digest.String()]).Infof("The underlying image is also used by non-deletable tags - skipping deletion")
				}
			} else {
				logger.WithField("tag", tag.Tag).Debug("Image under tag already deleted")
			}
		}
		delStats[entry] = tagsDeletionRatio{
			len(tagsData),
			numDeletedTags,
		}
	}
	return delStats
}

func setupAuthentication(uname string, passwd string, insecure bool) {
	if uname != "" && passwd == "" {
		fmt.Println("Password:")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err == nil {
			stringPassword := string(bytePassword[:])
			passwd = stringPassword
		} else {
			fmt.Println("Could not read password. Quitting!")
			os.Exit(1)
		}
	}
	basicAuthTransport = util.NewBasicAuthTransport(*registryURL, uname, passwd, insecure)
}
