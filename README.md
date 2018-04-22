# running the server (local)
```bash
$ docker-compose up -d
```

# running unit tests
```bash
$ docker-compose run --rm app go test -v ./...
```

# running smoke test
```bash
$ docker-compose up -d
$ ./test/smoke.sh
```
