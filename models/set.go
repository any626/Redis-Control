package models

import (
	"github.com/garyburd/redigo/redis"
)

type Set struct {
	Key   string
	Value []string
	Pool  *redis.Pool
}

func (s *Set) GetKey() string {
	return s.Key
}

func (s *Set) SetKey(key string) {
	s.Key = key
}

func (s *Set) GetValue() ([]string, error) {

	conn := s.Pool.Get()
	defer conn.Close()

	cursor := int64(0)

	for {
		data, err := redis.Values(conn.Do("SSCAN", s.Key, cursor))
		if err != nil {
			return nil, err
		}

		cursor, err = redis.Int64(data[0], nil)
		if err != nil {
			return nil, err
		}

		scannedKeys, err := redis.Strings(data[1], nil)
		if err != nil {
			return nil, err
		}

		s.Value = append(s.Value, scannedKeys...)

		if cursor == 0 {
			break
		}
	}

	return s.Value, nil
}
