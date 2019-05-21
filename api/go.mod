module github.com/EnsicoinDevs/ensicoin-explorer/api

go 1.12

require (
	github.com/EnsicoinDevs/eccd v0.0.0-20190519221937-361dc6f1a950
	github.com/gin-gonic/gin v1.4.0
	github.com/golang/protobuf v1.3.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.3
	go.etcd.io/bbolt v1.3.2
	golang.org/x/net v0.0.0-20190514140710-3ec191127204
	google.golang.org/grpc v1.20.1
)

replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43
