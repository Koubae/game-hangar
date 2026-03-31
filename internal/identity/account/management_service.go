package account

import (
	"github.com/koubae/game-hangar/pkg/database"
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

func (s *ManagementService) Me(accountID string) (*Account, error) {
	return nil, nil
}
