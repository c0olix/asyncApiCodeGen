# asyncApiCodeGen

This project aims to generate usable code from an asyncApi-spec (see https://www.asyncapi.com/).

## Features

* Create Kafka Code to be used in:
    * Go
    * Java

## Limits

* General: only Kafka as messaging backend
* Go: for now only a specific flavor of the Kafka package is used which relies on a private repository, so only the types may be usable if you don't
  have access to the private repository

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
| minimum            | :x:                | :heavy_check_mark: | :heavy_check_mark: |
| maximum            | :x:                | :heavy_check_mark: | :heavy_check_mark: |
| exclusiveMinimum   | :x:                | :heavy_check_mark: | :heavy_check_mark: |
| exclusiveMaximum   | :x:                | :heavy_check_mark: | :heavy_check_mark: |
| minLength          | :x:                | :heavy_check_mark: | :heavy_check_mark: |
| maxLength          | :x:                | :heavy_check_mark: | :heavy_check_mark: |
| minItems           | :x:                | :heavy_check_mark: | :heavy_check_mark: |
| maxItems           | :x:                | :heavy_check_mark: | :heavy_check_mark: |
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