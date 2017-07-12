package models

import (
	"github.com/garyburd/redigo/redis"
)

type HashElement struct {
	Field string
	Value string
}

type Hash struct {
	Key   string
	Value []HashElement
	Pool  *redis.Pool
}

func (h *Hash) GetKey() string {
	return h.Key
}

func (h *Hash) SetKey(key string) {
	h.Key = key
}

func (h *Hash) GetValue() ([]HashElement, error) {

	conn := h.Pool.Get()
	defer conn.Close()

	cursor := int(0)

	for {
		data, err := redis.Values(conn.Do("HSCAN", h.Key, cursor))
		if err != nil {
			return nil, err
		}

		cursor, err = redis.Int(data[0], nil)
		if err != nil {
			return nil, err
		}

		fieldAndValues, err := redis.Strings(data[1], nil) // passing nil error
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(fieldAndValues); i = i + 2 {
			f := fieldAndValues[i]
			v := fieldAndValues[i+1]

			h.Value = append(h.Value, HashElement{Field: f, Value: v})
		}

		if cursor == 0 {
			break
		}
	}

	return h.Value, nil

}
