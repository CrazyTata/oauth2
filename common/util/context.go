package util

import (
	"context"
	"strconv"

	"oauth2/common/xerr"
)

// GetUserIdFromContext 从上下文中提取 user_id
func GetUserIdFromContext(ctx context.Context) (int64, error) {
	userId, err := GetUserStrFromContext(ctx)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(userId, 10, 64)
}

// GetUserStrFromContext 从上下文中提取 user_id
func GetUserStrFromContext(ctx context.Context) (string, error) {
	userId, ok := ctx.Value("user_id").(string)
	if !ok {
		return "", xerr.NewErrCode(xerr.LoginMiss)
	}
	return userId, nil
}
