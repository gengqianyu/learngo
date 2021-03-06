package main

//一个非常重要的概念，channel 是负责 goroutine 和 goroutine 通信的，调度器是负责怎样把大量 goroutine 调度到内核线程上去执行。
//channel 本身和 goroutine 调度器没有半毛钱关系
//单线程多线程调度器都是基于全局锁来实现的效率不高
/*
任务窃取调度器 #
2012 年 Google 的工程师 Dmitry Vyukov 在 Scalable Go Scheduler Design Doc 中指出了现有多线程调度器的问题并在多线程调度器上提出了两个改进的手段：

	1.在当前的 G-M 模型中引入了处理器 P，增加中间层；
	2.在处理器 P 的基础上实现基于工作窃取的调度器；

基于任务窃取的 Go 语言调度器使用了沿用至今的 G-M-P 模型，我们能在 runtime: improved scheduler 提交中找到任务窃取调度器刚被实现时的源代码，
调度器的 runtime.schedule:779c45a 在这个版本的调度器中反而更简单了：

static void schedule(void) {
    G *gp;
 top:
    if(runtime·gcwaiting) {
        gcstopm();
        goto top;
    }

    gp = runqget(m->p);
    if(gp == nil)
        gp = findrunnable();

    ...

    execute(gp);
}

	1.如果当前运行时在等待垃圾回收，调用 runtime.gcstopm:779c45a 函数；
	2.调用 runtime.runqget:779c45a 和 runtime.findrunnable:779c45a 从本地或者全局的运行队列中获取待执行的 Goroutine；
	3.调用 runtime.execute:779c45a 在当前线程 M 上运行 Goroutine；

当前处理器 p 本地的运行队列中没有 Goroutine 时，调用 runtime.findrunnable:779c45a 会触发工作窃取，从其它的处理器 p 的队列中随机获取一些 Goroutine。

运行时 G-M-P 模型中引入的处理器 P, 是线程 m 和 Goroutine 的中间层，我们从它的结构体中就能看到处理器 p 与 M 和 G 的关系：

struct P {
	Lock;

	uint32	status;
	P*	link;
	uint32	tick;
	M*	m;
	MCache*	mcache;

	G**	runq;
	int32	runqhead;
	int32	runqtail;
	int32	runqsize;

	G*	gfree;
	int32	gfreecnt;
};

处理器 p 持有一个由可运行的 Goroutine 组成的环形的运行队列 runq，还反向持有一个线程。
调度器在调度时会从处理器的队列中选择队列头的 Goroutine 放到线程 M 上执行。
如下所示的图片展示了 Go 语言中的线程 M、处理器 P 和 Goroutine 的关系。

		M<--p<--G-G-G-G
			↑
			G
基于工作窃取的多线程调度器将每一个线程绑定到了独立的 CPU 上，这些线程会被不同处理器 P 管理，
不同的处理器 P 通过工作窃取对任务进行再分配实现任务的平衡，也能提升调度器和 Go 语言程序的整体性能，今天所有的 Go 语言服务都受益于这一改动。

抢占式调度器 #
对 Go 语言并发模型的修改提升了调度器的性能，但是 1.1 版本中的调度器仍然不支持抢占式调度，程序只能依靠 Goroutine 主动让出 CPU 资源才能触发调度。
Go 语言的调度器在 1.2 版本4中引入基于协作的抢占式调度解决下面的问题5：

	1.某些 Goroutine 可以长时间占用线程，造成其它 Goroutine 的饥饿；
	2.垃圾回收需要暂停整个程序（Stop-the-world，STW），最长可能需要几分钟的时间，导致整个程序无法工作；

1.2 版本的抢占式调度虽然能够缓解这个问题，但是它实现的抢占式调度是基于协作的，在之后很长的一段时间里 Go 语言的调度器都有一些无法被抢占的边缘情况，
例如：for 循环或者垃圾回收长时间占用线程，这些问题中的一部分直到 1.14 才被基于信号的抢占式调度解决。

基于协作的抢占式调度的工作原理：
	1.编译器会在调用函数前插入 runtime.morestack；
	2.Go 语言运行时会在垃圾回收暂停程序、系统监控发现 Goroutine 运行超过 10ms 时发出抢占请求 StackPreempt；
	3.当发生函数调用时，可能会执行编译器插入的 runtime.morestack，它调用的 runtime.newstack 会检查 Goroutine 的 stackguard0 字段是否为 StackPreempt；
	4.如果 stackguard0 是 StackPreempt，就会触发抢占让出当前线程；
这种实现方式虽然增加了运行时的复杂度，但是实现相对简单，也没有带来过多的额外开销，总体来看还是比较成功的实现，也在 Go 语言中使用了 10 几个版本。
因为这里的抢占是通过编译器插入函数实现的，还是需要函数调用作为入口才能触发抢占，所以这是一种协作式的抢占式调度。

目前的抢占式调度也只会在垃圾回收扫描任务时触发（抢占运行很长时间产生很多垃圾的goroutine的执行权），我们可以梳理一下上述代码实现的抢占式调度过程：

	1.程序启动时，在 runtime.sighandler 中注册 SIGURG 信号的处理函数 runtime.doSigPreempt；(如果不注册操作系统收到SIGURG信号会执行挂起)
	2.在触发垃圾回收的栈扫描时会调用 runtime.suspendG 挂起 Goroutine，该函数会执行下面的逻辑：
	3.将*全部*处于 _Grunning 状态的 Goroutine 标记成可以被抢占，即将 preemptStop 设置成 true；调用 runtime.preemptM 触发抢占；
	4.runtime.preemptM 会调用 runtime.signalM 向线程发送信号 SIGURG；
	5.操作系统收到信号会中断正在运行的线程并执行预先注册的信号处理函数 runtime.doSigPreempt；
	6.runtime.doSigPreempt 函数会处理抢占信号，获取当前的 SP 和 PC 寄存器并调用 runtime.sigctxt.pushCall；
runtime.sigctxt.pushCall 会修改寄存器并在程序回到用户态时执行 runtime.asyncPreempt；
汇编指令 runtime.asyncPreempt 会调用运行时函数 runtime.asyncPreempt2；
runtime.asyncPreempt2 会调用 runtime.preemptPark；
runtime.preemptPark 会修改当前 Goroutine 的状态到 _Gpreempted 并调用 runtime.schedule 让当前函数陷入休眠并让出线程，调度器会选择其它的 Goroutine 继续执行；
上述 9 个步骤展示了基于信号的抢占式调度的执行过程。除了分析抢占的过程之外，我们还需要讨论一下抢占信号的选择，提案根据以下的四个原因选择 SIGURG 作为触发异步抢占的信号7；

该信号需要被调试器透传；
该信号不会被内部的 libc 库使用并拦截；
该信号可以随意出现并且不触发任何后果；
我们需要处理多个平台上的不同信号；
STW 和栈扫描是一个可以抢占的安全点（Safe-points），所以 Go 语言会在这里先加入抢占功能8。基于信号的抢占式调度只解决了垃圾回收和栈扫描时存在的问题，它到目前为止没有解决所有问题，但是这种真抢占式调度是调度器走向完备的开始，相信在未来我们会在更多的地方触发抢占。

GO 语言运行时调度器
运行时调度器的三个重要组成部分 — 线程 M、Goroutine G 和处理器 P

	G — 表示 Goroutine，它是一个待执行的任务；
	M — 表示操作系统的线程，它由操作系统的调度器调度和管理；
	P — 表示处理器，它可以被看做运行在线程上的用户本地调度器；

我们会在这一节中分别介绍不同的结构体，详细介绍它们的作用、数据结构以及在运行期间可能处于的状态。

G #
Goroutine 是 Go 语言调度器中待执行的任务，它在运行时调度器中的地位与线程在操作系统中差不多，但是它占用了更小的内存空间，也降低了上下文切换的开销。

Goroutine 只存在于 Go 语言的运行时，它是 Go 语言在用户态提供的线程，作为一种粒度更细的资源调度单元，如果使用得当能够在高并发的场景下更高效地利用机器的 CPU。
Goroutine 与我们在前面章节提到的 defer 和 panic 也有千丝万缕的联系，每一个 Goroutine 上都持有两个分别存储 defer 和 panic 对应结构体的链表：
这就是为什么 panic 只会触发 当前 goroutine 中的 defer 的原因。
type g struct {
	_panic       *_panic // 最内侧的 panic 结构体        用于递归panic
	_defer       *_defer // 最内侧的延迟函数结构体			用于递归 defer
}

最后，我们再节选一些作者认为比较有趣或者重要的字段：

type g struct {
	m              *m
	sched          gobuf
	atomicstatus   uint32
	goid           int64
}

	.m — 当前 Goroutine 占用的线程，可能为空；
	.atomicstatus — Goroutine 的状态；
	.sched — 存储 Goroutine 的调度相关的数据；
	.goid — Goroutine 的 ID，该字段对开发者不可见，Go 团队认为引入 ID 会让部分 Goroutine 变得更特殊，从而限制语言的并发能力10；

上述四个字段中，我们需要展开介绍 sched 字段的 runtime.gobuf 结构体中包含哪些内容：

type gobuf struct {
	sp   uintptr
	pc   uintptr
	g    guintptr
	ret  sys.Uintreg
	...
}

	.sp — 栈指针；
	.pc — 程序计数器；
	.g — 持有 runtime.gobuf 的 Goroutine；
	.ret — 系统调用的返回值；发生系统调用 goroutine 就会去等待队列

这些内容会在调度器保存或者恢复上下文的时候用到，其中的栈指针和程序计数器会用来存储或者恢复寄存器中的值，
。

结构体 runtime.g 的 atomicstatus 字段存储了当前 Goroutine 的状态。

Goroutine 的状态

我们可以将这些不同的状态聚合成三种：等待中、可运行、运行中，运行期间会在这三种状态来回切换：

	.等待中：Goroutine 正在等待某些条件满足，例如：系统调用结束等，包括 _Gwaiting、_Gsyscall 和 _Gpreempted 几个状态；
	.可运行：Goroutine 已经准备就绪，可以在线程上运行，如果当前程序中有非常多的 Goroutine，每个 Goroutine 就可能会等待更多的时间，即 _Grunnable；
	.运行中：Goroutine 正在某个线程上运行，即 _Grunning；

M #
Go 语言并发模型中的 M 是操作系统线程。调度器最多可以创建 10000 个线程，但是其中大多数的线程都不会执行用户代码（可能陷入系统调用），最多只会有 GOMAXPROCS 个活跃线程能够正常运行。

在默认情况下，运行时会将 GOMAXPROCS 设置成当前机器的核数，我们也可以在程序中使用 runtime.GOMAXPROCS 来改变最大的活跃线程数。

在默认情况下，一个四核机器会创建四个活跃的操作系统线程，每一个线程都对应一个运行时中的 runtime.m 结构体。
默认的设置不会频繁触发操作系统的线程调度和上下文切换，所有的调度都会发生在用户态，由 Go 语言调度器触发，能够减少很多额外开销。

Go 语言会使用私有结构体 runtime.m 表示操作系统线程，这个结构体也包含了几十个字段，这里先来了解几个与 Goroutine 相关的字段：

type m struct {
	g0   *g
	curg *g
	...
}
Go
其中 g0 是持有调度栈的 Goroutine，curg 是在当前线程上运行的用户 Goroutine，这也是操作系统线程唯一关心的两个 Goroutine。
g0 是一个运行时中比较特殊的 Goroutine，它会深度参与运行时的调度过程，包括 Goroutine 的创建、大内存分配和 CGO 函数的执行。在后面的小节中，我们会经常看到 g0 的身影。

runtime.m 结构体中还存在三个与处理器相关的字段，它们分别是:
	1.表示正在运行代码的处理器 p、
	2.暂存的处理器 nextp
	3.执行系统调用之前使用线程的处理器 oldp,G执行系统调用会导致 P 和 M 暂时分离，当前线程会陷入系统调用等待返回，在锁被释放后，会有其他 Goroutine 抢占处理器资源。
	  当系统调用结束后，会调用操作系统的退出系统调用的函数 runtime.exitsyscall 为当前 Goroutine 重新分配资源

type m struct {
	p             puintptr
	nextp         puintptr
	oldp          puintptr
}

P #
调度器中的处理器 P 是线程和 Goroutine 的中间层，
它能提供线程 M 需要的上下文环境，
也会负责调度线程上的等待队列，通过处理器 P 的调度，每一个内核线程都能够执行多个 Goroutine，
它能在 Goroutine 进行一些 I/O 操作陷入等待时及时让出计算资源，调度其他 goroutine 的去执行，提高线程的利用率。

因为调度器在启动时就会创建 GOMAXPROCS 个处理器，所以 Go 语言程序的处理器数量一定会等于 GOMAXPROCS，这些处理器会绑定到不同的内核线程上。

触发调度
主动挂起 — runtime.gopark -> runtime.park_m
系统调用 — runtime.exitsyscall -> runtime.exitsyscall0
协作式调度 — runtime.Gosched -> runtime.gosched_m -> runtime.goschedImpl
系统监控 — runtime.sysmon -> runtime.retake -> runtime.preemptone
我们在这里介绍的调度时间点不是将线程的运行权直接交给其他任务，而是通过调度器的 runtime.schedule 重新调度。
*/
