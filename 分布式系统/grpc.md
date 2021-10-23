## grpc是什么？
## grpc相比http的优势是什么？
## 如何使用grpc？
1. 安装plugins
```
go get -u google.golang.org/protobuf/cmd/protoc-gen-go
go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
```
0. buf 用于代替 protoc 进行生成代码，可以避免使用复杂的 protoc 命令，避免 protoc 各种失败问题
1. 安装buf `https://docs.buf.build/installation/`
2. 安装问题
```
如果遇到timeout
执行vim /opt/homebrew/Library/Taps/bufbuild/homebrew-buf/Formula/buf.rb
编辑，在下面代码后添加两行
ENV["GOPATH"] = HOMEBREW_CACHE/"go_cache"
天剑
ENV['GOPROXY'] = 'https://goproxy.cn'
ENV['GO111MODULE'] = 'on'
```
3. buf参数解释
breaking
lint
build
## 如何写proto文件
`https://colobu.com/2019/10/03/protobuf-ultimate-tutorial-in-go/#proto%E6%95%99%E7%A8%8B`   


参考：
https://www.cnblogs.com/hacker-linner/p/14618862.html