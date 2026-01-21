# serrors (service errors)

Package for creating errors with added information about the calling function.
All functions add a prefix in the format `{pkg.[type.]functionName.line}` to the error message.

## Functions

### Errorf

```go
func Errorf(format string, a ...any) error
```

Analog of `fmt.Errorf()`, but with added information about the calling function in the prefix.

### Error

```go
func Error(err error) error
```

Simplified version of `serrors.Errorf("%w", err)`. Wraps an existing error, adding information about the calling function.

### Join

```go
func Join(err1, err2 error) error
```

Analog of `errors.Join(err1, err2)`, but with added information about the calling function to `err2`.

### Joinf

```go
func Joinf(err error, format string, a ...any) error
```

Combines an existing error with a new formatted error, adding information about the calling function.
