package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/blang/semver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var major *bool
var minor *bool
var patch *bool
var meta *string
var ver *string
var git *bool
var push *bool
var (
	rType = regexp.MustCompile(`\+(major|minor|patch)`)
	rMeta = regexp.MustCompile(`\+meta=([0-9A-Za-z-\.]+)`)
)
var multiBumpError = errors.New("You specified more than one bump type!")

type bumpType int

const (
	bumpPatch bumpType = iota
	bumpMinor
	bumpMajor
)

type bumpOptions struct {
	major bool
	minor bool
	patch bool
	meta  string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "semver",
	Short: "Bumps up the semver accordingly",
	Long: `Bumps up the semver, given a current semver and a type to update. 

Examples:
* "semver 0.1.2 --minor" will return 0.2.0,
* "semver 0.1.2 " will return 0.1.3 (patch is the default bump type)
* "semver 0.1.2+metadata --meta newMetaData" will return 0.1.2+newMetaData

If "--git" is specified at runtime, the new semver is determined by using:
- The latest tag (on the current branch of the repo where it is being ran) AND
- The bump type(s) found in the commit messgaes from the last tag on the current branch up to HEAD. 
 * If anywhere in the messages, a string of types "+major", "+minor" or "+patch" is found, "major", 
 "minor" or "patch" will be used as the bump type, respectively. 
 * If more than one bump type is found in the commit messages, the highest one will be used as the 
bump type (major > minor > patch)

Similarly, if the commit message contains +meta=<some_meta_data> the meta data will be 
attached to the new sematic version when specifying --git

For example, assume that the latest tag was 1.3.4+6.8 and that the last commit message was

"
Adding new functionality in a backwards-compatible manner

New feature to do something cool!

+minor
+meta=7.3

That's all folks
"

In this case, running 

* "semver --git" will print 1.4.0+7.3

If you want to create a tag using this semver and push it, then run:
* "semver --git --push" will print 1.4.0+7.3 and it will also push the tag to the repo.`,

	Args: func(cmd *cobra.Command, args []string) error {
		if *git {
			if len(args) > 0 {
				return errors.New("A <semver> arg and --git cannot be specified at the same time!")
			}
		} else if len(args) != 1 {
			return errors.New("Either a <semver> arg or --git is required")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var ver string

		if *git {
			fetchTagErr := exec.Command("git", "fetch", "origin", `+refs/tags/*:refs/tags/*`).Run()
			if fetchTagErr != nil {
				return fetchTagErr
			}

			tag, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
			if err != nil {
				return err
			}

			ver = strings.TrimSpace(string(tag))

			prevMessages, prevMessagesErr := exec.Command("git", "log", "--format=%B", ver+"..HEAD").Output()
			if prevMessagesErr != nil {
				return prevMessagesErr
			}

			message := string(prevMessages)
			parentBumps := readBumpTypesFromString(message)

			if len(parentBumps) == 0 {
				*patch = true
			} else {
				switch coalescedBump := findMaxBump(parentBumps); coalescedBump {
				case bumpMajor:
					*major = true
				case bumpMinor:
					*minor = true
				case bumpPatch:
					*patch = true
				default:
					return errors.New("No valid bump type found!")
				}
			}

			if rMeta.MatchString(message) {
				*meta = rMeta.FindStringSubmatch(message)[1]
			}
		} else {
			ver = args[0]
		}

		thisSemver, err := semver.Make(ver)
		if err != nil {
			return err
		}

		err = bump(&thisSemver, bumpOptions{major: *major, minor: *minor, patch: *patch, meta: *meta})
		if err != nil {
			return err
		}

		fmt.Println(thisSemver.String())

		if *git && *push {
			tagErr := exec.Command("git", "tag", thisSemver.String()).Run()
			if tagErr != nil {
				return tagErr
			}

			pushErr := exec.Command("git", "push", "origin", "--tags").Run()
			if pushErr != nil {
				return pushErr
			}
		}
		return nil
	},
}

func bump(v *semver.Version, opt bumpOptions) error {
	if !(opt.major || opt.minor) {
		v.Patch++
	} else if opt.major && opt.minor || opt.major && opt.patch || opt.minor && opt.patch {
		return multiBumpError
	} else if opt.major {
		v.Major++
		v.Minor = 0
		v.Patch = 0
	} else if opt.minor {
		v.Minor++
		v.Patch = 0
	} else {
		return errors.New("No valid bump type was specified!")
	}

	if opt.meta != "" {
		var buildErr error

		metaData := strings.Split(opt.meta, ".")

		for i := range metaData {
			v.Build[i], buildErr = semver.NewBuildVersion(metaData[i])
			if buildErr != nil {
				return errors.New("[" + opt.meta + "] is not a valid metadata format! Check https://semver.org/#spec-item-10")
			}
		}
	}

	return nil
}

func readBumpTypesFromString(message string) []bumpType {
	var bumps []bumpType

	if rType.MatchString(message) {
		bumpType := rType.FindAllStringSubmatch(message, -1)

		for _, match := range bumpType {
			switch match[1] {
			case "major":
				bumps = append(bumps, bumpMajor)
			case "minor":
				bumps = append(bumps, bumpMinor)
			case "patch":
				bumps = append(bumps, bumpPatch)
			}
		}
	}

	return bumps
}

func findMaxBump(bumps []bumpType) bumpType {
	var maxBump bumpType
	maxBump = bumpPatch

	for _, bump := range bumps {
		if bump > maxBump {
			maxBump = bump
		}
	}

	return maxBump
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.semver.yaml)")
	major = rootCmd.PersistentFlags().BoolP("major", "", false, "Bump major version")
	minor = rootCmd.PersistentFlags().BoolP("minor", "", false, "Bump minor version")
	patch = rootCmd.PersistentFlags().BoolP("patch", "", false, "Bump patch version")
	meta = rootCmd.PersistentFlags().StringP("meta", "", "", "Metadata to be attached to the semantic version")
	git = rootCmd.PersistentFlags().BoolP("git", "", false, "Use git to get the current version and bump type")
	push = rootCmd.PersistentFlags().BoolP("push", "", false, "Push new semver as a tag to the repo")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".semver") // name of config file (without extension)
	viper.AddConfigPath("$HOME")   // adding home directory as first search path
	viper.AutomaticEnv()           // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func main() {
	Execute()
}
