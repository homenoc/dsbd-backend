package chat

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
)

type ResultDatabase struct {
	Err  error
	Chat []core.Chat
}
