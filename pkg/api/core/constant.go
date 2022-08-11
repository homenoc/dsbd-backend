package core

//
// constant
//

// Membership
var MemberTypeStandard = ConstantMembership{ID: 1, Name: "一般会員"}
var MemberTypeStudent = ConstantMembership{ID: 70, Name: "学生会員"}
var MemberTypeCommittee = ConstantMembership{ID: 90, Name: "運営委員"}
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
