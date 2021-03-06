package data_resources

import (
	"github.com/noah-blockchain/coinExplorer-tools/models"
	"github.com/noah-blockchain/noah-explorer-api/internal/helpers"
	"github.com/noah-blockchain/noah-explorer-api/internal/resource"
)

type BuyCoin struct {
	CoinToBuy          string `json:"coin_to_buy"`
	CoinToSell         string `json:"coin_to_sell"`
	ValueToBuy         string `json:"value_to_buy"`
	ValueToSell        string `json:"value_to_sell"`
	MaximumValueToSell string `json:"maximum_value_to_sell"`
}

func (BuyCoin) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*models.BuyCoinTxData)
	model := params[0].(models.Transaction)

	return BuyCoin{
		CoinToBuy:          data.CoinToBuy,
		CoinToSell:         data.CoinToSell,
		ValueToBuy:         helpers.QNoahStr2Noah(data.ValueToBuy),
		ValueToSell:        helpers.QNoahStr2Noah(model.Tags["tx.return"]),
		MaximumValueToSell: helpers.QNoahStr2Noah(data.MaximumValueToSell),
	}
}
