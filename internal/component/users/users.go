package users

import (
	"context"
	"net/mail"
	"time"

	"casino_loyalty_reward_system/internal/store"
	"casino_loyalty_reward_system/internal/types"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Provider interface {
	Register(ctx context.Context, user types.User) (types.User, string, error)
	Login(ctx context.Context, req types.User) (types.User, string, error)
	Auth(ctx context.Context, token string, path string, method string) (types.User, error)
	GetUsers(ctx context.Context) types.User
	GetUser(ctx context.Context, userID uuid.UUID) types.User
	UpdateUser(ctx context.Context, user types.User) types.User
	DeleteUser(ctx context.Context, userID uuid.UUID) types.User
}

type component struct {
	persistent  store.Persistent
	jwtKey      []byte
	jwtDuration time.Duration
}

// DeleteUser implements Provider.
func (c *component) DeleteUser(ctx context.Context, userID uuid.UUID) types.User {
	panic("unimplemented")
}

// GetUser implements Provider.
func (c *component) GetUser(ctx context.Context, userID uuid.UUID) types.User {
	panic("unimplemented")
}

// GetUsers implements Provider.
func (c *component) GetUsers(ctx context.Context) types.User {
	panic("unimplemented")
}

// UpdateUser implements Provider.
func (c *component) UpdateUser(ctx context.Context, user types.User) types.User {
	panic("unimplemented")
}

var _ Provider = (*component)(nil)

func New(persistent store.Persistent, jwtKey []byte, jwtDuration time.Duration) *component {
	return &component{
		persistent:  persistent,
		jwtKey:      jwtKey,
		jwtDuration: jwtDuration,
	}
}

func (c *component) Register(ctx context.Context, user types.User) (types.User, string, error) {
	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return types.User{}, "", err
	}

	user.Password, err = hashPassword(user.Password)
	if err != nil {
		return types.User{}, "", err
	}

	user.ID = uuid.New()

	now := time.Now()
	user.Created = now
	user.Updated = now

	createdUser, err := c.persistent.UserCreate(ctx, user)
	if err != nil {
		return types.User{}, "", err
	}

	authClaims := types.AuthClaims{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(c.jwtDuration)),
		},
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, authClaims).SignedString(c.jwtKey)

	return createdUser, tokenString, err
}

func (c *component) Login(ctx context.Context, req types.User) (types.User, string, error) {
	user, err := c.persistent.UserGetBy(ctx, types.UserFilter{ByEmail: &req.Email})
	if store.IsErrNotFound(err) {
		return types.User{}, "", types.ErrUnauthorized
	}
	if err != nil {
		return types.User{}, "", err
	}

	match, err := comparePasswords(user.Password, req.Password)
	if err != nil {
		return types.User{}, "", types.ErrUnauthorized
	}

	if !match {
		return types.User{}, "", types.ErrUnauthorized
	}

	authClaims := types.AuthClaims{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(c.jwtDuration)),
		},
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, authClaims).SignedString(c.jwtKey)

	user.Password = ""

	return user, tokenString, err
}

func (c *component) Auth(ctx context.Context, token string, path string, method string) (types.User, error) {
	var authClaims types.AuthClaims
	_, err := jwt.ParseWithClaims(token, &authClaims, func(token *jwt.Token) (interface{}, error) {
		return c.jwtKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
	if err != nil {
		return types.User{}, err
	}

	return types.User{
		ID:    authClaims.ID,
		Email: authClaims.Email,
		Name:  authClaims.Name,
	}, nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func comparePasswords(hashedPassword string, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil, err
}
