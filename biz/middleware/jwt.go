package middleware

import (
	"LearnShare/biz/model/user"
	"LearnShare/biz/pack"
	"LearnShare/biz/service"
	"LearnShare/pkg/constants"
	"LearnShare/pkg/errno"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hertz-contrib/jwt"

	"github.com/cloudwego/hertz/pkg/app"
)

var (
	IdentityKey               = "userid"
	AccessTokenJwtMiddleware  *jwt.HertzJWTMiddleware
	RefreshTokenJwtMiddleware *jwt.HertzJWTMiddleware
)

func AccessTokenJwt() {
	var err error
	AccessTokenJwtMiddleware, err = jwt.New(&jwt.HertzJWTMiddleware{
		Realm:                       "LS",
		Key:                         []byte("AccessToken_key"),
		Timeout:                     time.Hour,
		MaxRefresh:                  time.Hour,
		WithoutDefaultTokenHeadName: true,
		TokenLookup:                 "header: Access-Token",
		IdentityKey:                 IdentityKey,

		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(int64); ok {
				return jwt.MapClaims{
					AccessTokenJwtMiddleware.IdentityKey: v,
				}
			}
			return jwt.MapClaims{}
		},

		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			return claims[AccessTokenJwtMiddleware.IdentityKey]
		},

		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			pack.BuildFailResponse(c, errno.AuthNoToken)
		},

		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			c.Set("Access-Token", token)
		},

		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			var loginStruct user.LoginInReq
			if err := c.BindAndValidate(&loginStruct); err != nil {
				return nil, err
			}
			users, err := service.NewUserService(ctx, c).LoginIn(&loginStruct)
			if err != nil {
				return nil, err
			}

			c.Set(constants.ContextUid, users.UserId)
			return users.UserId, nil
		},
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
}

func RefreshTokenJwt() {
	var err error
	RefreshTokenJwtMiddleware, err = jwt.New(&jwt.HertzJWTMiddleware{
		Realm:                       "LS",
		Key:                         []byte("refresh_token_key"),
		Timeout:                     time.Hour * 72,
		WithoutDefaultTokenHeadName: true,
		TokenLookup:                 "header: Refresh-Token",
		IdentityKey:                 IdentityKey,

		PayloadFunc: func(data interface{}) jwt.MapClaims {

			return jwt.MapClaims{
				RefreshTokenJwtMiddleware.IdentityKey: data,
			}

		},

		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)

			return claims[RefreshTokenJwtMiddleware.IdentityKey]
		},

		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			pack.BuildFailResponse(c, errno.AuthNoToken)
		},

		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			c.Set("Refresh-Token", token)
		},

		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			uid, exist := c.Get("userid")
			if !exist {
				return nil, err
			}

			return uid, nil

		},
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
}

func GenerateAccessToken(c *app.RequestContext) {
	data := service.GetUidFormContext(c)
	tokenString, _, _ := AccessTokenJwtMiddleware.TokenGenerator(data)
	c.Header("New-Access-Token", tokenString)

}

func IsAccessTokenAvailable(ctx context.Context, c *app.RequestContext) bool {
	claims, err := AccessTokenJwtMiddleware.GetClaimsFromJWT(ctx, c)
	if err != nil {
		return false
	}
	switch v := claims["exp"].(type) {
	case nil:
		return false
	case float64:
		if int64(v) < AccessTokenJwtMiddleware.TimeFunc().Unix() {
			return false
		}
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			return false
		}
		if n < AccessTokenJwtMiddleware.TimeFunc().Unix() {
			return false
		}
	default:
		return false
	}
	c.Set("JWT_PAYLOAD", claims)
	identity := AccessTokenJwtMiddleware.IdentityHandler(ctx, c)
	if identity != nil {
		c.Set(AccessTokenJwtMiddleware.IdentityKey, identity)
	}
	if !AccessTokenJwtMiddleware.Authorizator(identity, ctx, c) {
		return false
	}

	return true

}

func IsRefreshTokenAvailable(ctx context.Context, c *app.RequestContext) bool {

	claims, err := RefreshTokenJwtMiddleware.GetClaimsFromJWT(ctx, c)
	if err != nil {
		return false
	}

	switch v := claims["exp"].(type) {
	case nil:
		return false
	case float64:
		if int64(v) < RefreshTokenJwtMiddleware.TimeFunc().Unix() {
			return false
		}
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			return false
		}
		if n < RefreshTokenJwtMiddleware.TimeFunc().Unix() {
			return false
		}
	default:
		return false
	}

	c.Set("JWT_PAYLOAD", claims)
	identity := RefreshTokenJwtMiddleware.IdentityHandler(ctx, c)
	if identity != nil {
		c.Set(RefreshTokenJwtMiddleware.IdentityKey, identity)
	}
	if !RefreshTokenJwtMiddleware.Authorizator(identity, ctx, c) {
		return false
	}

	return true
}

func Init() {
	AccessTokenJwt()
	RefreshTokenJwt()
	errInit := AccessTokenJwtMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("AccessTokenJwtMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	errInit = RefreshTokenJwtMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Fatal("RefreshTokenJwtMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}
}
