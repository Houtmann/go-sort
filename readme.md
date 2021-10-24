## go-sort


Create plugin for goling-ci

`go build -buildmode=plugin ./plugin/`

In .golangci.yml
```yaml
linters-settings:
  custom:
    example:
      path: go-sort.so
      description: go-sort linter check if struct fields are sorted
      original-url: github.com/Houtmann/go-sort
```
and add

````yaml
linters:
  disable-all: true
  enable:
    - go-sort
````