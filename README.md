# 中数文联盟链Demo
## GO SDK 教程

### 1. 按照启动链教程启动本地联盟链
https://chaindocs.zhongshuwen.com/docs/intro
### 2. Clone + Go Get

```bash
git clone https://github.com/zhongshuwen/zswchain-go-demo
cd zswchain-go-demo
go get -u github.com/zhongshuwen/zswchain-go@1.1.0
cp ./.env.example.sh ./.env.sh
```
### 3. 编辑.env.sh
这个教程会需要四个密钥，一个你已经生成在启动链的时候（
```bash
#!/bin/sh
 # 这是docker-compose设置的“ZSW_ADMIN_PRIVATE_KEY” ENV VAR
export ZSW_CONTENT_REVIEW_PRIVATE_KEY=""

# 下面三个需要生成密钥（生成工具：https://chaintools.zhongshuwen.com/zsw-key-generator.html）
export KEXIN_JIEDIAN_A_PRIVATE_KEY="" 
export USER_A_PRIVATE_KEY=""
export USER_B_PRIVATE_KEY=""
```

### 4. Run！
```bash
./build-run-demo.sh
```