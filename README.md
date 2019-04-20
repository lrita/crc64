# CRC64

CRC64 variant with "Jones" which is used by Redis. Be optimized by slice-by-8.

# Benchmark

On CPU `2.7 GHz Intel Core i5`:
```
BenchmarkCRC64-4      	10000000	       203 ns/op	1474.16 MB/s
```
