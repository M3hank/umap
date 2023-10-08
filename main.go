package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var (
	urlMappings   = make(map[string]map[string][]map[string]string)
	seenParameters = make(map[string]struct{})
	seenPatterns   = make(map[string]struct{})
	reInteger      = regexp.MustCompile(`/\d+([?/]|$)`)
	reContentCheck = regexp.MustCompile(`(post|article|blog)s?|docs|system|support/|/(\d{4}|pages?)/\d+/`)
	staticExtensions = map[string]struct{}{
	    "css": {}, "svg": {}, "png": {}, "mp3": {}, "jpg": {}, "pdf": {}, "woff2": {}, "bmp": {}, "ico": {},
		"mp4": {}, "woff": {}, "jpeg": {}, "ttf": {}, "avi": {}, "webp": {}, "eot": {}, "otf": {}, "gif": {},
	}
)

var parametersFlag = flag.Bool("params", false, "Only output URLs with parameters")

func parametersToDictionary(params string) map[string]string {
	res := make(map[string]string)
	for _, pair := range strings.Split(params, "&") {
		parts := strings.Split(pair, "=")
		if len(parts) == 2 {
			res[parts[0]] = parts[1]
		}
	}
	return res
}

func dictionaryToParameters(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}
	pairs := make([]string, 0, len(params))
	for k, v := range params {
		pairs = append(pairs, k+"="+v)
	}
	return "?" + strings.Join(pairs, "&")
}

func compareParameters(originalParams []map[string]string, newParams map[string]string) bool {
	originalKeys := make(map[string]struct{})
	for _, params := range originalParams {
		for k := range params {
			originalKeys[k] = struct{}{}
		}
	}
	for k := range newParams {
		if _, exists := originalKeys[k]; !exists {
			return true
		}
	}
	return false
}

func hasBadExtension(path string) bool {
	parts := strings.Split(path, ".")
	if len(parts) < 2 {
		return false
	}
	extension := parts[len(parts)-1]
	_, exists := staticExtensions[extension]
	return exists
}

func isContentPath(path string) bool {
	for _, part := range strings.Split(path, "/") {
		if strings.Count(part, "-") > 3 {
			return true
		}
	}
	return false
}

func main() {
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		parsed, err := url.Parse(line)
		if err != nil {
			continue
		}
		host := parsed.Scheme + "://" + parsed.Host
		if _, exists := urlMappings[host]; !exists {
			urlMappings[host] = make(map[string][]map[string]string)
		}
		params := parametersToDictionary(parsed.RawQuery)
		path := parsed.Path
		if hasBadExtension(path) || reContentCheck.MatchString(path) || isContentPath(path) {
			continue
		}
		if _, exists := urlMappings[host][path]; !exists {
			urlMappings[host][path] = make([]map[string]string, 0)
		}
		if compareParameters(urlMappings[host][path], params) {
			urlMappings[host][path] = append(urlMappings[host][path], params)
		}
	}

	for host, paths := range urlMappings {
		for path, allParams := range paths {
			if *parametersFlag && len(allParams) == 0 {
				continue
			}
			for _, params := range allParams {
				fmt.Println(host + path + dictionaryToParameters(params))
			}
			if !*parametersFlag && len(allParams) == 0 {
				fmt.Println(host + path)
			}
		}
	}
}
