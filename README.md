# TSDB Tester

This is a golang implementation of `TSDB Test Framework`, ref to [TiDB MySQL Tester](https://github.com/pingcap/mysql-tester).

## Requirements

- All the tests should be put in [`t`](./t), take [t/metric_meta.test](./t/metric_meta.test) as an example.
- All the expected test results should be put in the same path/file name with the corresponding test file, but with a `.result` file suffix, take [t/metric_meta.result](./t/metric_meta.result) as an examle.

## How to use

Build the `tsdb-tester` binary:
```sh
make build
```

Basic usage:
```
Usage of ./tsdb-tester:
  -url string
        The api url of the LinDb broker. (default "http://127.0.0.1:9000/api/v1/exec")
```

By default, it connects to the LinDB broker at `127.0.0.1:9000`:

```sh
./bin/tsdb-tester # run all the tests
```

Test file:

```sql
# setting test database name
@use _internal;

# test tag key not found(test case desc)
select * from cpu group by host;
```

Result:

```sh
=== RUN   group/group
=== RUN   group/group/test tag key not found
--- PASS: group/group (0.00 sec)
    --- PASS: group/group/test tag key not found (0.00 sec)
ok ./t/group/group.test 
=== RUN   metric_meta
=== RUN   metric_meta/show metrics
=== RUN   metric_meta/select * from cpu ...
=== RUN   metric_meta/select * from cpu
--- FAIL: metric_meta (0.01 sec)
    --- FAIL: metric_meta/show metrics (0.00 sec)
	Statement: show metrics
	Not equal: 
	   expected: "{}"
	   actual  : "{\"type\":\"metric\",\"values\":[\"cpu\",\"lindb.broker.family.writete\",\"lindb.concurrent.limit\",\"lindb.concurrent.limitite\",\"lindb.concurrent.pool\",\"lindb.concurrent.pool.tasks_waiting_durationon\",\"lindb.coordinator.state_managerting_duration\",\"lindb.http.ingest_durationnager\",\"lindb.ingestion.flatration\",\"lindb.ingestion.infl\",\"lindb.ingestion.proto\",\"lindb.ingestion.protox\",\"lindb.kv.table\",\"lindb.kv.table.cacheion\",\"lindb.kv.table.reade\",\"lindb.kv.table.writ\",\"lindb.master.control\",\"lindb.master.shard.lead\",\"lindb.monitor.native_push\",\"lindb.monitor.system.cpu_st\",\"lindb.monitor.system.disk_ino\",\"lindb.monitor.system.disk_usage_statss\",\"lindb.monitor.system.mem_statge_stats\",\"lindb.monitor.system.net_stat\",\"lindb.query\",\"lindb.queryor.system.net_stat\",\"lindb.storage\",\"lindb.storage.que\",\"lindb.storage.repli\",\"lindb.storage.replicator.ru\",\"lindb.storage.wallicator.runner\",\"lindb.task.transp\",\"lindb.task.transport\",\"lindb.traffic.grpc_client.stream\",\"lindb.traffic.grpc_client.stream.sent_duration\",\"lindb.traffic.grpc_client.stream.sent_durationtion\",\"lindb.traffic.grpc_server\",\"lindb.traffic.grpc_server.stream\",\"lindb.traffic.grpc_server.stream.sent_durationtion\",\"lindb.traffic.tcp\",\"lindb.traffic.tcpc_server.stream.sent_duration\",\"lindb.tsdb.matabase\",\"lindb.tsdb.memdb\",\"lindb.tsdb.memdbase.metadb_flush_duration\",\"lindb.tsdb.shard\",\"lindb.tsdb.shard.memdb_flush_duration\",\"lindb.tsdb.shard.memdb_flush_durationon\"]}"
    --- PASS: metric_meta/select * from cpu ... (0.00 sec)
    --- FAIL: metric_meta/select * from cpu (0.00 sec)
	Statement: select * from cpu
	Not equal: 
	   expected: "[\"ddd\"]"
	   actual  : "{\"metricName\":\"cpu\",\"startTime\":1702514340000,\"endTime\":1702517940000,\"interval\":10000}"
fail ./t/metric_meta.test 
Summary: total 4, success 2, fail 2
```

## License

TSDB Tester is under the Apache 2.0 license. See the [LICENSE](./LICENSE) file for details.
