# tinysyslog
[![Go Report Card](http://goreportcard.com/badge/alexferl/tinysyslog)](http://goreportcard.com/report/alexferl/tinysyslog)

A tiny and simple syslog server with log rotation. tinysyslog was born out of the need for a tiny, easy to set up and 
use syslog server that simply writes every incoming log (RFC5424 format) to a file that is automatically rotated, 
to stdout or stderr (mostly for Docker) and or to Elasticsearch.
tinysyslog is based on [go-syslog](https://github.com/mcuadros/go-syslog) and [lumberjack](https://github.com/natefinch/lumberjack).

## Quickstart
```shell
git clone https://github.com/alexferl/tinysyslog.git
cd tinysyslog
make run
```

If tinysyslog started properly you should see:
```shell
2023-08-30T18:38:09-04:00 INF server.go:52 > tinysyslog starting
2023-08-30T18:38:09-04:00 INF server.go:63 > tinysyslog listening on 127.0.0.1:5140
```
You can take make sure logs are processed by the server by entering the following in a terminal:
```shell
nc -w0 -u 127.0.0.1 5140 <<< '<165>1 2016-01-01T12:01:21Z hostname appname 1234 ID47 [exampleSDID@32473 iut="9" eventSource="test" eventID="123"] message'
```

You should then see the following output in your terminal:
```shell
Jan  1 12:01:21 hostname appname[1234]: message
```

## Docker Quickstart
Download the image:
```shell
docker pull admiralobvious/tinysyslog
```
    
Start the container:
```shell
docker run --rm --name tinysyslog -p 5140:5140/udp -d admiralobvious/tinysyslog
```

Send a log:
```shell
nc -w0 -u 127.0.0.1 5140 <<< '<165>1 2016-01-01T12:01:21Z hostname appname 1234 ID47 [exampleSDID@32473 iut="9" eventSource="test" eventID="123"] message'
```

Confirm the container received it:
```shell
docker logs tinysyslog
```

```shell
2023-08-30T22:46:06Z INF build/server.go:52 > tinysyslog starting
2023-08-30T22:46:06Z INF build/server.go:63 > tinysyslog listening on 0.0.0.0:5140
Jan  1 12:01:21 hostname appname[1234]: message
```

## Kubernetes Quickstart
Apply the manifest to your cluster:
```shell
kubectl apply -f kubernetes/tinysyslog.yaml
```

Make sure the container is running:
```shell
kubectl get pods | grep tinysyslog
```

```shell
tinysyslog-6c85886f65-q9cxw          1/1       Running   0          1m
```

Confirm the pod started properly:

```shell
kubectl logs tinysyslog-6c85886f65-q9cxw
```

```shell
2023-08-30T22:46:06Z INF build/server.go:52 > tinysyslog starting
2023-08-30T22:46:06Z INF build/server.go:63 > tinysyslog listening on 0.0.0.0:5140
```

You can now send logs from your app(s) to `tinysyslog:5140`.

## Configuration
```
Usage of ./tinysyslogd:
      --app-name string                           The name of the application. (default "tinysyslog")
      --bind-addr string                          IP and port to listen on. (default "127.0.0.1:5140")
      --env-name string                           The environment of the application. Used to load the right configs file. (default "PROD")
      --filter string                             Filter to filter logs with. Valid filters are: regex.
      --filter-regex string                       Regex to filter with.
      --log-level string                          The granularity of log outputs. Valid levels: 'PANIC', 'FATAL', 'ERROR', 'WARN', 'INFO', 'DEBUG', 'TRACE', 'DISABLED' (default "INFO")
      --log-output string                         The output to write to. 'stdout' means log to stdout, 'stderr' means log to stderr. (default "stdout")
      --log-writer string                         The log writer. Valid writers are: 'console' and 'json'. (default "console")
      --mutator string                            Mutator type to use. Valid mutators are: text, json. (default "text")
      --sink-console-output string                Console to output too. Valid outputs are: stdout, stderr. (default "stdout")
      --sink-elasticsearch-addresses strings      Elasticsearch server address. (default [http://127.0.0.1:9200])
      --sink-elasticsearch-api-key string         Elasticsearch api key.
      --sink-elasticsearch-cloud-id string        Elasticsearch cloud id.
      --sink-elasticsearch-index-name string      Elasticsearch index name. (default "tinysyslog")
      --sink-elasticsearch-password string        Elasticsearch password.
      --sink-elasticsearch-service-token string   Elasticsearch service token.
      --sink-elasticsearch-username string        Elasticsearch username.
      --sink-filesystem-filename string           File to write incoming logs to. (default "syslog.log")
      --sink-filesystem-max-age int               Maximum age (in days) before a log is deleted. (default 30)
      --sink-filesystem-max-backups int           Maximum backups to keep. (default 10)
      --sink-filesystem-max-size int              Maximum log size (in megabytes) before it's rotated. (default 100)
      --sinks strings                             Sinks to save syslogs to. Valid sinks are: console, elasticsearch and filesystem. (default [console])
      --socket-type string                        Type of socket to use, TCP or UDP. If no type is specified, both are used.
```

## Benchmarks
Nothing scientific here but with a simple client consisting of a for loop sending large messages as fast as possible over UDP:

`iostat -d 5`
```
    KB/t tps  MB/s
  127.61 585 72.95
  127.66 592 73.74
  126.41 591 72.98
  126.36 590 72.76
  124.76 615 74.95
```
