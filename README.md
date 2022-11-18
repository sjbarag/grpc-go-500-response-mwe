# HTTP 500 Response for Bad Requests instead of 400

## Description

[grpc-go](https://github.com/grpc/grpc-go) (aka `google.golang.org/grpc`) returns an HTTP 500 Internal Server Error in response to malformed HTTP requests when used as an HTTP handler, when HTTP 4XX responses more- closely follows the HTTP spec as defined by [RFC9110](https://httpwg.org/specs/ rfc9110.html#status.4xx). More practically, returning HTTP 500 artificially increases the number of server-side errors reported via most monitoring tools, since these aren't actually errors.

While only one case is demonstrated here, there are several HTTP request conditions that incorrectly return HTTP 500 responses:

1. Requests made with any protocol other than HTTP/2 (e.g. HTTP/1.1 or HTTP/3)
2. Requests with any method other than `POST` (e.g. `GET` or `HEAD`)
3. Requests with any `content-type` header value that doesn't include the content type `application/grpc` (e.g. `text/plain` or missing headers)

Respectively, these are most appropriately modeled as:

1. [HTTP 426 Upgrade Required](https://httpwg.org/specs/rfc9110.html#status.426)
2. [HTTP 405 Method Not Allowed](https://httpwg.org/specs/rfc9110.html#status.405)
3. [HTTP 415 Unsupported Media Type](https://httpwg.org/specs/rfc9110.html#status.415)

but a general [HTTP 400 Bad Request](https://httpwg.org/specs/rfc9110.html#status.400) in any of these cases appears semantically appropriate.

## Reproducing
1. Clone this repo
2. Start the demo Greeter service (modeled after the one in [grpc-go's examples](https://github.com/grpc/grpc-go/tree/master/examples)):
```sh
$ go run main.go
2022/11/17 17:07:36 server listening on :55123
```
3. In another terminal, make an HTTP-only healthcheck request:
```sh
$ curl --include localhost:55123/health
HTTP/1.1 200 OK
Date: Fri, 18 Nov 2022 01:07:58 GMT
Content-Length: 2
Content-Type: text/plain; charset=utf-8

ok
```
4. Then make an HTTP/1.1 request to the gRPC Greeter service:
```sh
$ curl --http1.1 --include localhost:55123/greet/foo
HTTP/1.1 500 Internal Server Error
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 18 Nov 2022 01:09:29 GMT
Content-Length: 21

gRPC requires HTTP/2
```

While the "gRPC requires HTTP/2" message is expected here, an HTTP/1.1 500 Internal Server Error response code is not.
