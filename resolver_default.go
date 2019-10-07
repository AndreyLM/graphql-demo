package graphql_demo

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/andreylm/graphql-demo/api/errors"
)

type contextKey string

var (
	// UserIDCtxKey - user context id key
	UserIDCtxKey contextKey = "userID"
	// UserRoleCtxKey - user context role key
	UserRoleCtxKey contextKey = "userRole"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

// Resolver -resolver
type Resolver struct {
	db *sql.DB
}

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }

// NewRootResolvers - creates root resolvers
func NewRootResolvers(db *sql.DB) Config {
	c := Config{
		Resolvers: &Resolver{
			db: db,
		},
	}

	c.Directives.IsAuthenticated = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		ctxUserID := ctx.Value(UserIDCtxKey)
		if ctxUserID != nil {
			return next(ctx)
		}
		return nil, errors.UnauthorisedError
	}

	c.Directives.HasRole = func(ctx context.Context, obj interface{}, next graphql.Resolver, role Role) (res interface{}, err error) {
		ctxUserID := ctx.Value(UserIDCtxKey).(string)
		if ctxUserID == "" {
			return nil, errors.ForbiddenError
		}

		userID, err := strconv.Atoi(ctxUserID)
		if err != nil {
			return nil, errors.ForbiddenError
		}

		if !hasRole(db, userID, role.String()) {
			return nil, errors.ForbiddenError
		}

		return next(ctx)
	}

	countCompexity := func(childComplexity int, limit *int, offset *int) int {
		return *limit * childComplexity
	}

	c.Complexity.Query.Videos = countCompexity
	c.Complexity.Video.Related = countCompexity

	return c
}

// NewResolver - creates new resolver
func NewResolver(db *sql.DB) *Resolver {
	return &Resolver{db: db}
}

// Mutation - mutation resolver
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// Query - query resolver
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

// Video - video resolver
func (r *Resolver) Video() VideoResolver {
	return &videoResolver{r}
}

// User - user resolver
func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}
