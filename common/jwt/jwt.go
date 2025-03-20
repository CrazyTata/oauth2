package jwt

import (
	"context"
	"oauth2/common/redis"
	"oauth2/common/util"
	"oauth2/common/xerr"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type Jwt struct {
	logx.Logger
	Key string
	Ctx context.Context
}

func NewJwt(ctx context.Context, key string) *Jwt {
	return &Jwt{
		Key:    key,
		Ctx:    ctx,
		Logger: logx.WithContext(ctx),
	}
}

func (p *Jwt) GetJwtToken(user string, iat, seconds int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["user_id"] = user
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(p.Key))
}

// ParseToken 解析token
func (p *Jwt) ParseToken(tokenString string) (string, error) {
	InvalidError := util.ReturnError(xerr.InvalidToken)
	if p.CheckBlack(tokenString) {
		p.Logger.Errorf("Jwt CheckBlack: 无效的token token:%s", tokenString)
		return "", InvalidError
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the token signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			p.Logger.Errorf("Jwt  unexpected signing method: %v", token.Header["alg"])
			return nil, InvalidError
		}
		// Return the secret key for verification
		return []byte(p.Key), nil
	})
	if err != nil {
		p.Logger.Errorf("Jwt Parse error : %+v", err)
		return "", InvalidError
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := claims["user"].(string)
		// 将过期时间转化为时间戳
		exp := time.Unix(int64(claims["exp"].(float64)+claims["iat"].(float64)), 0)

		// 判断是否过期
		if !time.Now().Before(exp) {
			p.Logger.Errorf("Jwt Parse error : Token has expired")

			return "", InvalidError
		}

		return user, nil
	}
	p.Logger.Errorf("Jwt token.Claims error : Handle invalid token")
	return "", InvalidError

}

// AddBlack 加入到黑名单
func (p *Jwt) AddBlack(token string) (err error) {
	p.Logger.Infof("Jwt AddBlack token:%s", token)
	redis.Rdb.SAdd(p.Ctx, "logout", token)
	return
}

// CheckBlack 检查token是否在黑名单
func (p *Jwt) CheckBlack(token string) bool {
	isExist, err := redis.Rdb.SIsMember(p.Ctx, "logout", token).Result()
	if err == nil && isExist {
		return true
	}
	return false
}

//
////@author: [piexlmax](https://github.com/piexlmax)
////@function: GetRedisJWT
////@description: 从redis取jwt
////@param: userName string
////@return: redisJWT string, err error
//
//func (jwtService *JwtService) GetRedisJWT(userName string) (redisJWT string, err error) {
//	redisJWT, err = global.GVA_REDIS.Get(context.Background(), userName).Result()
//	return redisJWT, err
//}
//
////@author: [piexlmax](https://github.com/piexlmax)
////@function: SetRedisJWT
////@description: jwt存入redis并设置过期时间
////@param: jwt string, userName string
////@return: err error
//
//func (jwtService *JwtService) SetRedisJWT(jwt string, userName string) (err error) {
//	// 此处过期时间等于jwt过期时间
//	dr, err := utils.ParseDuration(global.GVA_CONFIG.JWT.ExpiresTime)
//	if err != nil {
//		return err
//	}
//	timer := dr
//	err = global.GVA_REDIS.Set(context.Background(), userName, jwt, timer).Err()
//	return err
//}
//
//func LoadAll() {
//	var data []string
//	err := global.GVA_DB.Model(&system.JwtBlacklist{}).Select("jwt").Find(&data).Error
//	if err != nil {
//		global.GVA_LOG.Error("加载数据库jwt黑名单失败!", zap.Error(err))
//		return
//	}
//	for i := 0; i < len(data); i++ {
//		global.BlackCache.SetDefault(data[i], struct{}{})
//	} // jwt黑名单 加入 BlackCache 中
//}
