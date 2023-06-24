# UrlMap
**A sophisticated URL filtering tool, written in Go, that effectively identifies and categorizes URLs based on intriguing file extensions and query parameters**

### Installation
```
go install github.com/M3hank/umap@latest
```
### Basic Usage
**At its most basic, UrlMap can accept a list of URLs via stdin (Standard Input), when used in its most basic mode without any command-line arguments, UrlMap will filter URLs based on interesting file extensions, query parameters & performs de-duplication.:** 
```
cat urls.txt | umap
```
**Umap can be easily integrated with other tools.**
```
echo "www.example.com" | waybackurls | umap
```
---
### Available Options

UrlMap provides the following command-line options:

```-vuln```: This option filters URLs based on known vulnerable query parameters such as 'id', 'file', 'page', 'path', etc. This can help in quickly identifying potentially risky URLs.

```
cat urls.txt | umap -vuln
```
![umap -vuln](https://github.com/M3hank/umap/assets/70057473/e3703dc7-64c9-4a2f-a7ce-161f1e82f648)
---

```-params```: Use this option to display only URLs that contain query parameters. This can be handy when you're interested in web pages that are likely to display different content based on the parameters.
```
cat urls.txt | umap -params
```
![umap -params](https://github.com/M3hank/umap/assets/70057473/4f74671c-9774-4162-aced-522ca702a5a8)
---

```-ext```: With this option, UrlMap will display only the URLs that have an interesting file extension (such as .php, .js, .sql, etc.) in their path. It's useful when you're interested in identifying scripts, configuration files, or other potentially interesting resources.
```
cat urls.txt | umap -ext
```
![umap -ext](https://github.com/M3hank/umap/assets/70057473/44782c46-19ca-4a50-aa0a-c7c5929a51df)
---



## Contributing

Contributions are always welcome!
