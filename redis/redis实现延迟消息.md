## 使用zset和list实现

### zset负责将任务时间作为score排序，这样携程A可以定期取出过期的任务

### 将过期的任务按照score依次放入list，携程B定期获取实现FIFO顺序消费
可利用list的BLPOP LIST1 TIMEOUT命令，避免频繁的扫描

### 如果要多机消费，则需支持分布式锁