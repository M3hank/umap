package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"
)

var (
	urlMappings      = make(map[string]map[string]map[string]string)
	staticExtensions = map[string]struct{}{
		"css": {}, "svg": {}, "png": {}, "mp3": {}, "jpg": {}, "pdf": {},
		"woff2": {}, "bmp": {}, "ico": {}, "mp4": {}, "woff": {}, "jpeg": {},
		"ttf": {}, "avi": {}, "webp": {}, "ppt": {}, "eot": {}, "otf": {}, "gif": {},
	}
)

var parametersFlag = flag.Bool("params", false, "Only output URLs with parameters")

func parametersToNameSet(params string) map[string]struct{} {
	res := make(map[string]struct{})
	for _, pair := range strings.Split(params, "&") {
		if strings.Contains(pair, "=") {
			parts := strings.SplitN(pair, "=", 2)
			key := parts[0]
			if key != "" {
				res[key] = struct{}{}
			}
		}
	}
	return res
}

func paramNamesToString(paramNames map[string]struct{}) string {
	if len(paramNames) == 0 {
		return ""
	}
	keys := make([]string, 0, len(paramNames))
	for k := range paramNames {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, "&")
}

func hasBadExtension(pathStr string) bool {
	extension := strings.TrimPrefix(path.Ext(pathStr), ".")
	if extension == "" {
		return false
	}
	_, exists := staticExtensions[extension]
	return exists
}

func isContentPath(pathStr string) bool {
	for _, part := range strings.Split(pathStr, "/") {
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
		if parsed.Scheme == "" {
			parsed.Scheme = "http"
		}
		host := parsed.Scheme + "://" + parsed.Host
		if _, exists := urlMappings[host]; !exists {
			urlMappings[host] = make(map[string]map[string]string)
		}
		pathStr := parsed.Path
		if hasBadExtension(pathStr) || isContentPath(pathStr) {
			continue
		}
		paramNames := parametersToNameSet(parsed.RawQuery)
		paramNamesStr := paramNamesToString(paramNames)
		if *parametersFlag && paramNamesStr == "" {
			// Skip URLs without parameters when -params flag is set
			continue
		}
		if _, exists := urlMappings[host][pathStr]; !exists {
			urlMappings[host][pathStr] = make(map[string]string)
		}
		// Use paramNamesStr as key for deduplication
		if _, exists := urlMappings[host][pathStr][paramNamesStr]; !exists {
			// Store the full URL (including parameter values) for output
			fullURL := parsed.String()
			urlMappings[host][pathStr][paramNamesStr] = fullURL
		}
	}

	for _, paths := range urlMappings {
		for _, paramMap := range paths {
			for _, fullURL := range paramMap {
				fmt.Println(fullURL)
			}
		}
	}
}
