package graphql_demo

import (
	"context"

	"github.com/andreylm/graphql-demo/api/dal"
	"github.com/andreylm/graphql-demo/api/errors"
)

func (r *mutationResolver) AddVideoRelation(ctx context.Context, input VideoRelation) (bool, error) {
	_, err := dal.LogAndQuery(
		r.db,
		"INSERT INTO related_videos(first_id, second_id) VALUES($1, $2)",
		input.FirstID, input.SecondID)
	if err != nil {
		errors.DebugPrintf(err)
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) RemoveVideoRelation(ctx context.Context, input VideoRelation) (bool, error) {
	_, err := dal.LogAndQuery(
		r.db,
		"DELETE FROM related_videos WHERE ( first_id = $1 AND second_id = $2 ) OR ( first_id = $2 AND second_id = $1 )",
		input.FirstID, input.SecondID)
	if err != nil {
		errors.DebugPrintf(err)
		return false, err
	}

	return true, nil
}
