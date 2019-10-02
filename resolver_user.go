package graphql_demo

import (
	"context"

	"github.com/andreylm/graphql-demo/api"
	// "github.com/andreylm/graphql-demo/api/dal"
	// "github.com/andreylm/graphql-demo/api/errors"
)

type userResolver struct{ *Resolver }

func (r *queryResolver) Users(ctx context.Context, limit *int, offset *int) ([]*api.User, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateUser(ctx context.Context, input NewUser) (*api.User, error) {
	panic("Not implemented yet")
}

func (r *userResolver) Videos(ctx context.Context, obj *api.User, limit *int, offset *int) ([]*api.Video, error) {
	panic("not implemented")
}
