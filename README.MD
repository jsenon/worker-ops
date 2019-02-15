# Worker Ops

This is a tool to report usage of worker node

## Func

- Generate static report with command line
- Launch Api server 
    - With report generation
    - With Custom Prometheus metrics
- Send report to slack

Instrumentation:

- Prometheus GoLang 
- OpenTracing
- well-known and health endpoint

### Build

Locally 
```sh
make build_local
```

Build container
```sh
make build
```

You can have a look on [Makefile](Makefile),in order to launch linter or swagger file generation

### Usage

Launch command-line
```sh
worker-ops --help
```

```sh
An Utility to launch reporter for worker node, or launch a web server API
to export usage to prometheus, and send report to slack

Usage:
  worker-ops [command]

Available Commands:
  help        Help about any command
  report      Launch reporting on stdout
  server      Launch Worker Ops Server

Flags:
      --debug              Set log level to Debug
  -h, --help               help for worker-ops
      --jaegerurl string   Set jaegger collector endpoint
  -t, --toggle             Help message for toggle

Use "worker-ops [command] --help" for more information about a command.
```

Launch static report generation
```sh
worker-ops report
```

Launch Server, and open [http://127.0.0.1:8080](http://127.0.0.1:8080)
```sh
worker-ops server
```

For Slack reporting don't miss to export env `SLACK_URL` containing your slack webhook

### Doc

You can find package doc

```sh
godoc .
```

```sh
COMMAND DOCUMENTATION

    Worker-ops generate can print static report, or serve an api with slack
    and metrics functionnality

SUBDIRECTORIES

[...]
 ```


### API 

You can have a look on the [swagger](swagger.yml)

Api are:

- /metrics for prometheus export
- /report for reporting
- /healthz for health check probe
- /send for sending report (will be used to trigger report generation from scheduled job)


## TODO

- [ ] Change url and port static to env
- [x] Link with slack webhook
- [x] Link with aws, made the calculation
- [ ] AWS env from vault