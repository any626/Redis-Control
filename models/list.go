package models

import (
	"github.com/garyburd/redigo/redis"
)

type List struct {
	Key   string
	Value []string
	Pool  *redis.Pool
}

func (l *List) GetKey() string {
	return l.Key
}

func (l *List) SetKey(key string) {
	l.Key = key
}

func (l *List) GetValue() ([]string, error) {

	conn := l.Pool.Get()
	defer conn.Close()

	var err error

	l.Value, err = redis.Strings(conn.Do("LRANGE", l.Key, 0, -1))

	return l.Value, err
}
