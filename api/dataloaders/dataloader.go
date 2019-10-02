package dataloaders

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/andreylm/graphql-demo/api"
	"github.com/andreylm/graphql-demo/api/dal"
	"github.com/andreylm/graphql-demo/api/errors"
	"github.com/jmoiron/sqlx"
)

type ctxKeyType struct{ name string }

// CtxKey - dataloader middleware context key
var CtxKey = ctxKeyType{name: "dataloaderctx"}

// DataloaderMiddleware - dataloader middleware
func DataloaderMiddleware(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userLoader := UserLoader{
			wait:     1 * time.Millisecond,
			maxBatch: 100,
			fetch: func(ids []int) ([]*api.User, []error) {
				sqlQuery := getQuery(ids)
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

		ctx := context.WithValue(r.Context(), CtxKey, &userLoader)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getQuery(ids []int) string {
	if len(ids) == 1 {
		return "SELECT id, name, email FROM users WHERE id = ?"
	}
	return "SELECT id, name, email FROM users WHERE id IN (?)"
}
