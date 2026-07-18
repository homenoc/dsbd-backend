package core

import "fmt"

// GroupStatus is the onboarding state of a group's application. It is currently
// DERIVED (not a stored column) from Pass/AddAllow/service state; these constants
// exist so the Slack transition strings are generated, not hand-written.
type GroupStatus uint

const (
	StatusServiceInput    GroupStatus = 1 // ネットワーク情報記入段階(User)
	StatusExamination     GroupStatus = 2 // 審査中
	StatusConnectionInput GroupStatus = 3 // 接続情報記入段階(User)
	StatusOpening         GroupStatus = 4 // 開通作業中
)

// Label returns the Japanese label for a status.
func (s GroupStatus) Label() string {
	switch s {
	case StatusServiceInput:
		return "ネットワーク情報記入段階(User)"
	case StatusExamination:
		return "審査中"
	case StatusConnectionInput:
		return "接続情報記入段階(User)"
	case StatusOpening:
		return "開通作業中"
	default:
		return ""
	}
}

// TransitionText renders the "1[label] =>2[label]" history string used in Slack
// status notifications, replacing the hand-written literals in the handlers.
func TransitionText(from, to GroupStatus) string {
	return fmt.Sprintf("%d[%s] =>%d[%s]", from, from.Label(), to, to.Label())
}
