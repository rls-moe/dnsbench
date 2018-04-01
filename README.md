# dnsbench

Simple DNS Benchmarking Tool

## Usage

```
$ ./dnsbench --help
usage: dnsbench [<flags>] [<dns servers>...]

Flags:
      --help              Show context-sensitive help (also try --help-long and --help-man).
  -i, --i=1000            Number of Requests per DNS server
  -e, --e="google.com"    DNS Wildcard Endpoint to benchmark against
      --wait-request=100  Time to wait between requests in milliseconds
      --anti-cache        Prepend randomized subdomains to query to prevent some caching. THIS
                          REQUIRES A WILDCARD DNS ENTRY!

Args:
  [<dns servers>]  DNS Servers to ping
```

It is recommended to change the benchmark endpoint.

## License

This tool is licensed under MIT License. See LICENSE for details
