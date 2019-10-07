package graphql_demo

import (
	"context"
	"time"

	"github.com/andreylm/graphql-demo/api"
	"github.com/andreylm/graphql-demo/api/dal"
	"github.com/andreylm/graphql-demo/api/dataloaders"
	"github.com/andreylm/graphql-demo/api/errors"
	"github.com/andreylm/graphql-demo/api/utils"
)

type videoResolver struct{ *Resolver }

func (r *videoResolver) User(ctx context.Context, obj *api.Video) (*api.User, error) {
	return ctx.Value(dataloaders.CtxKeyUser).(*dataloaders.UserLoader).Load(obj.UserID)
}

func (r *videoResolver) Screenshots(ctx context.Context, obj *api.Video) ([]*api.Screenshot, error) {
	return nil, nil
}

func (r *videoResolver) Related(ctx context.Context, obj *api.Video, limit *int, offset *int) ([]*api.Video, error) {
	var videos []*api.Video

	rows, err := dal.LogAndQuery(
		r.db,
		"SELECT id, name, url, created_at, user_id FROM videos RIGHT JOIN related_videos ON first_id = id OR second_id = id WHERE ( first_id = $1 OR second_id = $1 ) AND id != $1 ORDER BY created_at desc LIMIT $2 offset $3",
		obj.ID, utils.GetInteger(limit, 10), utils.GetInteger(offset, 0))
	if err != nil {
		errors.DebugPrintf(err)
		return nil, errors.InternalServerError
	}
	defer rows.Close()

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

func (r *mutationResolver) CreateVideo(ctx context.Context, input NewVideo) (*api.Video, error) {
	newVideo := api.Video{
		URL:         input.URL,
		Name:        input.Name,
		Description: input.Description,
		CreatedAt:   time.Now().UTC(),
		UserID:      input.UserID,
	}

	rows, err := dal.LogAndQuery(r.db, "INSERT INTO videos (name, url, description, user_id, created_at) VALUES($1, $2, $3, $4, $5) RETURNING id",
		input.Name, input.URL, input.Description, input.UserID, newVideo.CreatedAt)

	if err != nil {
		errors.DebugPrintf(err)
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	if err := rows.Scan(&newVideo.ID); err != nil {
		errors.DebugPrintf(err)
		if errors.IsForeignKeyError(err) {
			return nil, errors.UserNotExist
		}
		return nil, errors.InternalServerError
	}

	for _, observer := range videoPublishedChannel {
		observer <- &newVideo
	}

	return &newVideo, nil
}

func (r *mutationResolver) RemoveVideo(ctx context.Context, input int) (bool, error) {
	_, err := dal.LogAndQuery(
		r.db,
		"DELETE FROM videos WHERE id = $1",
		input,
	)
	if err != nil {
		errors.DebugPrintf(err)
		return false, err
	}

	return true, nil
}

func (r *queryResolver) Videos(ctx context.Context, limit *int, offset *int) ([]*api.Video, error) {
	var videos []*api.Video

	rows, err := dal.LogAndQuery(r.db, "SELECT id, name, url, created_at, user_id FROM videos ORDER BY created_at desc LIMIT $1 offset $2",
		utils.GetInteger(limit, 10), utils.GetInteger(offset, 0))
	if err != nil {
		errors.DebugPrintf(err)
		return nil, errors.InternalServerError
	}
	defer rows.Close()

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
