package jwt

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	jwtWare "github.com/golang-jwt/jwt/v5"
	"hack/internal/auth/domain"
	"hack/internal/shared"
	"hack/pkg"
	"hack/pkg/secure"
	"strconv"
	"time"
)

type AuthUserParams struct {
	UserId         int64
	HashedPassword string
	PlainPassword  string
	RequestId      int64
	Jwt            *Config
}

func AccessTokenCookie() *fiber.Cookie {
	return &fiber.Cookie{
		Name:     shared.AccessTokenCookie,
		Path:     "/",
		MaxAge:   int(shared.JwtTokenExp.Seconds()),
		Expires:  time.Now().Add(shared.JwtTokenExp),
		Secure:   true,
		HTTPOnly: true,
	}
}

func CreateAccessToken(userId, requestId int64, cfg *Config) (string, error) {
	key := []byte(cfg.Secret)
	encryptedId, err := secure.Encrypt(strconv.FormatInt(userId, 10), cfg.EncryptionKey)
	if err != nil {
		return "", err
	}

	ts := pkg.GetLocalTime().Unix()
	hashedRequestId := secure.GenerateHash(fmt.Sprintf("%d%d", ts, requestId))

	claims := jwtWare.MapClaims{
		shared.UserIdField:    encryptedId,
		shared.TimestampField: ts,
		shared.RequestIdField: hashedRequestId,
		shared.ExpTokenField:  pkg.GetLocalTime().Add(shared.JwtTokenExp).Unix(),
	}
	accessToken := jwtWare.NewWithClaims(jwtWare.SigningMethodHS256, claims)
	return accessToken.SignedString(key)
}

func VerifyAccessToken(accessTokenString string, secret string) (*jwtWare.Token, error) {
	accessToken, err := jwtWare.Parse(accessTokenString, func(token *jwtWare.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtWare.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	return accessToken, err
}

func AuthenticateUser(params AuthUserParams) (domain.AuthResp, error) {
	var resp domain.AuthResp
	if !secure.VerifyPassword(params.HashedPassword, params.PlainPassword) {
		return resp, shared.ErrInvalidPassword
	}

	accessToken, err := CreateAccessToken(params.UserId, params.RequestId, params.Jwt)
	if err != nil {
		return resp, err
	}

	resp.AccessToken = accessToken
	resp.Type = shared.TokenTypeBearer
	return resp, nil
}
