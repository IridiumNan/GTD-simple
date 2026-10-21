# Design

## Gen

A counter which record the generation by uint64 is better then a bool flag (maybe `flushed`).

As we a writing data into disk and new write request coming. The status is vague for `flushed`.

For generation counter, when a write request reach, it will wait for the write lock then guarantee gen++ which ensure all change status record.

For this system store, it hold three gen counter

```go
type JSONStore struct {
    // Other fields
    // ...

    gen uint64 // +1 when RAM data updated
    flushedGen uint64 // then last flushed gen
    statGen uint64 // the gen for last stat cache
}
```

When reach the flush interval, check by `if flushedGen != gen`. So that we can reduce the flush count while ensuring the data safety.
