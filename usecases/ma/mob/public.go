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
	return e.value(entity, model.MeasurementTypeTypeEnergy, model.ScopeTypeTypeCharge)
}

// return the total energy discharged from the battery, in Wh
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MOB) EnergyDischarged(entity spineapi.EntityRemoteInterface) (float64, error) {
	return e.value(entity, model.MeasurementTypeTypeEnergy, model.ScopeTypeTypeDischarge)
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

func (e *MOB) value(
	entity spineapi.EntityRemoteInterface,
	measurementType model.MeasurementTypeType,
	scope model.ScopeTypeType,
) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}
	return internal.MeasurementValueForFilter(e.LocalEntity, entity, filterFor(measurementType, scope), "")
}
