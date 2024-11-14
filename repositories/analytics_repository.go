package repositories

import (
	"github.com/go-redis/redis"
	"log"
)

type AnalyticsRepository struct {
	client *redis.Client
}

func NewAnalyticsRepository(client *redis.Client) *AnalyticsRepository {
	return &AnalyticsRepository{client: client}
}

func (r *AnalyticsRepository) CountViews() (int64, error) {
	count, err := r.client.Incr("visit_count").Result()
	if err != nil {
		log.Print(err)
	}
	return count, nil
}
