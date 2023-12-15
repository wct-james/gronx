# wct-james/gronx
Forked from [adhocore/gronx](https://github.com/wct-james/gronx)

[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE)

> The reason for this fork is to make the implementation granular to minutes, rather than seconds. 
> 
> This is so that the package aligns more closely to the [Posix crontab specification](https://pubs.opengroup.org/onlinepubs/9699919799/utilities/crontab.html)

`gronx` is Golang [cron expression](#cron-expression) parser ported from [adhocore/cron-expr](https://github.com/adhocore/php-cron-expr). You can also use gronx to find the next (`NextTick()`) or previous (`PrevTick()`) run time of an expression from any arbitrary point of time.

- Zero dependency.
- Very **fast** because it bails early in case a segment doesn't match.


Find gronx in [pkg.go.dev](https://pkg.go.dev/github.com/wct-james/gronx).

## Installation

```sh
go get -u github.com/wct-james/gronx
```

## Usage

```go
import (
	"time"

	"github.com/wct-james/gronx"
)

gron := gronx.New()
expr := "* * * * *"

// check if expr is even valid, returns bool
gron.IsValid(expr) // true

// check if expr is due for current time, returns bool and error
gron.IsDue(expr) // true|false, nil

// check if expr is due for given time
gron.IsDue(expr, time.Date(2021, time.April, 1, 1, 1, 0, 0, time.UTC)) // true|false, nil
```

### Batch Due Check

If you have multiple cron expressions to check due on same reference time use `BatchDue()`:
```go
gron := gronx.New()
exprs := []string{"* * * * *", "0 */5 * * * *"}

// gives []gronx.Expr{} array, each item has Due flag and Err enountered.
dues := gron.BatchDue(exprs)

for _, expr := range dues {
    if expr.Err != nil {
        // Handle err
    } else if expr.Due {
        // Handle due
    }
}

// Or with given time
ref := time.Now()
gron.BatchDue(exprs, ref)
```

### Next Tick

To find out when is the cron due next (in near future):
```go
allowCurrent = true // includes current time as well
nextTime, err := gron.NextTick(expr, allowCurrent) // gives time.Time, error

// OR, next tick after certain reference time
refTime = time.Date(2022, time.November, 1, 1, 1, 0, 0, time.UTC)
allowCurrent = false // excludes the ref time
nextTime, err := gron.NextTickAfter(expr, refTime, allowCurrent) // gives time.Time, error
```

### Prev Tick

To find out when was the cron due previously (in near past):
```go
allowCurrent = true // includes current time as well
prevTime, err := gron.PrevTick(expr, allowCurrent) // gives time.Time, error

// OR, prev tick before certain reference time
refTime = time.Date(2022, time.November, 1, 1, 1, 0, 0, time.UTC)
allowCurrent = false // excludes the ref time
nextTime, err := gron.PrevTickBefore(expr, refTime, allowCurrent) // gives time.Time, error
```

> The working of `PrevTick*()` and `NextTick*()` are mostly the same except the direction.
> They differ in lookback or lookahead.

### Standalone Daemon

In a more practical level, you would use this tool to manage and invoke jobs in app itself and not
mess around with `crontab` for each and every new tasks/jobs.

In crontab just put one entry with `* * * * *` which points to your Go entry point that uses this tool.
Then in that entry point you would invoke different tasks if the corresponding Cron expr is due.
Simple map structure would work for this.

Check the section below for more sophisticated way of managing tasks automatically using `gronx` daemon called `tasker`.

---
### Cron Expression

Cron expression is most commonly made of 5 segments. These segments are interpreted as:
```
<minute> <hour> <day> <month> <weekday>
```

In a 6 segments expression, if 6th segment matches `<year>` (i.e 4 digits at least) it will be interpreted as:
```
<minute> <hour> <day> <month> <weekday> <year>
```

For each segments you can have **multiple choices** separated by comma:
> Eg: `0,30 * * * *` means either 0th or 30th minute.

To specify **range of values** you can use dash:
> Eg: `10-15 * * * *` means 10th, 11th, 12th, 13th, 14th and 15th minute.

To specify **range of step** you can combine a dash and slash:
> Eg: `10-15/2 * * * *` means every 2 minutes between 10 and 15 i.e 10th, 12th and 14th minute.

For the `<day>` and `<weekday>` segment, there are additional [**modifiers**](#modifiers) (optional).

And if you want, you can mix the multiple choices, ranges and steps in a single expression:
> `5,12-20/4,55 * * * *` matches if any one of `5` or `12-20/4` or `55` matches the minute.

### Real Abbreviations

You can use real abbreviations (3 chars) for month and week days. eg: `JAN`, `dec`, `fri`, `SUN`

### Tags

Following tags are available and they are converted to real cron expressions before parsing:

- *@yearly* or *@annually* - every year
- *@monthly* - every month
- *@daily* - every day
- *@weekly* - every week
- *@hourly* - every hour
- *@5minutes* - every 5 minutes
- *@10minutes* - every 10 minutes
- *@15minutes* - every 15 minutes
- *@30minutes* - every 30 minutes
- *@always* - every minute

```go
// Use tags like so:
gron.IsDue("@hourly")
gron.IsDue("@5minutes")
```

### Modifiers

Following modifiers supported

- *Day of Month / 3rd of 5 segments / 4th of 6+ segments:*
    - `L` stands for last day of month (eg: `L` could mean 29th for February in leap year)
    - `W` stands for closest week day (eg: `10W` is closest week days (MON-FRI) to 10th date)
- *Day of Week / 5th of 5 segments / 6th of 6+ segments:*
    - `L` stands for last weekday of month (eg: `2L` is last monday)
    - `#` stands for nth day of week in the month (eg: `1#2` is second sunday)

---
## License

> &copy; [MIT](./LICENSE) | 2021-2099 Will James

## Credits

This project is forked from [adhocore/gronx](https://github.com/adhocore/gronx)
