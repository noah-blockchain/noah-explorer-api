package addresses

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/noah-blockchain/coinExplorer-tools/models"
	"github.com/noah-blockchain/noah-explorer-api/internal/address"
	"github.com/noah-blockchain/noah-explorer-api/internal/aggregated_reward"
	"github.com/noah-blockchain/noah-explorer-api/internal/chart"
	"github.com/noah-blockchain/noah-explorer-api/internal/core"
	"github.com/noah-blockchain/noah-explorer-api/internal/delegation"
	"github.com/noah-blockchain/noah-explorer-api/internal/errors"
	"github.com/noah-blockchain/noah-explorer-api/internal/events"
	"github.com/noah-blockchain/noah-explorer-api/internal/helpers"
	"github.com/noah-blockchain/noah-explorer-api/internal/resource"
	"github.com/noah-blockchain/noah-explorer-api/internal/reward"
	"github.com/noah-blockchain/noah-explorer-api/internal/slash"
	"github.com/noah-blockchain/noah-explorer-api/internal/tools"
	"github.com/noah-blockchain/noah-explorer-api/internal/transaction"
	validatorMeta "github.com/noah-blockchain/noah-explorer-api/internal/validator/meta"
)

const (
	precision = 100
)

type GetAddressRequest struct {
	Address string `uri:"address" binding:"noahAddress"`
}

type GetAddressesRequest struct {
	Addresses []string `form:"addresses[]" binding:"required,noahAddress,max=50"`
}

// TODO: replace string to int
type FilterQueryRequest struct {
	StartBlock *string `form:"startblock" binding:"omitempty,numeric"`
	EndBlock   *string `form:"endblock"   binding:"omitempty,numeric"`
	Page       *string `form:"page"       binding:"omitempty,numeric"`
}

type StatisticsQueryRequest struct {
	StartTime *string `form:"startTime" binding:"omitempty,timestamp"`
	EndTime   *string `form:"endTime"   binding:"omitempty,timestamp"`
}

type AggregatedRewardsQueryRequest struct {
	StartTime *string `form:"startTime" binding:"omitempty,timestamp"`
	EndTime   *string `form:"endTime"   binding:"omitempty,timestamp"`
}

//Get list of addresses ranked by balance
func GetTopAddresses(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	//fetch address
	pagination := tools.NewPagination(c.Request)
	addresses := explorer.AddressRepository.GetPaginatedAddresses(&pagination)

	c.JSON(http.StatusOK, gin.H{
		"data": address.ResourceTopAddresses{}.TransformCollection(addresses, pagination),
	})

}

// Get list of addresses
func GetAddresses(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var request GetAddressesRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// remove Noah wallet prefix from each address
	noahAddresses := make([]string, len(request.Addresses))
	for key, addr := range request.Addresses {
		noahAddresses[key] = helpers.RemoveNoahPrefix(addr)
	}

	// fetch addresses
	addresses := explorer.AddressRepository.GetByAddresses(noahAddresses)

	// extend the model array with empty model if not exists
	if len(addresses) != len(noahAddresses) {
		for _, item := range noahAddresses {
			if isModelsContainAddress(item, addresses) {
				continue
			}

			addresses = append(addresses, *makeEmptyAddressModel(item, explorer.Environment.BaseCoin))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resource.TransformCollection(addresses, address.Resource{}),
	})
}

