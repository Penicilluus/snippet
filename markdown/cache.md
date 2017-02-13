###### 缓存更新的一般套路

Cache Aside Pattern

失效：程序先从cache中获取数据，如果没有则从数据库获取，成功后更新缓存

命中：从cache中获取到数据，取到后返回

更新：先把数据存到数据库，成功后，让缓存失效（为什么不是更新数据库后，更新缓存，主要存在两个并发写操作会导致脏数据）

问题：一个获取缓存操作，查找缓存，没有找到，从数据库中查找。一更新操作，修改数据库中的值并是缓存失效，这时 第一个操作，get操作将数据库中获取的数据写到cache中，导致数据不一致。

Read/Write Through Pattern

把更新数据库（Repository）的操作由缓存自己代理

后端就是一个单一的存储，而存储自己维护自己的Cache

read through：cache aside是应用程序需要对没有在缓存中的数据，写入缓存，而read through则把写入缓存操作直接交给缓存服务，相对于应用程序是透明的

Write Behind Caching Pattern

Write Back：在更新数据的时候，只更新缓存，不更新数据库，而我们的缓存会异步地批量更新数据库。

###### 驱逐策略

- 随机替换
- 最近最少使用（将访问最少的淘汰）LFU 是一样的思路，但是是根据访问次数而不是时间
- W-TinyLFU 为了解决LFU的缓存数据是访问 过去访问次数多的数据，加入了一个时间窗口的机制



###### 缓存的并发访问

- 使用锁，但是粒度不好划分，读写锁？
- 将大块的cache map划分为若干个区
- 使用commit logs记录，在读写时不是真正改变数据而是记录，另外一个进程处理log



###### 分布式缓存

memcache的特点

- 基于C/S架构，协议简单
- 基于libevent的事件处理
- 自主内存储存处理
- 基于客户端的memcache分布式

数据存储方式：Slab Allocation，按照预先规定的大小，将分配的内存划分为固定大小的块，解决内存问题。

将分配的内层划分各种尺寸的块，把尺寸大小相同的块归类为组chunk （有点像golang里的内存分配策略）

数据过期方式：Lazy Expiration + LRU

lazy Expiration：内部不会监视缓存的过期时间，只是在get缓存是查看缓存的时间，这个被称为lazy Expiration

LRU：最近最少使用，一般使用 队列 和 hashmap，map主要是为了加速查找，队列则是为了快速添加和删除元素

基于客户端的分布式：表示memcache服务端仅仅是存储数据，存储策略，访问哪一个server需要客户端指定特定的算法

###### 缓存策略

- 缓存数据库查询 A hashed version of your query is the cache key，但是针对复杂的查询条件可能不是很适用，一些细微的数据变化也会导致，大量的缓存失效，不是很实用
- 缓存对象，将查询之后的结果缓存起来



有持久化需求、数据同步或者对数据结构和处理有高级要求的应用，选择redis，其他简单的key/value存储，选择memcache

###### 参考

- [High-Concurrency-LRU-Caching](http://openmymind.net/High-Concurrency-LRU-Caching/)
- [Design-of-a-modern-cache](http://highscalability.com/blog/2016/1/25/design-of-a-modern-cache.html)
- [Why does Facebook use delete to remove the key-value pair in Memcached instead of updating the Memcached during write request to the backend?](https://www.quora.com/Why-does-Facebook-use-delete-to-remove-the-key-value-pair-in-Memcached-instead-of-updating-the-Memcached-during-write-request-to-the-backend)
- [缓存更新的策略](http://coolshell.cn/articles/17416.html)
- [Scalability for Dummies - Part 3: Cache](http://www.lecloud.net/post/9246290032/scalability-for-dummies-part-3-cache)

