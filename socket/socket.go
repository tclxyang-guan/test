package socket

import (
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	log "github.com/sirupsen/logrus"
	"test/middleware"
)

type websocketController struct {
	//注意你也可以使用匿名字段，无所谓，binder会找到它。
	//这是当前的websocket连接，每个客户端都有自己的*websocketController实例。
	Conn websocket.Conn
}

const namespace = "default"

var serverEvents = websocket.Namespaces{
	namespace: websocket.Events{
		websocket.OnNamespaceConnected: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			// with `websocket.GetContext` you can retrieve the Iris' `Context`.
			ctx := websocket.GetContext(nsConn.Conn)

			log.Printf("[%s] connected to namespace [%s] with IP [%s]",
				nsConn, msg.Namespace,
				ctx.RemoteAddr())
			return nil
		},
		websocket.OnNamespaceDisconnect: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			log.Printf("[%s] disconnected from namespace [%s]", nsConn, msg.Namespace)
			return nil
		},
		"chat": func(nsConn *websocket.NSConn, msg websocket.Message) error {
			// room.String() returns -> NSConn.String() returns -> Conn.String() returns -> Conn.ID()
			log.Printf("[%s] sent: %s", nsConn, string(msg.Body))

			// Write message back to the client message owner with:
			// nsConn.Emit("chat", msg)
			// Write message to all except this client with:
			nsConn.Conn.Server().Broadcast(nsConn, msg)
			return nil
		},
	},
}

func WebSocketHandler() *neffos.Server {
	websocketServer := websocket.New(websocket.DefaultGorillaUpgrader, serverEvents)
	websocketServer.OnConnect = func(c *websocket.Conn) error {
		ctx := websocket.GetContext(c)
		if err := middleware.GetJWT().CheckJWT(ctx); err != nil {
			// will send the above error on the client
			// and will not allow it to connect to the websocket server at all.
			return err
		}
		user := ctx.Values().Get("jwt").(*jwt.Token)
		// or just: user := j.Get(ctx)
		log.Printf("This is an authenticated request\n")
		log.Printf("Claim content:")
		log.Printf("%#+v\n", user.Claims)
		log.Printf("[%s] connected to the server", c.ID())
		return nil
	}
	return websocketServer
}
