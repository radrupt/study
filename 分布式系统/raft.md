## 什么是raft协议？
raft是一种分布式一致性协议，解决工程上的最终一致性，去中心化，高可用的性问题的一种协议方案。raft是一种共识算法，即多个节点对某个事情达成一致看法，即使部分节点出现故障，网络延时，网络分割的情况下，大部分节点仍然是能达成一致看法。

## raft协议是如何工作的呢？
raft如何工作，即raft协议是如何实现最终一致性，去中心化，高可用的。
### 如何实现去中心化？
#### 什么是去中心化？
> (摘自百度百科)在一个分布有众多节点的系统中，每个节点都具有高度自治的特征。节点之间彼此可以自由连接，形成新的连接单元。任何一个节点都可能成为阶段性的中心，但不具备强制性的中心控制功能。节点与节点之间的影响，会通过网络而形成非线性因果关系。这种开放式、扁平化、平等性的系统现象或结构，我们称之为去中心化。

raft协议通过给节点定义状态（leader，follower，candidate），结合选举算法每次选举得到阶段性的任期中心节点，从而实现去中心化
#### 定义节点三种状态    
1. Follower State
2. Candidate State
3. Leader State
#### 如何选出主节点
1. 一开始所有所有节点都是follower State
2. 系统运行中，处于Follower State的节点在一定时间内没有收到Leader 节点的心跳，则判断Leader节点出现了故障，该follower节点启动选主过程, 当前term（term是逻辑时钟，全局递增）任期+1，向其它节点发起投票请求RPC，该节点会变成candidate状态，直到选主结束，变成Follower or Leader 节点。
3. 通过prevote可以解决在主节点没问题，但是follower由于自身原因长时间和其它节点无法通信，导致自己变为candidate后term不断增加，重新加入集群时，会中断集群的问题（term比其他节点都大，不一致，中断重新一致）。方法是，通过确认节点是否能赢得集群中大多数节点的投票，如果能，则term+1，否则重新变为Follower State.
4. 考虑正常情况下，由于leader节点发生故障，那么follower节点在随机的leader超时时间范围内，将状态变为candidate，并发起投票RPC
	```
	其他节点统一投票的条件是
	1. 没有收到有效领导的心跳，至少有一次选举超时。
	2. Candidate的日志足够新（Term更大，或者Term相同raft index更大）。

	```
5. Candidate如何发送投票RPC
	```
	1. 自增当前节点的任期号
	2. 给自己投票
	3. 重置选举超时计时器
	4. 发送请求投票的RPC给其他服务器
	```
6. 节点收到请求投票的RPC如何处理
	```
	1. 判断当前Term和请求投票参数中的Term
		如果更大，则拒绝投票
		否则更新更新当前Term为请求投票参数中的Term，并将自身状态设为Follower
	2. 检测当前节点的投票状态
		如果当前节点没有给其他节点投过票，或者是投给自己过，那么继续检查日志的匹配状态
		否则，拒绝投票，因为投票采取先到先得策略
	3. 检测候选人的日志是否比当前节点日志新，通过比较候选人的lastLogIndex和lastLogIterm和当前节点日志，确保新选举出来的Leader不会丢失已经提交的日志
		如果候选人的任期比当前节点任期高，或者任期相同，但是候选人的日志比当前节点日志新，则给候选人投票
		否则，拒绝投票
	```
7. candidate收到请求投票的响应会如何处理？  
每一个候选人在每一个任期term内都会发起一轮投票，如果指定时间内，收到N/2+1个节点的同意投票，即投票成为，变为Leader节点。  
当然也可能存在该任期term内，没有选出Leader节点，即由于网络问题，或者部分节点得到了相同的投票从而导致投票失败。  
另外也可能由于其他节点发起了更新的任期投票，那么当前节点的选举失败，节点变为Follower状态.
如果收到的响应term < 当前节点term，表示这是一个过期的term，不处理，否则如果相同，投票数+1，并判断是否达到n/2+1,达到则变为哦Leader。  
至此，阶段性中心节点产生。
#### 最终一致性和高可用行
这部分主要讲最终一致性，基于最终一致性和去中心化的选主策略即可实现高可用性。  
一致性是构建具有容错性的分布式系统的基础。在一个具有一致性性质的集群中，同一时刻所有几点对存储在其中的某个值都有相同的结果，即对其共享的存储保持一致。集群具有自动恢复的性质，少数节点失效不影响集群的正常工作，当大多数节点失效，集群停止服务，而不会返回错误数据。  
一致性协议就是来提供这样能力的，一致性协议通常基于replicated State machines，即所有节点都从同一个state出发，都经过同样的一些操作序列，最后到达相同的state。  
#### 如何实现一致性协议
1. 每个节点需要有三个组件
	```
	1. 状态机：状态机会从log中取出所有的命令，然后执行一遍，得到的结果就是我们对外提供的保证了一致性的结果
	2. log：保存了所有的修改结果
	3.一致性模块：一致性模块算法就是用来保证写入的log命令一致性
	```
