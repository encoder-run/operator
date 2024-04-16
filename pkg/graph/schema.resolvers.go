package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"

	"github.com/encoder-run/operator/pkg/graph/model"
	"github.com/encoder-run/operator/pkg/graph/resolvers/models"
	"github.com/encoder-run/operator/pkg/graph/resolvers/repositories"
	"github.com/encoder-run/operator/pkg/graph/resolvers/storage"
)

// AddModel is the resolver for the addModel field.
func (r *mutationResolver) AddModel(ctx context.Context, input model.AddModelInput) (*model.Model, error) {
	return models.Add(ctx, input)
}

// DeployModel is the resolver for the deployModel field.
func (r *mutationResolver) DeployModel(ctx context.Context, input model.DeployModelInput) (*model.Model, error) {
	panic(fmt.Errorf("not implemented: DeployModel - deployModel"))
}

// DeleteModel is the resolver for the deleteModel field.
func (r *mutationResolver) DeleteModel(ctx context.Context, id string) (*model.Model, error) {
	return models.Delete(ctx, id)
}

// AddRepository is the resolver for the addRepository field.
func (r *mutationResolver) AddRepository(ctx context.Context, input model.AddRepositoryInput) (*model.Repository, error) {
	return repositories.Add(ctx, input)
}

// DeleteRepository is the resolver for the deleteRepository field.
func (r *mutationResolver) DeleteRepository(ctx context.Context, id string) (*model.Repository, error) {
	return repositories.Delete(ctx, id)
}

// AddStorage is the resolver for the addStorage field.
func (r *mutationResolver) AddStorage(ctx context.Context, input model.AddStorageInput) (*model.Storage, error) {
	return storage.Add(ctx, input)
}

// DeleteStorage is the resolver for the deleteStorage field.
func (r *mutationResolver) DeleteStorage(ctx context.Context, id string) (*model.Storage, error) {
	return storage.Delete(ctx, id)
}

// Models is the resolver for the models field.
func (r *queryResolver) Models(ctx context.Context) ([]*model.Model, error) {
	return models.List(ctx)
}

// GetModel is the resolver for the getModel field.
func (r *queryResolver) GetModel(ctx context.Context, id string) (*model.Model, error) {
	return models.Get(ctx, id)
}

// Repositories is the resolver for the repositories field.
func (r *queryResolver) Repositories(ctx context.Context) ([]*model.Repository, error) {
	return repositories.List(ctx)
}

// GetRepository is the resolver for the getRepository field.
func (r *queryResolver) GetRepository(ctx context.Context, id string) (*model.Repository, error) {
	return repositories.Get(ctx, id)
}

// Storages is the resolver for the storages field.
func (r *queryResolver) Storages(ctx context.Context) ([]*model.Storage, error) {
	return storage.List(ctx)
}

// GetStorage is the resolver for the getStorage field.
func (r *queryResolver) GetStorage(ctx context.Context, id string) (*model.Storage, error) {
	return storage.Get(ctx, id)
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
