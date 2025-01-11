# 说明

go 中框架太多了，没有大一统框架，杂七杂八的都用过，没有特别熟练的，jwt 这种大概原理都清楚，真要写还是要查查怎么做，我用的 go-zero 不巧的是网上的教程都说的不那么明白，内心又着急，越着急越写不对，最后时间都一点点浪费了

我觉得这个面试题挺好的，虽然现场没做完，我还是希望能完成它，我觉得这个题能很好的展示自己思考全面性，虽然我写的慢了，但是我写好了就不用改了，哈哈

llm 我搜了不少，还是没找到免费的，我觉得这个也不是重点，所以只使用 google-translate 的 api 做了实现，但实际上未付费调用都会失败

# 设计思路

## 核心

llm 是收费服务，价格昂贵，尽量减少 llm 调用，这是第一优先级，贯穿整个实际

我们的设计，几乎没有针对 llm 多余的调用

## 

1. 幂等
   1. 创建用户：使用 user name 数据库 unique key 实现
   2. 创建任务：对文件内容进行 md5 存入 redis 中 5s 过期，需要原子操作，代码中简单实现未进行原子操作，仅用于防止短时间内因网络导致的重试，业务上本身允许同样的文件多次翻译
   3. 执行任务：使用数据库类似乐观锁的机制
2. 高可用
   1. monitor 任务可以使用多个实例，利用分布式锁实现主备模式，但是考虑到并非真实场景，只做说明，未实现
3. 并发问题
   1. 幂等设计规避了 创建用户，创建任务的并发问题
   2. monitor 重试任务与用户执行任务时候存在并发调用 llm 的问题，通过对状态的处理，消除了这种情况
4. 完整性
   1. 因为网络原因，此类服务一定要有一个 monitor 服务做兜底，对因为各种原因导致失败的任务进行重试

# monitor 完整性校验服务说明

接口设计的有点奇怪，创建任务和执行任务分成了两个接口，分成两个而不是用一个说明应该是异步处理任务，既然是异步的那么执行接口暴露给用户，似乎又失去了异步的意义，现场做题过程思考了很久，总是感觉哪里不对劲，any way，我就把他当成异步任务来处理，并且 llm 服务也当成是异步的，这样更又有趣，因为异步任务，所以就一定要有这个 monitor 的服务对失败的任务进行重试，因为有了这个重试，所以可能与用户手动执行存在并发行为

### task state

1. wait：上传任务，等待执行的状态
2. doing：交给 llm 进行翻译，正在翻译中的状态
3. success：翻译成功状态
4. failed：翻译失败状态

### 用户行为：

用户针对 wait，failed 状态可以手动点击重试

### monitor 服务

主要用来针对失败的任务，进行重试，重试也会有个次数限制，超过就不在重试，人工处理（这时候认为是一些未知原因，导致的失败，重试多少次都会失败）

##### 说明

1. wait 状态，按照我的理解，业务上这个状态是需要用户手动触发执行的，我们系统不能对其做处理
2. failed：出现了 failed，说明用户手动点击过了执行任务，这种情况我们需要帮助用户进行重新执行任务
3. doing：这个状态从业务逻辑上来讲是 llm 正在处理，但是由于分网络原因，有可能虽然显示 doing，但是实际上任务可能是未执行，正在执行，失败结束，成功结束，只是这个状态更新失败了
4. doing 状态的这种不确定性原因跟代码实现有关，但是无论怎么实现，一定会存在这种不确定性的状态，刚好我得实现方式这个不确定性的状态是 doing

##### 处理逻辑

1. failed： 这是一个明确的失败状态，不存在类似 doing 的不明确状态，所以monitor 发现这种状态会重新执行任务，用户也可以手动重新触发这个状态的任务，那么这里就存在并发问题，代码里在数据库层使用了类似乐观锁的技术，消除了并发，保证了只有一个操作会被执行
2. doing：因为 doing 的不确定性，所以做如下处理
   1. 实际是 success：通过检查是否产生结果文件（我们要的翻译后的文件），可以确定实际是 success，如果是这种情况，我们直接就更新 数据库，不在调用 llm
   2. 实际是 failed：同 doing
   3. 实际确实在 doing：failed 和 doing 我们没准确的办法区分（假设 llm 未提供查询状态服务，如果提供了实际也可以区分，我们这里认为没提供），这种情况只能根据时间戳与阈值比较，比如根据经验认为超过 60 s，就认为任务已经 failed





# 测试方法：

1. 测试文件在 etc/test.csv 下

2. db,redis 等配置在 etc/babelfish-api.yaml 中

3. 手动在创建数据库 create database babel_fish
4. 创建任务，执行任务，执行任务过程中 llm 随机成功失败，失败后，monitor 会对失败任务发起重试，观察数据库task 的 state

4. state 说明：
   1. 0：任务等待用户手动执行
   2. 1：任务正在执行中
   3. 2：任务成功
   4. 3：任务失败（会重试）
   5. 4：任务永久失败（不会重试，因为达到最大 llm 调用次数）
5. 具体测试方法
   1. 快速连续创建任务，通过数据库，可以观察到只会创建一条任务（去重）
   2. 创建多条任务，并且手动调用执行任务接口触发（llm 随机失败），使用 select id, state, llm_call_count,created_at, updated_at from tasks ;  观察 state，llm_call_count 会发现最终 state = 2 或者 4，对于 4 的情况，llm_call_count 一定是等于 配置中配置的最大值 
   3. 对于网络原因造成的问题，可以手动修改 state，删除 task/result 下面对应的文件模拟



创建账号

curl --location 'http://0.0.0.0:8888/auth/users' \

--header 'Content-Type: application/json' \

--data '{

​    "user_name": "likun",

​    "user_pwd": "123"

}'

登录账号

curl --location 'http://0.0.0.0:8888/auth/login' \

--header 'Content-Type: application/json' \

--data '{

​    "user_name": "likun",

​    "user_pwd": "123"

}'



创建任务

curl --location 'http://0.0.0.0:8888/tasks' \

--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzY2Mjk1MjcsImlhdCI6MTczNjYyNTkyNywidXNlcklEIjoxfQ.2OYf91wyehp35oEBmmGLpL8SRLii-f06uqlGHRsqLC4' \

--form 'file=@"/Users/zhengwei/Downloads/test.csv"'

执行任务

curl --location --request POST 'http://0.0.0.0:8888/tasks/34/translate' \

--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzY2Mjk1MjcsImlhdCI6MTczNjYyNTkyNywidXNlcklEIjoxfQ.2OYf91wyehp35oEBmmGLpL8SRLii-f06uqlGHRsqLC4' \

--data ''