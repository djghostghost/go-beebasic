package b_rediscache

import (
	"github.com/beego/beego/v2/core/logs"
	"testing"
	"time"
)

func TestGetRedisClient_Cluster(t *testing.T) {
	logs.Info("test cluster")
	RedisTestSutupWithConf("redis/cluster")
	GetRedisClient().Del("testMode")
	GetRedisClient().Set("testMode", "1", 10 * time.Second)
	v, err := GetRedisClient().Get("testMode").Result()
	if err != nil {
		t.Error(err)
	}
	if v != "1" {
		t.Fail()
	}
}
func TestGetRedisClient_Single(t *testing.T) {
	logs.Info("test single")
	RedisTestSutupWithConf("redis/single")
	GetRedisClient().Del("testMode")
	GetRedisClient().Set("testMode", "1", 10 * time.Second)
	v, err := GetRedisClient().Get("testMode").Result()
	if err != nil {
		t.Error(err)
	}
	if v != "1" {
		t.Fail()
	}
}
func TestGetRedisClient_Sentinel(t *testing.T) {
	logs.Info("test sentinelg")
	RedisTestSutupWithConf("redis/sentinel")
	GetRedisClient().Del("testMode")
	GetRedisClient().Set("testMode", "1", 10 * time.Second)
	v, err := GetRedisClient().Get("testMode").Result()
	if err != nil {
		t.Error(err)
	}
	if v != "1" {
		t.Fail()
	}
}