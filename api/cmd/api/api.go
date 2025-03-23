package api

import (
	"fmt"
	"go-api/internal/env"
	"go-api/service/routers"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ApiServer struct {
	addr string
	db   *gorm.DB
}

func NewApiServer(addr string, db *gorm.DB) *ApiServer {
	return &ApiServer{
		addr: addr,
		db:   db,
	}
}

func (s *ApiServer) Init(version string) *gin.Engine {
	if !checkVersion(version) {
		panic("version must be in format v1")
	}

	gin.SetMode(env.GetString("GIN_MODE", gin.ReleaseMode))

	// Middlewares

	// logger, recover, cors, requestID
	r := gin.Default()

	// groups
	versionRouter := r.Group(fmt.Sprintf("/api/%s", version))
	log.Printf("API version: %s", version)

	// routers
	routers.NewHealthRouter().RegisterRouter(versionRouter)

	shortenerRouter := routers.NewShortenerRouter(s.db)
	shortenerRouter.RegisterBaseRoutes(r)
	shortenerRouter.RegisterRouter(versionRouter)

	routers.NewAuthRouter(s.db).RegisterRouter(versionRouter)

	return r
}

func (s *ApiServer) Start(r *gin.Engine) error {
	log.Printf("Starting API server on %s", s.addr)
	server := &http.Server{
		Addr:           s.addr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func checkVersion(version string) bool {
	if len(version) < 2 {
		return false
	}

	if version[0] != 'v' {
		return false
	}

	v, err := strconv.Atoi(version[1:])
	if err != nil {
		return false
	}

	if v < 1 {
		return false
	}

	return true
}
