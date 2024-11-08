package kvstore

import "github.com/gomodule/redigo/redis"

var RedisPool *redis.Pool

//model一般返回struct就好了，可能会有&obj用法
//其他情况就按上面吧
