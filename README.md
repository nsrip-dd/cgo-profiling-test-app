This is a simple app to test cgo-related profiling.
It records people's favorite colors, backed by an SQLite database.

To use it:

```
go build
./sqliteapp &
curl -X POST 'localhost:8765/update?name=bob&color=blue'
curl 'localhost:8765/get?name=bob'
```

## Load testing

This can be tested with something like [`vegeta`](https://github.com/tsenart/vegeta), e.g.:

```
$ ./sqliteapp &
$ vegeta attack -connections 10 -rate=2000 -duration=30s -output vegeta.bin << EOF 
POST http://localhost:8765/update?name=foo&color=blue
POST http://localhost:8765/update?name=bar&color=green
POST http://localhost:8765/update?name=baz&color=purple
POST http://localhost:8765/update?name=foo&color=yello
POST http://localhost:8765/update?name=bar&color=red
POST http://localhost:8765/update?name=baz&color=gray
GET http://localhost:8765/get?name=foo
GET http://localhost:8765/get?name=bar
GET http://localhost:8765/get?name=baz
GET http://localhost:8765/get?name=foo
GET http://localhost:8765/get?name=bar
GET http://localhost:8765/get?name=baz
GET http://localhost:8765/get?name=foo
GET http://localhost:8765/get?name=bar
GET http://localhost:8765/get?name=baz
GET http://localhost:8765/get?name=foo
GET http://localhost:8765/get?name=bar
GET http://localhost:8765/get?name=baz
EOF
```

Here's some results I get:

* C allocation profiling disabled:
```
Requests      [total, rate, throughput]  60000, 2000.06, 2000.04
Duration      [total, attack, wait]      29.999449323s, 29.999107064s, 342.259µs
Latencies     [mean, 50, 95, 99, max]    857.807µs, 446.388µs, 2.764953ms, 7.823032ms, 26.041527ms
Bytes In      [total, mean]              1164392, 19.41
Bytes Out     [total, mean]              0, 0.00
Success       [ratio]                    100.00%
Status Codes  [code:count]               200:60000  
Error Set:
```
* C allocation enabled, at default 2MB sampling rate:
```
Requests      [total, rate, throughput]  60000, 1999.99, 1999.92
Duration      [total, attack, wait]      30.001148889s, 30.000122093s, 1.026796ms
Latencies     [mean, 50, 95, 99, max]    2.7608ms, 629.752µs, 9.338833ms, 39.267642ms, 576.20971ms
Bytes In      [total, mean]              1167971, 19.47
Bytes Out     [total, mean]              0, 0.00
Success       [ratio]                    100.00%
Status Codes  [code:count]               200:60000  
Error Set:
```