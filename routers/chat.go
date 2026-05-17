package routers

import (
	"net/http"
	"twitta/pkg/response"
	"twitta/pkg/utils"
	"twitta/pkg/ws"
	"twitta/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//	@Summary	chatList（历史消息）
//	@Tags		chat
//	@Produce	json
//	@Success	200		{object}	forms.ChatList
//	@Failure	400		{object}	response.Response
//	@Param		id		path		string	true	"对方用户ID"
//	@Param		token	header		string	true	"token"
//	@Router		/chats/{id} [get]
func chatList(c *gin.Context) {
	UIdForm := &utils.UIdForm{}
	if err := utils.GetValidUriParams(c, UIdForm); err != nil {
		response.Error(c, 2001, err)
		return
	}
	params := &utils.PageForm{}
	if err := utils.DefaultGetValidParams(c, params); err != nil {
		response.Error(c, 2001, err)
		return
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Size <= 0 {
		params.Size = 10
	}

	result, err := services.NewService().ChatList(c, UIdForm.Id, params)
	if err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.Success(c, result)
}

//	@Summary	chatWS（WebSocket 实时聊天）
//	@Tags		chat
//	@Param		id		path	string	true	"对方用户ID"
//	@Param		token	header	string	true	"token"
//	@Router		/chats/{id}/ws [get]
func chatWS(c *gin.Context) {
	UIdForm := &utils.UIdForm{}
	if err := utils.GetValidUriParams(c, UIdForm); err != nil {
		response.Error(c, 2001, err)
		return
	}

	user := utils.GetUser(c)
	userID := user.ID.Hex()
	roomID := services.RoomID(userID, UIdForm.Id)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := ws.NewClient(ws.Default, conn, roomID, userID)
	ws.Default.Register(client)

	go client.WritePump()
	client.ReadPump(func(msg *ws.Message) {
		go services.SaveChatMessage(msg)
	})
}

// chatRead 将对方发给我的消息标记为已读
func chatRead(c *gin.Context) {
	UIdForm := &utils.UIdForm{}
	if err := utils.GetValidUriParams(c, UIdForm); err != nil {
		response.Error(c, 2001, err)
		return
	}
	user := utils.GetUser(c)
	if err := services.MarkMessagesRead(c, UIdForm.Id, user.ID.Hex()); err != nil {
		response.Error(c, 2001, err)
		return
	}
	response.Success(c, nil)
}

