package core

import (
	"fmt"
	"github.com/DVV-15324/witches-auth/internal/shared/utils"
	pkg "github.com/DVV-15324/witches-auth/pkg"
	w_redis "github.com/DVV-15324/witches-auth/pkg/redis"
	wcmd_utils "github.com/DVV-15324/witches/cmd/utils"
	w_db "github.com/DVV-15324/witches/pkg/core/database/connect/sql"
	"github.com/DVV-15324/witches/pkg/core/response/logger"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	redis_driver "github.com/ulule/limiter/v3/drivers/store/redis"
	"gorm.io/gorm"
	"log"
	"path/filepath"
	"strconv"
	"time"
)

type CoreServices struct {
	Config     *wcmd_utils.Config
	Logger     *logger.ModelLogger
	DB         *gorm.DB
	Redis      *w_redis.RedisClient
	Limiter    *limiter.Limiter
	Helper     *w_utils.Helper
	TxManager  *w_utils.TxManager
	Middleware gin.HandlerFunc
	// Add
	JWT     *pkg.JwtService
	Session *pkg.SessionService
	//Blacklist *pkg.BlacklistService
	Hash    *pkg.Hash
	HashLen int64
}

func NewCoreServices(cfg *wcmd_utils.Config) *CoreServices {
	//UID Bit
	utils.Bits = cfg.UIDBits
	// Logger
	logPath := filepath.Join(cfg.LogPath, "logs.log")
	logg, err := logger.NewFileLogger(logPath, 1, 20, 30)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	// Database
	db, err := w_db.NewDatabaseInstance(cfg, logg)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	if db == nil || db.DB == nil {
		log.Fatal("Database connection is nil")
	}
	// Redis
	address := fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort)
	redisClient, err := w_redis.NewRedisClient(address, cfg.RedisPassword)
	if err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
	}
	// Rate Limiter
	rate := limiter.Rate{
		Period: time.Duration(cfg.RateLimitPeriod) * time.Second,
		Limit:  cfg.RateLimitMax,
	}
	store, _ := redis_driver.NewStore(redisClient.GetClient())
	limiterInstance := limiter.New(store, rate)
	// Services
	REFRESH_TOKEN_TTL, _ := strconv.ParseInt(cfg.Others["REFRESH_TOKEN_TTL"], 10, 64)
	ACCESS_TOKEN_TTL, _ := strconv.ParseInt(cfg.Others["ACCESS_TOKEN_TTL"], 10, 64)
	SESSION_TTL, _ := strconv.ParseInt(cfg.Others["SESSION_TTL"], 10, 64)
	IDLE_TIMEOUT, _ := strconv.ParseInt(cfg.Others["IDLE_TIMEOUT"], 10, 64)
	//REVOKED_TTL, _ := strconv.ParseInt(cfg.Others["IDLE_TIMEOUT"], 10, 64)
	jwt := pkg.NewJwtService(cfg.Others["JWT_SECRET"], ACCESS_TOKEN_TTL, REFRESH_TOKEN_TTL)
	sessionService := pkg.NewSessionService(redisClient.GetClient(), SESSION_TTL, IDLE_TIMEOUT)
	//blacklistService := pkg.NewBlacklistService(redisClient.GetClient(), REVOKED_TTL)
	helper := w_utils.NewHelper(cfg)
	txManager := w_utils.NewTxManager(db.DB)
	hashLen := int64(16)
	// Hash Pass
	hash := new(pkg.Hash)
	return &CoreServices{
		Config:  cfg,
		Logger:  logg,
		DB:      db.DB,
		Redis:   redisClient,
		Limiter: limiterInstance,
		JWT:     jwt,
		HashLen: hashLen,
		Session: sessionService,
		//	Blacklist:  blacklistService,
		Helper:     helper,
		TxManager:  txManager,
		Hash:       hash,
		Middleware: gin.HandlerFunc(func(ctx *gin.Context) { ctx.Next() }),
	}
}
