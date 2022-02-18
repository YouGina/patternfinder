# Patternfinder
Find patterns in HTTP output based on regex string. Display occurrences.

## Installation
`go install github.com/YouGina/patternfinder@latest`

## Example usage
Default usage
`cat urls | patternfinder`

Specify pattern:
`cat urls | patternfinder -p "plugins/([[a-zA-Z0-9-_]+)/"`

## Example output
```
cat urls.txt | patternfinder | sort -k2 -n
stop-user-enumeration   1
woocommerce     1
woocommerce-gateway-authorize-net-cim   1
wordcamp-coming-soon-page       1
wp-accessibility        1
wp-google-maps  1
wp-google-maps-pro      1
wporg-gp-customizations 1
wp-timelines    1
ultimate-faqs   3
seo     4
camptix 8
virtual-embeds  8
wc-post-types   8
blocks  9
jetpack 9
tagregator      9
gutenberg       10
```
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