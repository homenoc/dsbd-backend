package core

//
// constant
//

// Membership
// 1-49 一般支払い
// 70-89 特別枠(無料)
// 90-99 特別枠(運営)
// 1: 一般会員
// 40: 運営委員(有償)
// 70: 学生会員
// 90: 運営委員(無償)
// 99: ""
var MemberTypes = []ConstantMembership{
	MemberTypeStandard,
	MemberTypeStudent,
	MemberTypeCommittee,
	MemberTypeCommitteeFree,
	MemberTypeDisable,
}
var MemberTypeStandard = ConstantMembership{ID: 1, Name: "一般会員"}
var MemberTypeCommittee = ConstantMembership{ID: 40, Name: "運営委員(有償)"}
var MemberTypeStudent = ConstantMembership{ID: 70, Name: "学生会員"}
var MemberTypeCommitteeFree = ConstantMembership{ID: 90, Name: "運営委員(無償)"}
var MemberTypeDisable = ConstantMembership{ID: 99, Name: ""}

// Payment Type
const PaymentMembership = 1
const PaymentDonate = 2

//
// constant struct
//

type ConstantMembership struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
