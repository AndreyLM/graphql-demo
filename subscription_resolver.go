package graphql_demo

import (
	"context"

	"github.com/andreylm/graphql-demo/api"
	"github.com/markbates/going/randx"
)

var (
	videoPublishedChannel map[string]chan *api.Video
)

func init() {
	videoPublishedChannel = map[string]chan *api.Video{}
}

// Subscription - subscriptions
func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) VideoPublished(ctx context.Context) (<-chan *api.Video, error) {
	id := randx.String(8)

	videoEvent := make(chan *api.Video, 1)
	go func() {
		<-ctx.Done()
	}()
	videoPublishedChannel[id] = videoEvent
	return videoEvent, nil
}
