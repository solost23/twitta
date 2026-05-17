package routers

import (
	"twitta/pkg/utils"
	"twitta/pkg/ws"
	"twitta/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// notifyWS 建立全局通知 WebSocket，用户登录后保持连接，实时接收新消息通知。
// 连接建立后立即推送离线期间未读的消息。
func notifyWS(c *gin.Context) {
	user := utils.GetUser(c)
	userID := user.ID.Hex()

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := ws.NewNotifyClient(ws.DefaultNotify, conn, userID)
	ws.DefaultNotify.Register(client)

	// 推送离线消息
	go func() {
		offline, err := services.GetAndMarkOfflineMessages(userID)
		if err != nil {
			zap.S().Warnf("notify: get offline messages error: %v", err)
			return
		}
		for _, n := range offline {
			client.SendNotification(n)
		}
	}()

	client.WritePump()
}
