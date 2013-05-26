# irclogs - save and scroll irc logs #

Archive messages from your irc channel, and scroll back through
them. Something like [botbot](http://botbot.me/5/log).

## Usage ##

Clone the repository. Compile/Build binary with Golang.

     git clone https://github.com/ciju/irclogs.git
     go get
     go build

Execute it. The options available are.

    $ ./irclogs -h
    Usage of ./irclogs:
      -c="#astest": channel to connect to
      -h=false: Print console options
      -l=".": log directory, also to serve
      -p="3001": port to serve assets and logs
      -s=30: page size, to be served
      ...

## Contribute ##

It works. But could be made much better. To contribute, follow the
same instructions as above. Change the code, build it, run it, test
it, and send pull requests.

Things that are in my mind. 
- better UX/UI
- more efficient serving of logs
- package as a single binary [go-bindata](https://github.com/jteeuwen/go-bindata)

## Technical Details ##

There are three parts to the server. From the `irclogs.go` file.

    go logIRCMessages(p, *channel, quit)
    go serveAssets("./assets")
    go serveLogs(p, *page_size, quit)

1) `logIRCMessages`: It logs the messages sent on the irc
   channel. Uses [goirc](http://github.com/fluffle/goirc/client) to
   listen to a channel, and writes messages on to a text file for the
   particular day.

2) `serveAssets`: This is the part serving the static content from the
   `./assets` directory. Front-end is pretty straightforward. Except
   for the reverse scroll, which is a behavior to the
   [infinite-scroll](https://github.com/paulirish/infinite-scroll). Its
   available at
   [reverse-infinite-scroll](https://github.com/ciju/reverse-infinite-scroll).

3) `serveLogs`: This part is responsible for reading the log files and
   serving the lines, requested by the api. The implementation could
   be improved.
   

## License ##

[MIT](https://raw.github.com/ciju/irclogs/master/LICENSE)

**Sponsored by [ActiveSphere](http://activesphere.com)**




