# Key value cache project **kvcache**

-- PROJECT IN PROGRESS --

TODO:
```
add example
add some simple logging solution with levels
add prometheus integration 

add openapi and swagger <- more work required
add dumping to filesystem <- more work required
add support for concurrent writes and reads <- more work required
```

### Contains four parts:
1. **kvcache** - core package, implements main functionality of storage, can be embedded into project
2. **client** - package provide way to connect to core kvcache running on different process/host
3. **cli** - tool which provides command line access to kvcache running on different process/host
4. **service** - standalone executable which provides kvcache through **rest api** or **grpc** 

<br />
<br />
### **protobuf generation:**

```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/model/kvcache.proto
```


### **run commands:**
```
go run cmd/cli/main.go set -p grpc -k testkey -v testvalue



go test -run=XXX -bench=. -benchmem -cpuprofile profile2.out -memprofile memprofile2.out
go test -race -run=TestRaceDetector 
```