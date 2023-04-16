# ChainCode-Electricity-Power
项目是基于该大佬的基础上进行修改，感谢大佬的分享：https://github.com/kevin-hf/education

> 项目同步到 gitee：https://gitee.com/KiceAmber/chaincode-power-trading.git

## 技术栈

-   Golang：后台逻辑开发
-   Vue3+Ts：前端界面开发
-   GoSDK：使用 Go 版本的超级账本开发链码

## 说明
需要编辑 `sudo vim /etc/hosts` 文件，添加下面的内容，保存退出
```
172.0.0.1 orderer.example.com
172.0.0.1 peer0.org1.example.com
172.0.0.1 peer1.org1.example.com
172.0.0.1 peer0.org2.example.com
172.0.0.1 peer1.org2.example.com
```

## 前置环境
- nodejs
- go 1.8.5
- docker
- docker-compose
- fabric v2.2 版本
- fabric-ca v1.4.7 版本

## 运行
### 后端及链码部分
```shell
# 需要将powerTrading项目移动到 ${GOPATH}/src/ 目录下，然后运行下面的命令
$ pwd 
${GOPATH}/src/powerTrading

$ ./clean_docker.sh # 运行后端以及链码部分
```

### 前端
```shell
$ pwd
${GOPATH}/src/chaincode-power-Trading/system

$ npm run dev # 运行前端项目
```

### 区块链浏览器
```shell
$ pwd
${GOPATH}/src/powerTrading/explorer

$ chmod u+x ./restart.sh
$ ./restart.sh
```

## 涉及到到加密算法
- 同态加密：https://github.com/tuneinsight/lattigo
- 环签名：https://github.com/zbohm/lirisi
- 零知识证明：https://pkg.go.dev/github.com/0xdecaf/zkrp#section-readme


