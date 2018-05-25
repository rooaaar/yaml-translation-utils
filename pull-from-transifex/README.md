# Pull translations from Transifex

this tool will Download all translations using Transifex API.

## Installation

```
go get github.com/sijad/yaml-translation-utils/pull-from-transifex
```

## Usage

```
Usage of pull-from-transifex:
  -api string
        Transifex API key
  -lang string
        Language code
  -out string
        output directory to save resouces
  -project-slug string
        Transifex API key
```

## example

```
pull-from-transifex --api=<TRANSIFEX_API_KEY> --project-slug=<TRANSIFEX_PROJECT_SLUG> --lang=<TRANSIFEX_TRANSLATION_CODE> --out=./out/
```
