package moi

import (
	"math"

	"github.com/enbility/eebus-go/api"
	internal "github.com/enbility/eebus-go/usecases/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

func electricity(measurementType model.MeasurementTypeType, scope model.ScopeTypeType) model.MeasurementDescriptionDataType {
	return model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(measurementType),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(scope),
	}
}

// return the momentary total active AC power of the inverter
//
//   - positive values: consumption, negative values: production
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MOI) Power(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	filter := electricity(model.MeasurementTypeTypePower, model.ScopeTypeTypeACPowerTotal)
	return internal.MeasurementValueForFilter(e.LocalEntity, entity, filter, model.EnergyDirectionTypeConsume)
}

// return the momentary active AC power per phase (ordered L1, L2, L3)
//
//   - positive values: consumption, negative values: production
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MOI) PowerPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	filter := electricity(model.MeasurementTypeTypePower, model.ScopeTypeTypeACPower)
	return internal.MeasurementValuesPerPhase(e.LocalEntity, entity, filter, model.EnergyDirectionTypeConsume, true)
}

// return the total AC energy the inverter consumed, in Wh
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MOI) EnergyConsumed(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	filter := electricity(model.MeasurementTypeTypeEnergy, model.ScopeTypeTypeACEnergyConsumed)
	return internal.MeasurementValueForFilter(e.LocalEntity, entity, filter, "")
}

// return the total AC energy the inverter produced, in Wh
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MOI) EnergyProduced(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	filter := electricity(model.MeasurementTypeTypeEnergy, model.ScopeTypeTypeACEnergyProduced)
	return internal.MeasurementValueForFilter(e.LocalEntity, entity, filter, "")
}

// return the total AC yield of the inverter, in Wh (magnitude)
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MOI) YieldTotal(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	filter := electricity(model.MeasurementTypeTypeEnergy, model.ScopeTypeTypeACYieldTotal)
	value, err := internal.MeasurementValueForFilter(e.LocalEntity, entity, filter, "")
	if err != nil {
		return 0, err
	}
	return math.Abs(value), nil
}

// return the momentary AC current per phase (ordered L1, L2, L3)
//
//   - positive values: consumption, negative values: production
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MOI) CurrentPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	filter := electricity(model.MeasurementTypeTypeCurrent, model.ScopeTypeTypeACCurrent)
	return internal.MeasurementValuesPerPhase(e.LocalEntity, entity, filter, model.EnergyDirectionTypeConsume, false)
}

// return the phase-to-neutral AC voltage per phase (ordered L1, L2, L3)
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MOI) VoltagePerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	filter := electricity(model.MeasurementTypeTypeVoltage, model.ScopeTypeTypeACVoltage)
	return internal.MeasurementValuesPerPhase(e.LocalEntity, entity, filter, "", true)
}

// return the AC frequency, in Hz
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MOI) Frequency(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	filter := electricity(model.MeasurementTypeTypeFrequency, model.ScopeTypeTypeACFrequency)
	return internal.MeasurementValueForFilter(e.LocalEntity, entity, filter, "")
}
