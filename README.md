# ShareKBM
ShareKBM is open source alternative software  inspired by applications like synergy and other such applications like ShareMouse which enables you to share your keyboard and mouse across multiple systems


## Build

```
> go mod tidy

> go get -v

> go build
```

## Run
```
> go run app.go -h
usage: sharekbm [-h|--help] -p|--port <integer> -a|--agent (s|c) -t|--target
                "<value>" -b|--buffer <integer> [-n|--name "<value>"]       

                Controls operation parameter for sharekbm app

Arguments:

  -h  --help    Print help information
  -p  --port    Port for application
  -a  --agent   Launch as clieant(c) or server(s)
  -t  --target  Host address for Client(only used in Client)
  -b  --buffer  bufferSize
  -n  --name    an optional client name. Default: client-15516

> go run app.go -p 3000 -a c -t 127.0.0.1 -n desk1

```