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

func newUserLoader(db *sql.DB, wait time.Duration, maxBatch int) *UserLoader {
	return &UserLoader{
		wait:     wait,
		maxBatch: maxBatch,
		fetch: func(ids []int) ([]*api.User, []error) {
			sqlQuery := getUsersQuery(ids)
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
			userByID := map[int]*api.User{}
			for rows.Next() {
				user := api.User{}
				if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
					errors.DebugPrintf(err)
					return nil, []error{errors.InternalServerError}
				}
				userByID[user.ID] = &user
			}
			users := make([]*api.User, len(ids))
			for i, id := range ids {
				users[i] = userByID[id]
				i++
			}
			return users, nil
		},
	}
}

func getUsersQuery(ids []int) string {
	if len(ids) == 1 {
		return "SELECT id, name, email FROM users WHERE id = ?"
	}
	return "SELECT id, name, email FROM users WHERE id IN (?)"
}
