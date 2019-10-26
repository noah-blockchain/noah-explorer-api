package aggregated_reward

import (
	"time"

	"github.com/noah-blockchain/CoinExplorer-BackEnd/helpers"
	"github.com/noah-blockchain/CoinExplorer-BackEnd/resource"
	validatorMeta "github.com/noah-blockchain/CoinExplorer-BackEnd/validator/meta"
	"github.com/noah-blockchain/noah-explorer-tools/models"
)

type Resource struct {
	TimeID        string             `json:"time_id"`
	Role          string             `json:"role"`
	Amount        string             `json:"amount"`
	Address       string             `json:"address"`
	Validator     string             `json:"validator"`
	ValidatorMeta resource.Interface `json:"validator_meta"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	reward := model.(models.AggregatedReward)

	return Resource{
		TimeID:        reward.TimeID.Format(time.RFC3339),
		Role:          reward.Role,
		Amount:        helpers.QNoahStr2Noah(reward.Amount),
		Address:       reward.Address.GetAddress(),
		Validator:     reward.Validator.GetPublicKey(),
		ValidatorMeta: new(validatorMeta.Resource).Transform(*reward.Validator),
	}
}
