package graphql_demo

import (
	"context"
	"database/sql"
	"time"

	"github.com/andreylm/graphql-demo/api"
	"github.com/andreylm/graphql-demo/api/dal"
	"github.com/andreylm/graphql-demo/api/errors"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

// Resolver -resolver
type Resolver struct {
	db *sql.DB
}

// NewResolver - creates new resolver
func NewResolver(db *sql.DB) *Resolver {
	return &Resolver{db: db}
}

// Mutation - Mutation
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// Query - Query
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

// Video - Video
func (r *Resolver) Video() VideoResolver {
	return &videoResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateVideo(ctx context.Context, input NewVideo) (*api.Video, error) {
	newVideo := api.Video{
		URL:       input.URL,
		Name:      input.Name,
		CreatedAt: time.Now().UTC(),
	}

	rows, err := dal.LogAndQuery(r.db, "INSERT INTO videos (name, url, user_id, created_at) VALUES($1, $2, $3, $4) RETURNING id",
		input.Name, input.URL, input.UserID, newVideo.CreatedAt)

	defer rows.Close()
	if err != nil || !rows.Next() {
		return nil, err
	}

	if err := rows.Scan(&newVideo.ID); err != nil {
		errors.DebugPrintf(err)
		if errors.IsForeignKeyError(err) {
			return nil, errors.UserNotExist
		}
		return nil, errors.InternalServerError
	}

	return &newVideo, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Videos(ctx context.Context, limit *int, offset *int) ([]*api.Video, error) {
	var videos []*api.Video

	rows, err := dal.LogAndQuery(r.db, "SELECT id, name, url, created_at, user_id FROM videos ORDER BY created_at desc LIMIT $1 offset $2",
		limit, offset)
	defer rows.Close()
	if err != nil {
		errors.DebugPrintf(err)
		return nil, errors.InternalServerError
	}

	for rows.Next() {
		video := api.Video{}
		if err := rows.Scan(&video.ID, &video.Name, &video.URL, &video.CreatedAt, &video.UserID); err != nil {
			errors.DebugPrintf(err)
			return nil, errors.InternalServerError
		}
		videos = append(videos, &video)
	}

	return videos, nil
}

type videoResolver struct{ *Resolver }

func (r *videoResolver) User(ctx context.Context, obj *api.Video) (*api.User, error) {
	panic("not implemented")
}
