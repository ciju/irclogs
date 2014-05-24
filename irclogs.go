package main

import (
	"flag"
	irc "github.com/fluffle/goirc/client"
	"github.com/golang/glog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

import (
	msglog "github.com/ciju/irclogs/ircfilelog"
	"github.com/ciju/irclogs/logserver"
)

// todo: directory to log to, as a flag
// todo: option for channel to connect to (list of channels?)

// serve the scroll back files and give api.
// the javascript files.
//

func logIRCMessages(root string, channel string) {
	c := irc.SimpleClient("logbot")
	c.EnableStateTracking()

	c.AddHandler(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		conn.Join(channel)
		glog.Infof("connecting to %s\n", channel)
	})

	c.AddHandler(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) {
		glog.Infoln("disconnecting")
		logIRCMessages(root, channel)
	})

	c.AddHandler("PRIVMSG", func(conn *irc.Conn, line *irc.Line) {
		// only log private messages to Nick's
		if line.Args[0] == channel {
			go msglog.LogLine(root, channel, line)
		}
	})

	if err := c.Connect("irc.freenode.net"); err != nil {
		glog.Error("Connection error: %s\n", err)
	}

}

func serveLogs(root string, page_size int) {
	http.Handle("/logs",
		http.HandlerFunc(
			logserver.LogServerHandler(root, page_size)))
}

func serveAssets(dir string) {
	if dir == "" {
		glog.Fatal("No directory given, to serve")
	}

	http.Handle("/", http.FileServer(http.Dir(dir)))
}

var (
	root      = flag.String("l", ".", "log directory, also to serve")
	page_size = flag.Int("s", 30, "page size, to be served")
	port      = flag.String("p", "3001", "port to serve assets and logs")
	channel   = flag.String("c", "#astest", "channel to connect to")
	help      = flag.Bool("h", false, "Print console options")
)

func main() {
	glog.Infoln("START. Use 'logbot -h' for command line options.")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	p, err := filepath.Abs(*root)
	if err != nil {
		glog.Fatal(err)
	}

	glog.Infof("Saving and serving logs to/from %s\n", p)

	os.MkdirAll(p, 0700)

	go logIRCMessages(p, *channel)
	go serveAssets("./assets")
	go serveLogs(p, *page_size)

	glog.Infof("Serving at port %s\n", *port)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		glog.Fatal("error", err)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-quit
	glog.Infoln("Quit.")
}
