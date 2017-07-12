package models

import (
	"github.com/garyburd/redigo/redis"
)

// redis string
type String struct {
	Key   string
	Value string
	Pool  *redis.Pool
}

func (s *String) GetKey() string {
	return s.Key
}

func (s *String) SetKey(key string) {
	s.Key = key
}

func (s *String) GetValue() (string, error) {

	conn := s.Pool.Get()
	defer conn.Close()

	var err error
	s.Value, err = redis.String(conn.Do("GET", s.Key))

	return s.Value, err
}
