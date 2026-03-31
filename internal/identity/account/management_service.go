package account

import (
	"context"
	"errors"
	"time"

	"github.com/koubae/game-hangar/internal/errs"
	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/database"
	"go.uber.org/zap"
)

type ManagementService struct {
	db         database.DBTX
	repository IAccountRepository
}

type ManagementServiceFactory func(d database.DBTX, r IAccountRepository) *ManagementService

func NewManagementService(
	d database.DBTX,
	r IAccountRepository,
) *ManagementService {
	return &ManagementService{
		db:         d,
		repository: r,
	}
}

func (s *ManagementService) GetAccount(ctx context.Context, accountID string) (*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logger := common.GetLogger()

	account, err := s.repository.GetAccount(ctx, s.db, accountID)
	if err != nil {
		if !errors.Is(err, errs.ResourceNotFound) {
			logger.Error(
				"[ManagementService] error while getting account",
				zap.String("accountID", accountID),
				zap.Error(err),
			)
		}
		return nil, err
	}

	return account, nil
}
