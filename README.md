```
        /'___\  /'___\           /'___\
       /\ \__/ /\ \__/  __  __  /\ \__/
       \ \ ,__\\ \ ,__\/\ \/\ \ \ \ ,__\
        \ \ \_/ \ \ \_/\ \ \_\ \ \ \ \_/
         \ \_\   \ \_\  \ \____/  \ \_\
          \/_/    \/_/   \/___/    \/_/
```

# ffuf - Fuzz Faster U Fool

A fast web fuzzer written in Go.

Heavily inspired by the great projects [gobuster](https://github.com/OJ/gobuster) and [wfuzz](https://github.com/xmendez/wfuzz).

## Installation

- [Download](https://github.com/ffuf/ffuf/releases/latest) a prebuilt binary from [releases page](https://github.com/ffuf/ffuf/releases/latest), unpack and run!
  or
- If you have go compiler installed: `go get github.com/ffuf/ffuf`

The only dependency of ffuf is Go 1.11. No dependencies outside of Go standard library are needed.

## Example usage

### Typical directory discovery

[![asciicast](https://asciinema.org/a/211350.png)](https://asciinema.org/a/211350)

By using the FUZZ keyword at the end of URL (`-u`):

```
ffuf -w /path/to/wordlist -u https://target/FUZZ
```

### Virtual host discovery (without DNS records)

[![asciicast](https://asciinema.org/a/211360.png)](https://asciinema.org/a/211360)

Assuming that the default virtualhost response size is 4242 bytes, we can filter out all the responses of that size (`-fs 4242`)while fuzzing the Host - header:

```
ffuf -w /path/to/vhost/wordlist -u https://target -H "Host: FUZZ" -fs 4242
```

### GET parameter fuzzing

GET parameter name fuzzing is very similar to directory discovery, and works by defining the `FUZZ` keyword as a part of the URL. This also assumes an response size of 4242 bytes for invalid GET parameter name.

```
ffuf -w /path/to/paramnames.txt -u https://target/script.php?FUZZ=test_value -fs 4242
```

If the parameter name is known, the values can be fuzzed the same way. This example assumes a wrong parameter value returning HTTP response code 401.

```
ffuf -w /path/to/values.txt -u https://target/script.php?valid_name=FUZZ -fc 401
```

### POST data fuzzing

This is a very straightforward operation, again by using the `FUZZ` keyword. This example is fuzzing only part of the POST request. We're again filtering out the 401 responses.

```
ffuf -w /path/to/postdata.txt -X POST -d "username=admin\&password=FUZZ" -u https://target/login.php -fc 401
```

### Using external mutator to produce test cases

For this example, we'll fuzz JSON data that's sent over POST. [Radamsa](https://gitlab.com/akihe/radamsa) is used as the mutator.

When `--input-cmd` is used, ffuf will display matches as their position. This same position value will be available for the callee as an environment variable `$FFUF_NUM`. We'll use this position value as the seed for the mutator. Files example1.txt and example2.txt contain valid JSON payloads. We are matching all the responses, but filtering out response code `400 - Bad request`:

```
ffuf --input-cmd 'radamsa --seed $FFUF_NUM example1.txt example2.txt' -H "Content-Type: application/json" -X POST -u https://ffuf.io.fi/ -mc all -fc 400
```

It of course isn't very efficient to call the mutator for each payload, so we can also pre-generate the payloads, still using [Radamsa](https://gitlab.com/akihe/radamsa) as an example:

```
# Generate 1000 example payloads
radamsa -n 1000 -o %n.txt example1.txt example2.txt

# This results into files 1.txt ... 1000.txt
# Now we can just read the payload data in a loop from file for ffuf

ffuf --input-cmd 'cat $FFUF_NUM.txt' -H "Content-Type: application/json" -X POST -u https://ffuf.io.fi/ -mc all -fc 400
```

## Usage

To define the test case for ffuf, use the keyword `FUZZ` anywhere in the URL (`-u`), headers (`-H`), or POST data (`-d`).

```
Usage of ffuf:
  -D	DirSearch style wordlist compatibility mode. Used in conjunction with -e flag. Replaces %EXT% in wordlist entry with each of the extensions provided by -e.
  -H "Name: Value"
    	Header "Name: Value", separated by colon. Multiple -H flags are accepted.
  -V	Show version information.
  -X string
    	HTTP method to use (default "GET")
  -ac
    	Automatically calibrate filtering options
  -acc value
    	Custom auto-calibration string. Can be used multiple times. Implies -ac
  -b "NAME1=VALUE1; NAME2=VALUE2"
    	Cookie data "NAME1=VALUE1; NAME2=VALUE2" for copy as curl functionality.
    	Results unpredictable when combined with -H "Cookie: ..."
  -c	Colorize output.
  -compressed
    	Dummy flag for copy as curl functionality (ignored) (default true)
  -cookie value
    	Cookie data (alias of -b)
  -d string
    	POST data
  -data string
    	POST data (alias of -d)
  -data-ascii string
    	POST data (alias of -d)
  -data-binary string
    	POST data (alias of -d)
  -debug-log string
    	Write all of the internal logging to the specified file.
  -e string
    	Comma separated list of extensions to apply. Each extension provided will extend the wordlist entry once. Only extends a wordlist with (default) FUZZ keyword.
  -fc string
    	Filter HTTP status codes from response. Comma separated list of codes and ranges
  -fl string
    	Filter by amount of lines in response. Comma separated list of line counts and ranges
  -fr string
    	Filter regexp
  -fs string
    	Filter HTTP response size. Comma separated list of sizes and ranges
  -fw string
    	Filter by amount of words in response. Comma separated list of word counts and ranges
  -i	Dummy flag for copy as curl functionality (ignored) (default true)
  -input-cmd value
    	Command producing the input. --input-num is required when using this input method. Overrides -w.
  -input-num int
    	Number of inputs to test. Used in conjunction with --input-cmd. (default 100)
  -k	TLS identity verification
  -mc string
    	Match HTTP status codes from respose, use "all" to match every response code. (default "200,204,301,302,307,401,403")
  -ml string
    	Match amount of lines in response
  -mode string
    	Multi-wordlist operation mode. Available modes: clusterbomb, pitchfork (default "clusterbomb")
  -mr string
    	Match regexp
  -ms string
    	Match HTTP response size
  -mw string
    	Match amount of words in response
  -o string
    	Write output to file
  -od string
    	Directory path to store matched results to.
  -of string
    	Output file format. Available formats: json, ejson, html, md, csv, ecsv (default "json")
  -p delay
    	Seconds of delay between requests, or a range of random delay. For example "0.1" or "0.1-2.0"
  -r	Follow redirects
  -s	Do not print additional information (silent mode)
  -sa
    	Stop on all error cases. Implies -sf and -se. Also stops on spurious 429 response codes.
  -se
    	Stop on spurious errors
  -sf
    	Stop when > 95% of responses return 403 Forbidden
  -t int
    	Number of concurrent threads. (default 40)
  -timeout int
    	HTTP request timeout in seconds. (default 10)
  -u string
    	Target URL
  -v	Verbose output, printing full URL and redirect location (if any) with the results.
  -w value
    	Wordlist file path and (optional) custom fuzz keyword, using colon as delimiter. Use file path '-' to read from standard input. Can be supplied multiple times. Format: '/path/to/wordlist:KEYWORD'
  -x string
    	HTTP Proxy URL
```

eg. `ffuf -u https://example.org/FUZZ -w /path/to/wordlist`


## License

ffuf is released under MIT license. See [LICENSE](https://github.com/ffuf/ffuf/blob/master/LICENSE).
