// This file is part of ezBastion.

//     ezBastion is free software: you can redistribute it and/or modify
//     it under the terms of the GNU Affero General Public License as published by
//     the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.

//     ezBastion is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU Affero General Public License for more details.

//     You should have received a copy of the GNU Affero General Public License
//     along with ezBastion.  If not, see <https://www.gnu.org/licenses/>.

// Package logmanager add helper for logrus
package logmanager

import (
	"fmt"
<<<<<<< HEAD
	"path"
	"runtime"
	"strings"

=======
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	ezbevent "github.com/ezbastion/ezb_lib/eventlogmanager"
>>>>>>> 6a5ce0a8749038afa3ec71a0628dbc7e36664633
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type callInfo struct {
	packageName string
	fileName    string
	funcName    string
	line        int
}

<<<<<<< HEAD
// SetLogLevel set logrus level
func SetLogLevel(LogLevel string, exPath string, fileName string, maxSize int, maxBackups int, maxAge int) error {
=======
// EventLog management
var eid int
var osspecific bool
var level string

func init() {
	osspecific = true
}

// SetLogLevel set logrus level
func SetLogLevel(LogLevel string, exPath string, fileName string, maxSize int, maxBackups int, maxAge int, interactive bool, reportcaller bool, jsontostdout bool) error {
>>>>>>> 6a5ce0a8749038afa3ec71a0628dbc7e36664633
	log.SetFormatter(&log.JSONFormatter{})
	level = LogLevel

	switch LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
		break
	case "info":
		log.SetLevel(log.InfoLevel)
		break
	case "warning":
		log.SetLevel(log.WarnLevel)
		break
	case "error":
		log.SetLevel(log.ErrorLevel)
		break
	case "critical":
		log.SetLevel(log.FatalLevel)
		break
	default:
		return fmt.Errorf("ezb_lib/logmanager/SetLogLevel() failed: Bad log level name")
	}
<<<<<<< HEAD
	log.SetOutput(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
	})
=======
	// Adding the method and line caller, easier to debug
	log.SetReportCaller(reportcaller)
	abspathfilename := exPath+string(os.PathSeparator)+fileName

	lj := &lumberjack.Logger{
		Filename:   abspathfilename,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
	}

	if jsontostdout {
	mWriter := io.MultiWriter(os.Stderr, lj)
		log.SetOutput(mWriter)
	} else {
		log.SetOutput(lj)
	}
	log.Info("Log system initialized.")
	return nil

}

func retrieveCallInfo() *callInfo {
	pc, file, line, _ := runtime.Caller(2)
	_, fileName := path.Split(file)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	packageName := ""
	funcName := parts[pl-1]

	if parts[pl-2][0] == '(' {
		funcName = parts[pl-2] + "." + funcName
		packageName = strings.Join(parts[0:pl-2], ".")
	} else {
		packageName = strings.Join(parts[0:pl-1], ".")
	}

	return &callInfo{
		packageName: packageName,
		fileName:    fileName,
		funcName:    funcName,
		line:        line,
	}
}

func WithFields(s1 string,s2 string) {
	log.WithFields(log.Fields{s1:s2})
}

func StartWindowsEvent(name string) {
	if osspecific == true {
		if ezbevent.Status == 0 {
				ezbevent.Open(name)
		}
	}
}

// Info logs an info event into the windows eventlog system
func Debug(logline string) error {
	log.Debugln(logline)
	if level == "debug" {
		if osspecific == true {
			if ezbevent.Status == 0 {
				ezbevent.Elog.Info(1, "DEBUG : "+logline)
			}
		}
	}

	return nil
}

func Info(logline string, forceStdout ...bool) error {
	log.Infoln(logline)
	output := false
	if len(forceStdout) > 0 {
		output = forceStdout[0]
	}

	if output {
		fmt.Println(logline)
	}
	if level == "debug" || (level == "info") {
		if osspecific == true {
			if ezbevent.Status == 0 {
				ezbevent.Elog.Info(1, logline)
			}
		}
	}
	return nil
}

// Error logs an error event into the windows eventlog system
func Error(logline string) error {
	log.Errorln(logline)
	if (level == "info") || (level == "warning") || (level == "error") || (level == "debug") {
		if osspecific == true {
			if ezbevent.Status == 0 {
				ezbevent.Elog.Error(1, logline)
			}
		}
	}
>>>>>>> 6a5ce0a8749038afa3ec71a0628dbc7e36664633
	return nil
}

// Warning logs an warning event into the windows eventlog system
func Warning(logline string) error {
	log.Warnln(logline)
	if (level == "debug") || (level == "info") || (level == "warning") {
		if osspecific == true {
			if ezbevent.Status == 0 {
				ezbevent.Elog.Warning(1, logline)
			}
		}
	}
	return nil
}

<<<<<<< HEAD
func retrieveCallInfo() *callInfo {
	pc, file, line, _ := runtime.Caller(2)
	_, fileName := path.Split(file)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	packageName := ""
	funcName := parts[pl-1]

	if parts[pl-2][0] == '(' {
		funcName = parts[pl-2] + "." + funcName
		packageName = strings.Join(parts[0:pl-2], ".")
	} else {
		packageName = strings.Join(parts[0:pl-1], ".")
	}

	return &callInfo{
		packageName: packageName,
		fileName:    fileName,
		funcName:    funcName,
		line:        line,
	}
}
=======
func Fatal(logline string) {
	log.Fatal(logline)
}

>>>>>>> 6a5ce0a8749038afa3ec71a0628dbc7e36664633