// Get address detail
func GetAddress(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	noahAddress, err := getAddressFromRequestUri(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch address
	model := explorer.AddressRepository.GetByAddress(*noahAddress)

	// if model not found
	if model == nil {
		model = makeEmptyAddressModel(*noahAddress, explorer.Environment.BaseCoin)
	}

	c.JSON(http.StatusOK, gin.H{"data": new(address.Resource).Transform(*model)})
}

// Get list of transactions by noah address
func GetTransactions(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	noahAddress, err := getAddressFromRequestUri(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// validate request query
	var requestQuery FilterQueryRequest
	err = c.ShouldBindQuery(&requestQuery)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	pagination := tools.NewPagination(c.Request)
	txs := explorer.TransactionRepository.GetPaginatedTxsByAddresses(
		[]string{*noahAddress},
		transaction.BlocksRangeSelectFilter{
			StartBlock: requestQuery.StartBlock,
			EndBlock:   requestQuery.EndBlock,
		}, &pagination)

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(txs, transaction.Resource{}, pagination))
}

// Get list of rewards by Noah address
func GetRewards(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	filter, pagination, err := prepareEventsRequest(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	rewards := explorer.RewardRepository.GetPaginatedByAddress(*filter, pagination)

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(rewards, reward.Resource{}, *pagination))
}

func GetAggregatedRewards(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	noahAddress, err := getAddressFromRequestUri(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	var requestQuery FilterQueryRequest
	if err := c.ShouldBindQuery(&requestQuery); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	pagination := tools.NewPagination(c.Request)
	rewards := explorer.RewardRepository.GetPaginatedAggregatedByAddress(aggregated_reward.SelectFilter{
		Address:   *noahAddress,
		StartTime: requestQuery.StartBlock,
		EndTime:   requestQuery.EndBlock,
	}, &pagination)

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(rewards, aggregated_reward.Resource{}, pagination))
}

// Get list of slashes by Noah address
func GetSlashes(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	filter, pagination, err := prepareEventsRequest(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	slashes := explorer.SlashRepository.GetPaginatedByAddress(*filter, pagination)

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(slashes, slash.Resource{}, *pagination))
}

// Get list of delegations by Noah address
func GetDelegations(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	noahAddress, err := getAddressFromRequestUri(c)
	if err != nil || noahAddress == nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	pagination := tools.NewPagination(c.Request)

	stakesSum, err := explorer.StakeRepository.GetSumInNoahValueByAddress(*noahAddress)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	stakes := explorer.StakeRepository.GetPaginatedByAddress(*noahAddress, &pagination)
	delegatedStakeList := make([]delegation.Resource, len(stakes))
	for i, stake := range stakes {

		yourMoney := helpers.NewFloat(0, precision)
		//if *stake.Validator.Commission < 100 {
		//	fmt.Println(*stake.Validator.Commission)
		//	//// get sum reward validator from time created >= stake.created_at
		//	sumReward := explorer.RewardRepository.GetSumRewardForValidator(stake.ValidatorID, stake.CreatedAt)
		//	sumRewardBigFloat, _ := helpers.NewFloat(0, precision).SetString(sumReward)
		//	log.Println("sumReward", sumReward)
		//	log.Println("sumRewardFloat", sumRewardBigFloat.String())
		//
		//	//// ((sum_reward-(sum_reward * commission_validator_%)) * stake_%) = profit
		//	validatorsMoney := helpers.NewFloat(0, precision)
		//	validatorsMoney = validatorsMoney.Mul(sumRewardBigFloat, big.NewFloat(float64(*stake.Validator.Commission)/100))
		//	log.Println("commission", big.NewFloat(float64(*stake.Validator.Commission)/100).String())
		//	log.Println("validatorsMoney", validatorsMoney.String())
		//
		//	delegationsMoney := validatorsMoney.Sub(sumRewardBigFloat, validatorsMoney)
		//	log.Println("delegationsMoney", delegationsMoney.String())
		//
		//	//// get your stake_% from delegations total stake
		//	percentYourStake, _ := helpers.NewFloat(0, precision).SetString(stake.NoahValue)
		//	percentYourStake = percentYourStake.Mul(percentYourStake, big.NewFloat(100))
		//	percentYourStake = percentYourStake.Quo(percentYourStake, delegationsMoney)
		//	log.Println("percentYourStake", percentYourStake.String())
		//
		//	yourMoney = delegationsMoney.Mul(delegationsMoney, percentYourStake)
		//	yourMoney = yourMoney.Quo(yourMoney, big.NewFloat(100))
		//}

		delegatedStakeList[i] = delegation.Resource{
			Coin:           stake.Coin.Symbol,
			PubKey:         stake.Validator.GetPublicKey(),
			Value:          helpers.QNoahStr2Noah(stake.Value),
			NoahValue:      helpers.QNoahStr2Noah(stake.NoahValue),
			ProfitReceived: helpers.QNoahStr2Noah(yourMoney.String()),
			ValidatorMeta:  new(validatorMeta.Resource).Transform(*stake.Validator),
		}
	}

	additionalFields := map[string]interface{}{
		"total_delegated_noah_value": helpers.QNoahStr2Noah(stakesSum),
	}

	c.JSON(http.StatusOK, resource.TransformPaginatedCollectionWithAdditionalFields(
		delegatedStakeList,
		delegation.Resource{},
		pagination,
		additionalFields,
	))
}

// Get rewards statistics by noah address
func GetRewardsStatistics(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	noahAddress, err := getAddressFromRequestUri(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	var requestQuery StatisticsQueryRequest
	if err := c.ShouldBindQuery(&requestQuery); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	chartData := explorer.RewardRepository.GetAggregatedChartData(aggregated_reward.SelectFilter{
		Address:   *noahAddress,
		EndTime:   requestQuery.EndTime,
		StartTime: requestQuery.StartTime,
	})

	c.JSON(http.StatusOK, gin.H{
		"data": resource.TransformCollection(chartData, chart.RewardResource{}),
	})
}

func prepareEventsRequest(c *gin.Context) (*events.SelectFilter, *tools.Pagination, error) {
	noahAddress, err := getAddressFromRequestUri(c)
	if err != nil {
		return nil, nil, err
	}

	var requestQuery FilterQueryRequest
	if err := c.ShouldBindQuery(&requestQuery); err != nil {
		return nil, nil, err
	}

	pagination := tools.NewPagination(c.Request)

	return &events.SelectFilter{
		Address:    *noahAddress,
		StartBlock: requestQuery.StartBlock,
		EndBlock:   requestQuery.EndBlock,
	}, &pagination, nil
}

// Get noah address from current request uri
func getAddressFromRequestUri(c *gin.Context) (*string, error) {
	var request GetAddressRequest
	if err := c.ShouldBindUri(&request); err != nil {
		return nil, err
	}

	noahAddress := helpers.RemoveNoahPrefix(request.Address)
	return &noahAddress, nil
}

// Return model address with zero base coin
func makeEmptyAddressModel(noahAddress string, baseCoin string) *models.Address {
	return &models.Address{
		Address: noahAddress,
		Balances: []*models.Balance{{
			Coin: &models.Coin{
				Symbol: baseCoin,
			},
			Value: "0",
		}},
	}
}

// Check that array of address models contain exact noah address
func isModelsContainAddress(noahAddress string, models []models.Address) bool {
	for _, item := range models {
		if item.Address == noahAddress {
			return true
		}
	}

	return false
}
