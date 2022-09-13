# RedisPrefixFilter
The RedisPrefixFilter module provides the prefix filter data structure.
PF is used to determine, with a high degree of certainty, whether an element is  a member of a set.



## Prerequisites
- Install Python 3.7
- For building the prefix filter, install the Prerequisites mentioned [here.](https://github.com/TomerEven/Prefix-Filter#prerequisites)
- Get redis server with the bloom filter module. For example, you can install from [here.](https://hub.docker.com/r/redislabs/rebloom/)
- Install golang 1.18 (for running the benchmarks)

## Installation
### Building and Loading RedisPrefixFilter

- clone:
```
$ git clone https://github.com/eshelyaron/RedisPrefixFilter
```
- make:
```
$ cd RedisPrefixFilter/
$ /RedisPrefixFilter$ make
```
- verify module.so created:
```
$ /RedisPrefixFilter$ ls -la module.so
-rwxr-xr-x 1 user user 893384 Sep 13 20:06 module.so
```
### Use RedisPrefixFilter with redis-cli
- Run cli command:
```
$ /RedisPrefixFilter$ redis-cli 
```

- Load the module:
```
127.0.0.1:6379> module load module.so
OK
```

- Create a new prefix  filter:
```
127.0.0.1:6379> pf.reserve my_table 1024
OK
```

- Add items to the filter:
```
127.0.0.1:6379> pf.madd my_table foo bar baz spam
1) (integer) 1
2) (integer) 1
3) (integer) 1
4) (integer) 1
```

- Find out whether the items exists in the filter:
```
127.0.0.1:6379> pf.mexists my_table foo bar what spam
1) (integer) 1
2) (integer) 1
3) (integer) 0
4) (integer) 1
```
## Benchmarks
To benchmark prefix filter against Redis Bloom and cuckoo filters, run the following commands:
```
$ cd benchmark 
$ go build # build benchmark executable
$ ./benchmark 
$ cd visualisations/
$ python3 main.py # generate graphs
$ cd ../results/
$ ls -la *.png 
-rw-r--r-- 1 user user 31444 Sep 13 22:01 testExistsPerNumberOfParalleledTests.png
-rw-r--r-- 1 user user 28756 Sep 13 22:01 testMAddPerNumberOfItems.png
-rw-r--r-- 1 user user 27084 Sep 13 22:01 testMAddPerNumberOfParalleledTests.png
-rw-r--r-- 1 user user 31073 Sep 13 22:01 testMExistsPerNumberOfItemsAlwaysNegative.png
-rw-r--r-- 1 user user 33571 Sep 13 22:01 testMExistsPerNumberOfItems.png
```











