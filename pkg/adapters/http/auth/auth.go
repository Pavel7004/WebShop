package auth

import (
	"crypto/ecdsa"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"

	"github.com/Pavel7004/WebShop/pkg/domain"
)

var (
	ErrNoAuthHeader         = errors.New("Header Authorization not found")
	ErrInvalidSigningMethod = errors.New("JWT Signing method is different from server's")
	ErrInvalidToken         = domain.NewError(
		http.StatusUnauthorized,
		"invalid_token",
		"Received token is invalid",
	)
)

type (
	Auth struct {
		key *ecdsa.PrivateKey
	}

	claims struct {
		Username string `json:"username"`
		jwt.RegisteredClaims
	}
)

func New(key *ecdsa.PrivateKey) *Auth {
	return &Auth{
		key: key,
	}
}

func (a *Auth) Middleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, c.Error(ErrNoAuthHeader))
		return
	}

	if !strings.HasPrefix(tokenString, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusForbidden, c.Error(ErrNoAuthHeader))
		return
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	claims := &claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return a.key.Public(), nil
		},
	)
	log.Info().Msgf("Parsed token: %v", token)

	if err != nil || !token.Valid {
		log.Error().Err(err).Bool("Token valid", token.Valid).Send()
		c.AbortWithStatusJSON(http.StatusUnauthorized, c.Error(ErrInvalidToken))
		return
	}

	c.Set("username", claims.Username)
	c.Next()
}

// Authentication godoc
// @Summary     Login user
// @Description	Login user
// @Tags        Users
// @Accept		json
// @Produce     json
// @Param       req  body  domain.LoginUserRequest	true  "Request to login user"
// @Success      200  {object}  string
// @Failure      400  {object}  domain.Error
// @Failure      404  {object}  domain.Error
// @Failure      500  {object}  domain.Error
// @Router       /login [post]
func (a *Auth) Authentication(c *gin.Context) {
	refreshTokenString := c.GetHeader("X-Refresh-Token")
	if refreshTokenString != "" {
		claims := &claims{}
		token, err := jwt.ParseWithClaims(
			refreshTokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				return a.key, nil
			},
		)

		if err == nil && token.Valid {
			token, err := a.createAccessToken(claims)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}

			c.Header("X-Access-Token", token)
			c.Status(http.StatusOK)
			return
		}
	}

	var req domain.LoginUserRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if req.Name != "admin" || req.Password != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	refreshToken, err := a.createRefreshToken(req.Name)
	if err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &claims{
		Username: req.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	AccessToken := jwt.NewWithClaims(jwt.SigningMethodES512, claims)
	AccessTokenString, err := AccessToken.SignedString(a.key)
	if err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Header("X-Access-Token", AccessTokenString)
	c.Header("X-Refresh-Token", refreshToken)
	c.Status(http.StatusOK)
}

func (a *Auth) createRefreshToken(name string) (string, error) {
	refreshTokenExpirationTime := time.Now().Add(24 * time.Hour)
	refreshClaims := &claims{
		Username: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshTokenExpirationTime),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodES512, refreshClaims)
	return refreshToken.SignedString(a.key)
}

func (a *Auth) createAccessToken(refreshClaims *claims) (string, error) {
	expTime := time.Now().Add(5 * time.Minute)
	claims := &claims{
		Username: refreshClaims.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES512, claims)
	return token.SignedString(a.key)
}
