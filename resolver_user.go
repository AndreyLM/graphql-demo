package graphql_demo

import (
	"context"

	"github.com/andreylm/graphql-demo/api"
	"github.com/andreylm/graphql-demo/api/dal"
	"github.com/andreylm/graphql-demo/api/dataloaders"
	"github.com/andreylm/graphql-demo/api/utils"

	// "github.com/andreylm/graphql-demo/api/dal"
	"github.com/andreylm/graphql-demo/api/errors"
)

type userResolver struct{ *Resolver }

func (r *queryResolver) Users(ctx context.Context, limit *int, offset *int) ([]*api.User, error) {
	var users []*api.User

	rows, err := dal.LogAndQuery(r.db, "SELECT id, name, email FROM users ORDER BY id LIMIT $1 offset $2",
		utils.GetInteger(limit, 10), utils.GetInteger(offset, 0))
	if err != nil {
		errors.DebugPrintf(err)
		return nil, errors.InternalServerError
	}
	defer rows.Close()

	for rows.Next() {
		user := api.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			errors.DebugPrintf(err)
			return nil, errors.InternalServerError
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input NewUser) (*api.User, error) {
	newUser := &api.User{
		Name:  input.Name,
		Email: input.Email,
	}

	rows, err := dal.LogAndQuery(
		r.db,
		"INSERT INTO users(name, email) VALUES($1, $2) RETURNING id",
		input.Name, input.Email)
	if err != nil {
		errors.DebugPrintf(err)
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}

	if err := rows.Scan(&newUser.ID); err != nil {
		errors.DebugPrintf(err)
		return nil, errors.InternalServerError
	}

	return newUser, nil
}

func (r *userResolver) Videos(ctx context.Context, obj *api.User, limit *int, offset *int) ([]*api.Video, error) {
	return ctx.Value(dataloaders.CtxKeyVideo).(*dataloaders.VideoLoader).Load(obj.ID)

	// var videos []*api.Video

	// rows, err := dal.LogAndQuery(r.db, "SELECT id, name, description, url FROM videos WHERE user_id =  $1 LIMIT $2 OFFSET $3",
	// 	obj.ID, utils.GetInteger(limit, 10), utils.GetInteger(offset, 0))
	// if err != nil {
	// 	errors.DebugPrintf(err)
	// 	return nil, errors.InternalServerError
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	video := api.Video{UserID: obj.ID}
	// 	if err := rows.Scan(&video.ID, &video.Name, &video.Description, &video.URL); err != nil {
	// 		errors.DebugPrintf(err)
	// 		return nil, errors.InternalServerError
	// 	}

	// 	videos = append(videos, &video)
	// }

	// return videos, nil
}
