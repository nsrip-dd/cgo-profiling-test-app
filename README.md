This is a simple app to test cgo-related profiling.
It records people's favorite colors, backed by an SQLite database.

To use it:

```
go build
./sqliteapp &
curl -X POST 'localhost:8765/update?name=bob&color=blue'
curl 'localhost:8765/get?name=bob'
```
