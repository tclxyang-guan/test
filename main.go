package main

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/websocket"
	"github.com/onrik/logrus/filename"
	"strings"
	"test/config"
	"test/datasource"
	"test/middleware"
	"test/socket"
	"test/web/route"
	"time"
	//"github.com/iris-contrib/middleware/cors"
	//"github.com/onrik/logrus/sentry"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetLevel(log.InfoLevel)
	//aa := log.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"} // 设置时间输出格式
	aa := log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"}
	log.SetFormatter(&aa)
}

func LoginNew3() *os.File {
	filenameHook := filename.NewHook()
	filenameHook.Field = "line" // Customize source field name
	log.AddHook(filenameHook)

	/*sentryHook,_ := sentry.NewHook(nil, log.PanicLevel, log.FatalLevel, log.ErrorLevel)

	log.AddHook(sentryHook)*/
	log.SetOutput(os.Stdout)
	filename := config.Sysconfig.Basic.ErrorLogPath + "error_" + time.Now().Format("20060102") + ".txt"
	//打开一个输出文件，如果重新启动服务器，它将追加到今天的文件中
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		if !config.Sysconfig.Basic.IsLocal {
			log.SetOutput(file)
		}
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	return file
}

func main() {
	file := LoginNew3()
	defer file.Close()
	app := iris.New().Configure(iris.WithConfiguration(iris.Configuration{DisableBodyConsumptionOnUnmarshal: true}))
	r, close := newRequestLogger()
	defer close()
	app.Use(r)
	app.Logger().SetOutput(datasource.NewLogFile())
	//静态文件路径
	app.HandleDir("/static", config.Sysconfig.Basic.StaticPath)
	//跨域资源访问
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // 允许通过的主机名称 *，表示接受任意域名的请求
		AllowCredentials: true,          //该字段可选。它的值是一个布尔值，表示是否允许发送Cookie true，即表示服务器明确许可，Cookie可以包含在请求中，一起发给服务器
		AllowedHeaders:   []string{"*"},
	})
	//jwt
	//app.Use(middleware.GetJWT().Serve)
	app.Use(iris.Gzip, logger.New(), crs)
	app.AllowMethods(iris.MethodOptions)
	app.Get("/test", func(ctx context.Context) {
		ctx.WriteString("123456")
	})
	//注册路由
	route.Route(app)
	//websocket  检查token
	websocketRoute := app.Get("/websocket", websocket.Handler(socket.WebSocketHandler(), nil))
	if config.Sysconfig.Basic.SocketCheckToken {
		websocketRoute.Use(middleware.GetJWT().Serve)
	}
	//从任何恐慌中恢复，如果有恐慌，则写入500
	//app.Use(recover.New())
	err := app.Run(iris.Addr(":49111"))
	if err != nil {
		log.Fatal("启动失败", err.Error())
		return
	}

}

var excludeExtensions = [...]string{
	".js",
	".css",
	".jpg",
	".png",
	".ico",
	".svg",
}

func newRequestLogger() (h iris.Handler, close func() error) {
	close = func() error { return nil }
	c := logger.Config{
		Status:  true,
		IP:      true,
		Method:  true,
		Path:    true,
		Columns: true,
	}
	logFile := os.Stdout
	if !config.Sysconfig.Basic.IsLocal {
		logFile = datasource.NewLogFile()
	}
	c.LogFunc = func(now time.Time, latency time.Duration, status, ip, method, path string, message interface{}, headerMessage interface{}) {
		output := logger.Columnize(now.Format("2006/01/02 - 15:04:05"), latency, status, ip, method, path, message, headerMessage)
		logFile.Write([]byte(output))
	}
	//我们不想使用记录器，一些静态请求等
	c.AddSkipper(func(ctx iris.Context) bool {
		path := ctx.Path()
		for _, ext := range excludeExtensions {
			if strings.HasSuffix(path, ext) {
				return true
			}
		}
		return false
	})
	h = logger.New(c)
	return
}
