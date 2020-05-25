# Genid Uniform Random Id generator

Genid is a Go and Redis based uniform random id generator with support
for numeric and alphanumeric ids

## Setup 

### Requirements
- Go
- Redis 

### MacOS
```
$ go mod download
$ brew install redis
$ redis-server
```

## Run

```
$ go run genid.go --help

usage: main [-h|--help] [-c|--calibrate] [-r|--refill] [-i|--info] [--initial
            <integer>] [--total <integer>] [-t|--type "<value>"]

            Configuration provided for core genid

Arguments:

  -h  --help       Print help information
  -c  --calibrate  Initiate calibration
  -r  --refill     Refill ids
  -i  --info       Show info of app
      --initial    Initial starting number for id. Default: 100
      --total      Total number stored each refill. Default: 100
  -t  --type       Type of id to generate (alphanum, num). Default alphanum.
                   Default: alphanum

$ go run genid.go --info

Genid Id Generation System
Status: Calibrated
Version: 0.1
Next Id: 1njci1
Initial Id: 100000000
Id Type: alphanum
Offset: 3
Refill Amount: 15

```

### Initial calibration
```
go run genid.go -c --initial 100000000 --total 1000 --type num
```

### Generate Id
```
-- If type num
$ go run genid.go 

Genid Id Genertion System
id generated is 100000061

-- If type alphanum
$ go run genid.go 

Genid Id Genertion System
id generated is 1njck4
```