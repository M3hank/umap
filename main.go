package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"
)

var (
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
		if pair == "" {
			continue
		}
		parts := strings.SplitN(pair, "=", 2)
		key := parts[0]
		if key != "" {
			res[key] = struct{}{}
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

// hasBadExtension checks if the URL path has a static file extension.
func hasBadExtension(pathStr string) bool {
	base := path.Base(pathStr)
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

// normalizePathForDedup replaces UUIDs and long slugs in the path with a placeholder
// so that similar endpoints (e.g., /user/UUID1 and /user/UUID2) deduplicate down
// to a single representative URL instead of completely dropping them.
func normalizePathForDedup(pathStr string) string {
	parts := strings.Split(pathStr, "/")
	modified := false
	for i, part := range parts {
		// If a path segment has 4 or more hyphens, treat it as a UUID or slug
		if strings.Count(part, "-") > 3 {
			parts[i] = "{slug}"
			modified = true
		}
	}
	if !modified {
		return pathStr
	}
	return strings.Join(parts, "/")
}

func main() {
	flag.Parse()

	// Use bufio.Reader so arbitrarily long lines are processed without buffer size limits or drops.
	reader := bufio.NewReader(os.Stdin)

	// Preserve insertion order for deterministic output across runs
	var outputURLs []string
	seenKeys := make(map[string]struct{})

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "umap: error reading input: %v\n", err)
			os.Exit(1)
		}

		line = strings.TrimSpace(line)
		if line != "" {
			// Handle scheme-less and protocol-relative URLs
			if strings.HasPrefix(line, "//") {
				line = "http:" + line
			} else if !strings.Contains(line, "://") {
				line = "http://" + line
			}

			parsed, parseErr := url.Parse(line)
			if parseErr == nil && parsed.Host != "" {
				host := parsed.Scheme + "://" + parsed.Host
				if parsed.Path == "" {
					parsed.Path = "/"
				}
				pathStr := parsed.Path

				if !hasBadExtension(pathStr) {
					paramNames := parametersToNameSet(parsed.RawQuery)
					paramNamesStr := paramNamesToString(paramNames)

					if !(*parametersFlag && paramNamesStr == "") {
						dedupPath := normalizePathForDedup(pathStr)
						dedupKey := host + "|" + dedupPath + "|" + paramNamesStr
						if _, exists := seenKeys[dedupKey]; !exists {
							seenKeys[dedupKey] = struct{}{}
							outputURLs = append(outputURLs, parsed.String())
						}
					}
				}
			}
		}

		if err == io.EOF {
			break
		}
	}

	for _, fullURL := range outputURLs {
		fmt.Println(fullURL)
	}
}
