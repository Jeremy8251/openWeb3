## 以太坊数据持久化

源码：https://github.com/ethereum/go-ethereum/tree/release/1.15/ethdb

目录结构如下所示:

```properties
├── dbtest   // 对数据库测试
├── leveldb  // 开源k-v存储数据库
├── memorydb // 内存数据库实际是map
├── pebble   // 并发型k-v键值存储库
├── remotedb  //远程数据库进行键值
├── batch.go  //批量处理操作数据
├── database.go //定义数据存储处理数据操作接口
└── iterator.go //定义数据存储迭代处理数据操作接口

```

### 一、数据接口

database.go

```go
//Database接口位置: go-ethereum/ethdb/database.go
type Database interface {
	KeyValueStore
	AncientStore
}
type KeyValueStore interface {
	KeyValueReader//定义了读取键值对的方法
	KeyValueWriter//定义了写入键值对的方法
	KeyValueStater//定义了获取键值对状态（存在、删除等）的方法
	KeyValueRangeDeleter//范围删除
	Batcher   // 批处理接口
	Iteratee  // 迭代接口
	Compacter // 压缩功能
	io.Closer // 数据库关闭功能
}

// AncientReader 扩展的旧数据读取接口（支持原子读操作）
type AncientReader interface {
    AncientReaderOp

    // 在无写入操作的情况下执行读取函数
    ReadAncients(fn func(AncientReaderOp) error) (err error)
}

// AncientReaderOp 接口定义不可变历史数据（Ancient Data）的读取操作
type AncientReaderOp interface {
    // 检查指定类型和编号的旧数据是否存在
    HasAncient(kind string, number uint64) (bool, error)

    // 获取指定类型和编号的旧数据
    Ancient(kind string, number uint64) ([]byte, error)

    // 范围获取旧数据：
    // - 从 start 开始获取最多 count 个条目
    // - 如果设置了 maxBytes，则返回总大小不超过该值的数据
    AncientRange(kind string, start, count, maxBytes uint64) ([][]byte, error)

    // 获取存储的最高旧数据编号
    Ancients() (uint64, error)

    // 获取已删除的最小编号（即当前存储的最小可用编号）
    Tail() (uint64, error)

    // 获取指定类型的旧数据总大小
    AncientSize(kind string) (uint64, error)
}

//上述中都是接口,那么ethdb的数据库必须要实现上述接口中的所有方法.
```

**底层数据库支持leveldb(kv键值存储库),memorydb(内存数据库)以及pebble(kv键值存储库)三种数据库.** leveldb是最常用和成熟的选项,而memorydb则适用于对速度要求较高且数据不需要持久化的场景.

### 二、LevelDB**‌与**Pebble

**设计差异**

- ‌**LevelDB**‌
  - 由Google开发，基于LSM-Tree（日志结构合并树）的键值存储引擎，适用于顺序写密集场景‌。
  - 采用`单进程设计`，依赖全局锁进行并发控制，在单线程写入性能上表现优异，但多线程环境下锁竞争可能导致性能下降‌。
- ‌**Pebble**‌
  - Cockroach Labs开发的LevelDB改进版，保留LSM-Tree结构但针对现代硬件优化。
  - 引入更细粒度的锁机制和并发控制算法，提升`多线程`环境下的吞吐量‌。

**适配场景**‌

1. ‌**LevelDB**‌
   - 适合轻量级节点或低并发场景，例如个人开发者调试或小型网络部署‌。
   - 成熟稳定，社区支持广泛，适合无需高吞吐的链上数据存储（如区块头、交易日志）‌。
2. ‌**Pebble**‌
   - 适用于全节点或企业级应用，需要处理`高并发`读写请求（如交易所节点、区块浏览器后端）‌。
   - 在统计账户数量、批量读取区块头等批量操作中表现更优‌

`LevelDB`满足基础需求与稳定性，内存分配固定，成熟度高，兼容性强，默认集成于Geth客户端，适合轻节点、开发者调试等资源受限场景‌

`Pebble`针对高并发场景优化性能，动态内存管理减少碎片，资源利用率更高

### 三、leveldb

google开发开源k-v存储数据库

源码：https://github.com/syndtr/goleveldb

特点：

1. leveldb是一个持久化存储的KV系统，与redis相比，leveldb是将大部分数据存储到磁盘中。而redis是一个内存型的KV存储系统，会吃内存。
2. leveldb在存储数据时，是有序存储的，也就是相邻的key值在存储文件中是按照顺序存储的
3. 与其它KV系统一样，levelDb操作接口简单，基本操作也只包括增、删、改、查。也支持批量操作
4. levelDb支持数据快照(snapshot)功能，可以使得读取操作不受到写操作的影响
5. levelDb支持数据压缩，可以很好的减少存储空间，提高IO效率。
6. 支持布隆过滤器，Bloom Filter主要用于减少磁盘I/O操作。当查询一个键时，LevelDB会检查包含该键的SSTable文件是否包含该键的Bloom Filter。如果Bloom Filter判断该键可能存在，LevelDB才会进一步读取SSTable文件进行确认；如果Bloom Filter判断该键不存在，则直接返回不存在，`从而大大减少了不必要的磁盘I/O操作‌`

* 限制
  1. 非关系型数据库，不支持sql查询，不支持索引
  2. 一次只允许一个进程访问一个特定的数据库

**Metircs**

系统性能度量框架，如果我们需要为某个系统或者服务做监控、统计等，就可以用到它。通常有5种类型：

1. Meters：监控一系列事件发生的速率，在以太坊最大的作用就是监控TPS，Meters会统计最近1min，5min，15min以及全部时间的速率。
2. Gauge：最简单的度量指标、统计瞬时状态，只有一个简单的返回值。
3. Histogram：统计数据的分布情况，比如最小值、最大值、中间值、中位数
4. Timers和meters类似，他是meters和histogram结合，histogram统计耗时，meters统计TPS

**代码分析**

