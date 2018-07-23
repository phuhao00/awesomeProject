package ErrorCode

const (
	None = 0

	//chat 模块
	ChatUserIdNotExsitInGroup  = 75000 //用户ID不存在于所维护的聊天频道组中
	ChaGroupIsClosing          = 75001 //聊天频道已经关闭了
	ChatIsNotSatisfy           = 75002 //没有达到聊天的条件限制
	ChatCDIsNoRingUp           = 75003 //聊天的CD没有到时间
	ChatConsumeItemIsNotEnough = 75004 //聊天需要消耗的道具不够
	ChatConsumeItemFailed      = 75005 //聊天消耗道具操作失败
)
