/**
 * @Author: 夜央 Oh oh oh oh oh oh (https://github.com/togettoyou)
 * @Email: zoujh99@qq.com
 * @Date: 2020/3/4 2:00 下午
 * @Description: 定义的数据结构体、常量
 */
package lib

//账户，虚拟管理员和若干业主账号
type Account struct {
	AccountId string  `json:"accountId "` //账号ID
	UserName  string  `json:"userName"`   //账号名
	Balance   float64 `json:"balance"`    //余额
}

//房地产可以作为担保(当Encumbrance为true时)出售、捐赠或质押。
type RealEstate struct {
	RealEstateID string  `json:"realEstateId"` //房地产ID
	Proprietor   string  `json:"proprietor"`   //所有者
	Encumbrance  bool    `json:"encumbrance"`  //是否作为担保
	TotalArea    float64 `json:"totalArea"`    //总面积
	LivingSpace  float64 `json:"livingSpace"`  //生活空间
	Owner        string  `json:"owner"`        //业主(业主AccountId)
}

//销售要约
//需要确定ObjectOfSale是否属于Seller
//买家初始为空
//可以将Seller和Buyer一起作为复合键
type Selling struct {
	ObjectOfSale  string  `json:"objectOfSale"`  //销售对象(正在出售的房地产RealEstateID)
	Seller        string  `json:"seller"`        //卖家(卖家AccountId)
	Buyer         string  `json:"buyer"`         //买家(买家AccountId)
	Price         float64 `json:"price"`         //价格
	SalePeriod    int     `json:"salePeriod"`    //智能合约的有效期(单位为天)
	SellingStatus string  `json:"sellingStatus"` //销售状态
}

//销售状态
var SellingStatusConstant = func() map[string]string {
	return map[string]string{
		"saleStart": "销售中", //正在销售状态,等待买家光顾
		"cancelled": "已取消", //被卖家取消销售或销售期限到期或买家退款操作导致取消
		"delivery":  "交付中", //买家买下并付款,处于等待卖家确认收款状态,如若卖家未能确认收款，买家可以取消并退款
		"done":      "完成",  //卖家确认接收资金，交易完成
	}
}

//捐赠要约
//需要确定ObjectOfDonating是否属于Donor
//需要指定受赠人Grantee，并等待受赠人同意接收
type Donating struct {
	ObjectOfDonating string `json:"objectOfDonating"` //捐赠对象(正在捐赠的房地产RealEstateID)
	Donor            string `json:"donor"`            //捐赠人(捐赠人AccountId)
	Grantee          string `json:"grantee"`          //受赠人(受赠人AccountId)
	DonatingStatus   string `json:"donatingStatus"`   //捐赠状态
}

//捐赠状态
var DonatingStatusConstant = func() map[string]string {
	return map[string]string{
		"donatingStart": "捐赠中", //捐赠人发起捐赠合约，等待受赠人确认受赠
		"cancelled":     "已取消", //捐赠人在受赠人确认受赠之前取消捐赠或受赠人取消接收受赠
		"done":          "完成",  //受赠人确认接收，交易完成
	}
}

//质押要约
//业主将房地产质押给质押人
//需要确定ObjectOfPledge是否属于Owner
//需要指定质押人Pledge，并等待质押人转移质押金
//质押人转移质押金,业主确认收到资金,质押期开始
//如果期限到期,且业主尚未返还资金,质押对象将成为质押人的财产
//业主可以随时取消报价,直到他们确认收到资金
//质押人可在业主确认收到上述资金之前,随时提取资金
type Pledge struct {
	ObjectOfPledge string  `json:"objectOfPledge"` //质押对象(正在质押的房地产RealEstateID)
	Owner          string  `json:"owner"`          //业主(业主AccountId)
	Pledge         string  `json:"pledge"`         //质押人(质押人AccountId)
	Money          float64 `json:"money"`          //质押金额
	PledgePeriod   int     `json:"pledgePeriod"`   //质押期(单位为天)
	PledgeStatus   string  `json:"pledgeStatus"`   //质押状态
}

//质押状态
var PledgeStatusConstant = func() map[string]string {
	return map[string]string{
		"pendingPayment":      "待转账", //业主发起质押合约，等待质押人转移质押金
		"pendingConfirmation": "待确认", //质押人转移资金，等待业主确认收到资金
		"pledgedStart":        "质押中", //业主确认收到资金，正在质押状态
		"cancelled":           "已取消", //业主和质押人可以在业主确认收到资金前取消质押，处于质押状态过程不可再取消
		"done":                "已完成", //质押期限到期，完成质押，如若业主尚未返还资金,质押对象将成为质押人的财产
	}
}
