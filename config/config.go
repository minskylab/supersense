package config

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

const appName = "ss" // supersense

// Config wraps all the necessary
type Config struct {
	Port              int64 `default:"8080" split_words:"true"`
	Debug             bool  `default:"false"`
	GraphQLPlayground bool  `envconfig:"GRAPHQL_PLAYGROUND" default:"false"`

	Persistence               bool   `split_words:"true" default:"false"`
	PersistenceBoltDBFilePath string `envconfig:"PERSISTENCE_BOLTDB_FILEPATH" default:"ss.db"`
	PersistenceRedisAddress   string `split_words:"true"`
	PersistenceRedisPassword  string `split_words:"true"`

	GithubToken string   `split_words:"true"`
	GithubRepos []string `split_words:"true"`

	DummyPeriod  string `split_words:"true"`
	DummyMessage string `split_words:"true"`

	TwitterConsumerKey    string   `split_words:"true"`
	TwitterConsumerSecret string   `split_words:"true"`
	TwitterAccessToken    string   `split_words:"true"`
	TwitterAccessSecret   string   `split_words:"true"`
	TwitterQuery          []string `split_words:"true"`

	Spokesman         bool   `split_words:"true" default:"false"`
	SpokesmanName     string `split_words:"true" default:"Spokesman"`
	SpokesmanUsername string `split_words:"true" default:"spokesman"`
	SpokesmanEmail    string `split_words:"true"`

	RootCredentialUsername string `split_words:"true" default:"root"`
	RootCredentialPassword string `split_words:"true" default:""`

	ObserverBuffer  int    `split_words:"true" default:"20"`
	ObserverTitle   string `split_words:"true" default:"Hello World"`
	ObserverHashtag string `split_words:"true" default:"#opensource"`
	ObserverBrand   string `split_words:"true" default:"SUPERSENSE"`
}

func load(appName string) (*Config, error) {
	conf := new(Config)
	if err := envconfig.Process(appName, conf); err != nil {
		return nil, errors.WithStack(err)
	}

	conf.ObserverBrand = strings.ReplaceAll(conf.ObserverBrand, "\"", "")
	conf.ObserverTitle = strings.ReplaceAll(conf.ObserverTitle, "\"", "")
	conf.ObserverHashtag = strings.ReplaceAll(conf.ObserverHashtag, "\"", "")

	return conf, nil
}

// LoadDefault load a minimal default configuration
func LoadDefault() (*Config, error) {
	return load(appName)
}
