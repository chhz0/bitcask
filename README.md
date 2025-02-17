# caskv
Efficient log-structured key-value storage engine, designed for low-latency read and write, high throughput and high reliability

```
/caskv
  ├── cmd
  │   └── caskv-cli            # 命令行工具入口（可选）
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
  │   │   ├── shard_map.go     # 分配哈希表索引实现
  │   │   └── loader.go         # 数据文件索引重建
  │   ├── record                # 数据记录结构
  │   │   ├── record.go         # 记录结构定义 文件记录结构和内存索引定义
  │   │   └── validator.go      # CRC校验逻辑
  │   ├── codec                 # 编码解码
  │   │   ├── encoder.go        # 二进制编码
  │   │   └── decoder.go        # 二进制解码
  │   ├── merger                # 数据合并模块
  │   │    └── compact.go        # 合并旧文件，清理失效数据
  │   ├── wal.go                 #
  │   └── engine.go             # 引擎接口（Put/Get/Delete）
  ├── tests                     # 单元测试和压力测试
  ├── betch.go // 批处理
  ├── caskv.go // 引擎初始化，对外暴露接口
  ├── errors.go // 错误定义
  ├── hint.go
  ├── iter.go // 迭代器
  ├── merge.go // 合并逻辑
  ├── options.go // 引擎配置
  ├── go.mod
  └── README.md
```

## Record

文件日记记录：
```
header[19] = CRC(4) + Timestamp(8) + KeySize(2) + ValueSize(4) + Deleted(1)

| header  | Key(N) | Value(M) |

// LogRecord 文件日志记录结构
type LogRecord struct {
	CRC       uint32 // crc校验码(Header + Key + Value)
	Timestamp uint64 // 时间戳 (unix时间戳)
	Key       []byte
	Value     []byte
	Delete    bool
}
```

- `CRC`：校验范围除 CRC 字段以外的所有字段，用于校验数据是否损坏