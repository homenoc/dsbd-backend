package store

import "git.bgp.ne.jp/dsbd/backend/pkg/auth"

type userResult struct {
	user auth.User
	err  error
}

type allUserResult struct {
	user []auth.User
	err  error
}
