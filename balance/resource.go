package balance

import (
	"github.com/noah-blockchain/CoinExplorer-BackEnd/helpers"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/resource"
	"github.com/noah-blockchain/coinExplorer-tools/models"
)

type Resource struct {
	Coin   string `json:"coin"`
	Amount string `json:"amount"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	balance := model.(models.Balance)

	return Resource{
		Coin:   balance.Coin.Symbol,
		Amount: helpers.QNoahStr2Noah(balance.Value),
	}
}
