# panicker

> pattern for handling panics in a preemptable application

## usage 
```
$ git clone git@github.com:frankgreco/panicker.git
$ cd panicker
$ PANICKER=YES go run ./...
2020/08/17 19:54:22 starting panicker
2020/08/17 19:54:22 starting parent process monitor
2020/08/17 19:54:22 starting Chuck printer
2020/08/17 19:54:22 starting Bass printer
Bass
Chuck
Bass
Chuck
Bass
Chuck
Bass
Chuck
Bass
Chuck
2020/08/17 19:54:27 recovered from panic: the panicker has panicked
2020/08/17 19:54:27 closing parent process monitor
2020/08/17 19:54:27 closing Chuck printer
2020/08/17 19:54:27 closing Bass printer
2020/08/17 19:54:27 closing panicker
```