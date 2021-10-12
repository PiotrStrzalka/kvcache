# Key value cache project **kvcache**

-- PROJECT IN PROGRESS --

Contains four parts:
1. **kvcache** - core package, implements main functionality of storage, can be embedded into project
2. **client** - package provide way to connect to core kvcache running on different process/host
3. **cli** - tool which provides command line access to kvcache running on different process/host
4. **service** - standalone executable which provides kvcache through **rest api** or **grpc** 


protobuf generation:

```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/model/kvcache.proto
```

TODO:
```
add example
add dumping to filesystem
add grpc access type
add comments
add cache status on some endpoint
```
