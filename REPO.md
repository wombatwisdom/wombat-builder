# Repository

## Building
```shell
task build:ww
```

## Managing Libraries
```shell
target/ww library add -m github.com/redpanda-data/benthos/v4 redpanda-benthos
target/ww library add -m github.com/redpanda-data/connect/v4 redpanda-connect
```

## Managing Versions
```shell
target/ww version add redpanda-connect v4.32.1
target/ww version add redpanda-benthos v4.33.0
```

## Managing Packages
```shell
target/ww package add --path public/components/nats redpanda-connect v4.32.1 nats
target/ww package add --path public/components/aws redpanda-connect v4.32.1 aws
target/ww package add --path public/components/io redpanda-benthos v4.33.0 io
target/ww package add --path public/components/pure redpanda-benthos v4.33.0 pure
```

## Managing Package Builds
```shell
target/ww package build --goos darwin --goarch arm64 redpanda-benthos v4.33.0 pure
```