# portscan

A simple portscan tool written in golang.

# Installation

```
go get github.com/hiroakis/portscan
```

# Usage

* scan tcp ports from 1-1024 on local machine

```
portscan
```

* scan tcp ports from 80-11211 on remote machine

```
portscan -host=REMOTE_MACHINE -lower=80 -upper=11211
```

If you get `too many open files` error freaquentry, increase file descriptors using `ulimit -n 65536`.

# LICENSE

MIT.