package dataloaders

import (
	"database/sql"
	"log"
	"time"

	"github.com/andreylm/graphql-demo/api"
	"github.com/andreylm/graphql-demo/api/dal"
	"github.com/andreylm/graphql-demo/api/errors"
	"github.com/jmoiron/sqlx"
)

func newVideoLoader(db *sql.DB, wait time.Duration, maxBatch int) *VideoLoader {
	return &VideoLoader{
		wait:     wait,
		maxBatch: maxBatch,
		fetch: func(ids []int) ([][]*api.Video, []error) {
			sqlQuery := getVideosQuery(ids)
			sqlQuery, arguments, err := sqlx.In(sqlQuery, ids)
			if err != nil {
				log.Println(err)
			}
			sqlQuery = sqlx.Rebind(sqlx.DOLLAR, sqlQuery)
			rows, err := dal.LogAndQuery(db, sqlQuery, arguments...)
			defer rows.Close()
			if err != nil {
				log.Println(err)
			}
			allVideos := []*api.Video{}
			for rows.Next() {
				video := api.Video{}
				if err := rows.Scan(&video.ID, &video.Name, &video.Description, &video.UserID, &video.URL, &video.CreatedAt); err != nil {
					errors.DebugPrintf(err)
					return nil, []error{errors.InternalServerError}
				}
				allVideos = append(allVideos, &video)
			}

			videos := make([][]*api.Video, len(ids))
			for i, id := range ids {
				for _, v := range allVideos {
					if v.UserID == id {
						videos[i] = append(videos[i], v)
					}
				}
			}

			return videos, nil
		},
	}
}

func getVideosQuery(ids []int) string {
	if len(ids) == 1 {
		return "SELECT id, name, description, user_id, url, created_at FROM videos WHERE user_id = ?"
	}
	return "SELECT id, name, description, user_id, url, created_at FROM videos WHERE user_id IN (?)"
}
