# Basic Architecture
主要由以下几个服务组成。项目的目录命名以Service前面的英文名称进行命名。
* AccountService: 账号服务，用于用户账号管理。详细设计见`account_architecture.md`
* AuthService: 授权服务，用于管理用户的登录Token的生命周期.
* UserInfoService: 用户信息服务，用于用户信息的维护。
* CommodityService：商品服务，用于商品的管理。
* OrderService：订单服务，用于订单的管理。
* ChatService：客服服务，用于顾客与商家售后服务。
* EvaluateService：评价服务，用于顾客购物后评价的管理。
* ExpressService：快递服务，用于商品完成购买后快递过程的管理。