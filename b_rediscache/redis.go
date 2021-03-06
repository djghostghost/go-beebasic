package b_rediscache

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/go-redis/redis"
	"reflect"
	"strings"
	"sync"
	"time"
)

func Init() {
	BuildRedisClient()
}

var (
	mu   sync.Mutex
	mode string
	// NOTICE!!!
	RedisClient *redis.ClusterClient
	single      *redis.Client
	sentinel    *redis.Client
)

func BuildRedisClient() {
	mu.Lock()
	defer mu.Unlock()

	mode = web.AppConfig.DefaultString("redis.mode", "cluster")

	c := GetRedisClient()
	vc := reflect.ValueOf(c)
	if !vc.IsNil() {
		return
	}

	addrs,_ := web.AppConfig.String("redis.conn")
	if addrs == "" {
		logs.Info("there is not redis.conn config, skip build redis client")
		return
	}
	password,_ := web.AppConfig.String("redis.password")
	diaTimeout := time.Duration(web.AppConfig.DefaultInt64("redis.conn_timeout", 2000)) * 1e6
	readTimeout := time.Duration(web.AppConfig.DefaultInt64("redis.so_timeout", 4000)) * 1e6
	minIdleConns := web.AppConfig.DefaultInt("redis.min_idle_conns", 10)
	maxRetries := web.AppConfig.DefaultInt("redis.max_retries", 3)

	switch mode {
	case "single":
		options := &redis.Options{
			Addr:         addrs,
			Password:     password,
			DialTimeout:  diaTimeout,
			ReadTimeout:  readTimeout,
			WriteTimeout: readTimeout,
			MinIdleConns: minIdleConns,
			MaxRetries:   maxRetries,
		}
		single = redis.NewClient(options)

	case "sentinel":
		masterName,_ := web.AppConfig.String("redis.sentinel.master")
		if len(masterName) == 0 {
			panic("redis.sentinel.master must not empty when redis.mode=sentinel")
		}
		option := &redis.FailoverOptions{
			MasterName:    masterName,
			SentinelAddrs: strings.Split(addrs, ","),
			Password:      password,
			DialTimeout:   diaTimeout,
			ReadTimeout:   readTimeout,
			WriteTimeout:  readTimeout,
			MinIdleConns:  minIdleConns,
			MaxRetries:    maxRetries,
		}
		sentinel = redis.NewFailoverClient(option)
	default:
		RedisClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        strings.Split(addrs, ","),
			Password:     password,
			DialTimeout:  diaTimeout,
			ReadTimeout:  readTimeout,
			WriteTimeout: readTimeout,
			MinIdleConns: minIdleConns,
			MaxRetries:   maxRetries,
		})
	}

	_, err := GetRedisClient().Ping().Result()

	if err != nil {
		logs.Error("Redis Connection error. %s", err)
		RedisClient = nil
	}
}

func GetRedisClient() redis.Cmdable {
	switch mode {
	case "single":
		return single
	case "sentinel":
		return sentinel
	default:
		return RedisClient
	}
}
