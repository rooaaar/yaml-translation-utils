# Checks two yamls file are identical

this tool will check if given files both have same keys and depth.

## Installation

```
go get github.com/sijad/yaml-translation-utils/same-yaml
```

## Usage

```
Usage of same-yaml:
  -ref string
        reference file path
  -tra string
        translation file path
```

## example

```
same-yaml --ref en.yml --tra fa.yml
```
