package mapper

import (
	"github.com/4aykovski/effective_mobile_test_task/internal/service/carservice"
	"github.com/4aykovski/effective_mobile_test_task/internal/service/ownerservice"
	"github.com/4aykovski/effective_mobile_test_task/pkg/client/carinfo"
)

func CarInfoIntoCarAndOwner(carInfos map[string]carinfo.CarInfo) ([]carservice.AddNewCarInput, []ownerservice.AddNewOwnerInput) {
	var cars []carservice.AddNewCarInput
	for regNumber, carInfo := range carInfos {
		valid := true
		if carInfo == (carinfo.CarInfo{}) {
			valid = false
		}

		car := carservice.AddNewCarInput{
			RegistrationNumber: regNumber,
			Mark:               carInfo.Mark,
			Model:              carInfo.Model,
			Year:               carInfo.Year,
			OwnerName:          carInfo.Owner.Name,
			OwnerSurname:       carInfo.Owner.Surname,
			Valid:              valid,
		}
		cars = append(cars, car)
	}

	var owners []ownerservice.AddNewOwnerInput
	for _, carInfo := range carInfos {
		owner := ownerservice.AddNewOwnerInput{
			Name:       carInfo.Owner.Name,
			Surname:    carInfo.Owner.Surname,
			Patronymic: carInfo.Owner.Patronymic,
		}
		owners = append(owners, owner)
	}

	return cars, owners
}
