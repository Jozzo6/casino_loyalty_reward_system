package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/mail"
	"time"

	"casino_loyalty_reward_system/internal/store"
	"casino_loyalty_reward_system/internal/types"

	"github.com/coder/websocket"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Provider interface {
	Register(ctx context.Context, user types.User) (types.User, string, error)
	Login(ctx context.Context, req types.User) (types.User, string, error)
	Auth(ctx context.Context, token string, path string, method string) (types.User, error)
	GetUsers(ctx context.Context) ([]types.User, error)
	GetUser(ctx context.Context, userID uuid.UUID) (types.User, error)
	UpdateUser(ctx context.Context, user types.User) (types.User, error)
	UpdateUserBalance(ctx context.Context, user types.User, value float64, transacrionType types.TransactionType) (types.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	AddPromotion(ctx context.Context, userPromotion types.UserPromotion) (types.UserPromotion, error)
	GetUserPromotions(ctx context.Context, userID uuid.UUID) ([]types.UserPromotion, error)
	GetUserPromotionByID(ctx context.Context, userPromotionID uuid.UUID) (types.UserPromotion, error)
	ClaimPromotion(ctx context.Context, userPromotionID uuid.UUID) error
	DeleteUserPromotion(ctx context.Context, userPromotionID uuid.UUID) error
	ListenToNotifications(ctx context.Context, conn *websocket.Conn, userID uuid.UUID) error
	UserPing(ctx context.Context, conn *websocket.Conn, userID uuid.UUID) error
}

type component struct {
	persistent  store.Persistent
	pubsub      store.PubSub
	jwtKey      []byte
	jwtDuration time.Duration
}

var _ Provider = (*component)(nil)

func New(persistent store.Persistent, pubsub store.PubSub, jwtKey []byte, jwtDuration time.Duration) *component {
	return &component{
		persistent:  persistent,
		jwtKey:      jwtKey,
		jwtDuration: jwtDuration,
		pubsub: pubsub,
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

func (c *component) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return c.persistent.UserDelete(ctx, userID)
}

func (c *component) GetUser(ctx context.Context, userID uuid.UUID) (types.User, error) {
	return c.persistent.UserGetBy(ctx, types.UserFilter{ByID: uuid.NullUUID{UUID: userID, Valid: true}})
}

func (c *component) GetUsers(ctx context.Context) ([]types.User, error) {
	return c.persistent.GetUsers(ctx)
}

func (c *component) UpdateUser(ctx context.Context, user types.User) (types.User, error) {
	return c.persistent.UserUpdate(ctx, user)
}

func (c *component) UpdateUserBalance(ctx context.Context, user types.User, value float64, transacrionType types.TransactionType) (types.User, error) {
	var newBalance float64

	if transacrionType == types.Add {
		newBalance = user.Balance + value
	} else {
		newBalance = user.Balance - value
		if newBalance < 0 {
			return types.User{}, types.ErrInsufficientBalance
		}
	}

	newUser, err := c.persistent.UserBalanceUpdate(ctx, user.ID, newBalance)
	return newUser, err
}

func (c *component) AddPromotion(ctx context.Context, userPromotion types.UserPromotion) (types.UserPromotion, error) {
	userPromotion.ID = uuid.New()

	if userPromotion.StartDate.After(userPromotion.EndDate) {
		return types.UserPromotion{}, types.ErrStartAfterEndDate
	}

	promotion, err := c.persistent.PromotionGetByID(ctx, userPromotion.PromotionID)
	if err != nil {
		return types.UserPromotion{}, err
	}

	if !promotion.IsActive {
		return types.UserPromotion{}, types.ErrPromotionNoLongerActive
	}

	return c.persistent.AddPromotion(ctx, userPromotion)

}

func (c *component) ClaimPromotion(ctx context.Context, userPromotionID uuid.UUID) error {
	userPromotion, err := c.persistent.GetUserPromotionByID(ctx, userPromotionID)
	if err != nil {
		return err
	}

	if userPromotion.Claimed != nil {
		return types.ErrPromotionClaimed
	}

	if !userPromotion.Promotion.IsActive {
		return types.ErrPromotionNoLongerActive
	}

	if time.Now().Before(userPromotion.StartDate) {
		return types.ErrPromotionNotStarted
	}

	if time.Now().After(userPromotion.EndDate) {
		return types.ErrPromotionExpired
	}

	err = c.persistent.ClaimPromotion(ctx, userPromotion.ID)
	if err != nil {
		return err
	}

	user, err := c.persistent.UserGetBy(ctx, types.UserFilter{ByID: uuid.NullUUID{UUID: userPromotion.UserID, Valid: true}})
	if err != nil {
		return err
	}

	_, err = c.UpdateUserBalance(ctx, user, userPromotion.Promotion.Amount, types.Add)
	return err
}

func (c *component) DeleteUserPromotion(ctx context.Context, userPromotionID uuid.UUID) error {
	return c.persistent.DeleteUserPromotion(ctx, userPromotionID)
}

func (c *component) GetUserPromotionByID(ctx context.Context, userPromotionID uuid.UUID) (types.UserPromotion, error) {
	return c.persistent.GetUserPromotionByID(ctx, userPromotionID)
}

func (c *component) GetUserPromotions(ctx context.Context, userPromotionID uuid.UUID) ([]types.UserPromotion, error) {
	return c.persistent.GetUserPromotions(ctx, userPromotionID)
}

func (c *component) ListenToNotifications(ctx context.Context, conn *websocket.Conn,  userID uuid.UUID) error {
	sub := c.pubsub.Subscribe(ctx, fmt.Sprintf("notifications:%s", userID))

	defer sub.Close()

	ch := sub.Channel()

	for msg := range ch {
		err := conn.Write(ctx, websocket.MessageText, []byte(msg.Payload))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *component) UserPing(ctx context.Context, conn *websocket.Conn, userID uuid.UUID) error {
	for {
		msgType, rdr, err := conn.Reader(ctx)
		if err == nil {
			if msgType == websocket.MessageText {
				var msg types.LiveMessage
				err := json.NewDecoder(rdr).Decode(&msg)
				if err != nil {
					conn.Close(websocket.StatusNormalClosure, "")
				}

				if msg.MessageType == "user_ping" {
					c.pubsub.Publish(
						ctx,
						fmt.Sprintf("notifications:%s", userID),
						types.LiveMessage{
							MessageType: "user_ping",
							SentByID:    userID,
						},
					)
				}
			}
		}
	}
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
