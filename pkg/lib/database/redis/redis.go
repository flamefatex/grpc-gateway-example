package redis

import (
	"fmt"
	"strings"
	"time"

	"github.com/FZambia/sentinel"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/redigo"
	"github.com/gomodule/redigo/redis"
)

type RedisClient interface {
	Close() error
	Ping() error
	Get(key string) (interface{}, error)
	MGet(keys []string) (interface{}, error)
	Set(key string, value interface{}) error
	SetEx(key string, ttl int32, value interface{}) error
	MSet(vals map[string]interface{}) error
	Del(key string) error
	MDel(keys []string) error
	INCR(key string, addValue int64) (int64, error)
	Expire(key string, ttl int32) error
	ExpireAt(key string, expireTime time.Time) error
	HSet(key string, field string, value interface{}) error
	HMSet(key string, fields map[string]interface{}) error
	HDel(key string, field string) error
	HMDel(key string, fields []string) (interface{}, error)
	HGet(key string, field string) (interface{}, error)
	HMGet(key string, fields []string) (interface{}, error)
	HGetall(key string) (interface{}, error)
	Hkeys(key string) (interface{}, error)
	HIncr(key string, field string, addValue int64) (int64, error)
	Keys(pattern string) (interface{}, error)
	Scan(pattern string, count int64) ([]interface{}, error)
	DoScript(keyCount int, src string, keysAndArgs ...interface{}) (interface{}, error)
	NewMutex(name string, options ...redsync.Option) *redsync.Mutex
	UnlockMutex(mutex *redsync.Mutex) error
	Rename(oldKey string, newKey string) error
	RenameEx(oldKey string, newKey string) error
	Lpush(key string, val interface{}) error
	Rpop(key string) (interface{}, error)
	Llen(key string) (int64, error)
}

type RedisMode = string

const (
	ModeNormalDsn  RedisMode = "NormalDsn"
	ModeNormalAddr RedisMode = "NormalAddr"
	ModeSentinel   RedisMode = "Sentinel"
)

// Config Redis的配置
type Config struct {
	Mode              RedisMode `default:"NormalDsn"` // 模式,支持NormalDsn、NormalAddr、Sentinel
	IsNeedFillDefault bool      //是否需要填充默认参数

	// dial
	Dsn            string        // 用于NormalDsn模式
	Addr           string        // redis地址，用于NormalAddr模式
	Password       string        // redis密码
	Db             int           // db库
	ConnectTimeout time.Duration `default:"3s"`
	ReadTimeout    time.Duration `default:"3s"`
	WriteTimeout   time.Duration `default:"3s"`

	// pool
	MaxIdle     int           `default:"3"`
	MaxActive   int           `default:"64"`
	IdleTimeout time.Duration `default:"240s"`

	// sentinel
	Sentinel *Sentinel // 用于Sentinel
}

type Sentinel struct {
	Addrs          []string
	MasterName     string
	Password       string        // sentinel密码
	ConnectTimeout time.Duration `default:"3s"`
	ReadTimeout    time.Duration `default:"3s"`
	WriteTimeout   time.Duration `default:"3s"`
}

// Pool redis connection struct
type Pool struct {
	conf *Config
	pool *redis.Pool
}

func NewRedisClient(conf *Config) (RedisClient, error) {
	pool := NewPool(conf)
	if pool == nil {
		return nil, fmt.Errorf("fail to create redis pool")
	}

	return pool, nil
}
func fillDefaultConfig(conf *Config) {
	if conf.ConnectTimeout == 0 {
		conf.ConnectTimeout = 3 * time.Second
	}
	if conf.ReadTimeout == 0 {
		conf.ReadTimeout = 3 * time.Second
	}
	if conf.WriteTimeout == 0 {
		conf.WriteTimeout = 3 * time.Second
	}

	if conf.MaxIdle == 0 {
		conf.MaxIdle = 3
	}
	if conf.MaxActive == 0 {
		conf.MaxActive = 64
	}
	if conf.IdleTimeout == 0 {
		conf.IdleTimeout = 240 * time.Second
	}

	if conf.Mode == ModeSentinel && conf.Sentinel != nil {
		if conf.Sentinel.ConnectTimeout == 0 {
			conf.Sentinel.ConnectTimeout = 3 * time.Second
		}
		if conf.Sentinel.ReadTimeout == 0 {
			conf.Sentinel.ReadTimeout = 3 * time.Second
		}
		if conf.Sentinel.WriteTimeout == 0 {
			conf.Sentinel.WriteTimeout = 3 * time.Second
		}
	}
}

func getRedisTimeoutOpts(conf *Config) []redis.DialOption {
	timeoutOpts := make([]redis.DialOption, 0)
	if conf.ConnectTimeout != 0 {
		timeoutOpts = append(timeoutOpts, redis.DialConnectTimeout(conf.ConnectTimeout))
	}
	if conf.ReadTimeout != 0 {
		timeoutOpts = append(timeoutOpts, redis.DialReadTimeout(conf.ReadTimeout))
	}
	if conf.ReadTimeout != 0 {
		timeoutOpts = append(timeoutOpts, redis.DialWriteTimeout(conf.WriteTimeout))
	}
	return timeoutOpts
}

