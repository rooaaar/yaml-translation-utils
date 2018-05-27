# Lint yaml files

this tool will check yaml file to make sure:

* all keys sorted alphabetically
* all keys are in lower case
* all key only contains letter and `_`

## Installation

```
go get github.com/sijad/yaml-translation-utils/yamllinter
```

## Usage

```
Usage of yamllinter:
  -file string
        yaml file path
  -level int
        minimal level to check alphabetically sorted
```

## example

```
yamllinter --file file.yml --level 2
```
