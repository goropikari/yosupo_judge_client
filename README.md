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

$ yosupocl probinfo https://judge.yosupo.jp/problem/aplusb | jq .
{
  "title": "A + B",
  "source_url": "https://github.com/yosupo06/library-checker-problems/tree/master/sample/aplusb",
  "time_limit": 2,
  "version": "970b9dc1dca7858d5bb9d9b06ad79fd741a211c754a5b67b3448032906835138",
  "testcases_version": "c190c9571890cf3710f989430d14d54d73dbddebfe5d184bf87b5e687688e10c"
}

$ yosupocl download-test https://judge.yosupo.jp/problem/aplusb outdir
example_00.in
example_01.in
example_00.out
example_01.out


$ yosupocl submit https://judge.yosupo.jp/problem/aplusb sample/aplusb.cpp cpp
https://judge.yosupo.jp/submission/271165
```
