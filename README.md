# Patternfinder
Find patterns in HTTP output based on regex string. Display occurrences.

## Installation
`go install github.com/YouGina/patternfinder@latest`

## Example usage
Default usage
`cat urls | patternfinder`

Specify pattern:
`cat urls | patternfinder -p "plugins/([[a-zA-Z0-9-_]+)/"`

## Command-line options
```
Find patterns in http output based on regex string. Display occurences.

Usage:
    patternfinder [options] < urls.txt

Options:
    -p,         --pattern <string>       Pattern to search for, default "plugins/([[a-zA-Z0-9-_]+)/"
    -d,         --depth <int>            Depth to crawl, default 1
    -c,         --concurrency <int>      Concurrency Level, default 2
    -s,         --secure                 Enable TLS verification, default false
    -dbg,       --debug                  Print all found patterns for debugging, default false
```


## Disclaimer
I've used some examples from [@tomnomnom](https://github.com/tomnomnom) and [@hakluke](https://github.com/hakluke)