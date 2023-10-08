# UrlMap

**UrlMap** is a sophisticated URL filtering tool crafted in Go. It intelligently filters urls discarding uninteresting content.

## Installation
```
go install github.com/M3hank/umap@latest
```

Feed **UrlMap** a list of URLs via stdin, and it will intelligently filter out repetitive URLs, emphasizing the more dynamic or intriguing ones.

**Before & After Filtering:**

| Total URLs | Filtered URLs |
|:-------------:|:-------------:|
| ![All URLs used for testing](https://github.com/M3hank/umap/assets/70057473/deb37664-b6ec-4103-8282-cf54ce4258c8) | ![URLs filtered using Umap](https://github.com/M3hank/umap/assets/70057473/07256b81-19e3-480d-90e5-30e53675e2e6) |

### Commands:
```
cat urls.txt | umap
```

For integrating with other tools:
```
echo "www.example.com" | waybackurls | umap
```

## Options

Spotlight URLs with parameters using the `-params` flag:
```
cat urls.txt | umap -params
```

**Parameter-focused URLs:**

![Parameters](https://github.com/M3hank/umap/assets/70057473/13c42cd9-673c-4afc-8345-bc73e52a4e34)

## Contributing

Contributions are always welcome! Enhance **UrlMap** by contributing code, do
