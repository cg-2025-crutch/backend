package service

import (
	"context"
)

func (s *FundsService) DeleteTransaction(ctx context.Context, id int64, userUID string) error {
	err := s.prod.Produce(ctx, []byte(userUID), []byte("update"))
	if err != nil {
		return err
	}
	return s.repo.DeleteTransaction(ctx, id, userUID)
}
