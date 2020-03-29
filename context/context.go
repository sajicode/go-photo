package context

import (
	"context"

	"github.com/sajicode/go-photo/models"
)

// * we do not want another app to overwrite our user key in context
// * context stores both the key and key type
const (
	userKey privateKey = "user"
)

type privateKey string

// WithUser sets a user on context
func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// User returns the User data stored in context
func User(ctx context.Context) *models.User {
	if temp := ctx.Value(userKey); temp != nil {
		if user, ok := temp.(*models.User); ok {
			return user
		}
	}
	return nil
}
