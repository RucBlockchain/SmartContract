# SmartContract
test1测试内容：用户购票后在规定时间内投保，并且航空公司在规定时间内提交保证金；行程结束后航班延误，用户调用索赔函数得到赔偿金。
test2测试内容：用户购票后未在规定时间内投保，则直接退出脚本。
test3测试内容：用户购票后在规定时间内投保，但保险公司在规定时间内未缴纳保证金，则退回用户的保费，退出脚本。

query.go 包含一些测试用的查找函数
    queryGuarantee函数：根据GuaranteeID准确查找某一份保单；
    queryOrder函数：根据OrderID准确查找某一机票订单；
    queryState函数：根据GuaranteeID查看当前保单状态；
    queryByObjectType函数：根据ObjectType查找某一类结构体信息。

fsmEvent.go 包含与状态变迁相关的函数
    clientDeposit函数：当保单状态为Bought并触发该函数时，保单状态会更改为Deposited；
    insurCompanyDeposit函数：当保单状态为Deposited并触发该函数时，保单状态会更改为Insurancde；
    clientWithdraw函数：当保单状态为Uninsuranced并触发该函数时，保单状态会更改为Client-Withdraw；
    changeOrder函数：当保单状态为Insuranced并触发该函数时，保单状态会更改为Delayed；（航班延误时调用）
    insurCompanyWithdraw函数：当保单状态为Undelayed并触发该函数时，保单状态会更改为InsurCompany-Withdrawn；
    indemnify函数：当保单状态为Delayed并触发该函数时，保单状态会更改为Success。
    
functions.go
    check_buyTicket函数：检查用户是否已经购票。
    check_delay函数：检查航班是否延误。
    clientRegist函数：用户登记函数，将用户的相关信息及公钥信息（Client结构体）存到链上，以“client”+ClientID拼接复合键作为Client的主键；
    insurCompanyRegist函数：保险公司登记函数，将保险公司信息及公钥信息（InsuranceCompany结构体）存到链上，以“company”+InsurompanyID拼接复合键作为InsuranceCompany的主键；
    writeOrder函数：当用户购票后，将用户的购票信息（Order结构体）存到链上，以“order”+OrderID拼接复合键作为Order的主键；
    initGuarantee函数：设置保单Guarantee的初始状态为Bought；
    verifySignature函数：验证签名，验证成功返回true，验证失败返回false。
