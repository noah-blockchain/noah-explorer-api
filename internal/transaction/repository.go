package transaction

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/noah-blockchain/coinExplorer-tools/models"
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

// Get paginated list of transactions by address filter
func (repository Repository) GetPaginatedTxsByAddresses(addresses []string, filter BlocksRangeSelectFilter, pagination *tools.Pagination) []models.Transaction {
	var transactions []models.Transaction
	var err error

	pagination.Total, err = repository.db.Model(&transactions).
		Join("INNER JOIN index_transaction_by_address AS ind").
		JoinOn("ind.transaction_id = transaction.id").
		Join("INNER JOIN addresses AS a").
		JoinOn("a.id = ind.address_id").
		ColumnExpr("DISTINCT transaction.id").
		Column("transaction.*", "FromAddress.address", "GasCoin.symbol").
		Where("a.address IN (?)", pg.In(addresses)).
		Apply(filter.Filter).
		Apply(pagination.Filter).
		Order("transaction.id DESC").
		SelectAndCount()

	helpers.CheckErr(err)

	return transactions
}

// Get paginated list of transactions by select filter
func (repository Repository) GetPaginatedTxsByFilter(filter tools.Filter, pagination *tools.Pagination) []models.Transaction {
	var transactions []models.Transaction
	var err error

	pagination.Total, err = repository.db.Model(&transactions).
		Column("transaction.*", "FromAddress.address", "GasCoin.symbol").
		Apply(filter.Filter).
		Apply(pagination.Filter).
		Order("transaction.id DESC").
		SelectAndCount()

	helpers.CheckErr(err)

	return transactions
}

// Get transaction by hash
func (repository Repository) GetTxByHash(hash string) *models.Transaction {
	var transaction models.Transaction

	err := repository.db.Model(&transaction).Column("FromAddress", "GasCoin.symbol").Where("hash = ?", hash).Select()
	if err != nil {
		return nil
	}

	return &transaction
}

type TxCountChartData struct {
	Time  time.Time
	Count uint64
}

// Get list of transactions counts filtered by created_at
func (repository Repository) GetTxCountChartDataByFilter(filter tools.Filter) []TxCountChartData {
	var tx models.Transaction
	var data []TxCountChartData

	err := repository.db.Model(&tx).
		ColumnExpr("COUNT(*) as count").
		Apply(filter.Filter).
		Select(&data)

	helpers.CheckErr(err)

	return data
}

// Get total transaction count
func (repository Repository) GetTotalTransactionCount(startTime *string) int {
	var tx models.Transaction

	query := repository.db.Model(&tx)
	if startTime != nil {
		query = query.Column("Block._").Where("block.created_at >= ?", *startTime)
	}

	count, err := query.Count()
	helpers.CheckErr(err)

	return count
}

type Tx24hData struct {
	FeeSum float64
	Count  int
	FeeAvg float64
}

// Get transactions data by last 24 hours
func (repository Repository) Get24hTransactionsData() Tx24hData {
	var tx models.Transaction
	var data Tx24hData

	err := repository.db.Model(&tx).
		Column("Block._").
		ColumnExpr("COUNT(*) as count, SUM(gas * gas_price) as fee_sum, AVG(gas * gas_price) as fee_avg").
		Where("block.created_at >= ?", time.Now().AddDate(0, 0, -1).Format(time.RFC3339)).
		Select(&data)

	helpers.CheckErr(err)
	return data
}

// Get paginated list of transactions by coin
func (repository Repository) GetPaginatedTxsByCoin(coinSymbol string, pagination *tools.Pagination) []models.TransactionOutput {
	var transactionOutputs []models.TransactionOutput
	var err error

	pagination.Total, err = repository.db.Model(&transactionOutputs).
		Join("LEFT JOIN coins as c").
		JoinOn("c.id = transaction_output.coin_id").
		Where("c.symbol=?", coinSymbol).
		Column("Transaction", "transaction_output.value",
			"Transaction.FromAddress", "Transaction.GasCoin").
		Apply(pagination.Filter).
		Order("transaction_output.id DESC").
		SelectAndCount()

	helpers.CheckErr(err)
	return transactionOutputs
}
