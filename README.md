This is a simple app to test cgo-related profiling.
It records people's favorite colors, backed by an SQLite3 database.

To use it:

```
go build
./sqliteapp &
curl 'localhost:8765/update?name=bob&color=blue'
curl 'localhost:8765/get?name=bob'
```