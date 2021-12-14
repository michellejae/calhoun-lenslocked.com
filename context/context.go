package context

import (
	"context"

	"gitlab.com/michellejae/lenslocked.com/models"
)

const (
	userKey privateKey = "user"
)

// we create type privateKey so just in case another packate has type user passed into context, it will not be correct
// since that type of user will not be of type privateKey. touched on this in on the justforfunc golang YT channel
type privateKey string

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	if temp := ctx.Value(userKey); temp != nil {
		if user := temp.(*models.User); user != nil {
			return user
		}
	}
	return nil
}
