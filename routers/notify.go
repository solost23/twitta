package routers

import (
	"twitta/pkg/utils"
	"twitta/pkg/ws"

	"github.com/gin-gonic/gin"
)

// notifyWS 建立全局通知 WebSocket，用户登录后保持连接，实时接收新消息通知。
func notifyWS(c *gin.Context) {
	user := utils.GetUser(c)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := ws.NewNotifyClient(ws.DefaultNotify, conn, user.ID.Hex())
	ws.DefaultNotify.Register(client)
	client.WritePump()
}
