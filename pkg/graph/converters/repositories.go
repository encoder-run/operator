package converters

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/encoder-run/operator/api/cloud/v1alpha1"
	"github.com/encoder-run/operator/pkg/graph/model"
)

func RepositoryCRDToModel(repo *v1alpha1.Repository) (*model.Repository, error) {
	var repoType model.RepositoryType
	switch repo.Spec.Type {
	case v1alpha1.RepositoryTypeGithub:
		repoType = model.RepositoryTypeGithub
	case v1alpha1.RepositoryTypeGitlab:
		repoType = model.RepositoryTypeGitlab
	case v1alpha1.RepositoryTypeBitbucket:
		repoType = model.RepositoryTypeBitbucket
	default:
		return nil, fmt.Errorf("unknown repository type: %s", repo.Spec.Type)
	}

	if repo.Spec.Type == v1alpha1.RepositoryTypeGithub {
		return &model.Repository{
			ID:          repo.Name,
			Name:        repo.Spec.Github.Name,
			Owner:       repo.Spec.Github.Owner,
			Type:        repoType,
			URL:         repo.Spec.Github.URL,
			DisplayName: fmt.Sprintf("%s/%s", repo.Spec.Github.Owner, repo.Spec.Github.Name),
		}, nil
	}

	return nil, fmt.Errorf("unsupported repository type: %s", repo.Spec.Type)
}

func SplitRepositoryURL(url string) (v1alpha1.RepositoryType, string, string, error) {
	// Regular expression to match the URL patterns
	regex := regexp.MustCompile(`(?:https?://)?(?:www\.)?(?P<type>github|gitlab|bitbucket)\.com/(?P<owner>[^/]+)/(?P<name>[^/.]+)(?:\.git)?`)

	matches := regex.FindStringSubmatch(url)
	if matches == nil {
		return "", "", "", errors.New("invalid repository URL")
	}

	// Extracting matches
	repositoryType := matches[regex.SubexpIndex("type")]
	var repoType v1alpha1.RepositoryType
	switch repositoryType {
	case "github":
		repoType = v1alpha1.RepositoryTypeGithub
	case "gitlab":
		repoType = v1alpha1.RepositoryTypeGitlab
	case "bitbucket":
		repoType = v1alpha1.RepositoryTypeBitbucket
	default:
		return "", "", "", fmt.Errorf("unsupported repository type: %s", repositoryType)
	}

	owner := matches[regex.SubexpIndex("owner")]
	name := matches[regex.SubexpIndex("name")]

	return repoType, owner, name, nil
}

func RepositoryURL(repoType model.RepositoryType, owner, name string) string {
	switch repoType {
	case model.RepositoryTypeGithub:
		return fmt.Sprintf("https://github.com/%s/%s", owner, name)
	case model.RepositoryTypeGitlab:
		return fmt.Sprintf("https://gitlab.com/%s/%s", owner, name)
	case model.RepositoryTypeBitbucket:
		return fmt.Sprintf("https://bitbucket.org/%s/%s", owner, name)
	default:
		return ""
	}
}
