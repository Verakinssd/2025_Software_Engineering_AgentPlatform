package middleware

import (
	"backend/config"
	"backend/exceptions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

var secretKey = []byte(config.Config.Jwt.Secret)

func GenerateToken(userId uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"role":   role,
		"exp":    time.Now().Add(time.Second * time.Duration(config.Config.Jwt.Token.Expiration)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func Auth(roleRequired []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie(config.Config.Jwt.Token.Name)
		if tokenString == "" || err != nil {
			c.JSON(http.StatusUnauthorized, &exceptions.ErrTokenMissing)
			c.Abort()
			return
		}

		// 解析 Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, &exceptions.ErrTokenInvalid)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, &exceptions.ErrTokenClaimsInvalid)
			c.Abort()
			return
		}

		// 续签 Token
		if time.Now().Add(time.Second*time.Duration(config.Config.Jwt.Token.RefreshBeforeExpiration)).Unix() > int64(claims["exp"].(float64)) {
			newToken, err := GenerateToken(uint(claims["userId"].(float64)), claims["role"].(string))
			if err != nil {
				c.JSON(http.StatusUnauthorized, &exceptions.ErrTokenGenerate)
				c.Abort()
				return
			}
			c.SetCookie(config.Config.Jwt.Token.Name, newToken, int(config.Config.Jwt.Token.Expiration), "/", "", false, true)
		}

		// 读取用户信息
		userId := uint(claims["userId"].(float64))
		role := claims["role"].(string)
		c.Set("userId", userId)
		c.Set("role", role)

		// 角色鉴权
		for _, r := range roleRequired {
			if role == r {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, &exceptions.ErrPermissionDenied)
		c.Abort()
		return
	}
}

// OperationAuth checks operations like commenting under workflows, publishing new agents, etc.
func OperationAuth(operation Operation) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, &exceptions.ErrTokenMissing)
			c.Abort()
			return
		}
		err := operationMap[operation](userId.(uint))
		if err != nil {
			c.JSON(http.StatusForbidden, err)
			c.Abort()
			return
		}
		c.Next()
		return
	}
}
