# mqtt-practice

## MQTT v3.1

Lib: https://github.com/eclipse/paho.mqtt.golang

### Run publisher

```
go run main.go --qos {qos_level}
```

### Run subscriber

Note: id must be unique number

```
go run subscriber/main.go --clientNo {id}
```

## MQTT v5

Lib: https://github.com/eclipse/paho.golang

### Run publisher

```
go run mqttv5/main.go --qos {qos_level}
```

### Run subscriber

Note: id must be unique number

```
go run mqttv5/subscriber/main.go --clientNo {id}
```
