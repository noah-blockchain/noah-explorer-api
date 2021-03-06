package slash

import (
	"github.com/go-pg/pg"
	"github.com/noah-blockchain/coinExplorer-tools/models"
	"github.com/noah-blockchain/noah-explorer-api/internal/events"
	"github.com/noah-blockchain/noah-explorer-api/internal/helpers"
	"github.com/noah-blockchain/noah-explorer-api/internal/tools"
)

type Repository struct {
	db *pg.DB
}

func NewRepository(db *pg.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (repository Repository) GetPaginatedByAddress(filter events.SelectFilter, pagination *tools.Pagination) []models.Slash {
	var slashes []models.Slash
	var err error

	pagination.Total, err = repository.db.Model(&slashes).
		Column("Coin.symbol", "Address.address", "Validator.public_key", "Block.created_at").
		Column("Validator.name", "Validator.description", "Validator.icon_url", "Validator.site_url").
		Apply(filter.Filter).
		Apply(pagination.Filter).
		Order("block_id DESC").
		SelectAndCount()

	helpers.CheckErr(err)

	return slashes
}
