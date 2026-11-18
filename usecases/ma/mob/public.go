package mob

import (
	"errors"

	"github.com/enbility/eebus-go/api"
	internal "github.com/enbility/eebus-go/usecases/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

func filterFor(measurementType model.MeasurementTypeType, scope model.ScopeTypeType) model.MeasurementDescriptionDataType {
	return model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(measurementType),
		ScopeType:       util.Ptr(scope),
	}
}

// return the momentary total power of the battery
//
//   - positive values: charging, negative values: discharging
//   - uses the AC total power, falls back to the DC power
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MOB) Power(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	value, err := internal.MeasurementValueForFilter(e.LocalEntity, entity,
		filterFor(model.MeasurementTypeTypePower, model.ScopeTypeTypeACPowerTotal), model.EnergyDirectionTypeConsume)
	if err == nil || errors.Is(err, api.ErrDataInvalid) {
		return value, err
	}
	return internal.MeasurementValueForFilter(e.LocalEntity, entity,
		filterFor(model.MeasurementTypeTypePower, model.ScopeTypeTypeDCPower), model.EnergyDirectionTypeConsume)
}

// return the total energy charged into the battery, in Wh
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MOB) EnergyCharged(entity spineapi.EntityRemoteInterface) (float64, error) {
	return e.value(entity, model.MeasurementTypeTypeEnergy,
		model.ScopeTypeTypeCharge, model.ScopeTypeTypeDCChargeEnergy)
}

// return the total energy discharged from the battery, in Wh
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MOB) EnergyDischarged(entity spineapi.EntityRemoteInterface) (float64, error) {
	return e.value(entity, model.MeasurementTypeTypeEnergy,
		model.ScopeTypeTypeDischarge, model.ScopeTypeTypeDCDischargeEnergy)
}

// return the state of charge, in percent
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MOB) StateOfCharge(entity spineapi.EntityRemoteInterface) (float64, error) {
	return e.value(entity, model.MeasurementTypeTypePercentage, model.ScopeTypeTypeStateOfCharge)
}

// return the state of health, in percent
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MOB) StateOfHealth(entity spineapi.EntityRemoteInterface) (float64, error) {
	return e.value(entity, model.MeasurementTypeTypePercentage, model.ScopeTypeTypeStateOfHealth)
}

// return the battery (component) temperature, in °C
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MOB) Temperature(entity spineapi.EntityRemoteInterface) (float64, error) {
	return e.value(entity, model.MeasurementTypeTypeTemperature, model.ScopeTypeTypeComponentTemperature)
}

// return the energy currently stored in the battery, in Wh
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MOB) StateOfEnergy(entity spineapi.EntityRemoteInterface) (float64, error) {
	return e.value(entity, model.MeasurementTypeTypeEnergy, model.ScopeTypeTypeStateOfEnergy)
}

// return the usable capacity of the battery, in Wh
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MOB) UseableCapacity(entity spineapi.EntityRemoteInterface) (float64, error) {
	return e.value(entity, model.MeasurementTypeTypeCapacity, model.ScopeTypeTypeUseableCapacity)
}

// return the number of load cycles of the battery
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MOB) LoadCycleCount(entity spineapi.EntityRemoteInterface) (float64, error) {
	return e.value(entity, model.MeasurementTypeTypeCount, model.ScopeTypeTypeLoadCycleCount)
}

// return the DC voltage of the battery, in V
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MOB) Voltage(entity spineapi.EntityRemoteInterface) (float64, error) {
	return e.value(entity, model.MeasurementTypeTypeVoltage, model.ScopeTypeTypeDCVoltage)
}

// return the DC current of the battery, in A
//
//   - positive values: charging, negative values: discharging
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MOB) Current(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}
	return internal.MeasurementValueForFilter(e.LocalEntity, entity,
		filterFor(model.MeasurementTypeTypeCurrent, model.ScopeTypeTypeDCCurrent), model.EnergyDirectionTypeConsume)
}

// value returns the first available value of the given scopes (in order).
func (e *MOB) value(
	entity spineapi.EntityRemoteInterface,
	measurementType model.MeasurementTypeType,
	scopes ...model.ScopeTypeType,
) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}
	err := api.ErrDataNotAvailable
	for _, scope := range scopes {
		var value float64
		value, err = internal.MeasurementValueForFilter(e.LocalEntity, entity, filterFor(measurementType, scope), "")
		if err == nil || !errors.Is(err, api.ErrDataNotAvailable) {
			return value, err
		}
	}
	return 0, err
}
