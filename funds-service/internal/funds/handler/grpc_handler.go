package handler

import (
	"github.com/cg-2025-crutch/backend/funds-service/internal/funds/service"
	pb "github.com/cg-2025-crutch/backend/funds-service/internal/grpc/gen"
	"github.com/cg-2025-crutch/backend/funds-service/internal/models"
)

type GRPCHandler struct {
	pb.UnimplementedFundsServiceServer
	service service.FundsServicer
}

func NewGRPCHandler(service service.FundsServicer) *GRPCHandler {
	return &GRPCHandler{
		service: service,
	}
}

func (h *GRPCHandler) transactionToProto(t *models.Transaction) *pb.Transaction {
	transaction := &pb.Transaction{
		Id:              t.ID,
		UserUid:         t.UserUID,
		CategoryId:      t.CategoryID,
		Type:            t.Type,
		Amount:          t.Amount,
		Title:           t.Title,
		Description:     t.Description,
		TransactionDate: t.TransactionDate.Format("2006-01-02"),
		CreatedAt:       t.CreatedAt.Unix(),
		UpdatedAt:       t.UpdatedAt.Unix(),
	}

	if t.Category != nil {
		transaction.Category = h.categoryToProto(t.Category)
	}

	return transaction
}

func (h *GRPCHandler) categoryToProto(c *models.Category) *pb.Category {
	return &pb.Category{
		Id:        c.ID,
		Name:      c.Name,
		Type:      c.Type,
		Icon:      c.Icon,
		CreatedAt: c.CreatedAt.Unix(),
	}
}

func (h *GRPCHandler) balanceToProto(b *models.UserBalance) *pb.UserBalance {
	balance := &pb.UserBalance{
		UserUid:      b.UserUID,
		TotalBalance: b.TotalBalance,
		TotalIncome:  b.TotalIncome,
		TotalExpense: b.TotalExpense,
		UpdatedAt:    b.UpdatedAt.Unix(),
	}

	if b.LastTransactionAt != nil {
		balance.LastTransactionAt = b.LastTransactionAt.Unix()
	}

	return balance
}
