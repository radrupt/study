mysql是公司维护的数据库
pg是社区维护的数据库

两者都是优秀的开源数据库。

mysql能做的事情pg都能做，但pg能做的事情mysql不一定能做。
比如大数据

pg免费，开源协议运行二次售卖
不仅仅是关系型数据库，可存储json，jsonb，array
支持地理信息处理扩展
索引更多：btree，gist，gin，部分索引，hash，联合索引，mysql只有B+tree索引、Hash索引，Full-text索引
对字符更友好，msyql需要设置utf8mb4才能显示emoji
事务隔离做的更好，虽然pg，mysql都是默认repeatable read，但是pg表自带version可保证并发更新的正确性

## 应用层面？

## 索引支持？

## 高可用实现

## 分库分表支持

## 存储数据类型