func getRedisOpts(conf *Config) []redis.DialOption {
	opts := make([]redis.DialOption, 0)
	if conf.Password != "" {
		opts = append(opts, redis.DialPassword(conf.Password))
	}
	if conf.Db != 0 {
		opts = append(opts, redis.DialDatabase(conf.Db))
	}
	timeoutOpts := getRedisTimeoutOpts(conf)
	opts = append(opts, timeoutOpts...)
	return opts
}

func getSentinelOpts(s *Sentinel) []redis.DialOption {
	opts := make([]redis.DialOption, 0)
	if s.Password != "" {
		opts = append(opts, redis.DialPassword(s.Password))
	}
	if s.ConnectTimeout != 0 {
		opts = append(opts, redis.DialConnectTimeout(s.ConnectTimeout))
	}
	if s.ReadTimeout != 0 {
		opts = append(opts, redis.DialReadTimeout(s.ReadTimeout))
	}
	if s.ReadTimeout != 0 {
		opts = append(opts, redis.DialWriteTimeout(s.WriteTimeout))
	}
	return opts
}

// NewPool 初始化redis连接池结构 c.Addr like: redis://user:secret@localhost:6379/0?foo=bar&qux=baz
func NewPool(conf *Config) *Pool {
	// 判断，填充默认值
	if conf.IsNeedFillDefault {
		fillDefaultConfig(conf)
	}

	p := &Pool{
		conf: conf,
	}
	var dial func() (redis.Conn, error)
	var testOnBorrow func(c redis.Conn, t time.Time) error

	switch conf.Mode {
	case ModeNormalAddr:
		opts := getRedisOpts(conf)
		dial = func() (redis.Conn, error) {
			return redis.Dial("tcp", conf.Addr, opts...)
		}
		testOnBorrow = func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		}
	case ModeSentinel:
		sntnlOpts := getSentinelOpts(conf.Sentinel)
		sntnl := &sentinel.Sentinel{
			Addrs:      conf.Sentinel.Addrs,
			MasterName: conf.Sentinel.MasterName,
			Dial: func(addr string) (redis.Conn, error) {
				c, err := redis.Dial("tcp", addr, sntnlOpts...)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
		}

		opts := getRedisOpts(conf)
		dial = func() (redis.Conn, error) {
			// 获取master地址
			masterAddr, err := sntnl.MasterAddr()
			if err != nil {
				return nil, err
			}
			// 连接master
			c, err := redis.Dial("tcp", masterAddr, opts...)
			if err != nil {
				return nil, err
			}
			if !sentinel.TestRole(c, "master") {
				c.Close()
				return nil, fmt.Errorf("%s is not redis master", masterAddr)
			}
			return c, nil
		}
		testOnBorrow = func(c redis.Conn, t time.Time) error {
			if !sentinel.TestRole(c, "master") {
				return fmt.Errorf("Role check failed")
			} else {
				return nil
			}
		}
	default: // 默认是NormalDsn
		var dsn = p.conf.Dsn
		if !strings.HasPrefix(dsn, "redis://") {
			dsn = "redis://" + dsn
		}
		timeoutOpts := getRedisTimeoutOpts(conf)
		dial = func() (redis.Conn, error) {
			return redis.DialURL(dsn, timeoutOpts...)
		}
		testOnBorrow = func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		}
	}

	p.pool = &redis.Pool{
		MaxIdle:      conf.MaxIdle,
		MaxActive:    conf.MaxActive,
		IdleTimeout:  conf.IdleTimeout,
		Dial:         dial,
		TestOnBorrow: testOnBorrow,
	}
	return p
}

// Close closes the connection pool.
func (p *Pool) Close() error {
	err := p.pool.Close()
	p.pool = nil
	p.conf = nil
	return err
}

// Pool 获取redis连接池
func (p *Pool) Pool() *redis.Pool {
	return p.pool
}

func (p *Pool) NewMutex(name string, options ...redsync.Option) *redsync.Mutex {
	rs := redsync.New(redigo.NewPool(p.pool))
	return rs.NewMutex(name, options...)
}

func (p *Pool) UnlockMutex(mutex *redsync.Mutex) (err error) {
	if ok, err := mutex.Unlock(); !ok {
		if err != nil {
			return fmt.Errorf("redis mutex unlock failed, err: %s", err.Error())
		}
		return fmt.Errorf("redis mutex unlock failed")
	}

	return
}

// Set set []byte value to redis key
func (p *Pool) Ping() error {
	conn := p.pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", "ping", "pong")
	return err
}

func (p *Pool) execute(cmd string, args ...interface{}) (reply interface{}, err error) {
	conn := p.pool.Get()
	defer conn.Close()
	return conn.Do(cmd, args...)
}

