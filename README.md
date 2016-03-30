# A go-ndn ping Program

NDN is named-data networking, and go-ndn is NDN implemented in go. But it still have no tools to ping ndn nodes. So, I write a ping program to test
the connection between two nodes.

Welcome to folk it and use go-ndn to do your research.

## Install
First ,your should install [go-ndn](https://github.com/go-ndn/example#step-0-install-go),and then build it.

```bash
go get github.com/qhsong/go-ndn-ping
cd $GOPATH/src/github.com/qhsong/go-ndn-ping
cd pingClient
go install 
cd ../pingServer
go install
```


## Server
Run it `./pingServer`
Params:
``` 
    -p Listen Path, default is /ndn/ping
    -s Nfd server, default is :6363
    -k Key path, default it key/default.pri
```
## Client
Run it `./pingClient`
```
    -i ping interval in seconds, default is 1.0
    -c count number,
    -n set starting number, default is random
    -p set sending path, default is /ndn/ping/
    -s set nfd Server, default is :6363
    -k set key path, default is key/default.pri
```

##News

###2016-03-25
- First Commit
- Can't got rtt static 
