// logserver Receives request for a range of lines. Serves them from
// the files in the logged directory.
package logserver

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func reverse(a []string) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

func readLines(path string) ([]string, error) {
	var lines []string

	file, err := os.Open(path)
	if err != nil {
		return lines, err
	}
	defer file.Close()

	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			break
		}
		lines = append(lines, line)
	}

	return lines, err
}

type Log struct {
	path        string
	name        string
	lastModTime time.Time
	Lines       []string
}

func NewLog(path string) (*Log, error) {
	var lines []string
	var err error
	var file os.FileInfo
	if lines, err = readLines(path); err != nil {
		fmt.Println("cound't read lines for ", path)
	}

	reverse(lines)

	if file, err = os.Stat(path); err != nil {
		fmt.Println("Coundn't open file for Stat", err)
		return nil, err
	}
	return &Log{path: path, Lines: lines, lastModTime: file.ModTime()}, err
}

func (l Log) HasUpdate() bool {
	file, err := os.Stat(l.path)
	if err != nil {
		fmt.Println("Coundn't open file for Stat", err)
		return false
	}
	return file.ModTime().After(l.lastModTime)
}

func (l *Log) Update() (err error) {
	var file os.FileInfo
	if l.Lines, err = readLines(l.path); err != nil {
		fmt.Println("Couldn't open the file for update", err)
		return err
	}

	if file, err = os.Stat(l.path); err != nil {
		fmt.Println("Coundn't open file for Stat", err)
		return err
	}
	l.lastModTime = file.ModTime()
	return nil
}

func (l Log) String() string {
	return fmt.Sprintf("linecnt: %-8d  - path: %s", len(l.Lines), l.path)
}

type Logs struct {
	Logs []*Log
	root string
}

// Range If logs were put together serially, range will return the
// lines between start and end (non included). The index starts from
// 0.
func (lgs *Logs) Range(start, end int) []string {
	var lns []string
	var lstart, lend int

	for _, l := range lgs.Logs {
		ln := len(l.Lines)
		if ln == 0 { // empty file
			continue
		}
		if end <= 0 { // range already finished
			break
		}

		if start >= ln { // range not in this file
			start -= ln
			end -= ln
			continue
		}

		// start is in this file
		lstart = start
		start = 0

		if end < ln {
			lend = end
			end = 0
		} else {
			lend = ln
			end -= ln
		}

		fmt.Println("range", l.Lines[lstart:lend])
		lns = append(lns, l.Lines[lstart:lend]...)
	}

	return lns
}

func (a *Logs) String() string {
	var stra []string
	for i := range a.Logs {
		stra = append(stra, a.Logs[i].String())
	}
	return strings.Join(stra, "\n")
}

// todo: better way to do this.
func (a *Logs) Update() {
	fls, err := ioutil.ReadDir(a.root)
	if err != nil {
		fmt.Println("readError", err)
	}

	a.Logs = []*Log{}
	for i := len(fls) - 1; i >= 0; i-- {
		p := filepath.Join(a.root, fls[i].Name())
		l, err := NewLog(p)
		if err != nil {
			continue
		}
		a.Logs = append(a.Logs, l)
	}
}

func NewLogs(root string) (*Logs, error) {
	fls, err := ioutil.ReadDir(root)
	if err != nil {
		fmt.Println("readError", err)
		return nil, err
	}

	lgs := &Logs{root: root}

	for i := len(fls) - 1; i >= 0; i-- {
		p := filepath.Join(root, fls[i].Name())
		l, err := NewLog(p)
		if err != nil {
			continue
		}
		lgs.Logs = append(lgs.Logs, l)
	}

	return lgs, nil
}

func logPages(w http.ResponseWriter, r *http.Request, lgs *Logs, size int) {
	page := r.FormValue("page")

	p, err := strconv.Atoi(page)
	if err != nil {
		fmt.Println("page param not int", err)
		fmt.Fprintf(w, fmt.Sprintln("page param not int", err))
		return
	}

	// update if first page
	if p == 1 {
		lgs.Update()
	}

	start := size * (p - 1)
	end := start + size
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(w, "<div id='content'>")

	rng := lgs.Range(start, end)
	reverse(rng)
	for _, c := range rng {
		fmt.Fprint(w, "<div class='entry'>"+c+"</div>")
	}
	if len(rng) == size {
		p = p + 1
	}
	nxt_page := strconv.Itoa(p)
	fmt.Println("   page: ", page, " next", nxt_page, " - ")
	fmt.Fprintf(w, "<a id='next' href='/logs?page="+nxt_page+"'>next</a></div>")
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, *Logs, int), root string, size int) http.HandlerFunc {
	lgs, err := NewLogs(root)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("initializing")
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("new request")
		fn(w, r, lgs, size)
	}
}

func LogServerHandler(root string, size int) http.HandlerFunc {
	return makeHandler(logPages, root, size)
}
