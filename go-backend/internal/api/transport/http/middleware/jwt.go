package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	jwtWare "github.com/golang-jwt/jwt/v5"
	"hack/internal/api"
	"hack/internal/auth"
	"hack/internal/shared"
	"hack/pkg/jwt"
	"hack/pkg/secure"
	"strconv"
)

type apiMiddleware struct {
	tokenRepo auth.TokenRepo
	cfg       *jwt.Config
}

func New(tokenRepo auth.TokenRepo, cfg *jwt.Config) api.Middleware {
	return &apiMiddleware{
		tokenRepo: tokenRepo,
		cfg:       cfg,
	}
}

func (m *apiMiddleware) Jwt() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		accessTokenString := ctx.Cookies(shared.AccessTokenCookie)
		if accessTokenString == "" {
			return shared.ErrJwtMissing
		}

		token, err := jwt.VerifyAccessToken(accessTokenString, m.cfg.Secret)
		if err != nil {
			return shared.ErrJwtInvalid
		}

		claims, ok := token.Claims.(jwtWare.MapClaims)
		if !ok {
			return shared.ErrJwtMalformed
		}

		encryptedId, ok := claims[shared.UserIdField]
		if !ok {
			return shared.ErrJwtMalformed
		}

		encryptedIdString, ok := encryptedId.(string)
		if !ok {
			return shared.ErrJwtMalformed
		}

		userIdString, err := secure.Decrypt(encryptedIdString, m.cfg.EncryptionKey)
		if err != nil {
			return shared.ErrJwtMalformed
		}

		userId, err := strconv.ParseInt(userIdString, 10, 64)
		if err != nil {
			return shared.ErrJwtMalformed
		}

		gotRequestId, ok := claims[shared.RequestIdField]
		if !ok {
			return shared.ErrJwtMalformed
		}

		requestId, err := m.tokenRepo.GetIdRequest(ctx.Context(), userId)
		if err != nil {
			return shared.ErrFindRecord
		}

		tsRaw, ok := claims[shared.TimestampField]
		if !ok {
			return shared.ErrJwtMalformed
		}

		ts, ok := tsRaw.(float64)
		if !ok {
			return shared.ErrJwtMalformed
		}

		hashedRequestId := secure.GenerateHash(fmt.Sprintf("%d%d", int64(ts), requestId))
		if hashedRequestId != gotRequestId {
			return shared.ErrJwtInvalid
		}

		if err = m.tokenRepo.SetIdRequest(ctx.Context(), userId, requestId+1); err != nil {
			return shared.ErrInsertRecord
		}

		accessToken, err := jwt.CreateAccessToken(userId, requestId+1, m.cfg)
		if err != nil {
			return shared.ErrJwtCreate
		}

		cookie := jwt.AccessTokenCookie()
		cookie.Value = accessToken
		ctx.Cookie(cookie)

		ctx.Locals(shared.UserIdField, userId)
		return ctx.Next()
	}
}
