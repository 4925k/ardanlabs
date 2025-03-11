# EP 1
- fmt runs before go run.
- run debugger during development
- `go build .` build binaries
    - default output matches the directory it is in
- `go run` will compile and run the program
- go binaries are static executables
    - they are compiled to machine level. no need for any environment.
    - can be cross compiled.
- `main()` does not accept parameters nor returns
- go is mostly written in go 
    - runtime has some parts in go assembly
        - go assembly converts to assembly based on machine

# EP 2
- string -> struct -> pointer to string + length of string
- utf8 -> most popular character encoding 
- byte (uint8)
- rune (int32)

# EP 3
- headers are case insensitive
- servers can respond based on whether the request generates from a user or a machine
    - users can get pretty json while machines may get compact json

# EP 4
- defer file.Close() -> to avoid hitting system file description limit
- variable shadowing
- io.Reader has one function Read()

# EP 5
- integers in go is values, not pointers

# EP 6
- zero vs missing value is hard to distinguish
    - cant tell if it was initialized to zero or user put zero
- local variables are stored on stack
- global variables are stored on heap
- ```go build -gcflags=-m``` to show memory usage

# EP 7
- bigger the interface, smaller the abstraction
- Rule of thumb: DONT USE ANY

# EP 8
- don't panic
- alawys comment regexp and maps with what they are used for
- zero vs nil values in maps by ```value, ok := map[key]```
- maps are hash maps not ordered

# EP 9
- using time.Sleep means you're doing something in a bad way
- channel semantics
    - send & recieve will block until opposite operation (*)
    - recieve from closed channel will return the zero value without blocking
    - sending to a closed channel will panic
    - closing a closed channel will panic
    - send/receive to a nil channel will block forever

# EP 10
- go routines and channels
- channels can be <-chan or chan<- or chan. 

# EP 11
- waitgroups
- sync.Once to run something only one time
- cost of race detection is very high

# EP 12
- contexts
- go routine leaks when context cancel is not being handled properly
    - use buffered channel to avoid go routine leaks
- use contexts on network requests to cancel them

# EP 13 = project engineering
- go project basics, docs
- look into .doc file of packages to see how to use docs for yourself
- examples and test
    - examples need to have ```Output: [expected]``` to match the expected output
    - test needs to have ```func TestName(t *testing.T)```
- go vendor to avoid dependency issues
- testify pacakge for making tests easier

# EP 14 - testing
- table testing 
    - toml for loading table test cases
- code structuring and modules

# EP 15 - project structure
- fuzz testing
- hey for load testing endpoints