package invalid_transaction

import (
	"github.com/go-pg/pg"
	"github.com/noah-blockchain/coinExplorer-tools/models"
)

type Repository struct {
	db *pg.DB
}

func NewRepository(db *pg.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Get invalid transaction by hash
func (repository Repository) GetTxByHash(hash string) *models.InvalidTransaction {
	var transaction models.InvalidTransaction

	err := repository.db.Model(&transaction).Column("FromAddress").Where("hash = ?", hash).Select()
	if err != nil {
		return nil
	}

	return &transaction
}
