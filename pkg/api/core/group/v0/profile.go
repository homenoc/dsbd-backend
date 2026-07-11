package v0

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	"github.com/homenoc/dsbd-backend/pkg/api/middleware"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
)

// Get returns the caller's group profile and member list (GET /group).
func Get(c *gin.Context) {
	u := middleware.CurrentUser(c)

	// No group yet: the member list is just the user themselves.
	if u.GroupID == nil {
		c.JSON(http.StatusOK, gin.H{
			"group":     group.Profile{},
			"user_list": []user.Profile{user.ProfileFrom(u)},
		})
		return
	}

	groupResult := dbGroup.GetByID(*u.GroupID)
	if groupResult.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: groupResult.Err.Error()})
		return
	}
	if len(groupResult.Group) == 0 {
		c.JSON(http.StatusOK, gin.H{"group": group.Profile{}, "user_list": []user.Profile{}})
		return
	}
	g := groupResult.Group[0]

	membership, err := core.GetMembershipTypeID(g.MemberType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	// isExpired (課金確認)
	isExpired := false
	if core.IsPaidMemberType(g.MemberType) && g.MemberExpired != nil {
		jst, err := time.LoadLocation("Asia/Tokyo")
		if err != nil {
			panic(err)
		}
		nowJST := time.Now().In(jst)
		if nowJST.Unix() > g.MemberExpired.Add(time.Hour*24).Unix() {
			isExpired = true
		}
	} else if core.IsPaidMemberType(g.MemberType) && g.MemberExpired == nil {
		isExpired = true
	}

	// isStripeID
	isStripeID := true
	if g.StripeCustomerID == nil || *g.StripeCustomerID == "" ||
		g.StripeSubscriptionID == nil || *g.StripeSubscriptionID == "" {
		isStripeID = false
	}

	couponID := ""
	if g.CouponID != nil {
		couponID = *g.CouponID
	}

	profile := group.Profile{
		ID:            g.ID,
		Pass:          g.Pass,
		ExpiredStatus: g.ExpiredStatus,
		IsExpired:     isExpired,
		IsStripeID:    isStripeID,
		MemberTypeID:  membership.ID,
		MemberType:    membership.Name,
		MemberExpired: g.MemberExpired,
		CouponID:      couponID,
	}
	if core.CanManageServices(u.Level) {
		profile.Agree = g.Agree
		profile.Question = g.Question
		profile.Org = g.Org
		profile.OrgEn = g.OrgEn
		profile.PostCode = g.PostCode
		profile.Address = g.Address
		profile.AddressEn = g.AddressEn
		profile.Tel = g.Tel
		profile.Country = g.Country
		profile.Contract = g.Contract
		profile.AddAllow = g.AddAllow
	}

	var userList []user.Profile
	if core.CanViewGroup(u.Level) {
		for _, member := range g.Users {
			userList = append(userList, user.ProfileFrom(member))
		}
	}

	c.JSON(http.StatusOK, gin.H{"group": profile, "user_list": userList})
}
