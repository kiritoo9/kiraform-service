package middlewares

import (
	"fmt"
	"kiraform/src/infras/configs"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func VerifyToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// get authorization header
		token := c.Request().Header.Get("authorization")
		if token == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing auhtorization header")
		}

		// check if token contain bearer
		hasBearer := strings.HasPrefix(token, "Bearer ")
		if !hasBearer {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing bearer token")
		}
		tokenArr := strings.Split(token, " ")
		if len(tokenArr) < 1 {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing bearer token")
		}
		token = tokenArr[1]

		// check valid token
		jwtSecret := []byte(configs.Environment().SECRET_KEY)
		decode, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}

		if claims, ok := decode.Claims.(jwt.MapClaims); ok && decode.Valid {
			for key, val := range claims {
				if key == "id" {
					// convert id as user id to prevent ambigous naming
					key = "user_id"
				}
				c.Set(key, fmt.Sprintf("%v", val))
			}
		} else {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid access token")
		}

		// allow this request
		return next(c)
	}
}
