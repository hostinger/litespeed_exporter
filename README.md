# LiteSpeed exporter for Prometheus
[![Go Report Card](https://goreportcard.com/badge/github.com/hostinger/litespeed_exporter)](https://goreportcard.com/report/github.com/hostinger/litespeed_exporter)

This is a LiteSpeed exporter for Prometheus that runs on the same machine as
the target LiteSpeed server and exports its metrics from the automatically-generated `.rtreport` files.

## Usage
```
./litespeed_exporter [<flags>]
```

#### Flags
Name                        | Description
----------------------------|-----------------
help                        | Show usage help
version                     | Show application version.
log.level                   | Only log messages with the given severity or above. One of: [debug, info, warn, error]
log.format                  | Output format of log messages. One of: [logfmt, json]
web.telemetry-path          | HTTP path to metrics
web.listen-address          | HTTP address to listen on for web interface and telemetry
litespeed.scrape-pattern    | Pattern of files to scrape LiteSpeed metrics from
litespeed.exclude-metrics   | Comma-separated list of metrics to exclude. Available options: `AVAILCONN, AVAILSSL, BPS_IN, BPS_OUT, EXTAPP_CMAXCONN, EXTAPP_EMAXCONN, EXTAPP_IDLE_CONN, EXTAPP_INUSE_CONN,EXTAPP_POOL_SIZE, EXTAPP_REQ_PER_SEC, EXTAPP_TOT_REQS, EXTAPP_WAITQUE_DEPTH, IDLECONN, MAXCONN, MAXSSL_CONN, PLAINCONN, REQ_RATE_PRIVATE_CACHE_HITS_PER_SEC REQ_RATE_PUB_CACHE_HITS_PER_SEC, REQ_RATE_REQ_PER_SEC, REQ_RATE_REQ_PROCESSING REQ_RATE_STATIC_HITS_PER_SEC, REQ_RATE_TOTAL_PRIVATE_CACHE_HITS, REQ_RATE_TOTAL_PUB_CACHE_HITS REQ_RATE_TOTAL_STATIC_HITS, REQ_RATE_TOT_REQS, SSLCONN, SSL_BPS_IN, SSL_BPS_OUT`
litespeed.req-rates-by-host | Export Request Rates by host
litespeed.metrics-by-core   | Export metrics by core filename
litespeed.exclude-extapp    | Exclude EXTAPP metrics altogether

## Builds

#### Pre-built binaries
For already built binaries, check the [releases](https://www.github.com/hostinger/litespeed_exporter/releases).

#### Building locally
```sh
make build
```

## Compatibility
 - Go 1.15+
 - LSWS 6.0 | 5.4 | 5.3
