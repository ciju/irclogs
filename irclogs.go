package main

import (
	"flag"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

func logIRCMessages(root string, channel string, quit chan bool) {
	c := irc.SimpleClient("logbot")
	c.EnableStateTracking()

	c.AddHandler(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		conn.Join(channel)
		fmt.Println("connecting to ", channel)
	})

	c.AddHandler(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) {
		fmt.Println("disconnecting")
		quit <- true
	})

	c.AddHandler("PRIVMSG", func(conn *irc.Conn, line *irc.Line) {
		// only log private messages to Nick's
		if line.Args[0] == channel {
			go msglog.LogLine(root, channel, line)
		}
	})

	if err := c.Connect("foonetic.net"); err != nil {
		fmt.Printf("Connection error: %s\n", err)
	}

}

func serveLogs(root string, page_size int, quit chan bool) {
	http.Handle("/logs",
		http.HandlerFunc(
			logserver.LogServerHandler(root, page_size)))
}

func serveAssets(dir string) {
	if dir == "" {
		log.Fatal("No directory given, to serve")
	}

	http.Handle("/", http.FileServer(http.Dir(dir)))
}

var (
	root      = flag.String("l", ".", "log directory, also to serve")
	page_size = flag.Int("s", 30, "page size, to be served")
	port      = flag.String("p", "3001", "port to serve assets and logs")
	channel   = flag.String("c", "#astest", "channel to connect to")
)

func main() {
	flag.Parse()

	quit := make(chan bool)

	p, err := filepath.Abs(*root)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Saving and serving logs to/from %s\n", p)

	os.MkdirAll(p, 0700)

	go logIRCMessages(p, *channel, quit)
	go serveAssets("./assets")
	go serveLogs(p, *page_size, quit)

	fmt.Println("Serving at port", *port)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		log.Fatal("error", err)
	}

	<-quit
}
