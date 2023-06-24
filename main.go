package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
)

var (
	interestingExtensions = map[string]bool{
		".js": true, ".html": true,".php": true, ".py": true, ".rb": true, ".sql": true, ".bak": true, ".cfg": true,
		".yml": true, ".env": true, ".asp": true, ".aspx": true, ".jsp": true, ".xls": true, ".mdb": true,
		".db": true, ".log": true, ".tmp": true, ".swp": true, ".ini": true, ".conf": true, ".htaccess": true, ".htpasswd": true,
		".git": true, ".svn": true, ".tar": true, ".zip": true, ".rar": true, ".7z": true, ".htm": true, ".txt": true,
		".cfm": true, ".hta": true, ".inc": true, ".jhtml": true, ".nsf": true, ".pcap": true, ".php2": true,
		".php3": true, ".php4": true, ".php5": true, ".php6": true, ".php7": true, ".phps": true, ".pht": true, ".phtml": true,
		".sh": true, ".shtml": true, ".swf": true, ".xml": true, ".bzip2": true, ".bz2": true, ".gz": true,
	}
	vulnerableParams = map[string]bool{
		"file": true, "view": true, "home": true, "route": true, "load": true, "name": true, "doc": true,
		"page": true, "path": true, "template": true, "lang": true, "redirect": true, "content": true,
		"include": true, "id": true, "user": true, "search": true, "login": true, "value": true,
		"page_id": true, "author": true, "number": true, "query": true, "date": true, "order": true,
		"info": true, "prod": true, "article": true, "entry": true, "cmd": true, "exec": true,
		"run": true, "execute": true, "shell": true, "command": true, "debug": true, "code": true,
		"arg": true, "bash": true, "payload": true, "input": true, "test": true, "data": true,
		"ref": true, "comment": true, "msg": true, "url": true, "link": true, "next": true,
		"rurl": true, "callback": true, "out": true, "goto": true, "jump": true, "return": true,
		"continue": true, "back": true, "forward": true, "dest": true, "img_url": true, "source": true,
		"host": true, "uri": true, "target": true, "site": true, "proxy": true, "refer": true,
		"account": true, "member": true, "usergroup": true, "key": true, "filetype": true, "dir": true,
	}
	vuln     = flag.Bool("vuln", false, "Vulnerable parameters")
	param    = flag.Bool("params", false, "Url with Parameters")
	onlyPath = flag.Bool("ext", false, "Show urls with interesting extensions")
	urlChan  = make(chan string)
	visited  = make(map[string]bool)
	visitedMu = &sync.RWMutex{}
	wg       = &sync.WaitGroup{}
	interestingExtensionsRegex *regexp.Regexp
)

func init() {
	pattern := "("
	for ext := range interestingExtensions {
		pattern += `\` + ext + "|"
	}
	pattern = strings.TrimSuffix(pattern, "|")
	pattern += ")(\\?|$)"
	interestingExtensionsRegex = regexp.MustCompile(pattern)
}

func main() {
	flag.Parse()
	scanner := bufio.NewScanner(os.Stdin)
	const concurrency = 1000

	if !*vuln && !*param && !*onlyPath {
		*vuln = true
		*param = true
		*onlyPath = true
	}

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			for u := range urlChan {
				if *vuln {
					handleVulnerableParams(u)
				}
				if *param {
					handleParams(u)
				}
				if *onlyPath {
					handlePath(u)
				}
			}
			wg.Done()
		}()
	}

	for scanner.Scan() {
		u := scanner.Text()
		urlChan <- u
	}

	close(urlChan)
	wg.Wait()

	if scanner.Err() != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", scanner.Err())
	}
}

func handleVulnerableParams(u string) {
	uParsed, err := url.Parse(u)
	if err != nil {
		return
	}

	for p := range uParsed.Query() {
		if vulnerableParams[p] {
			printURL(u)
			break
		}
	}
}

func handleParams(u string) {
	uParsed, err := url.Parse(u)
	if err != nil {
		return
	}

	if uParsed.RawQuery != "" {
		printURL(u)
	}
}

func handlePath(u string) {
	uParsed, err := url.Parse(u)
	if err != nil {
		return
	}

	if interestingExtensionsRegex.FindStringSubmatch(uParsed.Path) != nil {
		printURL(u)
	}
}

func printURL(u string) {
	visitedMu.Lock()
	defer visitedMu.Unlock()

	if _, seen := visited[u]; !seen {
		visited[u] = true
		fmt.Println(u)
	}
}
