package repositories

import (
	"context"
	"fmt"

	"github.com/encoder-run/operator/api/cloud/v1alpha1"
	"github.com/encoder-run/operator/pkg/common"
	"github.com/encoder-run/operator/pkg/graph/converters"
	"github.com/encoder-run/operator/pkg/graph/model"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func List(ctx context.Context) ([]*model.Repository, error) {
	// Get the controller-runtime client from the context.
	ctrlClient, ok := ctx.Value(common.AdminClientKey).(client.Client)
	if !ok {
		return nil, fmt.Errorf("controller-runtime client not found in context")
	}

	// List all the repositories.
	repoList := &v1alpha1.RepositoryList{}
	if err := ctrlClient.List(ctx, repoList, &client.ListOptions{Namespace: "default"}); err != nil {
		return nil, err
	}

	// Convert the repositories to the model.
	repos := make([]*model.Repository, 0, len(repoList.Items))
	for _, repo := range repoList.Items {
		r, err := converters.RepositoryCRDToModel(&repo)
		if err != nil {
			return nil, err
		}
		repos = append(repos, r)
	}
	return repos, nil
}

func Get(ctx context.Context, id string) (*model.Repository, error) {
	// Get the controller-runtime client from the context.
	ctrlClient, ok := ctx.Value(common.AdminClientKey).(client.Client)
	if !ok {
		return nil, fmt.Errorf("controller-runtime client not found in context")
	}

	// Get the repository by name.
	repo := &v1alpha1.Repository{}
	if err := ctrlClient.Get(ctx, client.ObjectKey{Name: id, Namespace: "default"}, repo); err != nil {
		return nil, err
	}

	// Convert the repository to the model.
	return converters.RepositoryCRDToModel(repo)
}

func Add(ctx context.Context, input model.AddRepositoryInput) (*model.Repository, error) {
	// Get the controller-runtime client from the context.
	ctrlClient, ok := ctx.Value(common.AdminClientKey).(client.Client)
	if !ok {
		return nil, fmt.Errorf("controller-runtime client not found in context")
	}

	// If URL is provided, create a new repository with the URL.
	if input.URL != nil {
		return addURL(ctx, ctrlClient, &input)
	}

	// Create a new repository with the type, owner and name.
	return add(ctx, ctrlClient, &input)
}

func Delete(ctx context.Context, id string) (*model.Repository, error) {
	// Get the controller-runtime client from the context.
	ctrlClient, ok := ctx.Value(common.AdminClientKey).(client.Client)
	if !ok {
		return nil, fmt.Errorf("controller-runtime client not found in context")
	}

	// Get the repository by name.
	repo := &v1alpha1.Repository{}
	if err := ctrlClient.Get(ctx, client.ObjectKey{Name: id, Namespace: "default"}, repo); err != nil {
		return nil, err
	}

	// Delete the repository.
	if err := ctrlClient.Delete(ctx, repo); err != nil {
		return nil, err
	}

	// Convert the repository to the model.
	return converters.RepositoryCRDToModel(repo)
}

func add(ctx context.Context, c client.Client, input *model.AddRepositoryInput) (*model.Repository, error) {
	// Check if the owner and name are provided and not empty.
	if input.Owner == nil || *input.Owner == "" {
		return nil, fmt.Errorf("owner cannot be empty")
	}
	if input.Name == nil || *input.Name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if input.Type == nil {
		return nil, fmt.Errorf("type cannot be empty")
	}
	if input.Token == nil {
		return nil, fmt.Errorf("token cannot be empty")
	}
	if input.Branch == nil {
		return nil, fmt.Errorf("branch cannot be empty")
	}
	// Build the url based on the type, owner, and name.
	url := converters.RepositoryURL(*input.Type, *input.Owner, *input.Name)
	if url == "" {
		return nil, fmt.Errorf("failed to build repository URL")
	}

	switch *input.Type {
	case model.RepositoryTypeGithub:
		// Create the repo
		repo := &v1alpha1.Repository{
			ObjectMeta: v1.ObjectMeta{
				GenerateName: "repo-",
				Namespace:    "default",
			},
			Spec: v1alpha1.RepositorySpec{
				Type: v1alpha1.RepositoryTypeGithub,
				Github: &v1alpha1.GithubRepositorySpec{
					URL:   url,
					Owner: *input.Owner,
					Name:  *input.Name,
					Branch: *input.Branch,
				},
			},
		}
		if err := c.Create(ctx, repo); err != nil {
			return nil, err
		}

		// Create the secret with the name of the repository.
		secret := &corev1.Secret{
			ObjectMeta: v1.ObjectMeta{
				Name:      repo.Name,
				Namespace: repo.Namespace,
			},
			Data: map[string][]byte{
				"token": []byte(*input.Token),
			},
		}

		if err := c.Create(ctx, secret); err != nil {
			return nil, err
		}

		return converters.RepositoryCRDToModel(repo)
	default:
		return nil, fmt.Errorf("unsupported repository type: %s", *input.Type)
	}
}

func addURL(ctx context.Context, c client.Client, input *model.AddRepositoryInput) (*model.Repository, error) {
	if *input.URL == "" {
		return nil, fmt.Errorf("URL cannot be empty")
	}
	// Split repository URL to get the type, owner and name.
	repoType, owner, name, err := converters.SplitRepositoryURL(*input.URL)
	if err != nil {
		return nil, err
	}

	// Strip .git from the url
	url := *input.URL
	if url[len(url)-4:] == ".git" {
		url = url[:len(url)-4]
	}

	switch repoType {
	case v1alpha1.RepositoryTypeGithub:
		repo := &v1alpha1.Repository{
			ObjectMeta: v1.ObjectMeta{
				GenerateName: "repo-",
				Namespace:    "default",
			},
			Spec: v1alpha1.RepositorySpec{
				Type: repoType,
				Github: &v1alpha1.GithubRepositorySpec{
					URL:   url,
					Owner: owner,
					Name:  name,
				},
			},
		}
		if err := c.Create(ctx, repo); err != nil {
			return nil, err
		}
		return converters.RepositoryCRDToModel(repo)
	default:
		return nil, fmt.Errorf("unsupported repository type: %s", repoType)
	}
}
