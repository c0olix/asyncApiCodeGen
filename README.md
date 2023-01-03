# asyncApiCodeGen

This project aims to generate usable code from an asyncApi-spec (see https://www.asyncapi.com/).

## Features

* Create Kafka Code to be used in:
    * Go
    * Java

## Limits

* General: only Kafka as messaging backend

## Datatypes, Format, Validations

Currently, these datatypes and formats are supported

| Datatype/Format     | Java               | Go                 | Validate           |
|---------------------|--------------------|--------------------|--------------------|
| integer             | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| int32 (format)      | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| int64 (format)      | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| number              | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| float (format)      | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| double (format)     | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| string              | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| email (format)      | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| binary (format)     | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| date (format)       | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| date-time (format)  | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| password (format)   | :x:                | :x:                | :heavy_check_mark: |

Further the generator support to add validation/featuress on the resulting types/classes for

| Validation/Feature | Java               | Go                 | Validate           |
|--------------------|--------------------|--------------------|--------------------|
| required           | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| email              | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| default            | :heavy_check_mark: | :x:                | :heavy_check_mark: |
| minimum            | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| maximum            | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| exclusiveMinimum   | :x:                | :heavy_check_mark: | :heavy_check_mark: |
| exclusiveMaximum   | :x:                | :heavy_check_mark: | :heavy_check_mark: |
| minLength          | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| maxLength          | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| minItems           | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| maxItems           | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| uniqueItems        | :x:                | :heavy_check_mark: | :heavy_check_mark: |
| anyOf              | :x:                | :x:                | :heavy_check_mark: |
| oneOf              | :x:                | :x:                | :heavy_check_mark: |
| allOf              | :x:                | :x:                | :heavy_check_mark: |
| not                | :x:                | :x:                | :heavy_check_mark: |
| pattern            | :x:                | :x:                | :heavy_check_mark: |
| enum               | :x:                | :x:                | :heavy_check_mark: |
| multipleOf         | :x:                | :x:                | :heavy_check_mark: |
| minProperties      | :x:                | :x:                | :heavy_check_mark: |
| maxProperties      | :x:                | :x:                | :heavy_check_mark: |
| externalDoc        | :x:                | :x:                | :heavy_check_mark: |
| nullable           | :x:                | :x:                | :heavy_check_mark: |
| readOnly           | :x:                | :x:                | :heavy_check_mark: |
| writeOnly          | :x:                | :x:                | :heavy_check_mark: |

## Usage
### RootCmd
```shell
$ asyncApiCodeGen -h
This CLI-Tool is used to generate code for given async api spec

Usage:
  asyncApiCodeGen [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  generate    Used to generate code from async api spec
  help        Help about any command
  validate    Validate given asyncApiSpec

Flags:
  -h, --help           help for asyncApiCodeGen
  -i, --input string   Where is the source spec located?
```
### Generate
```shell
$ asyncApiCodeGen generate -h
Used to generate code from async api spec. First argument is the spec and the second is the path to the output. 
        
        For example:
        asyncApiCodeGen generate -i in_spec.yaml -o out.go

Usage:
  asyncApiCodeGen generate [flags]

Flags:
  -c, --createDir            Should directory be created if not present (recursive)?
  -f, --flavor string        Which (if) flavor should be used?
  -h, --help                 help for generate
  -l, --lang string          What kind of code should be generated? (default "go")
  -o, --output string        Where should the generated code saved to? Attention: Go=File, Java=Dir!
  -p, --packageName string   Which package name should the generated code have?

Global Flags:
  -i, --input string   Where is the source spec located?
```
#### Possible options for flags
##### --createDir
* true
* false (default)
##### --flavor
* "" (Blank) - Generates a default
* "mosaic" - Generates code with the mosaic flavor, which includes a private repository
* "mqtt" - Generates a mqtt compatible api (go only for now)
##### --lang
* "java"
* "go"
##### --output
A path or file, where the generated code should be created.
>__Attention__: In case of go, a single file will be created, so output must be a file! In case of java multiple files will be created, so the output value must be a directory!
##### --packageName
The name of the package e.g. events or com.yourcompany.events
### Validate
```shell
$ asyncApiCodeGen validate -h
Validate given asyncApiSpec.

Usage:
  asyncApiCodeGen validate [flags]

Flags:
  -h, --help   help for validate

Global Flags:
  -i, --input string   Where is the source spec located?

```