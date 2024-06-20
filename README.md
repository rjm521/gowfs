[![Build Status](https://drone.io/github.com/rjm521/gowfs/status.png)](https://drone.io/github.com/rjm521/gowfs/latest)

## gowfs 介绍
 gowfs 是一个通过 WebHDFS 接口为 Hadoop HDFS 提供 Go 封装的库。支持像操作文件系统一样操作HDFS
 gowfs 遵循 WebHDFS JSON 协议, 详见：https://hadoop.apache.org/docs/current/hadoop-project-dist/hadoop-hdfs/WebHDFS.html

### 快速使用
```
go get github.com/rjm521/gowfs
```
```go
import github.com/rjm521/gowfs
...
fs, err := gowfs.NewFileSystem(gowfs.Configuration{Addr: "localhost:50070", User: "hdfs"})
if err != nil{
	log.Fatal(err)
}
checksum, err := fs.GetFileChecksum(gowfs.Path{Name: "location/to/file"})
if err != nil {
	log.Fatal(err)
}
fmt.Println (checksum)
```

## API 介绍
gowfs 通过两个结构体 FileSystem 和 FsShell 让你访问 HDFS 资源。使用 FileSystem 进行低级别调用。FsShell 提供了更多更高级的文件操作

### FileSystem API 介绍
#### Configuration{} 结构体介绍
使用 Configuration{} 结构体来指定文件系统的参数。你可以使用 Configuration{} 字面量或使用 NewConfiguration() 来创建默认配置。

```
conf := *gowfs.NewConfiguration()
conf.Addr = "localhost:50070"
conf.User = "hdfs"
conf.ConnectionTime = time.Second * 15
conf.DisableKeepAlives = false
```

#### FileSystem{} 结构体介绍
在调用任何函数之前，创建一个新的 FileSystem{} 结构体。通过传递一个 Configuration 指针来创建 FileSystem，如下所示。
```
fs, err := gowfs.NewFileSystem(conf)
```

#### 创建文件
`FileSystem.Create()` 在 HDFS 服务器上创建并存储一个远程文件。详见: https://godoc.org/github.com/rjm521/gowfs#FileSystem.Create
```
ok, err := fs.Create(
    bytes.NewBufferString("Hello webhdfs users!"),
	gowfs.Path{Name:"/remote/file"},
	false,
	0,
	0,
	0700,
	0,
)
```

#### 打开 HDFS 文件
 `FileSystem.Open()` 打开并读取 HDFS 上的远程文件。 详见: https://godoc.org/github.com/rjm521/gowfs#FileSystem.Open
```
data, err := fs.Open(gowfs.Path{Name:"/remote/file"}, 0, 512, 2048)
...
rcvdData, _ := ioutil.ReadAll(data)
fmt.Println(string(rcvdData))

```

#### 追加到文件
 `FileSystem.Append()` 追加一些内容到已经存在的文件中  详见： https://godoc.org/github.com/rjm521/gowfs#FileSystem.Append
```
ok, err := fs.Append(
    bytes.NewBufferString("Hello webhdfs users!"),
    gowfs.Path{Name:"/remote/file"}, 4096)
```

#### 重命名文件
 使用 FileSystem.Rename() 重命名 HDFS 文件名（文件夹名） 详见： https://godoc.org/github.com/rjm521/gowfs#FileSystem.Rename
```
ok, err := fs.Rename(gowfs.Path{Name:"/old/name"}, Path{Name:"/new/name"})
```

#### 删除HDFS资源
删除 HDFS 资源 要删除 HDFS 资源（文件/目录），使用 FileSystem.Delete()。详见 See https://godoc.org/github.com/rjm521/gowfs#FileSystem.Delete
```go
ok, err := fs.Delete(gowfs.Path{Name:"/remote/file/todelete"}, false)
```

#### 查看文件状态
你可以使用 FileSystem.GetFileStatus() 获取现有 HDFS 资源的状态。详见 https://godoc.org/github.com/rjm521/gowfs#FileSystem.GetFileStatus

```go
fileStatus, err := fs.GetFileStatus(gowfs.Path{Name:"/remote/file"})
```
gowfs 返回一个类型为 FileStatus 的值，这是一个包含远程文件信息的结构体。
```go
type FileStatus struct {
	AccesTime int64
    BlockSize int64
    Group string
    Length int64
    ModificationTime int64
    Owner string
    PathSuffix string
    Permission string
    Replication int64
    Type string
}
```
你可以使用 FileSystem.ListStatus() 获取文件状态列表。
```go
stats, err := fs.ListStatus(gowfs.Path{Name:"/remote/directory"})
for _, stat := range stats {
    fmt.Println(stat.PathSuffix, stat.Length)
}
```
### FsShell 使用示例
#### 创建一个FsShell
要创建 FsShell，需要有一个已经创建好的 FileSystem 实例。FsShell会依赖到它
```go
shell := gowfs.FsShell{FileSystem:fs}
```
#### FsShell.Put()
使用 put 将本地文件上传到 HDFS 文件系统。 详见 https://godoc.org/github.com/rjm521/gowfs#FsShell.PutOne
```go
ok, err := shell.Put("local/file/name", "hdfs/file/path", true)
```
#### FsShell.Get()
使用 Get 将远程 HDFS 文件下载到本地文件系统。 详见 https://godoc.org/github.com/rjm521/gowfs#FsShell.Get
```go
ok, err := shell.Get("hdfs/file/path", "local/file/name")
```

#### FsShell.AppendToFile()
将本地文件追加到远程 HDFS 文件或目录。 详见 https://godoc.org/github.com/rjm521/gowfs#FsShell.AppendToFile
```go
ok, err := shell.AppendToFile([]string{"local/file/1", "local/file/2"}, "remote/hdfs/path")
```

#### FsShell.Chown()
更改远程文件的所有者。  详见 https://godoc.org/github.com/rjm521/gowfs#FsShell.Chown.
```go
ok, err := shell.Chown([]string{"/remote/hdfs/file"}, "owner2")
```

#### FsShell.Chgrp()
更改远程 HDFS 文件的组。  详见 https://godoc.org/github.com/rjm521/gowfs#FsShell.Chgrp
```go
ok, err := shell.Chgrp([]string{"/remote/hdfs/file"}, "superduper")
```

#### FsShell.Chmod()
更改远程 HDFS 文件的权限。 详见 https://godoc.org/github.com/rjm521/gowfs#FsShell.Chmod
```go
ok, err := shell.Chmod([]string{"/remote/hdfs/file/"}, 0744)
```


### 参考资料
1. WebHDFS API - http://hadoop.apache.org/docs/current/hadoop-project-dist/hadoop-hdfs/WebHDFS.html
2. FileSystemShell - http://hadoop.apache.org/docs/current/hadoop-project-dist/hadoop-common/FileSystemShell.html#getmerge
