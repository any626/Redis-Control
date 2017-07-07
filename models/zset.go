package models

import (
    "github.com/garyburd/redigo/redis"
    // "strconv"
)

type ZSetElement struct {
    Value string
    Score float64
}

type ZSet struct {
    Key string
    Value []ZSetElement
    Pool *redis.Pool
}

func (zs *ZSet) GetKey() string {
    return zs.Key
}

func (zs *ZSet) SetKey(key string) {
    zs.Key = key
}

func (zs *ZSet) GetValue() ([]ZSetElement, error) {

    conn := zs.Pool.Get()
    defer conn.Close()

    cursor := int(0)

    for {
        data, err := redis.Values(conn.Do("ZSCAN", zs.Key, cursor))
        if err != nil {
            return nil, err
        }

        cursor, err = redis.Int(data[0], nil)
        if err != nil {
            return nil, err
        }

        valuesAndScores, err := redis.Values(data[1], nil)
        if err != nil {
            return nil, err
        }

        for i := 0; i < len(valuesAndScores); i = i+2 {
            value, err := redis.String(valuesAndScores[i], nil)
            if err != nil {
                return nil, err
            }

            score, err := redis.Float64(valuesAndScores[i+1], nil)
            if err != nil {
                return nil, err
            }

            zs.Value = append(zs.Value, ZSetElement{Value: value, Score: score})
        }

        if cursor == 0 {
            break;
        }
    }

    return zs.Value, nil
}