package core

// Domain enums replacing raw integer literals scattered across handlers.
// The underlying numeric values are FROZEN for DB/wire compatibility — only
// the comparisons in code change, not the stored or serialized numbers.

// UserLevel is User.Level. Lower value = more privilege within a group.
type UserLevel = uint

const (
	LevelMaster UserLevel = 1 // グループ内の申請・変更・閲覧が可能
	LevelMember UserLevel = 2 // 一般メンバー
)

// CanManageServices reports whether a level may create/modify services and
// connections. Replaces the raw `Level > 2` guard.
func CanManageServices(level UserLevel) bool {
	return level <= LevelMember
}

// ExpiredStatus is User.ExpiredStatus / Group.ExpiredStatus.
type ExpiredStatus = uint

const (
	ExpiredNone         ExpiredStatus = 0 // 有効
	ExpiredByMaster     ExpiredStatus = 1 // Masterアカウントによる停止
	ExpiredByCommittee  ExpiredStatus = 2 // 運営委員会による停止
	ExpiredReviewFailed ExpiredStatus = 3 // 審査不合格による停止
)

// ExpiredMessage returns the user-facing reason for a group's expired status.
// Mirrors the strings previously inlined in auth.GroupAuthorization.
func ExpiredMessage(s ExpiredStatus) string {
	switch s {
	case ExpiredByMaster:
		return "error: discontinued by Master Account"
	case ExpiredByCommittee:
		return "error: discontinuation by the steering committee"
	case ExpiredReviewFailed:
		return "error: discontinuation due to failed review"
	default:
		return ""
	}
}

// TokenTier is Token.Status (validity duration tier). Values preserved.
type TokenTier = uint

const (
	TokenInit TokenTier = 0  // initToken (30m)
	Token30m  TokenTier = 1  // 30m
	Token6h   TokenTier = 2  // 6h
	Token12h  TokenTier = 3  // 12h
	Token30d  TokenTier = 10 // 30d
	Token180d TokenTier = 11 // 180d
)

// MemoType is Memo.Type. Values preserved.
type MemoType = uint

const (
	MemoImportant MemoType = 1 // Important(Red)
	MemoComment1  MemoType = 2 // Comment1(Blue)
	MemoComment2  MemoType = 3 // Comment2(Gray)
)
