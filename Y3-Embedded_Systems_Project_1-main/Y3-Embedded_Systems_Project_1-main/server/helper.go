package main

import (
	"os"
	"strconv"
)

func getEnvBool(key string) (bool, error) {
	s := os.Getenv(key)
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return v, nil
}

func getEnvUint(key string) (uint64, error) {
	s := os.Getenv(key)
	v, err := strconv.ParseUint(s, 10, 0)
	if err != nil {
		return 15, err
	}
	return v, nil
}

func reverseMap(m map[string]string) map[string]string {
	n := make(map[string]string, len(m))
	for k, v := range m {
		n[v] = k
	}
	return n
}
