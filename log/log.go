/*
Copyright 2013 Petru Ciobanu, Francesco Paglia, Lorenzo Pierfederici

This file is part of Mapo.

Mapo is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 2 of the License, or
(at your option) any later version.

Mapo is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Mapo.  If not, see <http://www.gnu.org/licenses/>.
*/

/*
Package log contains a simple multi-level logger.
*/
package log

import (
	"fmt"
	"time"
)

// Available log levels
const (
	ERROR = iota
	INFO
	DEBUG
)

func init() {

    // create a global logger object
    l.level = debug
}

var l logger

// SetLevel sets the output level for the global logger
func SetLevel(level int) {
	l.level = level
}

func print(level int, format string, v ...interface{}) {
	if level <= l.level {
		var msgType string

		switch level {
		case ERROR:
			msgType = "ERROR"
		case INFO:
			msgType = "INFO"
		case DEBUG:
			msgType = "DEBUG"
		}

		msg := fmt.Sprintf(format, v...)
		t := time.Now().Format(time.RFC1123)
		fmt.Printf("%s [%s]: %s\n", t, msgType, msg)

	}
}

// Error logs a message at "ERROR" level
func Error(format string, v ...interface{}) {
	print(ERROR, format, v...)
}

// Info logs a message at "INFO" level
func Info(format string, v ...interface{}) {
	print(INFO, format, v...)
}

// Debug logs a message at "DEBUG" level
func Debug(format string, v ...interface{}) {
	print(DEBUG, format, v...)
}
