# bitcask

Bitcask 是一种日志结构哈希表的键值存储引擎，最初为 Riak 分布式数据库设计，具有通用性。其核心目标是实现低延迟读写、高吞吐量（尤其随机写入）、支持超内存数据集、崩溃友好（快速恢复且不丢数据）等

## 核心设计

Bitcask 实例对应一个目录, 仅允许一个进程打开该目录进行写入

- 文件管理
  - 目录中仅有一个`活跃文件`(active data file), 用于追加写入
  - 当`活跃文件`大小到达阈值, 会被关闭并创建新的活跃文件
  - 关闭后的文件(无论是主动关闭还是自动关闭)变为immutable(不可变), 不再进行写入
  - 数据目录格式

| crc | tstamp | ksz | value_sz | key | value |
| --- | ------ | --- | -------- | --- | ----- |
| 校验和 | 时间戳 | 键的大小 | 值的大小 | 键 | 值(删除为墓碑值) |

- keydir
  内存中的哈希表, 映射每个键到最近数据的元信息

| file_id | value_sz | value_pos | tstamp |
| ------- | -------- | --------- | ------ |
| 文件ID | 值大小 | 值位置 | 时间戳 |

keydir写入时原子更新, 确保读取时直接获取到最新数据位置, 避免扫盘; 确保读取仅需 1 次磁盘寻道.

- 读, 写, 合并过程
  - 读取过程: 在keydir中查询键对应的文件ID, 键位置和大小, 基于该元信息直接读取磁盘数据
  - 写入过程: 将键值条目追加到活跃文件, 原子更新keydir, 记录该键值对应的最新数据位置; 旧数据依旧在磁盘, 但其不会再被读取
  - 合并操作: 将清理掉 immutable(不可变) 文件中的旧数据和墓碑值, 仅保留每个键的最新版本, 产生新的合并数据文件, 以及hint文件(记录元信息, 加速后续启动)

<!--

目标:
性能：早期测试中，笔记本慢磁盘上随机写入吞吐量达5000-6000 次 / 秒，延迟中位数低于 1 毫秒；读取依赖 OS 缓存，效率高。
超内存支持：测试中数据集为内存的 10 倍以上，性能无退化。
崩溃恢复：数据文件即日志，无需回放；hint 文件可加速启动。
备份恢复：文件 immutable，支持系统级备份；恢复仅需将数据文件放入目标目录。
简洁性：代码和数据格式简单，易于理解和维护
 -->


```txt
/go-bitcask(待修改)
  ├── cmd
  │   └── go-bitcask-cli            # 命令行工具入口（可选）
  ├── internal
  │   ├── fio                   # 文件 I/O 管理
  │   │   ├── fio.go            # 文件操作接口定义
  │   │   ├── mmap.go           # MMAP 内存映射实现
  │   │   ├── asyncio.go       # 异步文件IO实现
  │   │   ├── manager.go        # 管理 ActiveFile 和 OldFiles
  │   │   └── sequential.go     # 顺序文件IO实现
  │   ├── index                 # 内存索引管理
  │   │   ├── btree.go          # B树索引实现
  │   │   ├── index.go          # 内存操作接口定义
  │   │   ├── shardmap.go      # 分片哈希表索引实现
  │   │   └── loader.go         # 数据文件索引重建
  │   ├── record                # 数据记录结构
  │   │   ├── record.go         # 记录结构定义 文件记录结构和内存索引定义
  │   │   └── validator.go      # CRC校验逻辑
  │   ├── codec                 # 编码解码
  │   │   ├── encoder.go        # 二进制编码
  │   │   └── decoder.go        # 二进制解码
  │   ├── merger                # 数据合并模块
  │   │    └── compact.go        # 合并旧文件，清理失效数据
  │   ├── keydir.go                 # 索引操作
  │   └── datafile.go             # 数据文件操作（Put/Get/Delete）
  ├── tests                     # 单元测试和压力测试
  ├── betch.go // 批处理
  ├── bitcask.go // 引擎初始化，对外暴露接口
  ├── errors.go // 错误定义
  ├── hint.go
  ├── iter.go // 迭代器
  ├── merge.go // 合并逻辑
  ├── options.go // 引擎配置
  ├── go.mod
  └── README.md
```


## ShardMap

ShardMap 是一个基于go原生map实现的分片哈希结构，用于快速定位数据记录

- 分片设计原理
  - `减少锁竞争` 全局map拆分成多个独立分片(shard)，每个分片持有自己的锁，减少锁颗粒度
  - `哈希定位` 每个分片使用一致性哈希算法，将key映射到分片，减少冲突
  - `读写分离` 使用sync.RWMutex读写锁，允许并发读，写互斥

> 分片map的相关资料
>
> - <https://github.com/orcaman/concurrent-map>
> - <https://github.com/HDT3213/godis/blob/master/datastruct/dict/concurrent.go>
> - <https://github.com/jianghushinian/blog-go-example/tree/main/sync/map/concurrent-map>

## Inspired

- [rosedblabs/mini-bitcask](https://github.com/rosedblabs/mini-bitcask.git)
- [prologic/bitcask](https://git.mills.io/prologic/bitcask)
- [ahmeducf/bitcask](https://github.com/ahmeducf/bitcask)
