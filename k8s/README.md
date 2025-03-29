## 创建私有镜像的密钥
```shell
kubectl create secret docker-registry regcred \
  --docker-server=your-registry.example.com \
  --docker-username=your-username \
  --docker-password=your-password \
  --docker-email=your-email@example.com
```