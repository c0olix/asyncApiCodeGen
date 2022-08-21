# asyncApiCodeGen
This project aims to generate usable code from an asyncApi-spec (see https://www.asyncapi.com/).

## Features
* Create Kafka Code to be used in:
  * Go
  * Java
## Limits
* General: only Kafka as messaging backend
* General: Spec has to be yaml
* Go: for now only a specific flavor of the Kafka package is used which relies on a private repository, so only the types may be usable if you don't have access to the private repository
## Usage
```shell
$ asyncApiCodeGen generate -h
Used to generate code from async api spec. First argument is the spec and the second is the path to the output. 
        
        For example:
        asyncApiCodeGen generate in_spec.yaml out.go

Usage:
  asyncApiCodeGen generate [flags]

Flags:
  -f, --flavor string   Which flavor should be used?? (default "mosaic")
  -h, --help            help for generate
  -l, --lang string     What kind of code should be generated? (default "go")

```