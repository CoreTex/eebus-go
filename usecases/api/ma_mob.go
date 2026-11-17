package api

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
)

// Actor: Monitoring Appliance
// UseCase: Monitoring of Battery
//
// Sign convention for power: positive values mean the battery is charging,
// negative values mean it is discharging (consistent with VABD). Values
// reported with the opposite positive energy direction are converted.
type MaMOBInterface interface {
	api.UseCaseInterface

	// return the momentary total power of the battery
	//
	// The AC total power is used if the battery reports it, otherwise its DC
	// power.
	//
	// parameters:
	//   - entity: the entity of the battery
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such value is (yet) available
	//   - ErrDataInvalid if the currently available data is invalid and should be ignored
	//   - and others
	Power(entity spineapi.EntityRemoteInterface) (float64, error)

	// return the total energy charged into the battery, in Wh
	//
	// parameters:
	//   - entity: the entity of the battery
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such value is (yet) available
	//   - ErrDataInvalid if the currently available data is invalid and should be ignored
	//   - and others
	EnergyCharged(entity spineapi.EntityRemoteInterface) (float64, error)

	// return the total energy discharged from the battery, in Wh
	//
	// parameters:
	//   - entity: the entity of the battery
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such value is (yet) available
	//   - ErrDataInvalid if the currently available data is invalid and should be ignored
	//   - and others
	EnergyDischarged(entity spineapi.EntityRemoteInterface) (float64, error)

	// return the state of charge, in percent
	//
	// parameters:
	//   - entity: the entity of the battery
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such value is (yet) available
	//   - ErrDataInvalid if the currently available data is invalid and should be ignored
	//   - and others
	StateOfCharge(entity spineapi.EntityRemoteInterface) (float64, error)

	// return the state of health, in percent
	//
	// parameters:
	//   - entity: the entity of the battery
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such value is (yet) available
	//   - ErrDataInvalid if the currently available data is invalid and should be ignored
	//   - and others
	StateOfHealth(entity spineapi.EntityRemoteInterface) (float64, error)

	// return the battery (component) temperature, in °C
	//
	// parameters:
	//   - entity: the entity of the battery
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such value is (yet) available
	//   - ErrDataInvalid if the currently available data is invalid and should be ignored
	//   - and others
	Temperature(entity spineapi.EntityRemoteInterface) (float64, error)
}
