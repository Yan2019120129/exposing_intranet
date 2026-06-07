## 生成CRUD数据模型

### 在线工具

管理员身份运行后，访问：http://127.0.0.1:8010/admin

### 使用命令行工具

```
adm generate -l cn -c adm.ini
```
打包linux
```shell
GOOS=linux CGO_ENABLED=0  GOARCH=amd64 go build -ldflags="-s -w" -o ./build/linux_service main.go
```

文件上传命令
```shell
rsync -avz --delete -e "ssh" ./build/linux_service server:/usr/local/go_project/car/
```