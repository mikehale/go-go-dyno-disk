go-go-dyno-disk
===============

Reads /proc/diskstats inside a heroku dyno and emits logs which can be
drained to l2met.

Example:

```bash
source=$APP-$DEPLOY.$DYNO count#xvda2.inflight=123 count#xvda2.weighted-io-time=456
```
