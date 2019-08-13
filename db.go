package main

import (
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"

	"github.com/go-redis/redis"
)

type DB struct {
	conn *redis.Client
}

func NewDB(uri string) *DB {
	opt, err := redis.ParseURL(uri)
	if err != nil {
		panic("unable to parse redis uri: " + err.Error())
	}

	return &DB{conn: redis.NewClient(opt)}
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) UpdateForUrlMatch(m *urlMatch, addr string) {
	if _, err := db.conn.SAdd(m.RedisKey(), addr).Result(); err != nil {
		panic("unable to SAdd for url: " + err.Error())
	}
	if _, err := db.conn.SAdd(m.id, m.RedisKey()).Result(); err != nil {
		panic("unable to SAdd for id: " + err.Error())
	}
}

func (db *DB) GetResultsForId(id string) (*LookUpResults, error) {
	client := NewIpInfoClient(os.Getenv("DNS_LEAK_IP_INFO_KEY"))

	results := &LookUpResults{Id: id}

	keys, err := db.conn.SMembers(id).Result()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get results for id")
	}

	var wg sync.WaitGroup
	var m sync.Mutex

	for _, key := range keys {
		wg.Add(1)

		go func(k string) {
			defer wg.Done()

			num, _ := strconv.Atoi(strings.Split(k, ".")[0])
			keyResult := LookUpResult{Number: num}

			addrs, err := db.conn.SMembers(key).Result()
			if err != nil {
				panic("unable to SMembers for key: " + err.Error())
			}

			for _, addr := range addrs {
				res, err := client.LookUpIp(addr)
				if err != nil {
					logrus.WithError(err).Errorf("unable to lookup ip %s", addr)
					res = nil // make sure it's nil
				}

				if res != nil && res.Bogon {
					logrus.Warningf("bogon ip: %s", addr)
					res = nil
				}

				keyResult.IPs = append(keyResult.IPs, LookUpResultIp{
					Address: addr,
					Info:    res,
				})
			}

			m.Lock()
			results.Results = append(results.Results, keyResult)
			m.Unlock()
		}(key)
	}

	wg.Wait()
	results.SortResults()

	return results, nil
}
