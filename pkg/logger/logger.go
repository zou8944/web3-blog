package logger

import "fmt"

func Info(args ...interface{}) {

}

func Infof(format string, args ...interface{}) {

}

func Warn(args ...interface{}) {

}

func Warnf(format string, args ...interface{}) {

}

func Error(args ...interface{}) {

}

func Errorf(format string, args ...interface{}) {
	fmt.Printf(format, args)
}

func Debug(args ...interface{}) {

}

func Debugf(format string, args ...interface{}) {

}