// Get get []byte value from redis key
func (p *Pool) Get(key string) (reply interface{}, err error) {
	return p.execute("GET", key)
}

func (p *Pool) MGet(keys []string) (reply interface{}, err error) {
	return p.execute("MGET", redis.Args{}.AddFlat(keys)...)
}

// Set set []byte value to redis key
func (p *Pool) Set(key string, val interface{}) error {
	_, err := p.execute("SET", key, val)
	return err
}

// setex
func (p *Pool) SetEx(key string, ttl int32, val interface{}) error {
	_, err := p.execute("SETEX", key, ttl, val)
	return err
}

// BatchSet
func (p *Pool) MSet(vals map[string]interface{}) error {
	_, err := p.execute("MSET", redis.Args{}.AddFlat(vals)...)
	return err
}

// Del redis key
func (p *Pool) Del(key string) error {
	_, err := p.execute("DEl", key)
	return err
}

// batch del
func (p *Pool) MDel(keys []string) error {
	_, err := p.execute("DEL", redis.Args{}.AddFlat(keys)...)
	return err
}

// INCR
func (p *Pool) INCR(key string, addValue int64) (int64, error) {
	return redis.Int64(p.execute("INCR", key, addValue))
}

// expire key
func (p *Pool) Expire(key string, ttl int32) error {
	_, err := p.execute("EXPIRE", key, ttl)
	return err
}

// expire at
func (p *Pool) ExpireAt(key string, expireTime time.Time) error {
	_, err := p.execute("EXPIREAT", key, expireTime.Unix())
	return err
}

// hset
func (p *Pool) HSet(key string, field string, value interface{}) error {
	_, err := p.execute("HSET", key, field, value)
	return err
}

// hmset
func (p *Pool) HMSet(key string, fields map[string]interface{}) error {
	_, err := p.execute("HMSET", redis.Args{}.Add(key).AddFlat(fields)...)
	return err
}

// hdel
func (p *Pool) HDel(key string, field string) error {
	_, err := p.execute("HDEL", key, field)
	return err
}

// hmdel
func (p *Pool) HMDel(key string, fields []string) (interface{}, error) {
	return p.execute("HDEL", redis.Args{}.Add(key).AddFlat(fields)...)
}

// hget
func (p *Pool) HGet(key string, field string) (interface{}, error) {
	return p.execute("HGet", key, field)
}

// hmget
func (p *Pool) HMGet(key string, fields []string) (interface{}, error) {
	return p.execute("HMGET", redis.Args{}.Add(key).AddFlat(fields)...)
}

// hgetall
func (p *Pool) HGetall(key string) (interface{}, error) {
	return p.execute("HGETALL", key)
}

// HINCR
func (p *Pool) HIncr(key string, field string, addValue int64) (int64, error) {
	return redis.Int64(p.execute("HINCR", key, field, addValue))
}

// KEYS
func (p *Pool) Keys(pattern string) (interface{}, error) {
	return p.execute("KEYS", pattern)
}

// DoScript
func (p *Pool) DoScript(keyCount int, src string, keysAndArgs ...interface{}) (interface{}, error) {
	conn := p.pool.Get()
	defer conn.Close()

	s := redis.NewScript(keyCount, src)
	return s.Do(conn, keysAndArgs...)
}

func (p *Pool) Scan(pattern string, count int64) ([]interface{}, error) {
	iter := 0
	//var cursor interface{}
	// this will store the keys of each iteration
	keys := make([]interface{}, 0)
	for {
		// we scan with our iter offset, starting at 0
		reply, exeErr := p.execute("SCAN", iter, "match", pattern, "count", count)
		if exeErr != nil {
			return keys, exeErr
		}
		arr, err := redis.Values(reply, exeErr)
		if err != nil {
			return keys, err
		} else {
			// now we get the iter and the keys from the multi-bulk reply
			iter, _ = redis.Int(arr[0], nil)
			keys = append(keys, arr[1])
		}
		// check if we need to stop...
		if iter == 0 {
			break
		}
	}
	return keys, nil
}

// RENAME
func (p *Pool) Rename(oldKey string, newKey string) error {
	_, err := p.execute("RENAME", oldKey, newKey)
	return err
}

// RENAMENX
func (p *Pool) RenameEx(oldKey string, newKey string) error {
	_, err := p.execute("RENAMENX", oldKey, newKey)
	return err
}

// HKEYS
func (p *Pool) Hkeys(key string) (interface{}, error) {
	return p.execute("HKEYS", key)
}

// Lpush LPUSH []byte value to redis key list
func (p *Pool) Lpush(key string, val interface{}) error {
	_, err := p.execute("LPUSH", key, val)
	return err
}

// Rpop RPOP []byte value from redis key list
func (p *Pool) Rpop(key string) (reply interface{}, err error) {
	return p.execute("RPOP", key)
}

// Llen LLEN redis list
func (p *Pool) Llen(key string) (int64, error) {
	return redis.Int64(p.execute("LLEN", key))
}
