package models


import (
    // "github.com/garyburd/redigo/redis"
)

type RedisType interface {
    GetKey() string
    SetKey(string)
    // GetValue(redis.Conn, bool) (interface{}, error)
}