### Task

Implement a recursive, mirroring web crawler . The crawler should be a command-line tool that accepts a starting URL and
a
destination directory. The crawler will then download the page at the URL, save it in the destination directory, and
then recursively proceed to any valid links in this page.

A valid link is the value of an href attribute in an ```<a>``` tag the resolves to urls that are children of the initial
URL.

For example, given initial URL https://start.url/abc , URLs that resolve to https://start.url/abc/foo
and https://start.url/abc/foo/bar are valid URLs, but ones that resolve to
https://another.domain/ or to https://start.url/baz are not valid URLs, and should be skipped.

Additionally, the crawler should:

- Correctly handle being interrupted by Ctrl-C
- Perform work in parallel where reasonable
- Support resume functionality by checking the destination directory for
  downloaded pages and skip downloading and processing where not necessary
- Provide “happy-path” test coverage

Some tips:

- If you’re not familiar with this kind of software, see ```wget --mirror``` for very
  similar functionality
- Document missing features and any other changes you would make if you had
  more time for the assignment implementation.

### Usage

Install `crawler` using the command below:

```shell
go install
```

To view available options and usage, run

```shell
crawler --help
```

To crawl a website, run

```
crawler -s https://example.com -d downloads
```

### Running tests and benchmarks

The project by default uses [github actions](https://github.com/features/actions) to run tests and benchmarks. To run
tests locally, please run the commands below.

Tests:

```shell
  go test -cover -race ./... -v
```

### Missing features

- The crawler should not exit if a link returns a 404. It should attempt to go back to the previous link and skip the
  missing link's URL.
- Keep track of the last crawled link and resume from it instead of starting afresh.  
