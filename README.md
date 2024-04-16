# RandomName
Restiki

## Execution

- Run this command in terminal
```Go
go run main.go
```
- Used a Struct for S, N, BOOL, NULL, M & L to unmarshall JSON input.
- Created a map[string]*interface{} for marshalling Output JSON in stdout. 
- Used 'omitempty' in all Struct to ignore empty.
- Used 'TrimSpace' to sanitize all the values of trailing and leading whitespace before processing.
- Used float64 in parsing integers and ignored errors. 
- Used Unix in parsing time and ignored errors.