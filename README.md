# running the server (local)
```bash
$ docker-compose up -d
```

# running unit tests
```bash
$ go test -v ./...
```

# running smoke test
```bash
$ docker-compose up -d
$ ./test/smoke.sh
```
