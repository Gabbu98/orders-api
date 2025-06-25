package strategy

import (
	"fmt"

	"github.com/Gabbu98/orders-api/repository/order"
)

type RepositoryStrategies struct {
	Strategies *map[string]order.OrderRepository
}

func (strats *RepositoryStrategies) GetStrategy(key string) (order.OrderRepository, error) {
	var repo order.OrderRepository = (*strats.Strategies)[key]
	if repo == nil {
		return nil, fmt.Errorf("Could not find repository with key %w", key)
	}

	return repo, nil
}
