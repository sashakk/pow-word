package config

import (
	"os"
	"strconv"
)

type Config struct {
	Host        string
	Port        int
	PuzzleSize  int
	TargetBits  int
	ConnTimeout int
}

func getIntVar(key string, defVal int) int {
	rawVal := os.Getenv(key)
	if val, err := strconv.Atoi(rawVal); err != nil {
		return defVal
	} else {
		return val
	}
}

func getStrVar(key string, defVal string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return defVal
	}
	return val
}

func ParseConfig() *Config {
	return &Config{
		Host:        getStrVar("HOST", "localhost"),
		Port:        getIntVar("PORT", 8080),
		PuzzleSize:  getIntVar("PUZZLE_SIZE", 20),
		TargetBits:  getIntVar("TARGET_BITS", 15),
		ConnTimeout: getIntVar("CONN_TIMEOUT", 500),
	}
}
