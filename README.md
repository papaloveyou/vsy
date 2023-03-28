# vsy

```shell
/proc/cpuinfo
/proc/diskstats
/proc/meminfo
/proc/stat
/proc/swaps
/proc/uptime
/proc/slabinfo
/sys/devices/system/cpu/online
```

```text
Filesystem     1MB-blocks     Used Available Use% Mounted on
/dev/sda1        949037MB 154611MB  746142MB  18% /usr/src/app/shared
```

## usage
```shell
docker run -it --name=test \
      -v /var/lib/vsy/meminfo:/proc/meminfo 
      debian /bin/bash
```

```shell
docker cp saturn-node:/bin/df /tmp/df0
docker cp /tmp/df0 saturn-node:/bin/
docker cp /tmp/df saturn-node:/bin/
```







