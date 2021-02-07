package store

import (
	"bytes"
	"log"

	"github.com/gomodule/redigo/redis"
)

const (
	longToShort = "longToShort"
	shortToLong = "shortToLong"
	lastIDKey   = "lastID"
)

// Storage keeps pool of connections for redis, number of saved URLs and channel required for generation of short
// aliases for new incoming URLs.
type Storage struct {
	Pool      *redis.Pool
	IDChannel chan int
	LastID    int
}

// New returns an instance of Storage.
func New(redisURL string, db int) (*Storage, error) {
	s := Storage{
		Pool:      newPool(redisURL, db),
		IDChannel: make(chan int),
	}

	conn := s.Pool.Get()

	lastID, err := s.retrieveLastID()
	if err != nil {
		return nil, err
	}
	s.LastID = lastID

	go func() {
		defer conn.Close()

		for {
			s.LastID++
			s.IDChannel <- s.LastID
			_, err = conn.Do("INCR", lastIDKey)
			if err != nil {
				log.Printf("failed to incr %v due to: %v", lastIDKey, err)
			}
		}
	}()

	return &s, nil
}

func (s *Storage) exists(key string) (bool, error) {
	conn := s.Pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		log.Printf("failed while checking if %v exists", key)
		return false, err
	}

	return exists, nil
}

func (s *Storage) retrieveLastID() (int, error) {
	exists, err := s.exists(lastIDKey)
	if err != nil {
		return 0, err
	}

	conn := s.Pool.Get()
	defer conn.Close()

	switch exists {
	case true:
		lastID, err := redis.Int(conn.Do("GET", lastIDKey))
		if err != nil {
			log.Printf("failed while getting %v", lastIDKey)
			return 0, err
		}
		return lastID, nil
	default:
		const lastID = 0

		log.Printf("key %v does not exist, creating it", lastIDKey)
		_, err = conn.Do("SET", lastIDKey, lastID)
		if err != nil {
			log.Printf("failed while setting %v", lastIDKey)
			return 0, err
		}
		log.Printf("created key %v, set value to %v", lastIDKey, lastID)

		return lastID, nil
	}
}

func newPool(redisURL string, db int) *redis.Pool {
	p := redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisURL, redis.DialDatabase(db))
		},
	}

	return &p
}

// LongByShort searches the original URL in Redis by given short alias.
func (s *Storage) LongByShort(short []byte) ([]byte, error) {
	conn := s.Pool.Get()
	defer conn.Close()

	long, err := redis.Bytes(conn.Do("HGET", shortToLong, short))
	if err != nil {
		if err == redis.ErrNil {
			return nil, err
		}
		return nil, err
	}

	return long, nil
}

// ShortByLong checks if the given URL has a short version saved earlier. If not, it saves it into Redis and returns
// a short alias for the given URL.
func (s *Storage) ShortByLong(longURL []byte) ([]byte, error) {
	conn := s.Pool.Get()
	defer conn.Close()

	short, err := redis.Bytes(conn.Do("HGET", longToShort, longURL))
	if err != nil {
		if err == redis.ErrNil {
			short, err = s.SaveFull(longURL)
			if err != nil {
				return nil, err
			}
			return short, nil
		}
		return nil, err
	}

	return short, nil
}

// SaveFull generates a unique short alias for the given URL, saves the match between this alias and the given
// URL into Redis and returns alias.
func (s *Storage) SaveFull(longURL []byte) ([]byte, error) {
	short := hash(s.IDChannel)
	conn := s.Pool.Get()
	defer conn.Close()

	if _, err := conn.Do("HSET", longToShort, longURL, short); err != nil {
		log.Printf("failed to save long link %q as short %q into Redis", longURL, short)
		return nil, err
	}

	if _, err := conn.Do("HSET", shortToLong, short, longURL); err != nil {
		log.Printf("failed to save short link %q as long %q into Redis", short, longURL)
		return nil, err
	}

	return short, nil
}

// hash generates the unique short alias for the incoming link
func hash(IDChannel <-chan int) []byte {
	const (
		allowedChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		lenChars     = len(allowedChars)
	)

	buf := new(bytes.Buffer)
	id := <-IDChannel
	for id > 0 {
		// the returned error here is always nil according to godoc
		_ = buf.WriteByte(allowedChars[id%lenChars])
		id /= lenChars
	}

	return buf.Bytes()
}

// Close closes all connections to Redis, releasing all resources.
func (s *Storage) Close() {
	err := s.Pool.Close()
	if err != nil {
		log.Printf("failed to close connections to Redis: %v", err)
	} else {
		log.Println("successfully disconnected from Redis")
	}
}
