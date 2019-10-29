package coins

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/helpers"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/tools"
	"github.com/noah-blockchain/noah-explorer-tools/models"
)

type Repository struct {
	DB             *pg.DB
	baseCoinSymbol string
}

func NewRepository(db *pg.DB, baseCoinSymbol string) *Repository {
	return &Repository{
		DB:             db,
		baseCoinSymbol: baseCoinSymbol,
	}
}

// Get list of coins
func (repository *Repository) GetCoins() []models.Coin {
	var coins []models.Coin

	err := repository.DB.Model(&coins).
		Column("crr", "volume", "reserve_balance", "name", "symbol", "updated_at", "a.address").
		Join("LEFT JOIN addresses AS a ON a.id = creation_address_id").
		Where("deleted_at IS NULL").
		Order("reserve_balance DESC").
		Select()

	helpers.CheckErr(err)

	return coins
}

// Get coin detail by symbol
func (repository *Repository) GetBySymbol(symbol string) []models.Coin {
	var coins []models.Coin

	err := repository.DB.Model(&coins).
		Column("crr", "volume", "reserve_balance", "name", "symbol").
		Where("symbol LIKE ?", fmt.Sprintf("%%%s%%", symbol)).
		Where("deleted_at IS NULL").
		Order("reserve_balance DESC").
		Select()
	helpers.CheckErr(err)

	return coins
}

type CustomCoinsStatusData struct {
	ReserveSum string
	Count      uint
}

// Get custom coins data for status page
func (repository *Repository) GetCustomCoinsStatusData() (CustomCoinsStatusData, error) {
	var data CustomCoinsStatusData

	err := repository.DB.
		Model(&models.Coin{}).
		ColumnExpr("SUM(reserve_balance) as reserve_sum, COUNT(*) as count").
		Where("symbol != ?", repository.baseCoinSymbol).
		Select(&data)

	return data, err
}

// Get paginated list of blocks
func (repository Repository) GetPaginated(pagination *tools.Pagination, field *string, orderBy *string) []models.Coin {
	var coins []models.Coin
	var err error
	fieldSql := "reserve_balance"
	orderBySql := "DESC"

	if field != nil {
		fieldSql = *field
	}

	if orderBy != nil {
		orderBySql = *orderBy
	}

	pagination.Total, err = repository.DB.Model(&coins).
		Column("coin.crr", "coin.volume", "coin.reserve_balance", "coin.name", "coin.symbol", "coin.updated_at", "a.address").
		Join("LEFT JOIN addresses AS a ON a.id = coin.creation_address_id").
		Where("coin.deleted_at IS NULL").
		Apply(pagination.Filter).
		Order(fmt.Sprintf("coin.%s %s", fieldSql, orderBySql)).
		SelectAndCount()

	helpers.CheckErr(err)

	return coins
}