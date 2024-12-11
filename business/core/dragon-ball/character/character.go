package character

import (
	"context"
	dragonball "dragonball/business/core/dragon-ball"
	v1 "dragonball/business/web/v1"
	"errors"
	"fmt"
)

func Run(ctx context.Context, db *dragonball.Core, name string) ([]Response, error) {
	URL := fmt.Sprintf("/characters?name=%s", name)
	service := "Character - Query - POST"

	fn := db.OnAfterResponse(ctx, service)
	fnErr := db.OnError(ctx, service)

	res, err := db.Client.Clone().OnAfterResponse(fn).OnError(fnErr).R().SetResult(&[]Response{}).Get(URL)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		err = errors.New("Service " + service + " returned a status " + res.Status())
		return nil, v1.NewRequestError(err, res.StatusCode())
	}

	result, ok := res.Result().(*[]Response)
	if !ok {
		return nil, errors.New("result is not ceclmvc.response type")
	}

	return *result, nil
}
