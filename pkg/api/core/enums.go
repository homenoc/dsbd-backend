package core

// Domain enums replacing raw integer literals scattered across handlers.
// The underlying numeric values are FROZEN for DB/wire compatibility — only
// the comparisons in code change, not the stored or serialized numbers.

// UserLevel is User.Level. Lower value = more privilege within a group.
// 1 is only ever assigned at account registration; 2-4 are what a master
// assigns when inviting users to the group (user.AddGroup validates 2..4,
// and the invite dialog labels them 追加・変更・閲覧(Master) / 閲覧のみ(User)
// / 通知のみ(Guest)).
type UserLevel = uint

const (
	LevelMaster UserLevel = 1 // 初期登録者(申請・変更・閲覧)
	LevelEditor UserLevel = 2 // Masterから割当された申請・変更・閲覧権限
	LevelViewer UserLevel = 3 // グループ内の情報閲覧のみ
	LevelGuest  UserLevel = 4 // 障害情報の通知のみ
)

// CanManageServices reports whether a level may create/modify services and
// connections (levels 1-2; replaces the raw `Level > 2` guard, minus its
// acceptance of the invalid zero value).
func CanManageServices(level UserLevel) bool {
	return LevelMaster <= level && level <= LevelEditor
}

// CanViewGroup reports whether a level may read group-scoped resources
// (levels 1-3; guests only receive notifications).
func CanViewGroup(level UserLevel) bool {
	return LevelMaster <= level && level <= LevelViewer
}

// ExpiredStatus is User.ExpiredStatus / Group.ExpiredStatus.
type ExpiredStatus = uint

// Meanings follow the writers (the admin UI's 廃止 menu and the Slack
// notifier's expiredStatusText): 1=審査落ち, 2=ユーザより廃止,
// 3=運営委員より廃止. The original auth error strings had 1 and 3
// cross-wired; ExpiredMessage below is the corrected reader.
const (
	ExpiredNone         ExpiredStatus = 0 // 有効
	ExpiredReviewFailed ExpiredStatus = 1 // 審査不合格による廃止
	ExpiredByMaster     ExpiredStatus = 2 // ユーザ(Master)による廃止
	ExpiredByCommittee  ExpiredStatus = 3 // 運営委員会による廃止
)

// ExpiredMessage returns the user-facing reason for a group's expired status.
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

// ExpiredLabel returns the Japanese admin-facing label for an expired status
// (previously duplicated as expiredStatusText in the group slack notifier and
// hardcoded in the admin frontend).
func ExpiredLabel(s ExpiredStatus) string {
	switch s {
	case ExpiredNone:
		return "0"
	case ExpiredReviewFailed:
		return "審査落ち"
	case ExpiredByMaster:
		return "ユーザより廃止"
	case ExpiredByCommittee:
		return "運営委員より廃止"
	default:
		return "status不明"
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
