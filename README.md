Parses yml files to retrieve the given value

## Installation

```
go install github.com/reactiveraven/ymlparse...
```

## Usage

```
> ls
docker-compose.yml docker-compose.override.yml
> cat docker-compose.yml
services:
   foo:
      bar: original value
> cat docker-compose.override.yml
services:
   foo:
      bar: overridden value
> ymlparse services foo bar
overridden value
```

options:

| flag | |
| ---- | - |
| i | input file to read from |
| override | override file to read from (attempts to find `input`.override.yml if not provided) |
