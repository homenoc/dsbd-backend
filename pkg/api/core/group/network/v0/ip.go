package v0

import (
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	"time"
)

func ipProcess(input network.Input) (*[]network.IP, error) {
	var net []network.IP

	for _, tmpIP := range *input.IP {
		if err := ipCheck(tmpIP); err != nil {
			return nil, err
		}

		startDate, _ := time.Parse("2006-01-02", tmpIP.StartDate)
		var endDate *time.Time = nil
		if tmpIP.EndDate != nil {
			tmpEndDate, _ := time.Parse("2006-01-02", *tmpIP.EndDate)
			endDate = &tmpEndDate
		}
		net = append(net, network.IP{
			Version:   tmpIP.Version,
			IP:        tmpIP.IP,
			Plan:      tmpIP.Plan,
			StartDate: startDate,
			EndDate:   endDate,
			UseCase:   tmpIP.UseCase,
			Open:      &[]bool{false}[0],
		})
	}

	return &net, nil
}
