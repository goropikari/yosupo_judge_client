# yosupo_judge_client

```sh
go install github.com/goropikari/yosupo_judge_client/cmd/yosupocl
```

```sh
aws s3 ls --recursive --endpoint-url http://localhost:9000 s3://testcase-public/v3/aplusb/testcase/

2025-02-02 15:23:56         10 v3/aplusb/testcase/c190c9571890cf3710f989430d14d54d73dbddebfe5d184bf87b5e687688e10c/in/example_00.in
2025-02-02 15:23:56         22 v3/aplusb/testcase/c190c9571890cf3710f989430d14d54d73dbddebfe5d184bf87b5e687688e10c/in/example_01.in
2025-02-02 15:23:56          5 v3/aplusb/testcase/c190c9571890cf3710f989430d14d54d73dbddebfe5d184bf87b5e687688e10c/out/example_00.out
2025-02-02 15:23:56         11 v3/aplusb/testcase/c190c9571890cf3710f989430d14d54d73dbddebfe5d184bf87b5e687688e10c/out/example_01.out
```

```sh
buf generate

yosupocl probinfo '{"name":"aplusb"}' | jq .
yosupocl submit '{"problem":"aplusb", "source": "long long solve(long long a, long long b) { return a+b;}", "lang":"cpp-func"}' | jq .
yosupocl download-test '{"name":"many_aplusb", "outdir": "test"}'

go run cmd/yosupocl/main.go probinfo '{"name":"aplusb"}' | jq .
go run cmd/yosupocl/main.go submit '{"problem":"aplusb", "source": "long long solve(long long a, long long b) { return a+b;}", "lang":"cpp-func"}' | jq .
go run cmd/yosupocl/main.go download-test '{"name":"many_aplusb", "outdir": "test"}'
```
