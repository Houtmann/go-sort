## go-sort


Create plugin for golangci-lint

`go build -o gosort.so -buildmode=plugin ./plugin/`

In .golangci.yml
```yaml
linters-settings:
  custom:
    gosort:
      path: gosort.so
      description: go-sort linter check if struct fields are sorted
      original-url: github.com/Houtmann/go-sort
```
and add

````yaml
linters:
  disable-all: true
  enable:
    - gosort
````