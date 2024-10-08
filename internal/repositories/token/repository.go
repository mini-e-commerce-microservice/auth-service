package token

import "github.com/redis/rueidis"

type repository struct {
	client rueidis.Client
}

func NewRepository(client rueidis.Client) *repository {
	return &repository{
		client: client,
	}
}
