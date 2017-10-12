# selpg

> 课程《服务计算》作业二：用 Go 实现 [selpg](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html)

## Install

```
go get github.com/mensu/selpg
```

## Example

- 从输入文件 ``input_file`` 中选取第 10 ~ 20 页（每页 2 行）打印

```
$GOPATH/bin/selpg -s11 -e20 -l2 input_file
```

- 查看帮助

```
$GOPATH/bin/selpg -h
```
