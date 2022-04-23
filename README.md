# rangetest

RangeTest is a tool for testing the compliance of web servers against the HTTP Range specification.

# Usage

## Webserver Setup

Spin up the web server you wish to test and ensure it can serve the sample file `data.txt` (included in this repository).

Then, to trigger the test suite:

```
./rangetest -u <absolute URL pointing to the data.txt file>
```

The URL value must be a fully qualified HTTP URL, for example: `http://localhost:8888/data.txt`. HTTPS is supported but
no validation is performed on the certificate.

## Example Output

```
[PASS] HEAD request
[PASS] Get all data without range
[PASS] Get all data with range
[PASS] Get single absolute range
[PASS] Get single relative range with start index
[PASS] Get single relative range with end index
```