2. Log Replication  
当Leader被选出来后，将开始接受客户端发来的请求，每个请求包含一条需要被replicated state machines执行的命令。leader会把它作为一个log entry append到日志中，然后给其他节点发送AppendEntriesRpc请求。当大部分Follower节点将该命令写入日志，就apply这条log entry到状态机中，然后返回结果给客户端。如果某个Follower宕机或者运行的很慢，或者网络丢包，则会一致给这个Follower发AppendEntriedRPC，直到日志一致。  
当一条日志是commited，即多数节点将该命令写入日志，Leader才会将它应用到状态机中。Raft保证一条commited的log entry已经持久化并被所有节点执行。  
3. Log Replcation如何让比较老的follower state 变的和leader一样？
leader会维护每个follower的nextIndex log，leader会给每个follower发送AppendEntriesRPC，携带log的Term和nextIndex，如果follower没有发现该log，则回复leader拒绝消息，leader会将nextIndex-1，并重新发送，直到回复ok，表示找到了同步的起始位置，并从改位置开始同步log数据，各节点的状态机再通过log来更新state
4. 什么情况下会出现stale data
设a，b，c，d，e5个节点，
	```
	1. a为leader
	2. c，d，e由于网络分割，开启新的选举，c为term+1的leader
	3. 若发送修改命令给a，由于a无法得到commited，写入失败
	4. 若发送修改命令给c，由于c得到commited，写入成功
	5. 读取a，会得到老数据
	6. 当5个节点恢复通信，由于a，b term 更小，a变成follower，更新a，b的term，并开启replcated state machines
	```

---

算法用途
原理
选举
normal Operation：稳定工作期
选主

* 问题分解
* 状态简化
问题分解是将"复制集中节点一致性"这个复杂的问题划分为数个可以被独立解释、理解、解决的子问题。在raft，子问题包括，leader election， log replication，safety，membership changes。  
而状态简化更好理解，就是对算法做出一些限制，减少需要考虑的状态数，使得算法更加清晰，更少的不确定性（比如，保证新选举出来的leader会包含所有commited log entry）
* leader election
 raft协议中，一个节点任一时刻处于以下三个状态之一：
1. leader 领导者
2. follower 跟随者
3. candidate 候选人
* log replication：日志副本，即如何保证主节点和从节点的状态是一致的
当有了leader，系统应该进入对外工作期了。客户端的一切请求来发送到leader，leader来调度这些并发请求的顺序，并且保证leader与followers状态的一致性。raft中的做法是，将这些请求以及执行顺序告知followers。leader和followers以相同的顺序来执行这些请求，保证状态一致。
共识算法的实现一般是基于复制状态机（Replicated state machines），何为复制状态机：
简单来说：相同的初识状态 + 相同的输入 = 相同的结束状态。引文中有一个很重要的词deterministic，就是说不同节点要以相同且确定性的函数来处理输入，而不要引入一下不确定的值，比如本地时间等。如何保证所有节点 get the same inputs in the same order，使用replicated log是一个很不错的注意，log具有持久化、保序的特点，是大多数分布式系统的基石。

没有leader的时候，选举leader失败的时候会如何？

election timeout （选举超时）：一般是多久？

一般选举过程需要花费多久？

majority：多数
term：任期，表示第几个任期
election： 选举
normal Operation：稳定工作期
deterministic：确定性的
stale leader：老的leader问题，网络分割
Election safety：选举安全
log matching：log匹配特性
State Machine Safety：？没看懂
leader crash：如何处理？
选举人必须比自己知道的更多？怎么判断？
log replication约束？

参考：  
https://zhuanlan.zhihu.com/p/27207160
https://www.cnblogs.com/xybaby/p/10124083.html