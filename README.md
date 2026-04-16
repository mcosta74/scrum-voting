# scrum-voting
Simple application that allows SCRUM teams to do tasks' estimation

## Configuration
The service can be configured using the following command line flags

| name | type | description | default value |
| ---- | ---- | ----------- | ------------- |
| `http.addr` | string | Listen address of the HTTP server | ":8080" |
| `http.timeout` | duration string | Graceful Shutdown Timeout of the HTTP server | "10s" |
| `log.json` | bool | Whether use JSON format for log messages | true |
| `log.level` | "INFO", "DEBUG", "WARN", "ERROR" | Application log level | "INFO" |
| `log.utc` | bool | Whether use UTC for log messages' timestamp | true |

## Development

To run the service in development mode you have to run 

```sh
go -C tools tool task run
```

You can pass command line parameters with
```sh
go -C tools tool task run -- <params>
```

For instance with
```sh
go -C tools tool task run -- -h
```
you get all available command line parameters

