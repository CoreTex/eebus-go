package mob

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntitiesScenarios` to get the current data
	UseCaseSupportUpdate api.EventType = "ma-mob-UseCaseSupportUpdate"

	// Battery power updated
	//
	// Use `Power` to get the current data
	DataUpdatePower api.EventType = "ma-mob-DataUpdatePower"

	// Battery charged energy updated
	//
	// Use `EnergyCharged` to get the current data
	DataUpdateEnergyCharged api.EventType = "ma-mob-DataUpdateEnergyCharged"

	// Battery discharged energy updated
	//
	// Use `EnergyDischarged` to get the current data
	DataUpdateEnergyDischarged api.EventType = "ma-mob-DataUpdateEnergyDischarged"

	// Battery state of charge updated
	//
	// Use `StateOfCharge` to get the current data
	DataUpdateStateOfCharge api.EventType = "ma-mob-DataUpdateStateOfCharge"

	// Battery state of health updated
	//
	// Use `StateOfHealth` to get the current data
	DataUpdateStateOfHealth api.EventType = "ma-mob-DataUpdateStateOfHealth"

	// Battery temperature updated
	//
	// Use `Temperature` to get the current data
	DataUpdateTemperature api.EventType = "ma-mob-DataUpdateTemperature"
)
