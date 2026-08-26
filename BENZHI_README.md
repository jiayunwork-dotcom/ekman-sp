# ekman-sp：Go 海洋表层 Ekman 螺线核算服务（HTTP API + 命令行）

给定风速应力、海水密度、科氏参数 f（或纬度）与涡粘系数 K，计算表层流速大小与方向、Ekman 深度 De、深度积分输运 Me，并输出随深度转向的 Ekman 螺线剖面。

## 构建 / 运行 / 测试

```text
go build ./...
go run . serve :8080
curl -s http://127.0.0.1:8080/health
curl -s -X POST http://127.0.0.1:8080/api/profile -d '{"tau":0.2,"latitude_deg":45,"K":0.05}'
go run . profile example/midlat-wind.json
go test ./...
```

`example/midlat-wind.json` 为 45°N 中纬西风算例；另有南半球（`example/south-hemi-wind.json`）与热带信风（`example/north-trade.json`）算例。非法输入（密度≤0、K≤0、赤道 f=0、有限水深≪De 等）会返回错误。

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -d -P --name ekman-sp-b14 <image-name>:latest
curl -s http://127.0.0.1:$(docker port ekman-sp-b14 8080 | cut -d: -f2)/health
docker rm -f ekman-sp-b14
```
