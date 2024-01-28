package main

import (
	repository "basic-go/webook/internal/repository"
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/web"
	"basic-go/webook/internal/web/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func main() {
	db := initDB()
	server := initWebServer()
	initUserHdl(db, server)

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
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
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
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}), func(context *gin.Context) {
		println("这是跨域middleware")
	})
	login := &middlewares.LoginMiddleWareBuilder{}
	//store := cookie.NewStore([]byte("secret"))
	//store := memstore.NewStore([]byte("qTzTTMzQcpXofciQynLVq1WbwRFeQrFn"), []byte("ik6K0pTEgJ8aqboo011NePKWmX837gxa"))
	store, err := redis.NewStore(16, "tcp", "localhost:6379", "", []byte("qTzTTMzQcpXofciQynLVq1WbwRFeQrFn"), []byte("ik6K0pTEgJ8aqboo011NePKWmX837gxa"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
	return server
}
