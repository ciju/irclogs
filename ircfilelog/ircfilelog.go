// ircfilelog Logs all channel messages to the channel, into daily log files.
package ircfilelog

import (
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"github.com/golang/glog"
	"os"
	"path/filepath"
	"strings"
)

func logFmtLine(line *irc.Line) string {
	return fmt.Sprintf("%s - %s - %s\n",
		line.Time.UTC().Format("06-01-02 15:04:05"),
		line.Nick,
		strings.Join(line.Args[1:], " "))
}

func logToFile(filename, msg string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		glog.Errorln("error while loggin", err)
	}
	defer f.Close()

	_, err = f.WriteString(msg)
	if err != nil {
		glog.Errorln("error while writing to log file", err)
	}
}

func logFileName(root, channel string, line *irc.Line) string {
	return filepath.Join(root,
		channel+"-"+line.Time.Format("06-01-02")+".txt")
}

func LogLine(root, channel string, line *irc.Line) {
	glog.V(2).Infof("path %s\n", logFileName(root, channel, line))
	logToFile(logFileName(root, channel, line), logFmtLine(line))
	glog.V(2).Infof("logging - %s\n", logFmtLine(line))
}
