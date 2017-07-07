package services

import (
    "fmt"
    "github.com/garyburd/redigo/redis"
    "github.com/any626/Redis-Control/models"
    "log"
    "errors"
)

// RedisConfig holds the config data for redis
type RedisConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Database int    `json:"database"`
}

// Get a new Redis Service
func NewRedisService(config *RedisConfig) *RedisService {
    rService := &RedisService{Config: config}
    rService.Pool = rService.getRedisPool()
    return rService
}

// GetRedisPool returns a redis pool provided the RedisConfig
func (rService *RedisService) getRedisPool() *redis.Pool {

    address := fmt.Sprintf("%s:%d", rService.Config.Host, rService.Config.Port)

    return &redis.Pool{
        MaxIdle:   10, // max idle connections
        MaxActive: 50, // max number of connections
        Dial: func() (redis.Conn, error) {
            c, err := redis.Dial("tcp", address)
            if err != nil {
                log.Fatalln(err.Error())
            }

            // Sets the database
            if rService.Config.Database != 0 {
                c.Do("SELECT", rService.Config.Database)
                if err != nil {
                    c.Close()
                    return nil, err
                }
            }
            return c, err
        },
    }
}

type RedisService struct {
    Config *RedisConfig
    Pool *redis.Pool
}

func (rService *RedisService) UpdatePool() {
    rService.Pool.Close()
    rService.Pool = rService.getRedisPool()
}

func (rService *RedisService) GetDatabaseCount() (int, error) {
    conn := rService.Pool.Get()
    defer conn.Close()

    data, err := redis.Values(conn.Do("CONFIG", "GET", "databases"))
    if err != nil {
        return 0, err
    }

    return redis.Int(data[1], nil)
}

func (rService *RedisService) GetKeys(keys chan string, keyErrors chan error) {

    conn := rService.Pool.Get()

    defer conn.Close()

    cursor := 0

    for {
        data, err := redis.Values(conn.Do("SCAN", cursor))
        if err != nil {
            keyErrors <- err
            close(keys)
            close(keyErrors)
            break
        }

        cursor, err = redis.Int(data[0], nil)
        if err != nil {
            keyErrors <- err
            close(keys)
            close(keyErrors)
            break
        }

        scannedKeys, err := redis.Strings(data[1], nil)
        if err != nil {
            keyErrors <- err
            close(keys)
            close(keyErrors)
            break
        }

        for _, v := range scannedKeys {
            keys <- v
        }

        if cursor == 0 {
            close(keys)
            close(keyErrors)
            break
        }
    }
}

func (rService *RedisService) GetModel(key string) (interface{}, error) {
    conn := rService.Pool.Get()
    defer conn.Close()

    rtype, err := redis.String(conn.Do("TYPE", key))
    if err != nil {
        return nil, err
    }

    switch rtype {
    case "string":
        return &models.String{Key: key, Pool: rService.Pool}, nil
    case "list":
        return &models.List{Key: key, Pool: rService.Pool}, nil
    case "set":
        return &models.Set{Key: key, Pool: rService.Pool}, nil
    case "hash":
        return &models.Hash{Key: key, Pool: rService.Pool}, nil
    case "zset":
        return &models.ZSet{Key: key, Pool: rService.Pool}, nil
    }

    return nil, errors.New("Unknown Type")

}

type TTLResult struct {
    Result int64
    Error error
}

func (rService *RedisService) GetTTL(key string, ch chan TTLResult) {
    conn := rService.Pool.Get()

    defer conn.Close()

    ttl, err := redis.Int64(conn.Do("TTL", key))
    ch <- TTLResult{Result: ttl, Error: err}
}