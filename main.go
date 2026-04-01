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

// hasBadExtension checks if the URL path (ignoring query/fragment) has a static file extension.
func hasBadExtension(pathStr string) bool {
	// Use only the path portion — strip any potential query string remnants
	// path.Ext operates on the last element, so we just need the clean path segment
	base := path.Base(pathStr)
	// Strip anything after '?' in case the path itself somehow carries it
	if idx := strings.Index(base, "?"); idx != -1 {
		base = base[:idx]
	}
	extension := strings.TrimPrefix(path.Ext(base), ".")
	extension = strings.ToLower(extension)
	if extension == "" {
		return false
	}
	_, exists := staticExtensions[extension]
	return exists
}

// isContentPath returns true if any path segment looks like a UUID / slug
// (more than 4 hyphens). A threshold of 3 was too aggressive and discarded
// legitimate endpoints like /wp-admin/admin-ajax.php.
func isContentPath(pathStr string) bool {
	for _, part := range strings.Split(pathStr, "/") {
		if strings.Count(part, "-") > 4 {
			return true
		}
	}
	return false
}

func main() {
	flag.Parse()

	// BUG FIX: Use a larger scanner buffer so very long lines are not silently dropped.
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // 1 MB per line

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// BUG FIX: Prepend scheme BEFORE parsing, not after.
		// url.Parse on a scheme-less string like "www.example.com/path"
		// treats the entire string as a path, leaving Host empty.
		if !strings.Contains(line, "://") {
			line = "http://" + line
		}

		parsed, err := url.Parse(line)
		if err != nil {
			continue
		}

		// BUG FIX: Skip URLs that have no host after parsing (e.g. bare relative paths).
		if parsed.Host == "" {
			continue
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

	// BUG FIX: Check scanner error — previously silently ignored read failures.
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "umap: error reading input: %v\n", err)
		os.Exit(1)
	}

	for _, paths := range urlMappings {
		for _, paramMap := range paths {
			for _, fullURL := range paramMap {
				fmt.Println(fullURL)
			}
		}
	}
}
