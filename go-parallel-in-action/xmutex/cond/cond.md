## 概念

Cond 是为等待 / 通知场景下的并发问题提供支持的。它提供了条件变量的三个基本方法 Signal、Broadcast 和 Wait，为并发的 goroutine 提供等待 / 通知机制。

- signal: 允许调用者 Caller 唤醒一个等待此 Cond 的 goroutine。如果此时没有等待的 goroutine，显然无需通知 waiter；如果 Cond 等待队列中有一个或者多个等待的 goroutine，则需要从等待队列中移除第一个 goroutine 并把它唤醒
- broadcast: 允许调用者 Caller 唤醒所有等待此 Cond 的 goroutine。如果此时没有等待的 goroutine，显然无需通知 waiter；如果 Cond 等待队列中有一个或者多个等待的 goroutine，则清空所有等待的 goroutine，并全部唤醒
- wait : 会把调用者 Caller 放入 Cond 的等待队列中并阻塞，直到被 Signal 或者 Broadcast 的方法从等待队列中移除并唤醒

Cond 有三点特性是 Channel 无法替代的：

- Cond 和一个 Locker 关联，可以利用这个 Locker 对相关的依赖条件更改提供保护。
- Cond 可以同时支持 Signal 和 Broadcast 方法，而 Channel 只能同时支持其中一种。
- Cond 的 Broadcast 方法可以被重复调用。等待条件再次变成不满足的状态后，我们又可以调用 Broadcast 再次唤醒等待的 goroutine。这也是 Channel 不能支持的，Channel 被 close 掉了之后不支持再 open。

Cond 可重用是 k8s 中很多地方使用 cond 的

cond 与 WaitGroup 的区别在于， WaitGroup 内已经将数量的统计和等待条件实现了，而 cond 需要开发者自己去控制

## 注意事项

- Wait 调用前需要加锁，因为调用 Wait 后，Wait 函数内就会先释放锁，然后进入阻塞队列等待，当唤醒后会再次去抢占加锁
- 被唤醒后一定要检查条件是否真的已经满足
