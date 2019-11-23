# Bench

A simple tool to benchmark remote APIs.

## Input

Describe the interfaces in a JSON file.
Each one can be given a name and then described
as a JSON representation of a [go http.Request object](https://golang.org/pkg/net/http/#Request).

### Example:
```
{
  "example test, with headers": {
    "Method": "GET",
    "URL": {
      "Scheme": "https",
      "Host": "www.example.com",
      "Path": "/test",
      "RawQuery": "param=test"
    },
    "Header": {
      "Accept-Encoding": ["gzip", "deflate"]
    }
  }
}
```

By default Bench looks for a file called `tests.json`.
The file can be specified with a `-tests=<path>` flag.

## Output

The results of the benchmark will be printed to standard err.
For more formatted results a template file can be provided with
`-template=<path>`.
The template file should be a valid go template.
A `map` of the results will be passed to the template with the names of the tests as the keys.

### Example:
```
{{ range $test, $result := . }}
{{ $test }} {{ $result }}
{{ end }}
```
