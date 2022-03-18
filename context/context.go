package context

import (
	"context"

	"github.com/google/uuid"
)


// GetContext 从pool中获取一个干净的Context，合并context.Context
func GetContext(c context.Context, UUID string) context.Context {
	if UUID != "" {
		c = context.WithValue(c, "uuid", UUID) //nolint
	} else {
		id, _ := uuid.NewUUID()
		c = context.WithValue(c, "uuid", id.String()) //nolint
	}
	return c
}
