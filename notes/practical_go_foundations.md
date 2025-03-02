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