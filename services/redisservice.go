package services

import (
	"io/ioutil"
	"errors"
	"fmt"
	"github.com/any626/Redis-Control/models"
	"github.com/any626/Redis-Control/utils"
	"github.com/garyburd/redigo/redis"
	"log"
	// "github.com/any626/Redis-Control/utils"
	"net"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"path/filepath"
	"os"
	// "strings"
)

// RedisConfig holds the config data for redis

// Get a new Redis Service
func NewRedisService(config *utils.Connection) *RedisService {
	rService := &RedisService{Config: config}
	rService.Pool = rService.getRedisPool()
	return rService
}

// GetRedisPool returns a redis pool provided the RedisConfig
func (rService *RedisService) getRedisPool() *redis.Pool {
	config := rService.Config

	address := fmt.Sprintf("%s:%d", config.Host, config.Port)

	var client *ssh.Client
	if config.SSHEnabled {
		hostKeyCallback, err := knownhosts.New(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
		if err != nil {
			log.Fatal(err)
		}

		privBuff, err := ioutil.ReadFile(config.SSHPrivateKey)
		if err != nil {
			log.Fatal(err)
		}

		var signer ssh.Signer
		if config.SSHPassword != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(privBuff, []byte(config.SSHPassword))
			if err != nil {
				log.Fatal(err)
			}
		} else {
			signer, err = ssh.ParsePrivateKey(privBuff)
			if err != nil {
				log.Fatal(err)
			}
		}

		sshConfig := &ssh.ClientConfig {
			User: config.SSHUser,
			Auth: []ssh.AuthMethod {
				ssh.PublicKeys(signer),
				ssh.Password(config.SSHPassword),
			},
			HostKeyCallback: hostKeyCallback,
		}

		sshAddress := fmt.Sprintf("%s:%d", config.SSHHost, config.SSHPort)
		netConn, err := net.Dial("tcp", sshAddress)
		if err != nil {
			log.Fatal(err)
		}

		clientConn, chans, reqs, err := ssh.NewClientConn(netConn, sshAddress, sshConfig)
		if err != nil {
			netConn.Close()
			log.Fatal(err)
		}
		client = ssh.NewClient(clientConn, chans, reqs)
	}

	return &redis.Pool{
		MaxIdle:   5, // max idle connections
		MaxActive: 5, // max number of connections
		Dial: func() (redis.Conn, error) {
			var dialOption redis.DialOption
			var c redis.Conn
			var err error
			if client != nil {
				dialOption = redis.DialNetDial(func(network, adr string) (net.Conn, error){
					conn, err := client.Dial("tcp", address)
					if err != nil {
						client.Close()
						log.Fatal(err)
					}
					return conn, nil
				})

				c, err = redis.Dial("tcp", address, dialOption)
				if err != nil {
					log.Fatalln(err.Error())
				}
			} else {
				c, err = redis.Dial("tcp", address)
				if err != nil {
					log.Fatalln(err.Error())
				}
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
	Config *utils.Connection
	Pool   *redis.Pool
}

func (rService *RedisService) UpdatePool() {
	rService.Pool.Close()
	rService.Pool = rService.getRedisPool()
}

func (rService *RedisService) GetDatabaseCount() (int, error) {
	conn := rService.Pool.Get()
	defer conn.Close()

	count := 0
	for {
		_, err := conn.Do("SELECT", count)
		if err != nil {
			if count == 0 {
				return 0, err
			}
			break
		}
		count++
	}

	return count, nil

	// if count 
	// data, err := redis.Values(conn.Do("CONFIG", "GET", "databases"))
	// if err != nil {
	// 	return 0, err
	// }

	// return redis.Int(data[1], nil)
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
	Error  error
}

func (rService *RedisService) GetTTL(key string, ch chan TTLResult) {
	conn := rService.Pool.Get()

	defer conn.Close()

	ttl, err := redis.Int64(conn.Do("TTL", key))
	ch <- TTLResult{Result: ttl, Error: err}
}
