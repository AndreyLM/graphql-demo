package dataloaders

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)

type ctxKeyType struct{ name string }

// CtxKeyUser - dataloader middleware context key for user
var CtxKeyUser = ctxKeyType{name: "dataloader_user"}

// CtxKeyVideo - dataloader middleware context key for video
var CtxKeyVideo = ctxKeyType{name: "dataloader_video"}

// DataloaderMiddleware - dataloader middleware
func DataloaderMiddleware(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userLoader := newUserLoader(db, 1*time.Millisecond, 100)
		videoLoader := newVideoLoader(db, 1*time.Millisecond, 100)

		ctx := context.WithValue(r.Context(), CtxKeyUser, userLoader)
		ctx = context.WithValue(ctx, CtxKeyVideo, videoLoader)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
