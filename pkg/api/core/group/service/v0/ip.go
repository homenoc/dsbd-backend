package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service"
	"time"
)

func ipProcess(restrict bool, input []service.IPInput) ([]core.IP, error) {
	var net []core.IP

	for _, tmpIP := range input {
		if err := ipCheck(restrict, tmpIP); err != nil {
			return nil, err
		}

		startDate, _ := time.Parse("2006-01-02", tmpIP.StartDate)
		var endDate *time.Time = nil
		if tmpIP.EndDate != nil {
			tmpEndDate, _ := time.Parse("2006-01-02", *tmpIP.EndDate)
			endDate = &tmpEndDate
		}

		net = append(net, core.IP{
			Version:   tmpIP.Version,
			Name:      tmpIP.Name,
			IP:        tmpIP.IP,
			Plan:      tmpIP.Plan,
			StartDate: startDate,
			EndDate:   endDate,
			UseCase:   tmpIP.UseCase,
			Open:      &[]bool{false}[0],
		})
	}
	return net, nil
}
