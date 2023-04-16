package web

import (
	"fmt"
	"net/http"

	"powerTrading/web/controller"
)

// 启动Web服务并指定路由信息
func WebStart(app controller.Application) {

	// fs := http.FileServer(http.Dir("web/static"))
	// http.Handle("/static/", http.StripPrefix("/static/", fs))
	// http.HandleFunc("/upload", app.UploadFile)

	http.HandleFunc("/register", app.UserRegister) // 用户注册 ✔️ ✔️
	http.HandleFunc("/login", app.UserLogin)       // 用户登陆 ✔️ ✔️

	// 涉及链码的操作
	http.HandleFunc("/queryUserBaseInfo", app.QueryUserBaseInfo)               // 查询用户基本信息(username, role, balance) ✔️ ✔️
	http.HandleFunc("/userTopUp", app.UserTopUp)                               // 用户充值余额 ✔️ ✔️
	http.HandleFunc("/publishBuyOrder", app.PublishBuyOrder)                   // 购电用户发布求购订单 ✔️ ✔️
	http.HandleFunc("/querySelfBuyOrder", app.QuerySelfBuyOrder)               // 查询购电用户自身发布的求购订单 ✔️ ✔️
	http.HandleFunc("/sellElectricity", app.SellElectricity)                   // 售卖电量资产 ✔️ ✔️
	http.HandleFunc("/queryAllOrder", app.QueryAllOrder)                       // 查询所有已发布交易信息,电力中心显示 ✔️ ✔️
	http.HandleFunc("/queryUserSelfElectricity", app.QueryUserSelfElectricity) // 查询用户订单资产信息️️️️️ ✔️  ✔️
	http.HandleFunc("/payOrder", app.PayOrder)                                 // 用户支付订单 ✔️ ✔️
	http.HandleFunc("/buyUserCancelOrder", app.BuyUserCancelOrder)             // 购电用户撤销交易 ✔️ ✔️
	http.HandleFunc("/sellUserCancelOrder", app.SellUserCancelOrder)           // 供电商撤销交易 ✔️ ✔️
	http.HandleFunc("/findOrderRecord", app.FindOrderRecord)                   // 购电用户查看交易记录 ✔️
	http.HandleFunc("/findAssetRecord", app.FindAssetRecord)                   // 供电商查看交易记录 ✔️

	// 管理员的操作
	http.HandleFunc("/queryAllUser", app.QueryAllUser)                                 // 查询区块网络上所有的用户 ✔️ ✔️
	http.HandleFunc("/queryUserOrder", app.QueryUserOrder)                             // 根据ID查询用户发布的订单 ✔️
	http.HandleFunc("/displayAllOrderWithTraceBack", app.DisplayAllOrderWithTraceBack) // 管理员显示所有交易(溯源) ✔️ ✔️
	http.HandleFunc("/traceBackOrder", app.TraceBackOrder)                             // 溯源操作 ✔️ ✔️

	http.HandleFunc("/queryAllRecords", app.QueryAllRecords) // 查询所有历史记录

	var port string = ":3000"
	fmt.Printf("启动Web服务, 监听端口号为%s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Web服务启动失败: %v\n", err)
	}
}
