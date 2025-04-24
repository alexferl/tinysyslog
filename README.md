# tinysyslog [![Go Report Card](https://goreportcard.com/badge/github.com/alexferl/tinysyslog)](https://goreportcard.com/report/github.com/alexferl/tinysyslog) [![codecov](https://codecov.io/gh/alexferl/tinysyslog/branch/master/graph/badge.svg)](https://codecov.io/gh/alexferl/tinysyslog)

A tiny and simple syslog server with log rotation. tinysyslog was born out of the need for a tiny, easy to set up and
use syslog server that simply writes every incoming log (in [RFC 5424](https://datatracker.ietf.org/doc/html/rfc5424) format **only**) to a file that is automatically rotated,
to stdout or stderr (mostly for Docker).
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
      --app-name string                   The name of the application. (default "tinysyslog")
      --bind-addr string                  IP and port to listen on. (default "127.0.0.1:5140")
      --env-name string                   The environment of the application. Used to load the right configs file. (default "PROD")
      --filter string                     Filter to filter logs with. Valid filters: [noop regex]
      --filter-regex string               Regex to filter with.
      --log-level string                  The granularity of log outputs. Valid levels: [PANIC FATAL ERROR WARN INFO DISABLED TRACE DISABLED] (default "INFO")
      --log-output string                 The output to write to. Valid outputs: [stdout stderr] (default "stdout")
      --log-writer string                 The log writer. Valid writers: [console json] (default "console")
      --mutator string                    Mutator type to use. Valid mutators: [text json] (default "text")
      --sink-console-output string        Console to output to. Valid outputs: [stdout stderr] (default "stdout")
      --sink-filesystem-filename string   File path to write incoming logs to. (default "syslog.log")
      --sink-filesystem-max-age int       Maximum age (in days) before a log is deleted. (default 30)
      --sink-filesystem-max-backups int   Maximum backups to keep. (default 10)
      --sink-filesystem-max-size int      Maximum log size (in megabytes) before it's rotated. (default 100)
      --sinks strings                     Sinks to save syslogs to. Valid sinks: [console filesystem] (default [console])
      --socket-type string                Type of socket to use, TCP or UDP. If no type is specified, both are used.
```
