lager
=====

lager is an opinionated logging library :beer:

## Why?

lager only performs string formatting for log levels that you set.
If you only want to log `Error` messages, you don't need to worry about all of your
`Trace` or `Debug` strings being (costly) formatted.

Here's how the `LogLager` compares to `log.Logger` when you want to only log one level
in your program instead of all.

```
BenchmarkLoggerOneAndAllLevels	  200000	      7842 ns/op	     605 B/op	      30 allocs/op
BenchmarkLogLagerOneLevel	 1000000	      2258 ns/op	     250 B/op	      15 allocs/op
BenchmarkLogLagerAllLevel	  200000	      7964 ns/op	     602 B/op	      35 allocs/op
```

## Usage

The main interface that defines a logger is `Lager`:

```
type Lager interface {
	Tracef(msg string, v ...interface{})
	Debugf(msg string, v ...interface{})
	Infof(msg string, v ...interface{})
	Warnf(msg string, v ...interface{})
	Errorf(msg string, v ...interface{})
}
```

`Lager` provides logging for five levels: `Trace`, `Debug`, `Info`, `Warn`, and `Error`.

Currently, `lager` provides three `Lager` implementations:
- `LogLager`: A performant logger for logging directly to `log.Logger`
- `BasicLager`: A basic logger for log messages with customizable output
- `ContextLager`: A context logger for adding context to log messages with customizable output

Log outputs are defined by `Drinker`:
```
type Drinker interface {
	Drink(e interface{}) error
}
```

Currently, there are two `Drinker` implementations:
- `LogDrinker`: logs messages using `log.Logger`
- `JSONDrinker`: logs messages using `json.Marshal`

For more usage, see the tests and benchmarks.
