```

### [（看书）线程安全的map](./thread/SafeMap/SafeMap.go)
### [（习题）线程安全的slice](./thread/SafeSlice/safeSlice.go)
### [（习题）多线程的IMG tag生成](./thread/imagetag/imagetag.go)
### [（习题＊抄）多线程html tag：img 的width height 补全](./thread/imgFill/imgFill.go)
### [（看书）支持多种格式的发票读写程序](./io/invoice/)
### [（看书）压缩文件](./io/pack/pack.go)
### [（看书）解压文件](./io/unpack/unpack.go)
### [（习题）支持bz2的解压程序](./io/my_unpack/unpack.go)
### [（习题）将utf16编码的文件转码为utf8并写入文件](./io/utf16-to-utf8/utf16-to-utf8.go)
### [（习题）提取url](./pkg/my_linkutil/my_linkutil.go)
### [（习题＊抄）检查页面链接是否存在］(./pkg/my_linkcheck/run.go)
### [（看书）websocket聊天室](./chat/main.go)
### [（看书) 生成关联词](./sprinkle/main.go)
### [（看书）过滤生成合适的域名](./domainify/main.go)
### [（看书）随机删除或复制一个元音字符](./coolify/main.go)
### [（看书）查询同义词](./synonyms/main.go)
### [（看书）查询域名是否可用，显然这个服务器不可信了](./available/main.go)
### [（看书）把前面几个串起来，用于查询域名, 其实用bash脚本更方便明了](./domainfinder/main.go)
### [（看书）分布式投票程序,主要是利用NSQ](./socialpoll/twittervotes/main.go)
### [（看书）RESTful风格的http API程序](./socialpoll/api/main.go)
### [（看书）一个备份程序的数据库用来存储要备份的路径](./backup/cmds/backup/main.go)
### [（看书）自动备份程序，从上一个数据库中获取要备份的路径，进行备份](./backup/cmds/backupd/main.go)
### [（看书）微服务](./vault/cmd/vaultd/main.go)
### [ (看博客) 微服务 https://jacobmartins.com/2016/03/14/web-app-using-microservices-in-go-part-1-design/](./micro_service/main.go)

发现一个很危险的事：
slice 在截取的时候，和append的时候，是基于原数组的，譬如下面这些
    ar := []int{1,2,3,4,5,6,7,8}
    ar = ar[2:3]
    ar = append(ar, 1)
    fmt.Printf("len: %d cap: %d val: %v\n", len(ar), cap(ar), ar)  // len: 2 cap: 6 val: [3 1]
    ar = ar[:6]
    fmt.Printf("len: %d cap: %d val: %v\n", len(ar), cap(ar), ar)  // len: 6 cap: 6 val: [3 1 5 6 7 8]
```