# topic-bleed

Tool for clearing a kafka topic of unwanted messages.

### Getting started
`go get github.com/daiLlew/topic-bleed`

### Configuration
Configure which topic(s) you wish to bleed by through a configuration `.yml` file in the root dir. Example config:
```yaml
topics:
  - name: myTopic1
    consumer_group: myConsumerGroupd1
    brokers:
      - localhost:9092
  - name: myTopic2
    consumer_group: myConsumerGroupd2
    brokers:
      - localhost:9092
      - localhost:9093
```
By default `topic-bleed` will attempt to load the default `config.yml` from the project root - you can update this as 
required or alternatively you can specify which configuration `.yml` file you wish to load with the `-config` flag
(useful if you want to set up multiple configs)

### Build
```go
go build topic-bleed
```

### Run
Using the default configuration file:
```
./topic-bleed
```
Specifying a configuration file:
```
./topic-bleed -config=my-config.yml
```
