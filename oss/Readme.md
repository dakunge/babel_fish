# 项目关键部分介绍

1. 项目编译运行
   1. docker build -t registry.light-field.tech/kun/fusion:v0.121 ./
   2. docker push registry.light-field.tech/kun/fusion:v0.121
   3. ssh yoke@192.168.62.212 (密码：yoke)
   4. cd kun
   5. 将 deployment.yaml image 字段修改为 步骤 2 的 image
   6. kubectl apply -f deployment.yaml
2. 项目中手动创建表的 sql 文件位于 internal/model/sql/all.sql
2. 采用 go zero 框架，DataFusionPlatform.api 文件中定义 api，使用 goctl api go --api DataFusionPlatform.api -dir ./ 生成代码，goctl api plugin -plugin goctl-swagger="swagger -filename data.json" -api DataFusionPlatform.api -dir . 生成 swagger 文档（将生成的 data.json 导入 apifox 供前端使用）
3. 为了接口返回值统一为 {"code":0, "msg": "ok", "data": {}} 模式，需要修改 goctl 代码生成模板中的 handle.tpl, 具体步骤为 
   1. goctl template init 
   2. 替换 handle.tpl
4. 代码 static 目录下存放的是 接口信息（swagger json），data.json 是步骤一生成的，其他的是通过 apifox 针对具体接口单独导出的，这些单独导出的文件用在 etc/data_release.csv 文件中的 文件接口字段

# 数据同步

### 外部 csv 文件

  外部数据源为 csv 文件，  csv 文件存在于 etc 目录下，目前有 company.csv，data_release.csv， profile_resident.csv， profile_visit.csv， 同步代码位于 internal/sync/ 下，基本就是读 csv，存入相应表中

### 外部接口数据

*外部数据源为接口，目前只有设备之家提供了接口，代码 internal/sync/sync_iot.go, 因为外部数据，我们需要同步到我们的系统中，这个过程需要创建对应的数据结构，数据表，目前外部接口只能同步全量更新数据，所以每次同步数据需要把就数据删除，另外数据接口比较多，所以相对来说有点繁琐，所以下面提供一个具体接口例子来解释*

接口： 

`curl --location --request GET 'http://183.47.50.83:3512/dataCenter/dataPlat/getColumnList/monthSrOrderData' \`

`--header 'Authorization: f8bb2e97e50c1a3c8d6201e352782d3f' \`

`--header 'Content-Type: application/json' \`

`--data '{`

​    `"page":1,`

​    `"size":1` 

`}'`

返回值：

`{`

​    `"status": 200,`

​    `"message": "success",`

​    `"serverTime": 1730740412778,`

​    `"record": {`

​        `"total": 1,`

​        `"list": [`

​            `{`

​                `"check_comp_qty": 45,`

​                `"req_comp_qty": 45,`

​                `"check_pass_qty": 42,`

​                `"month": 9,`

​                `"year": 2024,`

​                `"order_total": 50,`

​                `"check_pass_rate": 0.93,`

​                `"proj_id": 1,`

​                `"act_comp_qty": 45`

​            `}`

​        `],`

​        `"pageNum": 1,`

​        `"pageSize": 1,`

​        `"size": 1,`

​        `"startRow": 1,`

​        `"endRow": 1,`

​        `"pages": 1,`

​        `"prePage": 0,`

​        `"nextPage": 0,`

​        `"isFirstPage": true,`

​        `"isLastPage": true,`

​        `"hasPreviousPage": false,`

​        `"hasNextPage": false,`

​        `"navigatePages": 8,`

​        `"navigatepageNums": [`

​            `1`

​        `],`

​        `"navigateFirstPage": 1,`

​        `"navigateLastPage": 1`

​    `}`

`}`

具体同步步骤：

1. cd tools/generate

2.  `go run main.go iot MonthSrOrder '{"check_comp_qty":45,"req_comp_qty":45,"check_pass_qty":42,"month":9,"year":2024,"order_total":50,"check_pass_rate":0.93,"proj_id":1,"act_comp_qty":null}'`

   1. iot: 生成代码的 package
   2. MonthSrOrder：生成代码中 表对应的 struct 名称
   3. xxx：外部接口返回值数据

3. mv month_sr_order.go ../../internal/model/iot  (month_sr_order.go 是生成的代码，需要放到对应的目录下)

4. 生成地代码重点解释：

   `// copy to internal/svc/servicecontext.go`
   `// MonthSrOrderModel          iot.MonthSrOrderModel`
   `// monthSrOrderModel := iot.NewMonthSrOrderModel(db)`
   `// MonthSrOrderModel:          monthSrOrderModel,`
   `type MonthSrOrderModel interface {`
   	`// for sync data not for biz`
   	`AutoMigrateForSync(ctx context.Context, time string) error`
   	`CreateForSync(ctx context.Context, resp []byte, time string) error`
   	`RenameForSync(ctx context.Context, time string) error`
   `}`

   这几个接口是专门在 internal/sync/sync_iot.go 同步过程中适用的，结合代码重点看一下

   注释部分是为了方便 internal/svc/servicecontext.go 中适用，把这三行代码复制到 internal/svc/servicecontext.go 中，具体位置参照现有的代码很容易就能找到位置

5. 做完上述步骤后，在 internal/sync/sync_iot.go 文件中添加对应的 SyncTable ，这样每次部署得会后就会自动创建对应的表，以及数据同步了

### 遗留

目前每次重启会自动同步数据，同步数据接口的入口位于 main 函数中，因为还未确定数据同步周期，所以目前只在重新部署的会后同步数据，后续需要根据业务要求，不同接口采用不同周期定时同步，只需要在现有代码上加一层控制逻辑就可以啦
