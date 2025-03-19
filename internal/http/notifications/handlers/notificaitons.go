package handlers

import (
	"net/http"

	notifications "github.com/Jozzo6/casino_loyalty_reward_system/internal/component/notificaitons"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/types"

	"github.com/coder/websocket"
)

type notificationsRouter struct {
	component notifications.NotificationProvider
}

func NewNotificationsRouter(component notifications.NotificationProvider) *notificationsRouter {
	return &notificationsRouter{component: component}
}

// ListenToNotifications establishes a WebSocket connection for real-time notifications.
// @Summary Listen to notifications
// @Description Establishes a WebSocket connection to receive real-time notifications for the authenticated user.
// @Tags Notifications
// @Produce json
// @Success 101 {string} string "Switching Protocols - WebSocket connection established"
// @Failure 401 {object} string "Unauthorized - User not authenticated"
// @Failure 500 {object} string "Internal server error"
// @Router /api/v1/notifications [get]
func (nr *notificationsRouter) ListenToNotifications(w http.ResponseWriter, r *http.Request) {
	log := types.GetLoggerFromContext(r.Context())

	user, err := types.GetAccountFromContext(r.Context())
	if err != nil {
		log.Errorf("failed to get user from context: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	conn.CloseRead(r.Context())

	defer conn.CloseNow()

	err = nr.component.ListenToNotifications(r.Context(), conn, user.ID)
	if err != nil {
		log.Errorf("failed to listen to notificaitons: %s", err)
	}

}
