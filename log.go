package main

import (
	"fmt"
	"os"
)

func LogI(tag string, args ...interface{}) {
	LogLevel("I", tag, args...)
}

func LogD(tag string, args ...interface{}) {
	LogLevel("D", tag, args...)
}

func LogW(tag string, args ...interface{}) {
	LogLevel("W", tag, args...)
}

func LogE(tag string, args ...interface{}) {
	LogLevel("E", tag, args...)
}

func LogF(tag string, args ...interface{}) {
	LogLevel("E", tag, args...)
	os.Exit(1)
}

func LogV(tag string, args ...interface{}) {
	LogLevel("V", tag, args...)
}

func LogLevel(level string, tag string, args ...interface{}) {
	extra := make([]interface{}, 0)
	switch level {
	case "I":
		extra = append(extra, "\u001b[34m["+tag+"]")
	case "D":
		extra = append(extra, "\u001b[35m["+tag+"]")
	case "W":
		extra = append(extra, "\u001b[33m["+tag+"]")
	case "E":
		extra = append(extra, "\u001b[31m["+tag+"]")
	case "V":
		extra = append(extra, "["+tag+"]")
	}
	extra = append(extra, args...)
	extra = append(extra, "\u001b[0m")
	Log(extra...)
}

func Log(args ...interface{}) {
	fmt.Println(args...)
}
