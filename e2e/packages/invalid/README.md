# Ruby Logging Parser

This is a simple package containing a parser for parsing logs produced with the
[Ruby core library Logger](https://ruby-doc.org/stdlib-2.4.0/libdoc/logger/rdoc/Logger.html).

## Ruby on Rails

If you are using Ruby on Rails, we suggest you checkout the `humio/rails` package instead. It
has more sophisticated dashboards and queries tailored to Rails environments. 

## Package Content

This package contains:

_Alerts_
- `Errors` - An alert which triggers whenever an error occurs.
- `Exceptions` - An alert which triggers whenever an exception is thrown.

_Dashboards_
- `Monitoring` - A dashboard for monitoring uptime and errors

_Parsers_
- `ruby-logger` - A parser for parsing logs and stacktraces for logger.rb in the Ruby Std-lib.

## Well-defined fields

The parser produces the following well-defined fields.

- `#level` - INFO, ERROR, FATAL, DEBUG, UNKNOWN, WARN
- `msgtype` - exception, log, unknown
- `pid` - The PID of the process
- `progname` - The process name set in the Ruby program
- `message` - The log message

For exceptions, it also produces the following fields:

- `exception` - The name of the thrown exception (e.g. ArithmeticException)
- `linenum` - The line number from the stacktrace
- `file` - The name from the stacktrace
- `function` - The name of the function from the stacktrace

## Example Setup

This example assumes you are using filebeat to send your logs, we will

```yaml
filebeat.inputs:
- paths:
    - $PATH_TO_LOG_FILES/example.log
    - $PATH_TO_LOG_FILES/stderr.log
    - $PATH_TO_LOG_FILES/stdout.log
  encoding: utf-8

queue.mem:
  events: 8000
  flush.min_events: 100
  flush.timeout: 1s

output:
  elasticsearch:
    hosts: ["$HUMIO_SERVER/api/v1/ingest/elastic-bulk"]
    username: anything
    password: $INGEST_TOKEN
    compression_level: 5
    bulk_max_size: 200
    worker: 5
```

You should refer to the instructions at https://docs.humio.com/ref/filebeat for
more details on getting ingest to work.

Where is a simple example Ruby application that works with this package.

```ruby
require "logger"

logger = Logger.new File.new('example.log', 'w')
$stderr = File.new("stderr.log", "w")
$stdout = File.new("stdout.log", "w")

def throwException(x, logger)
    deepStacktrace(x, logger)
end

def deepStacktrace(x, logger)
    raise ArgumentError, 'Argument is not numeric' unless x.is_a? Numeric
rescue ArgumentError => error
    puts error.backtrace
end

# Generate some log output
i = 1
while(true) do
    logger.info("I am info. Controller: MySuperController User: Anders")

    if i % 2 == 0 then
        logger.warn("I am a warning. Foo: Bar")
    end

    if i % 3 == 0 then
        logger.error("I am an error. Controller: MySuperController User: Peter")
    end

    if i % 100 == 0 then
        throwException("dd", logger)
    end

    sleep 0.001
    i+=1
end
```

## Log Output

This package assumes you use the default format for Ruby logger,
which is as follows:

```
SeverityID, [DateTime #pid] SeverityLabel -- ProgName: message
```

If you change the logging format, you will need to modify the parser.


### Log Conventions

Many Ruby applications output key-value pairs separated to use colon like in this example:

```
UserId: 2133 Action: Purchase
```

If you use another separator, say `=` you should make sure to modify your parser
where `kvParse(seperator=":")` is used.


