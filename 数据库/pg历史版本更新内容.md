## postgresql 13更新

### 新特性
#### 逻辑复制支持分区表
PostgreSQL 10 版本开始支持逻辑复制，在12版本之前逻辑复制仅支持普通表，不支持分区表，如果需要对分区表进行逻辑复制，需单独对所有分区进行逻辑复制。
PostgreSQL 13 版本的逻辑复制新增了对分区表的支持，如下:
1. 可以显式地发布分区表，自动发布所有分区。
2. 从分区表中添加/删除分区将自动从发布中添加/删除。
具体是通过CREATE SUBSCRIPTION新增publish_via_partition_root选项支持异构分区表间的数据逻辑复制

#### 新增内置函数Gen_random_uuid()生成UUID数据
之前是需要创建外部扩展 uuid-ossp来生成，现在内置。

### 性能提升
#### Btree索引优化(引入Deduplication技术)
Deduplication技术的引入能够减少索引的存储空间和维护开销，同时提升查询效率。
#### 支持增量排序(Incremental Sorting)
PostgreSQL 13 版本的一个重要特性是支持增量排序(Incremental Sorting)，加速数据排序，例如以下SQL:  
`SELECT * FROM t ORDER BY a,b LIMIT 10;  `  
如果在字段a上建立了索引，由于索引是排序的，查询结果集的a字段是已排序的，这种场景下,PostgreSQL 13 的增量排序可以发挥重要作用，大幅加速查询，因为ORDER BY a,b中的字段a是已排序好的，只需要在此基础上对字段b进行批量排序即可。
#### 轻松连接分区表
通过分区，您可以使用范围，列表或哈希键对大型数据集进行分段，从而创建响应速度更快的数据库。分区会将一个表拆分为多个表，并以对客户端应用程序透明的方式完成，从而可以更快地访问所需的数据。以前，PostgreSQL仅在分区表具有匹配的分区边界时才允许您有效地连接分区表。  
在PostgreSQL 13中，分区连接允许你高效地连接表，即使它们的分区边界不完全匹配。这样做的好处是，连接分区表的速度更快，这鼓励了分区的使用，从而提高了数据库的响应速度。

## postgresql 12更新
### 新功能
#### 支持 SQL/JSON path 特性
可更灵活的查询jsonb数据
#### 支持 Generated Columns 特性
PostgreSQL 12 一个给力SQL特性是增加了对 Generated Columns 的支持，这个特性并不陌生，MySQL 已经支持这个特性。这个特性对分析类场景比较有用。
#### 新增 Pluggable Table Storage Interface
MySQL 支持多种存储引擎，例如 InnoDB、MyISAM、Memory 存储引擎等，现阶段 PostgreSQL 只提供 Heap 一种存储引擎。  
PostgreSQL 12 版本的一个重量级特性是引入了 Pluggable Table Storage Interface，为后续支持多种存储引擎奠定了基础，比如 zheap、Memory、columnar-oriented 等存储引擎。

### 性能优化
12 版本性能提升主要体现在分区表性能提升、CTE 支持 Inlined With Queries、Btree 索引性能提升等
#### B树增强
通过更有效地利用空间，多列索引大小最多可减少40％，从而节省了磁盘空间。
#### 分区表DML性能大辐提升
从 10 版本开始 PostgreSQL 的分区表在功能方面得到完善，到了 11 版本，在运维和开发方面得到了大幅增强，但是，当分区表分区数量较大时，分区表的DML性能并不好。  
PostgreSQL 12 版本的分区表在性能方面得到了大辐提升，尤其当分区表的分区数量非常多时，DML 性能提升更加明显。
#### 分区表数据导入性能提升
#### CTE 支持 Inlined With Queries
PostgreSQL 12 版本的 CTE 支持 Inlined WITH Queries 特性，由于 WITH 查询语句的条件可以外推到外层查询，避免中间结果数据产生，同时使用相关索引，从而大辐提升 CTE 性能。


## postgresql 11更新
### 提高分区的健壮性和性能
1. 提增加了通过 hash key 对数据进行分区的能力
2. PostgreSQL 11 为与分区键不匹配的数据引入了一个默认分区
3. 如果更新行的分区键，PostgreSQL 11 还支持自动将该行移动到正确的分区。
4. PostgreSQL 11 通过使用新的分区消除策略提高了从分区读取时的查询性能
5. 分区表支持创建主键、外键、索引
### 存储过程中支持事务
PostgreSQL 11 版本一个重量级新特性是对存储过程的支持，同时支持存储过程嵌入事务，存储过程是很多 PostgreSQL 从业者期待已久的特性，尤其是很多从Oracle转到PostgreSQL朋友。
### 查询并行性能提升
支持并行创建索引、并行Hash Join

## postgresql 10更新
### 内置分区表
PostgreSQL10 一个重量级新特性是支持分区表，在这之前，PostgreSQL不支持内置分区表。
### 逻辑复制
介绍Logical Replication之前，先介绍下Streaming Replication，中文常称之为流复制，流复制最早在 9.0 版本出现 ，生产环境使用非常普遍，常用在高可用、读写分离场景，流复制是基于 WAL 日志的物理复制，适用于实例级别的复制，而Logical Replication 属于逻辑复制，可基于表级别复制，是一种粒度可细的复制，主要用在以下场景：
1. 满足业务上需求，实现某些指定表数据同步；
2. 报表系统，采集报表数据；
3. PostgreSQL 跨版本数据同步
4. PostgreSQL 大版本升级
### 并行功能增强
### 全文检索支持JSON和JSONB数据类型


## 参考
https://postgres.fun/20190724143200.html