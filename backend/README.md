```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

```bash
go get backend/config
go get backend/db
go get backend/util
go get backend/router
go get backend/middleware
```

> add config.yaml

```bash
go run main.go
```