package main

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func watchTraces(pool *redis.Pool, ch chan trace) {
	size := 0
	for {
		traces, _ := getTraces(pool.Get())
		if len(traces) > size && size != 0 {
			for _, trace := range traces[size:] {
				ch <- trace
			}
		}
		size = len(traces)
		time.Sleep(time.Second)
	}
}

func getTraces(conn redis.Conn) ([]trace, error) {
	defer conn.Close()

	var traces []trace
	values, err := redis.Strings(conn.Do("HGETALL", "traces"))
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(values); i += 2 {
		traces = append(traces, trace{Id: values[i], Time_nanos: values[i+1]})
	}
	return traces, nil
}

func getTrace(conn redis.Conn, id string) ([]trace, error) {
	defer conn.Close()

	count, err := redis.Int(conn.Do("GET", id+"-count"))
	if err != nil {
		return nil, err
	}
	var (
		trc  = trace{}
		trcs []trace
	)
	for i := 1; i <= count; i++ {
		values, err := redis.Values(conn.Do("HGETALL", fmt.Sprintf("%s-%d", id, i)))
		if err != nil {
			return nil, err
		}
		if err := redis.ScanStruct(values, &trc); err != nil {
			return nil, err
		}
		trcs = append(trcs, trc)
	}
	return trcs, nil
}
