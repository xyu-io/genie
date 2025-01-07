package redispool

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

type (
	RedisClient struct {
		ctx       context.Context
		Option    Option
		originCli redis.UniversalClient
	}

	// Option is option for redis client
	Option struct {
		Network      string   `protobuf:"bytes,1,opt,name=address,proto3"  json:"network"` // 单机模式使用
		Timeout      int      `protobuf:"bytes,2,opt,name=timeout,proto3" json:"timeout"`
		DB           int      `protobuf:"bytes,3,opt,name=db,proto3" json:"db"`
		Password     string   `protobuf:"bytes,4,opt,name=password,proto3" json:"password"`
		PoolSize     int      `protobuf:"bytes,5,opt,name=pool_size,proto3" json:"pool_size"`
		AppName      string   `protobuf:"bytes,6,opt,name=app_name,proto3" json:"app_name"`
		ReadTimeout  *int     `protobuf:"bytes,7,opt,name=read_timeout,proto3" json:"read_timeout,omitempty"`
		WriteTimeout *int     `protobuf:"bytes,8,opt,name=write_timeout,proto3" json:"write_timeout,omitempty"`
		UseCluster   bool     `protobuf:"bytes,9,opt,name=use_cluster,proto3" json:"use_cluster" yaml:"use_cluster"`        // 是否使用集群模式
		ClusterAddrs []string `protobuf:"bytes,10,opt,name=cluster_addrs,proto3" json:"cluster_addrs" yaml:"cluster_addrs"` // 集群模式下的节点地址列表
	}
)

var (
	RedisInsMap sync.Map
)

// NewClient connects clients via option
// if one client connects error will throw panic
func MustNewClient(ctx context.Context, redisOption Option) (*RedisClient, error) {
	checkErr := checkRedisOption(redisOption)
	if checkErr != nil {
		return nil, checkErr
	}
	redisIns := GetRedisIns(redisOption.AppName)
	var rClient redis.UniversalClient
	if redisIns == nil {
		cli := redis.NewClient(&redis.Options{
			Addr:         redisOption.Network,
			DB:           redisOption.DB,
			Password:     redisOption.Password,
			PoolSize:     redisOption.PoolSize,
			ReadTimeout:  time.Second * time.Duration(redisOption.Timeout),
			WriteTimeout: time.Second * time.Duration(redisOption.Timeout),
			PoolTimeout:  time.Second * time.Duration(redisOption.Timeout),
		})
		if ctx == nil {
			ctx = context.Background()
		}
		if err := cli.Ping(ctx).Err(); err != nil {
			panic(err)
		}
		rClient = cli
		redisIns = &RedisClient{
			ctx:       ctx,
			Option:    redisOption,
			originCli: rClient,
		}
		RedisInsMap.Store(redisOption.AppName, redisIns)
	}
	return redisIns, nil
}

func MustNewClientCluster(ctx context.Context, redisOption Option) (*RedisClient, error) {
	checkErr := checkRedisOption(redisOption)
	if checkErr != nil {
		return nil, checkErr
	}
	redisIns := GetRedisIns(redisOption.AppName)
	var rClient redis.UniversalClient
	if redisIns == nil {
		cli := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        redisOption.ClusterAddrs,
			Password:     redisOption.Password,
			PoolSize:     redisOption.PoolSize,
			ReadTimeout:  time.Second * time.Duration(redisOption.Timeout),
			WriteTimeout: time.Second * time.Duration(redisOption.Timeout),
			PoolTimeout:  time.Second * time.Duration(redisOption.Timeout),
		})
		// 集群模式不支持ping
		if ctx == nil {
			ctx = context.Background()
		}

		rClient = cli
		redisIns = &RedisClient{
			ctx:       ctx,
			Option:    redisOption,
			originCli: rClient,
		}
		RedisInsMap.Store(redisOption.AppName, redisIns)
	}
	return redisIns, nil
}

func checkRedisOption(redisOption Option) error {
	if redisOption.Network == "" {
		return errors.New("no set redis addr")
	}
	if redisOption.AppName == "" {
		return errors.New("no set redis app name")
	}
	if len(redisOption.AppName) > 16 {
		return errors.New("app name is more than 16 characters")
	}
	//if ok, _ := regexp.MatchString("^[a-z\\d]+$", redisOption.AppName); !ok {
	//	panic(errors.New("app name must be lowercase letters or Numbers"))
	//}

	return nil
}

// Client gets redis client by name that maintaining in package
// if not exist,panic
func GetRedisIns(appName string) *RedisClient {
	v, ok := RedisInsMap.Load(appName)
	if ok {
		return v.(*RedisClient)
	}
	return nil
}

// To use the Baishan Redis cluster, you must use this method to get the specification key
func MakeBSKey(appName, key string) string {
	return appName + "_" + key
}
