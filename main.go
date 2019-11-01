package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/mvc"
	"github.com/onrik/logrus/filename"
	"strings"
	"test/datasource"
	"test/web/controllers"
	"time"

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
	filename := time.Now().Format("20060102") + ".txt"
	//打开一个输出文件，如果重新启动服务器，它将追加到今天的文件中
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		if !datasource.IsLocal {
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
	//logEn.Error("我拍的天啊")
	app := iris.New()
	r, close := newRequestLogger()
	defer close()
	app.Use(r)
	app.Logger().SetOutput(datasource.NewLogFile())
	mvc.New(app.Party("/test")).Handle(new(controllers.TestController))
	err := app.Run(iris.Addr(":8080"))
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
	if !datasource.IsLocal {
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
