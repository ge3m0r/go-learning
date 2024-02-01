package main

import (
	"basic-go/webook/config"
	repository "basic-go/webook/internal/repository"
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/web"
	"basic-go/webook/internal/web/middlewares"
	"basic-go/webook/pkg/ginx/middleware/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	redis "github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func main() {
	db := initDB()
	server := initWebServer()
	initUserHdl(db, server)
	//server := gin.Default()
	server.GET("/hello", func(context *gin.Context) {
		context.String(http.StatusOK, "hello,启动了")
	})
	server.Run(":8080")
}

func initUserHdl(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDAO(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	c := web.NewHandler(us)
	c.RegisterRoutes(server)
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}), func(context *gin.Context) {
		println("这是跨域middleware")
	})
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	//go get github.com/ulule/limiter/v3 另一个限流中间件
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
	useJWT(server)
	return server
}

func useSession(server *gin.Engine) {
	login := &middlewares.LoginMiddleWareBuilder{}
	store := cookie.NewStore([]byte("secret"))
	//store := memstore.NewStore([]byte("qTzTTMzQcpXofciQynLVq1WbwRFeQrFn"), []byte("ik6K0pTEgJ8aqboo011NePKWmX837gxa"))
	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "", []byte("qTzTTMzQcpXofciQynLVq1WbwRFeQrFn"), []byte("ik6K0pTEgJ8aqboo011NePKWmX837gxa"))
	//if err != nil {
	//panic(err)
	//}
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
}

func useJWT(server *gin.Engine) {
	loginJWT := &middlewares.LoginJWTMiddleWareBuilder{}
	server.Use(loginJWT.CheckLogin())
}