[go-ethereum/ethdb/leveldb/leveldb.go](https://github.com/ethereum/go-ethereum/tree/release/1.15/ethdb/leveldb/leveldb.go)

流程：

1. 创建数据库实例
2. 对数据库接口实现
   1. 单条数据操作
   2. 批量数据操作
3. 对eth服务的监听及数据统计



1. **引入相关包leveldb、metrics**

```go
package leveldb

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

```

2. **定义结构体 Database**

   封装 LevelDB 实例和文件路径，用于实现 Geth 的数据库接口。

```go
// Database is a persistent key-value store. Apart from basic data storage
// functionality it also supports batch writes and iterating over the keyspace in
// binary-alphabetical order.
// Database 封装了 LevelDB 实例，并实现 ethdb.Database 接口
type Database struct {
	fn string      // 数据库文件路径
	db *leveldb.DB // 数据库实例
    // 测量数据库性能的meter
	compTimeMeter       *metrics.Meter // Meter for measuring the total time spent in database compaction
	compReadMeter       *metrics.Meter // Meter for measuring the data read during compaction
	compWriteMeter      *metrics.Meter // Meter for measuring the data written during compaction
	writeDelayNMeter    *metrics.Meter // Meter for measuring the write delay number due to database compaction
	writeDelayMeter     *metrics.Meter // Meter for measuring the write delay duration due to database compaction
	diskSizeGauge       *metrics.Gauge // Gauge for tracking the size of all the levels in the database
	diskReadMeter       *metrics.Meter // Meter for measuring the effective amount of data read
	diskWriteMeter      *metrics.Meter // Meter for measuring the effective amount of data written
	memCompGauge        *metrics.Gauge // Gauge for tracking the number of memory compaction
	level0CompGauge     *metrics.Gauge // Gauge for tracking the number of table compaction in level0
	nonlevel0CompGauge  *metrics.Gauge // Gauge for tracking the number of table compaction in non0 level
	seekCompGauge       *metrics.Gauge // Gauge for tracking the number of table compaction caused by read opt
	manualMemAllocGauge *metrics.Gauge // Gauge to track the amount of memory that has been manually allocated (not a part of runtime/GC)

	levelsGauge []*metrics.Gauge // Gauge for tracking the number of tables in levels

	quitLock sync.Mutex      // Mutex protecting the quit channel access
	quitChan chan chan error // Quit channel to stop the metrics collection before closing the database

	log log.Logger // Contextual logger tracking the database path
}
```

3. **创建实例化数据库对象**

```go
// New returns a wrapped LevelDB object. The namespace is the prefix that the
// metrics reporting should use for surfacing internal stats.
// New 创建或打开一个 LevelDB 数据库
func New(file string, cache int, handles int, namespace string, readonly bool) (*Database, error) {
	return NewCustom(file, namespace, func(options *opt.Options) {
		// Ensure we have some minimal caching and file guarantees
		if cache < minCache {
			cache = minCache
		}
		if handles < minHandles {
			handles = minHandles
		}
		// Set default options
        // 配置 LevelDB 选项
		options.OpenFilesCacheCapacity = handles// 最大打开文件数
		options.BlockCacheCapacity = cache / 2 * opt.MiB// 块缓存大小
		options.WriteBuffer = cache / 4 * opt.MiB // 写缓冲区大小
		if readonly {
			options.ReadOnly = true
		}
	})
}

// NewCustom returns a wrapped LevelDB object. The namespace is the prefix that the
// metrics reporting should use for surfacing internal stats.
// The customize function allows the caller to modify the leveldb options.
func NewCustom(file string, namespace string, customize func(options *opt.Options)) (*Database, error) {
	options := configureOptions(customize)
	logger := log.New("database", file)
	usedCache := options.GetBlockCacheCapacity() + options.GetWriteBuffer()*2
	logCtx := []interface{}{"cache", common.StorageSize(usedCache), "handles", options.GetOpenFilesCacheCapacity()}
	if options.ReadOnly {
		logCtx = append(logCtx, "readonly", "true")
	}
	logger.Info("Allocated cache and file handles", logCtx...)

	// Open the db and recover any potential corruptions
    // 打开数据库
	db, err := leveldb.OpenFile(file, options)
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		db, err = leveldb.RecoverFile(file, nil)
	}
	if err != nil {
		return nil, err
	}
	// Assemble the wrapper with all the registered metrics
    // 返回封装后的 Database 实例
	ldb := &Database{
		fn:       file,
		db:       db,
		log:      logger,
		quitChan: make(chan chan error),
	}
	ldb.compTimeMeter = metrics.NewRegisteredMeter(namespace+"compact/time", nil)
	ldb.compReadMeter = metrics.NewRegisteredMeter(namespace+"compact/input", nil)
	ldb.compWriteMeter = metrics.NewRegisteredMeter(namespace+"compact/output", nil)
	ldb.diskSizeGauge = metrics.NewRegisteredGauge(namespace+"disk/size", nil)
	ldb.diskReadMeter = metrics.NewRegisteredMeter(namespace+"disk/read", nil)
	ldb.diskWriteMeter = metrics.NewRegisteredMeter(namespace+"disk/write", nil)
	ldb.writeDelayMeter = metrics.NewRegisteredMeter(namespace+"compact/writedelay/duration", nil)
	ldb.writeDelayNMeter = metrics.NewRegisteredMeter(namespace+"compact/writedelay/counter", nil)
	ldb.memCompGauge = metrics.NewRegisteredGauge(namespace+"compact/memory", nil)
	ldb.level0CompGauge = metrics.NewRegisteredGauge(namespace+"compact/level0", nil)
	ldb.nonlevel0CompGauge = metrics.NewRegisteredGauge(namespace+"compact/nonlevel0", nil)
	ldb.seekCompGauge = metrics.NewRegisteredGauge(namespace+"compact/seek", nil)
	ldb.manualMemAllocGauge = metrics.NewRegisteredGauge(namespace+"memory/manualalloc", nil)

	// Start up the metrics gathering and return
	go ldb.meter(metricsGatheringInterval, namespace)
	return ldb, nil
}

// configureOptions sets some default options, then runs the provided setter.
func configureOptions(customizeFn func(*opt.Options)) *opt.Options {
	// Set default options
	options := &opt.Options{
		Filter:                 filter.NewBloomFilter(10),
		DisableSeeksCompaction: true,
	}
	// Allow caller to make custom modifications to the options
	if customizeFn != nil {
		customizeFn(options)
	}
	return options
}
```

4. **关闭数据库**

```go
// Close stops the metrics collection, flushes any pending data to disk and closes
// all io accesses to the underlying key-value store.
func (db *Database) Close() error {
	db.quitLock.Lock()
	defer db.quitLock.Unlock()

	if db.quitChan != nil {
		errc := make(chan error)
		db.quitChan <- errc
		if err := <-errc; err != nil {
			db.log.Error("Metrics collection failed", "err", err)
		}
		db.quitChan = nil
	}
	return db.db.Close()
}
```

5. **调用 LevelDB增删改查、批量操作、迭代器**

```go
// Has retrieves if a key is present in the key-value store.
func (db *Database) Has(key []byte) (bool, error) {
	return db.db.Has(key, nil)
}

// Get retrieves the given key if it's present in the key-value store.
func (db *Database) Get(key []byte) ([]byte, error) {
	dat, err := db.db.Get(key, nil)
	if err != nil {
		return nil, err
	}
	return dat, nil
}

// Put inserts the given value into the key-value store.
func (db *Database) Put(key []byte, value []byte) error {
	return db.db.Put(key, value, nil)
}

// Delete removes the key from the key-value store.
func (db *Database) Delete(key []byte) error {
	return db.db.Delete(key, nil)
}

var ErrTooManyKeys = errors.New("too many keys in deleted range")

// DeleteRange deletes all of the keys (and values) in the range [start,end)
// (inclusive on start, exclusive on end).
// Note that this is a fallback implementation as leveldb does not natively
// support range deletion. It can be slow and therefore the number of deleted
// keys is limited in order to avoid blocking for a very long time.
// ErrTooManyKeys is returned if the range has only been partially deleted.
// In this case the caller can repeat the call until it finally succeeds.
func (db *Database) DeleteRange(start, end []byte) error {
	batch := db.NewBatch()
	it := db.NewIterator(nil, start)
	defer it.Release()

	var count int
	for it.Next() && bytes.Compare(end, it.Key()) > 0 {
		count++
		if count > 10000 { // should not block for more than a second
			if err := batch.Write(); err != nil {
				return err
			}
			return ErrTooManyKeys
		}
		if err := batch.Delete(it.Key()); err != nil {
			return err
		}
	}
	return batch.Write()
}

// NewBatch creates a write-only key-value store that buffers changes to its host
// database until a final write is called.
func (db *Database) NewBatch() ethdb.Batch {
	return &batch{
		db: db.db,
		b:  new(leveldb.Batch),
	}
}

// NewBatchWithSize creates a write-only database batch with pre-allocated buffer.
func (db *Database) NewBatchWithSize(size int) ethdb.Batch {
	return &batch{
		db: db.db,
		b:  leveldb.MakeBatch(size),
	}
}

// NewIterator creates a binary-alphabetical iterator over a subset
// of database content with a particular key prefix, starting at a particular
// initial key (or after, if it does not exist).
func (db *Database) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	return db.db.NewIterator(bytesPrefixRange(prefix, start), nil)
}

// Stat returns the statistic data of the database.
func (db *Database) Stat() (string, error) {
	var stats leveldb.DBStats
	if err := db.db.Stats(&stats); err != nil {
		return "", err
	}
	var (
		message       string
		totalRead     int64
		totalWrite    int64
		totalSize     int64
		totalTables   int
		totalDuration time.Duration
	)
	if len(stats.LevelSizes) > 0 {
		message += " Level |   Tables   |    Size(MB)   |    Time(sec)  |    Read(MB)   |   Write(MB)\n" +
			"-------+------------+---------------+---------------+---------------+---------------\n"
		for level, size := range stats.LevelSizes {
			read := stats.LevelRead[level]
			write := stats.LevelWrite[level]
			duration := stats.LevelDurations[level]
			tables := stats.LevelTablesCounts[level]

			if tables == 0 && duration == 0 {
				continue
			}
			totalTables += tables
			totalSize += size
			totalRead += read
			totalWrite += write
			totalDuration += duration
			message += fmt.Sprintf(" %3d   | %10d | %13.5f | %13.5f | %13.5f | %13.5f\n",
				level, tables, float64(size)/1048576.0, duration.Seconds(),
				float64(read)/1048576.0, float64(write)/1048576.0)
		}
		message += "-------+------------+---------------+---------------+---------------+---------------\n"
		message += fmt.Sprintf(" Total | %10d | %13.5f | %13.5f | %13.5f | %13.5f\n",
			totalTables, float64(totalSize)/1048576.0, totalDuration.Seconds(),
			float64(totalRead)/1048576.0, float64(totalWrite)/1048576.0)
		message += "-------+------------+---------------+---------------+---------------+---------------\n\n"
	}
	message += fmt.Sprintf("Read(MB):%.5f Write(MB):%.5f\n", float64(stats.IORead)/1048576.0, float64(stats.IOWrite)/1048576.0)
	message += fmt.Sprintf("BlockCache(MB):%.5f FileCache:%d\n", float64(stats.BlockCacheSize)/1048576.0, stats.OpenedTablesCount)
	message += fmt.Sprintf("MemoryCompaction:%d Level0Compaction:%d NonLevel0Compaction:%d SeekCompaction:%d\n", stats.MemComp, stats.Level0Comp, stats.NonLevel0Comp, stats.SeekComp)
	message += fmt.Sprintf("WriteDelayCount:%d WriteDelayDuration:%s Paused:%t\n", stats.WriteDelayCount, common.PrettyDuration(stats.WriteDelayDuration), stats.WritePaused)
	message += fmt.Sprintf("Snapshots:%d Iterators:%d\n", stats.AliveSnapshots, stats.AliveIterators)
	return message, nil
}

// Compact flattens the underlying data store for the given key range. In essence,
// deleted and overwritten versions are discarded, and the data is rearranged to
// reduce the cost of operations needed to access them.
//
// A nil start is treated as a key before all keys in the data store; a nil limit
// is treated as a key after all keys in the data store. If both is nil then it
// will compact entire data store.
func (db *Database) Compact(start []byte, limit []byte) error {
	return db.db.CompactRange(util.Range{Start: start, Limit: limit})
}

// Path returns the path to the database directory.
func (db *Database) Path() string {
	return db.fn
}

// batch is a write-only leveldb batch that commits changes to its host database
// when Write is called. A batch cannot be used concurrently.
type batch struct {
	db   *leveldb.DB
	b    *leveldb.Batch
	size int
}

// Put inserts the given value into the batch for later committing.
func (b *batch) Put(key, value []byte) error {
	b.b.Put(key, value)
	b.size += len(key) + len(value)
	return nil
}

// Delete inserts the key removal into the batch for later committing.
func (b *batch) Delete(key []byte) error {
	b.b.Delete(key)
	b.size += len(key)
	return nil
}

// ValueSize retrieves the amount of data queued up for writing.
func (b *batch) ValueSize() int {
	return b.size
}

// Write flushes any accumulated data to disk.
func (b *batch) Write() error {
	return b.db.Write(b.b, nil)
}

// Reset resets the batch for reuse.
func (b *batch) Reset() {
	b.b.Reset()
	b.size = 0
}

// Replay replays the batch contents.
func (b *batch) Replay(w ethdb.KeyValueWriter) error {
	return b.b.Replay(&replayer{writer: w})
}

// replayer is a small wrapper to implement the correct replay methods.
type replayer struct {
	writer  ethdb.KeyValueWriter
	failure error
}

// Put inserts the given value into the key-value data store.
func (r *replayer) Put(key, value []byte) {
	// If the replay already failed, stop executing ops
	if r.failure != nil {
		return
	}
	r.failure = r.writer.Put(key, value)
}

// Delete removes the key from the key-value data store.
func (r *replayer) Delete(key []byte) {
	// If the replay already failed, stop executing ops
	if r.failure != nil {
		return
	}
	r.failure = r.writer.Delete(key)
}

// bytesPrefixRange returns key range that satisfy
// - the given prefix, and
// - the given seek position
func bytesPrefixRange(prefix, start []byte) *util.Range {
	r := util.BytesPrefix(prefix)
	r.Start = append(r.Start, start...)
	return r
}
```

6. **启动性能检测**

```go
// meter periodically retrieves internal leveldb counters and reports them to
// the metrics subsystem.
func (db *Database) meter(refresh time.Duration, namespace string) {
	// Create the counters to store current and previous compaction values
	compactions := make([][]int64, 2)
	for i := 0; i < 2; i++ {
		compactions[i] = make([]int64, 4)
	}
	// Create storages for states and warning log tracer.
	var (
		errc chan error
		merr error

		stats           leveldb.DBStats
		iostats         [2]int64
		delaystats      [2]int64
		lastWritePaused time.Time
	)
	timer := time.NewTimer(refresh)
	defer timer.Stop()

	// Iterate ad infinitum and collect the stats
	for i := 1; errc == nil && merr == nil; i++ {
		// Retrieve the database stats
		// Stats method resets buffers inside therefore it's okay to just pass the struct.
		err := db.db.Stats(&stats)
		if err != nil {
			db.log.Error("Failed to read database stats", "err", err)
			merr = err
			continue
		}
		// Iterate over all the leveldbTable rows, and accumulate the entries
		for j := 0; j < len(compactions[i%2]); j++ {
			compactions[i%2][j] = 0
		}
		compactions[i%2][0] = stats.LevelSizes.Sum()
		for _, t := range stats.LevelDurations {
			compactions[i%2][1] += t.Nanoseconds()
		}
		compactions[i%2][2] = stats.LevelRead.Sum()
		compactions[i%2][3] = stats.LevelWrite.Sum()
		// Update all the requested meters
		if db.diskSizeGauge != nil {
			db.diskSizeGauge.Update(compactions[i%2][0])
		}
		if db.compTimeMeter != nil {
			db.compTimeMeter.Mark(compactions[i%2][1] - compactions[(i-1)%2][1])
		}
		if db.compReadMeter != nil {
			db.compReadMeter.Mark(compactions[i%2][2] - compactions[(i-1)%2][2])
		}
		if db.compWriteMeter != nil {
			db.compWriteMeter.Mark(compactions[i%2][3] - compactions[(i-1)%2][3])
		}
		var (
			delayN   = int64(stats.WriteDelayCount)
			duration = stats.WriteDelayDuration
			paused   = stats.WritePaused
		)
		if db.writeDelayNMeter != nil {
			db.writeDelayNMeter.Mark(delayN - delaystats[0])
		}
		if db.writeDelayMeter != nil {
			db.writeDelayMeter.Mark(duration.Nanoseconds() - delaystats[1])
		}
		// If a warning that db is performing compaction has been displayed, any subsequent
		// warnings will be withheld for one minute not to overwhelm the user.
		if paused && delayN-delaystats[0] == 0 && duration.Nanoseconds()-delaystats[1] == 0 &&
			time.Now().After(lastWritePaused.Add(degradationWarnInterval)) {
			db.log.Warn("Database compacting, degraded performance")
			lastWritePaused = time.Now()
		}
		delaystats[0], delaystats[1] = delayN, duration.Nanoseconds()

		var (
			nRead  = int64(stats.IORead)
			nWrite = int64(stats.IOWrite)
		)
		if db.diskReadMeter != nil {
			db.diskReadMeter.Mark(nRead - iostats[0])
		}
		if db.diskWriteMeter != nil {
			db.diskWriteMeter.Mark(nWrite - iostats[1])
		}
		iostats[0], iostats[1] = nRead, nWrite

		db.memCompGauge.Update(int64(stats.MemComp))
		db.level0CompGauge.Update(int64(stats.Level0Comp))
		db.nonlevel0CompGauge.Update(int64(stats.NonLevel0Comp))
		db.seekCompGauge.Update(int64(stats.SeekComp))

		for i, tables := range stats.LevelTablesCounts {
			// Append metrics for additional layers
			if i >= len(db.levelsGauge) {
				db.levelsGauge = append(db.levelsGauge, metrics.NewRegisteredGauge(namespace+fmt.Sprintf("tables/level%v", i), nil))
			}
			db.levelsGauge[i].Update(int64(tables))
		}

		// Sleep a bit, then repeat the stats collection
		select {
		case errc = <-db.quitChan:
			// Quit requesting, stop hammering the database
		case <-timer.C:
			timer.Reset(refresh)
			// Timeout, gather a new set of stats
		}
	}

	if errc == nil {
		errc = <-db.quitChan
	}
	errc <- merr
}
```



### 四、MemoryDB

go-ethereum中的MemoryDB是一个用于存储以太坊区块链数据的内存数据库。

MemoryDB在go-ethereum中主要用于存储区块和交易数据，以便快速访问这些信息。它通过内存中的数据结构来管理这些数据，使得查询和操作速度非常快，但缺点是重启服务后会丢失所有数据‌。

MemoryDB的工作原理和结构

MemoryDB在go-ethereum中通过内存中的数据结构来存储和管理区块链数据。它不依赖于磁盘存储，因此查询和操作速度非常快。然而，这也意味着所有数据在服务重启后会丢失，因此通常不用于长期存储‌。

**MemoryDB的优缺点**

‌**优点**‌：

- ‌**速度快**‌：由于数据存储在内存中，访问和操作速度非常快。
- ‌**高效**‌：适合需要快速读取和写入操作的场景。

‌**缺点**‌：

- ‌**数据不持久**‌：重启服务后，所有数据会丢失，不适合需要长期存储的场景。
- ‌**资源消耗大**‌：内存使用量大，可能不适合资源有限的系统。

[go-ethereum/ethdb/memorydb/memorydb.go](https://github.com/ethereum/go-ethereum/tree/release/1.15/ethdb/memorydb/memorydb.go)

1. **定义Database结构体**

实际创建了一个map

```go
package memorydb

import (
	"errors"
	"sort"
	"strings"
	"sync"//sync 包用于并发控制

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"//实现数据库接口
)
// Database 实现了 ethdb.Database 接口，表示一个内存中的键值数据库。
type Database struct {
	db   map[string][]byte//键值对存储实例
	lock sync.RWMutex//读写互斥锁
}
```

2. **调用 LevelDB增删改查、批量操作、迭代器**

```go
// New returns a wrapped map with all the required database interface methods
// implemented.
// New 创建一个新的内存数据库实例。
func New() *Database {
	return &Database{
		db: make(map[string][]byte),
	}
}

// NewWithCap returns a wrapped map pre-allocated to the provided capacity with
// all the required database interface methods implemented.
// New 创建一个size容量大小的内存数据库实例。
func NewWithCap(size int) *Database {
	return &Database{
		db: make(map[string][]byte, size),
	}
}
```

3. **关闭数据库，清空map**

```go
// Close deallocates the internal map and ensures any consecutive data access op
// fails with an error.
func (db *Database) Close() error {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.db = nil
	return nil
}
```

4. **实现增删改查、批量操作、迭代器等接口**

```go
// Has retrieves if a key is present in the key-value store.
func (db *Database) Has(key []byte) (bool, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if db.db == nil {
		return false, errMemorydbClosed
	}
	_, ok := db.db[string(key)]
	return ok, nil
}

// Get retrieves the given key if it's present in the key-value store.
func (db *Database) Get(key []byte) ([]byte, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if db.db == nil {
		return nil, errMemorydbClosed
	}
	if entry, ok := db.db[string(key)]; ok {
		return common.CopyBytes(entry), nil
	}
	return nil, errMemorydbNotFound
}

// Put inserts the given value into the key-value store.
func (db *Database) Put(key []byte, value []byte) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.db == nil {
		return errMemorydbClosed
	}
	db.db[string(key)] = common.CopyBytes(value)
	return nil
}

// Delete removes the key from the key-value store.
func (db *Database) Delete(key []byte) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.db == nil {
		return errMemorydbClosed
	}
	delete(db.db, string(key))
	return nil
}

// DeleteRange deletes all of the keys (and values) in the range [start,end)
// (inclusive on start, exclusive on end).
func (db *Database) DeleteRange(start, end []byte) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	if db.db == nil {
		return errMemorydbClosed
	}

	for key := range db.db {
		if key >= string(start) && key < string(end) {
			delete(db.db, key)
		}
	}
	return nil
}

// NewBatch creates a write-only key-value store that buffers changes to its host
// database until a final write is called.
func (db *Database) NewBatch() ethdb.Batch {
	return &batch{
		db: db,
	}
}

// NewBatchWithSize creates a write-only database batch with pre-allocated buffer.
func (db *Database) NewBatchWithSize(size int) ethdb.Batch {
	return &batch{
		db: db,
	}
}

// NewIterator creates a binary-alphabetical iterator over a subset
// of database content with a particular key prefix, starting at a particular
// initial key (or after, if it does not exist).
func (db *Database) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	db.lock.RLock()
	defer db.lock.RUnlock()

	var (
		pr     = string(prefix)
		st     = string(append(prefix, start...))
		keys   = make([]string, 0, len(db.db))
		values = make([][]byte, 0, len(db.db))
	)
	// Collect the keys from the memory database corresponding to the given prefix
	// and start
	for key := range db.db {
		if !strings.HasPrefix(key, pr) {
			continue
		}
		if key >= st {
			keys = append(keys, key)
		}
	}
	// Sort the items and retrieve the associated values
	sort.Strings(keys)
	for _, key := range keys {
		values = append(values, db.db[key])
	}
	return &iterator{
		index:  -1,
		keys:   keys,
		values: values,
	}
}

// Stat returns the statistic data of the database.
func (db *Database) Stat() (string, error) {
	return "", nil
}

// Compact is not supported on a memory database, but there's no need either as
// a memory database doesn't waste space anyway.
func (db *Database) Compact(start []byte, limit []byte) error {
	return nil
}

// Len returns the number of entries currently present in the memory database.
//
// Note, this method is only used for testing (i.e. not public in general) and
// does not have explicit checks for closed-ness to allow simpler testing code.
func (db *Database) Len() int {
	db.lock.RLock()
	defer db.lock.RUnlock()

	return len(db.db)
}

// keyvalue is a key-value tuple tagged with a deletion field to allow creating
// memory-database write batches.
type keyvalue struct {
	key    string
	value  []byte
	delete bool
}

// batch is a write-only memory batch that commits changes to its host
// database when Write is called. A batch cannot be used concurrently.
type batch struct {
	db     *Database
	writes []keyvalue
	size   int
}

// Put inserts the given value into the batch for later committing.
func (b *batch) Put(key, value []byte) error {
	b.writes = append(b.writes, keyvalue{string(key), common.CopyBytes(value), false})
	b.size += len(key) + len(value)
	return nil
}

// Delete inserts the key removal into the batch for later committing.
func (b *batch) Delete(key []byte) error {
	b.writes = append(b.writes, keyvalue{string(key), nil, true})
	b.size += len(key)
	return nil
}

// ValueSize retrieves the amount of data queued up for writing.
func (b *batch) ValueSize() int {
	return b.size
}

// Write flushes any accumulated data to the memory database.
func (b *batch) Write() error {
	b.db.lock.Lock()
	defer b.db.lock.Unlock()

	if b.db.db == nil {
		return errMemorydbClosed
	}
	for _, keyvalue := range b.writes {
		if keyvalue.delete {
			delete(b.db.db, keyvalue.key)
			continue
		}
		b.db.db[keyvalue.key] = keyvalue.value
	}
	return nil
}

// Reset resets the batch for reuse.
func (b *batch) Reset() {
	b.writes = b.writes[:0]
	b.size = 0
}

// Replay replays the batch contents.
func (b *batch) Replay(w ethdb.KeyValueWriter) error {
	for _, keyvalue := range b.writes {
		if keyvalue.delete {
			if err := w.Delete([]byte(keyvalue.key)); err != nil {
				return err
			}
			continue
		}
		if err := w.Put([]byte(keyvalue.key), keyvalue.value); err != nil {
			return err
		}
	}
	return nil
}

// iterator can walk over the (potentially partial) keyspace of a memory key
// value store. Internally it is a deep copy of the entire iterated state,
// sorted by keys.
type iterator struct {
	index  int
	keys   []string
	values [][]byte
}

// Next moves the iterator to the next key/value pair. It returns whether the
// iterator is exhausted.
func (it *iterator) Next() bool {
	// Short circuit if iterator is already exhausted in the forward direction.
	if it.index >= len(it.keys) {
		return false
	}
	it.index += 1
	return it.index < len(it.keys)
}

// Error returns any accumulated error. Exhausting all the key/value pairs
// is not considered to be an error. A memory iterator cannot encounter errors.
func (it *iterator) Error() error {
	return nil
}

// Key returns the key of the current key/value pair, or nil if done. The caller
// should not modify the contents of the returned slice, and its contents may
// change on the next call to Next.
func (it *iterator) Key() []byte {
	// Short circuit if iterator is not in a valid position
	if it.index < 0 || it.index >= len(it.keys) {
		return nil
	}
	return []byte(it.keys[it.index])
}

// Value returns the value of the current key/value pair, or nil if done. The
// caller should not modify the contents of the returned slice, and its contents
// may change on the next call to Next.
func (it *iterator) Value() []byte {
	// Short circuit if iterator is not in a valid position
	if it.index < 0 || it.index >= len(it.keys) {
		return nil
	}
	return it.values[it.index]
}

// Release releases associated resources. Release should always succeed and can
// be called multiple times without causing error.
func (it *iterator) Release() {
	it.index, it.keys, it.values = -1, nil, nil
}
```

### 五、Pebble

[Pebble](https://github.com/cockroachdb/pebble) 是一个受 LevelDB/RocksDB 启发的键值存储，专注于 CockroachDB 的性能和内部使用。Pebble 继承了 RocksDB 文件格式和一些扩展，例如范围删除墓碑、表级布隆过滤器和 MANIFEST 格式的更新。

Pebble 在 CockroachDB v20.1（2020 年 5 月发布）中作为 RocksDB 的替代存储引擎引入，并在当时成功投入生产。Pebble 在 CockroachDB v20.2（2020 年 11 月发布）中成为默认存储引擎。CockroachDB 用户正在大规模生产中使用 Pebble，它被认为是稳定的且可用于生产。

[pebble.go](https://github.com/ethereum/go-ethereum/blob/release/1.15/ethdb/pebble/pebble.go)

1. 引入第三方pebble包

```go
package pebble

import (
	"bytes"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/cockroachdb/pebble/bloom"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
)
```

2. 定义Database常量与结构体

```go
const (
	// minCache is the minimum amount of memory in megabytes to allocate to pebble
	// read and write caching, split half and half.
	minCache = 16

	// minHandles is the minimum number of files handles to allocate to the open
	// database files.
	minHandles = 16

	// metricsGatheringInterval specifies the interval to retrieve pebble database
	// compaction, io and pause stats to report to the user.
	metricsGatheringInterval = 3 * time.Second

	// degradationWarnInterval specifies how often warning should be printed if the
	// leveldb database cannot keep up with requested writes.
	degradationWarnInterval = time.Minute
)

// Database is a persistent key-value store based on the pebble storage engine.
// Apart from basic data storage functionality it also supports batch writes and
// iterating over the keyspace in binary-alphabetical order.
type Database struct {
	fn string     // filename for reporting
	db *pebble.DB // Underlying pebble storage engine

	compTimeMeter       *metrics.Meter // Meter for measuring the total time spent in database compaction
	compReadMeter       *metrics.Meter // Meter for measuring the data read during compaction
	compWriteMeter      *metrics.Meter // Meter for measuring the data written during compaction
	writeDelayNMeter    *metrics.Meter // Meter for measuring the write delay number due to database compaction
	writeDelayMeter     *metrics.Meter // Meter for measuring the write delay duration due to database compaction
	diskSizeGauge       *metrics.Gauge // Gauge for tracking the size of all the levels in the database
	diskReadMeter       *metrics.Meter // Meter for measuring the effective amount of data read
	diskWriteMeter      *metrics.Meter // Meter for measuring the effective amount of data written
	memCompGauge        *metrics.Gauge // Gauge for tracking the number of memory compaction
	level0CompGauge     *metrics.Gauge // Gauge for tracking the number of table compaction in level0
	nonlevel0CompGauge  *metrics.Gauge // Gauge for tracking the number of table compaction in non0 level
	seekCompGauge       *metrics.Gauge // Gauge for tracking the number of table compaction caused by read opt
	manualMemAllocGauge *metrics.Gauge // Gauge for tracking amount of non-managed memory currently allocated

	levelsGauge []*metrics.Gauge // Gauge for tracking the number of tables in levels

	quitLock sync.RWMutex    // Mutex protecting the quit channel and the closed flag
	quitChan chan chan error // Quit channel to stop the metrics collection before closing the database
	closed   bool            // keep track of whether we're Closed

	log log.Logger // Contextual logger tracking the database path

	activeComp    int           // Current number of active compactions
	compStartTime time.Time     // The start time of the earliest currently-active compaction
	compTime      atomic.Int64  // Total time spent in compaction in ns
	level0Comp    atomic.Uint32 // Total number of level-zero compactions
	nonLevel0Comp atomic.Uint32 // Total number of non level-zero compactions

	writeStalled        atomic.Bool  // Flag whether the write is stalled
	writeDelayStartTime time.Time    // The start time of the latest write stall
	writeDelayCount     atomic.Int64 // Total number of write stall counts
	writeDelayTime      atomic.Int64 // Total time spent in write stalls

	writeOptions *pebble.WriteOptions
}
```

3. **初始化函数，创建实例**

```go
// New returns a wrapped pebble DB object. The namespace is the prefix that the
// metrics reporting should use for surfacing internal stats.
func New(file string, cache int, handles int, namespace string, readonly bool) (*Database, error) {
	// Ensure we have some minimal caching and file guarantees
	if cache < minCache {
		cache = minCache
	}
	if handles < minHandles {
		handles = minHandles
	}
	logger := log.New("database", file)
	logger.Info("Allocated cache and file handles", "cache", common.StorageSize(cache*1024*1024), "handles", handles)

	// The max memtable size is limited by the uint32 offsets stored in
	// internal/arenaskl.node, DeferredBatchOp, and flushableBatchEntry.
	//
	// - MaxUint32 on 64-bit platforms;
	// - MaxInt on 32-bit platforms.
	//
	// It is used when slices are limited to Uint32 on 64-bit platforms (the
	// length limit for slices is naturally MaxInt on 32-bit platforms).
	//
	// Taken from https://github.com/cockroachdb/pebble/blob/master/internal/constants/constants.go
	maxMemTableSize := (1<<31)<<(^uint(0)>>63) - 1

	// Two memory tables is configured which is identical to leveldb,
	// including a frozen memory table and another live one.
	memTableLimit := 2
	memTableSize := cache * 1024 * 1024 / 2 / memTableLimit

	// The memory table size is currently capped at maxMemTableSize-1 due to a
	// known bug in the pebble where maxMemTableSize is not recognized as a
	// valid size.
	//
	// TODO use the maxMemTableSize as the maximum table size once the issue
	// in pebble is fixed.
	if memTableSize >= maxMemTableSize {
		memTableSize = maxMemTableSize - 1
	}
	db := &Database{
		fn:           file,
		log:          logger,
		quitChan:     make(chan chan error),
		writeOptions: &pebble.WriteOptions{Sync: false},
	}
	// 配置参数，此处省略
    
	// Disable seek compaction explicitly. Check https://github.com/ethereum/go-ethereum/pull/20130
	// for more details.
	opt.Experimental.ReadSamplingMultiplier = -1

	// Open the db and recover any potential corruptions
	innerDB, err := pebble.Open(file, opt)
	if err != nil {
		return nil, err
	}
	db.db = innerDB

	db.compTimeMeter = metrics.GetOrRegisterMeter(namespace+"compact/time", nil)
	//........配置性能检测

	// Start up the metrics gathering and return
	go db.meter(metricsGatheringInterval, namespace)
	return db, nil
}
```

4. **压缩（Compaction）事件监控**

- 当 Pebble 数据库触发压缩（Compaction）时，记录压缩操作的开始时间和类型。
- 当压缩结束时，更新活跃压缩计数 `activeComp`

```go
func (d *Database) onCompactionBegin(info pebble.CompactionInfo) {
    // 如果是第一个活跃的压缩操作，记录开始时间
	if d.activeComp == 0 {
		d.compStartTime = time.Now()
	}
    // 判断是否为 Level-0 的压缩
	l0 := info.Input[0]
	if l0.Level == 0 {
		d.level0Comp.Add(1)// Level-0 压缩计数
	} else {
		d.nonLevel0Comp.Add(1)// 非 Level-0 压缩计数
	}
	d.activeComp++ // 活跃压缩数递增
}

func (d *Database) onCompactionEnd(info pebble.CompactionInfo) {
     // 如果是最后一个活跃的压缩操作，累计总耗时
	if d.activeComp == 1 {
		d.compTime.Add(int64(time.Since(d.compStartTime)))
	} else if d.activeComp == 0 {
		panic("should not happen")// 逻辑错误：无活跃压缩却触发结束
	}
	d.activeComp--// 活跃压缩数递减
}
```

5. **写入延迟（Write Stall）事件监控**

* 当 Pebble 因写入压力过大进入暂停状态（Write Stall）
* 当写入暂停状态结束

```go
func (d *Database) onWriteStallBegin(b pebble.WriteStallBeginInfo) {
	d.writeDelayStartTime = time.Now()// 记录写入延迟开始时间
	d.writeDelayCount.Add(1)// 写入延迟次数 +1
	d.writeStalled.Store(true)// 标记写入暂停状态为 true
}

func (d *Database) onWriteStallEnd() {
	d.writeDelayTime.Add(int64(time.Since(d.writeDelayStartTime)))// 累加总延迟时间
	d.writeStalled.Store(false)// 标记写入暂停状态为 false
}
```

6. **日志处理（panicLogger）**

```go
// panicLogger is just a noop logger to disable Pebble's internal logger.
//
// TODO(karalabe): Remove when Pebble sets this as the default.
type panicLogger struct{}

func (l panicLogger) Infof(format string, args ...interface{}) {
}

func (l panicLogger) Errorf(format string, args ...interface{}) {
}

func (l panicLogger) Fatalf(format string, args ...interface{}) {
	panic(fmt.Errorf("fatal: "+format, args...))
}

```

7. 关闭数据库

```go
// Close stops the metrics collection, flushes any pending data to disk and closes
// all io accesses to the underlying key-value store.
func (d *Database) Close() error {
	d.quitLock.Lock()
	defer d.quitLock.Unlock()
	// Allow double closing, simplifies things
	if d.closed {
		return nil
	}
	d.closed = true
	if d.quitChan != nil {
		errc := make(chan error)
		d.quitChan <- errc
		if err := <-errc; err != nil {
			d.log.Error("Metrics collection failed", "err", err)
		}
		d.quitChan = nil
	}
	return d.db.Close()
}
```

8. **查找key值是否存在**

```go
// Has retrieves if a key is present in the key-value store.
func (d *Database) Has(key []byte) (bool, error) {
	d.quitLock.RLock()
	defer d.quitLock.RUnlock()
	if d.closed {
		return false, pebble.ErrClosed
	}
	_, closer, err := d.db.Get(key)
	if err == pebble.ErrNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}
	if err = closer.Close(); err != nil {
		return false, err
	}
	return true, nil
}
```

9. **实现接口增删改查、批量操作、迭代器、性能监控等**

```go
// Get retrieves the given key if it's present in the key-value store.
func (d *Database) Get(key []byte) ([]byte, error) {
	d.quitLock.RLock()
	defer d.quitLock.RUnlock()
	if d.closed {
		return nil, pebble.ErrClosed
	}
	dat, closer, err := d.db.Get(key)
	if err != nil {
		return nil, err
	}
	ret := make([]byte, len(dat))
	copy(ret, dat)
	if err = closer.Close(); err != nil {
		return nil, err
	}
	return ret, nil
}

// Put inserts the given value into the key-value store.
func (d *Database) Put(key []byte, value []byte) error {
	d.quitLock.RLock()
	defer d.quitLock.RUnlock()
	if d.closed {
		return pebble.ErrClosed
	}
	return d.db.Set(key, value, d.writeOptions)
}

// Delete removes the key from the key-value store.
func (d *Database) Delete(key []byte) error {
	d.quitLock.RLock()
	defer d.quitLock.RUnlock()
	if d.closed {
		return pebble.ErrClosed
	}
	return d.db.Delete(key, d.writeOptions)
}

// DeleteRange deletes all of the keys (and values) in the range [start,end)
// (inclusive on start, exclusive on end).
func (d *Database) DeleteRange(start, end []byte) error {
	d.quitLock.RLock()
	defer d.quitLock.RUnlock()
	if d.closed {
		return pebble.ErrClosed
	}
	return d.db.DeleteRange(start, end, d.writeOptions)
}

// NewBatch creates a write-only key-value store that buffers changes to its host
// database until a final write is called.
func (d *Database) NewBatch() ethdb.Batch {
	return &batch{
		b:  d.db.NewBatch(),
		db: d,
	}
}

// NewBatchWithSize creates a write-only database batch with pre-allocated buffer.
func (d *Database) NewBatchWithSize(size int) ethdb.Batch {
	return &batch{
		b:  d.db.NewBatchWithSize(size),
		db: d,
	}
}

// upperBound returns the upper bound for the given prefix
func upperBound(prefix []byte) (limit []byte) {
	for i := len(prefix) - 1; i >= 0; i-- {
		c := prefix[i]
		if c == 0xff {
			continue
		}
		limit = make([]byte, i+1)
		copy(limit, prefix)
		limit[i] = c + 1
		break
	}
	return limit
}

// Stat returns the internal metrics of Pebble in a text format. It's a developer
// method to read everything there is to read, independent of Pebble version.
func (d *Database) Stat() (string, error) {
	return d.db.Metrics().String(), nil
}

// Compact flattens the underlying data store for the given key range. In essence,
// deleted and overwritten versions are discarded, and the data is rearranged to
// reduce the cost of operations needed to access them.
//
// A nil start is treated as a key before all keys in the data store; a nil limit
// is treated as a key after all keys in the data store. If both is nil then it
// will compact entire data store.
func (d *Database) Compact(start []byte, limit []byte) error {
	// There is no special flag to represent the end of key range
	// in pebble(nil in leveldb). Use an ugly hack to construct a
	// large key to represent it.
	// Note any prefixed database entry will be smaller than this
	// flag, as for trie nodes we need the 32 byte 0xff because
	// there might be a shared prefix starting with a number of
	// 0xff-s, so 32 ensures than only a hash collision could touch it.
	// https://github.com/cockroachdb/pebble/issues/2359#issuecomment-1443995833
	if limit == nil {
		limit = bytes.Repeat([]byte{0xff}, 32)
	}
	return d.db.Compact(start, limit, true) // Parallelization is preferred
}

// Path returns the path to the database directory.
func (d *Database) Path() string {
	return d.fn
}

// meter periodically retrieves internal pebble counters and reports them to
// the metrics subsystem.
// 性能监控
func (d *Database) meter(refresh time.Duration, namespace string) {
	var errc chan error
	timer := time.NewTimer(refresh)
	defer timer.Stop()

	// Create storage and warning log tracer for write delay.
	var (
		compTimes  [2]int64
		compWrites [2]int64
		compReads  [2]int64

		nWrites [2]int64

		writeDelayTimes      [2]int64
		writeDelayCounts     [2]int64
		lastWriteStallReport time.Time
	)

	// Iterate ad infinitum and collect the stats
	for i := 1; errc == nil; i++ {
		var (
			compWrite int64
			compRead  int64
			nWrite    int64

			stats              = d.db.Metrics()
			compTime           = d.compTime.Load()
			writeDelayCount    = d.writeDelayCount.Load()
			writeDelayTime     = d.writeDelayTime.Load()
			nonLevel0CompCount = int64(d.nonLevel0Comp.Load())
			level0CompCount    = int64(d.level0Comp.Load())
		)
		writeDelayTimes[i%2] = writeDelayTime
		writeDelayCounts[i%2] = writeDelayCount
		compTimes[i%2] = compTime

		for _, levelMetrics := range stats.Levels {
			nWrite += int64(levelMetrics.BytesCompacted)
			nWrite += int64(levelMetrics.BytesFlushed)
			compWrite += int64(levelMetrics.BytesCompacted)
			compRead += int64(levelMetrics.BytesRead)
		}

		nWrite += int64(stats.WAL.BytesWritten)

		compWrites[i%2] = compWrite
		compReads[i%2] = compRead
		nWrites[i%2] = nWrite

		if d.writeDelayNMeter != nil {
			d.writeDelayNMeter.Mark(writeDelayCounts[i%2] - writeDelayCounts[(i-1)%2])
		}
		if d.writeDelayMeter != nil {
			d.writeDelayMeter.Mark(writeDelayTimes[i%2] - writeDelayTimes[(i-1)%2])
		}
		// Print a warning log if writing has been stalled for a while. The log will
		// be printed per minute to avoid overwhelming users.
		if d.writeStalled.Load() && writeDelayCounts[i%2] == writeDelayCounts[(i-1)%2] &&
			time.Now().After(lastWriteStallReport.Add(degradationWarnInterval)) {
			d.log.Warn("Database compacting, degraded performance")
			lastWriteStallReport = time.Now()
		}
		if d.compTimeMeter != nil {
			d.compTimeMeter.Mark(compTimes[i%2] - compTimes[(i-1)%2])
		}
		if d.compReadMeter != nil {
			d.compReadMeter.Mark(compReads[i%2] - compReads[(i-1)%2])
		}
		if d.compWriteMeter != nil {
			d.compWriteMeter.Mark(compWrites[i%2] - compWrites[(i-1)%2])
		}
		if d.diskSizeGauge != nil {
			d.diskSizeGauge.Update(int64(stats.DiskSpaceUsage()))
		}
		if d.diskReadMeter != nil {
			d.diskReadMeter.Mark(0) // pebble doesn't track non-compaction reads
		}
		if d.diskWriteMeter != nil {
			d.diskWriteMeter.Mark(nWrites[i%2] - nWrites[(i-1)%2])
		}
		// See https://github.com/cockroachdb/pebble/pull/1628#pullrequestreview-1026664054
		manuallyAllocated := stats.BlockCache.Size + int64(stats.MemTable.Size) + int64(stats.MemTable.ZombieSize)
		d.manualMemAllocGauge.Update(manuallyAllocated)
		d.memCompGauge.Update(stats.Flush.Count)
		d.nonlevel0CompGauge.Update(nonLevel0CompCount)
		d.level0CompGauge.Update(level0CompCount)
		d.seekCompGauge.Update(stats.Compact.ReadCount)

		for i, level := range stats.Levels {
			// Append metrics for additional layers
			if i >= len(d.levelsGauge) {
				d.levelsGauge = append(d.levelsGauge, metrics.GetOrRegisterGauge(namespace+fmt.Sprintf("tables/level%v", i), nil))
			}
			d.levelsGauge[i].Update(level.NumFiles)
		}

		// Sleep a bit, then repeat the stats collection
		select {
		case errc = <-d.quitChan:
			// Quit requesting, stop hammering the database
		case <-timer.C:
			timer.Reset(refresh)
			// Timeout, gather a new set of stats
		}
	}
	errc <- nil
}

// batch is a write-only batch that commits changes to its host database
// when Write is called. A batch cannot be used concurrently.
type batch struct {
	b    *pebble.Batch
	db   *Database
	size int
}

// Put inserts the given value into the batch for later committing.
func (b *batch) Put(key, value []byte) error {
	if err := b.b.Set(key, value, nil); err != nil {
		return err
	}
	b.size += len(key) + len(value)
	return nil
}

// Delete inserts the key removal into the batch for later committing.
func (b *batch) Delete(key []byte) error {
	if err := b.b.Delete(key, nil); err != nil {
		return err
	}
	b.size += len(key)
	return nil
}

// ValueSize retrieves the amount of data queued up for writing.
func (b *batch) ValueSize() int {
	return b.size
}

// Write flushes any accumulated data to disk.
func (b *batch) Write() error {
	b.db.quitLock.RLock()
	defer b.db.quitLock.RUnlock()
	if b.db.closed {
		return pebble.ErrClosed
	}
	return b.b.Commit(b.db.writeOptions)
}

// Reset resets the batch for reuse.
func (b *batch) Reset() {
	b.b.Reset()
	b.size = 0
}

// Replay replays the batch contents.
func (b *batch) Replay(w ethdb.KeyValueWriter) error {
	reader := b.b.Reader()
	for {
		kind, k, v, ok, err := reader.Next()
		if !ok || err != nil {
			return err
		}
		// The (k,v) slices might be overwritten if the batch is reset/reused,
		// and the receiver should copy them if they are to be retained long-term.
		if kind == pebble.InternalKeyKindSet {
			if err = w.Put(k, v); err != nil {
				return err
			}
		} else if kind == pebble.InternalKeyKindDelete {
			if err = w.Delete(k); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unhandled operation, keytype: %v", kind)
		}
	}
}

// pebbleIterator is a wrapper of underlying iterator in storage engine.
// The purpose of this structure is to implement the missing APIs.
//
// The pebble iterator is not thread-safe.
type pebbleIterator struct {
	iter     *pebble.Iterator
	moved    bool
	released bool
}

// NewIterator creates a binary-alphabetical iterator over a subset
// of database content with a particular key prefix, starting at a particular
// initial key (or after, if it does not exist).
func (d *Database) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	iter, _ := d.db.NewIter(&pebble.IterOptions{
		LowerBound: append(prefix, start...),
		UpperBound: upperBound(prefix),
	})
	iter.First()
	return &pebbleIterator{iter: iter, moved: true, released: false}
}

// Next moves the iterator to the next key/value pair. It returns whether the
// iterator is exhausted.
func (iter *pebbleIterator) Next() bool {
	if iter.moved {
		iter.moved = false
		return iter.iter.Valid()
	}
	return iter.iter.Next()
}

// Error returns any accumulated error. Exhausting all the key/value pairs
// is not considered to be an error.
func (iter *pebbleIterator) Error() error {
	return iter.iter.Error()
}

// Key returns the key of the current key/value pair, or nil if done. The caller
// should not modify the contents of the returned slice, and its contents may
// change on the next call to Next.
func (iter *pebbleIterator) Key() []byte {
	return iter.iter.Key()
}

// Value returns the value of the current key/value pair, or nil if done. The
// caller should not modify the contents of the returned slice, and its contents
// may change on the next call to Next.
func (iter *pebbleIterator) Value() []byte {
	return iter.iter.Value()
}

// Release releases associated resources. Release should always succeed and can
// be called multiple times without causing error.
func (iter *pebbleIterator) Release() {
	if !iter.released {
		iter.iter.Close()
		iter.released = true
	}
}
```

