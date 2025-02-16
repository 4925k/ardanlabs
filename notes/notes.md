# PRACTICAL GO FOUNDATIONS

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
