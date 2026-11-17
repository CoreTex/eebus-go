package api

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
)

// Actor: Monitoring Appliance
// UseCase: Monitoring of Inverter
//
// Sign convention for power and current: positive values mean the inverter
// consumes from its AC connection (e.g. charging a battery from the grid),
// negative values mean it produces (feeds into the AC side). Values reported
// with the opposite positive energy direction are converted.
type MaMOIInterface interface {
	api.UseCaseInterface

	// return the momentary total active AC power of the inverter
	//
	// parameters:
	//   - entity: the entity of the inverter
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such value is (yet) available
	//   - ErrDataInvalid if the currently available data is invalid and should be ignored
	//   - and others
	Power(entity spineapi.EntityRemoteInterface) (float64, error)

	// return the momentary active AC power per phase (ordered L1, L2, L3)
	//
	// parameters:
	//   - entity: the entity of the inverter
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such value is (yet) available
	//   - ErrDataInvalid if the currently available data is invalid and should be ignored
	//   - and others
	PowerPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error)

	// return the total AC energy the inverter consumed, in Wh
	//
	// parameters:
	//   - entity: the entity of the inverter
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such value is (yet) available
	//   - ErrDataInvalid if the currently available data is invalid and should be ignored
	//   - and others
	EnergyConsumed(entity spineapi.EntityRemoteInterface) (float64, error)

	// return the total AC energy the inverter produced, in Wh
	//
	// parameters:
	//   - entity: the entity of the inverter
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such value is (yet) available
	//   - ErrDataInvalid if the currently available data is invalid and should be ignored
	//   - and others
	EnergyProduced(entity spineapi.EntityRemoteInterface) (float64, error)

	// return the total AC yield of the inverter, in Wh
	//
	// Devices report the yield with the sign of their positive energy
	// direction; the magnitude is returned.
	//
	// parameters:
	//   - entity: the entity of the inverter
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such value is (yet) available
	//   - ErrDataInvalid if the currently available data is invalid and should be ignored
	//   - and others
	YieldTotal(entity spineapi.EntityRemoteInterface) (float64, error)

	// return the momentary AC current per phase (ordered L1, L2, L3)
	//
	// parameters:
	//   - entity: the entity of the inverter
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such value is (yet) available
	//   - ErrDataInvalid if the currently available data is invalid and should be ignored
	//   - and others
	CurrentPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error)

	// return the phase-to-neutral AC voltage per phase (ordered L1, L2, L3)
	//
	// Phase-to-phase voltages, which some devices report with the same
	// scope, are ignored.
	//
	// parameters:
	//   - entity: the entity of the inverter
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such value is (yet) available
	//   - ErrDataInvalid if the currently available data is invalid and should be ignored
	//   - and others
	VoltagePerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error)

	// return the AC frequency, in Hz
	//
	// parameters:
	//   - entity: the entity of the inverter
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such value is (yet) available
	//   - ErrDataInvalid if the currently available data is invalid and should be ignored
	//   - and others
	Frequency(entity spineapi.EntityRemoteInterface) (float64, error)
}
