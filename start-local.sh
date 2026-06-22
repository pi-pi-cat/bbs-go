对，那就不用脚本。前后端本地启动命令就是：

后端，在项目根目录：

```bash
go run -tags dev main.go
```

前端，另开一个终端：

```bash
cd web
corepack pnpm dev
```

第一次启动前可以先装依赖：

```bash
go mod download
cd web && corepack pnpm install --frozen-lockfile
```

前端地址：<http://localhost:3000/>  
后端地址：<http://127.0.0.1:8082>