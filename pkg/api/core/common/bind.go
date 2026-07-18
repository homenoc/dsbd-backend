package common

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
)

// dateKeys are JSON fields typed as time.Time on the domain models. The admin
// detail pages PUT the whole object back, and an empty-string date (e.g.
// "start_date":"") fails to unmarshal into time.Time, turning a status change
// into a 400. Dropping such empty date fields lets the bind succeed (the field
// stays zero and GORM's Updates skips it).
var dateKeys = []string{
	// gorm.Model fields have no json tags, so they serialize PascalCase.
	"CreatedAt", "UpdatedAt", "DeletedAt",
	"created_at", "updated_at", "deleted_at",
	"start_date", "end_date",
	"member_expired", "antisocial_check_at", "expired_at",
}

// BindJSONTolerant binds the request body into obj, first stripping empty-string
// values for known date fields so they don't break time.Time unmarshalling.
func BindJSONTolerant(c *gin.Context, obj any) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	var m map[string]json.RawMessage
	if json.Unmarshal(body, &m) == nil {
		changed := false
		for _, k := range dateKeys {
			if v, ok := m[k]; ok && string(v) == `""` {
				delete(m, k)
				changed = true
			}
		}
		if changed {
			if cleaned, mErr := json.Marshal(m); mErr == nil {
				body = cleaned
			}
		}
	}

	// Restore the body so any later reader still works, then bind.
	c.Request.Body = io.NopCloser(bytes.NewReader(body))
	return json.Unmarshal(body, obj)
}
