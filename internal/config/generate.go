package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

type GenerateConfig struct {
	Apps []GenerateApps
}

type GenerateApps struct {
	Name       string
	Repo       string
	RepoBranch string `mapstructure:"repo_branch"`
	RepoFile   string `mapstructure:"repo_file"`
	FileURL    string `mapstructure:"file_url"`
}

func (a *GenerateApps) Schema() (*ast.Source, error) {
	if a.RepoBranch == "" {
		a.RepoBranch = "main"
	}
	if a.RepoFile == "" {
		a.RepoFile = "schema.graphql"
	}

	if a.FileURL == "" {
		a.FileURL = fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", a.Repo, a.RepoBranch, a.RepoFile)
	}

	isUrl := strings.HasPrefix(a.FileURL, "http")

	if isUrl {
		// fetch from url and return it
		return nil, nil
	}

	b, err := os.ReadFile(a.FileURL)
	if err != nil {
		return nil, err
	}

	return &ast.Source{
		Name:  a.FileURL,
		Input: string(b),
	}, nil
}
