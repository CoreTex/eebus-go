package internal

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// MeasurementValueForFilter returns the value of the single measurement that
// matches the filter and carries a value.
//
// If energyDirection is set and the remote electrical connection of that
// measurement reports the opposite positive energy direction, the value is
// negated, so the result always follows the requested sign convention. If the
// remote provides no electrical connection information, the value is returned
// as reported.
//
// possible errors:
//   - ErrMetadataNotAvailable if the remote entity has no measurement feature
//   - ErrDataNotAvailable if there is no, or more than one, matching value
//   - ErrDataInvalid if the value state is not normal
func MeasurementValueForFilter(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface,
	filter model.MeasurementDescriptionDataType,
	energyDirection model.EnergyDirectionType,
) (float64, error) {
	measurement, err := client.NewMeasurement(localEntity, remoteEntity)
	if err != nil {
		return 0, api.ErrMetadataNotAvailable
	}

	data, err := measurement.GetDataForFilter(filter)
	if err != nil {
		return 0, api.ErrDataNotAvailable
	}

	var found *model.MeasurementDataType
	for i := range data {
		if data[i].Value == nil {
			continue
		}
		if found != nil {
			// ambiguous: the caller's filter is not specific enough
			return 0, api.ErrDataNotAvailable
		}
		found = &data[i]
	}
	if found == nil {
		return 0, api.ErrDataNotAvailable
	}

	if found.ValueState != nil && *found.ValueState != model.MeasurementValueStateTypeNormal {
		return 0, api.ErrDataInvalid
	}

	value := found.Value.GetValue()
	if energyDirection != "" && found.MeasurementId != nil &&
		isOppositeEnergyDirection(localEntity, remoteEntity, *found.MeasurementId, energyDirection) {
		value = -value
	}

	return value, nil
}

// MeasurementValuesPerPhase returns the phase specific values matching the
// filter, ordered by phase (a, b, c). Phases without a value are skipped.
//
// If toNeutral is true, only values measured against the neutral conductor
// (or without a reference) are used. This separates phase-to-neutral voltages
// from phase-to-phase voltages, which devices may report with the same scope
// and AcMeasuredPhases value.
//
// energyDirection works as in MeasurementValueForFilter.
//
// possible errors:
//   - ErrMetadataNotAvailable if the remote entity has no measurement or
//     electrical connection feature
//   - ErrDataNotAvailable if no phase specific value is available
//   - ErrDataInvalid if a value state is not normal
func MeasurementValuesPerPhase(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface,
	filter model.MeasurementDescriptionDataType,
	energyDirection model.EnergyDirectionType,
	toNeutral bool,
) ([]float64, error) {
	measurement, err := client.NewMeasurement(localEntity, remoteEntity)
	electricalConnection, err1 := client.NewElectricalConnection(localEntity, remoteEntity)
	if err != nil || err1 != nil {
		return nil, api.ErrMetadataNotAvailable
	}

	data, err := measurement.GetDataForFilter(filter)
	if err != nil || len(data) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	phases := []model.ElectricalConnectionPhaseNameType{
		model.ElectricalConnectionPhaseNameTypeA,
		model.ElectricalConnectionPhaseNameTypeB,
		model.ElectricalConnectionPhaseNameTypeC,
	}
	byPhase := make(map[model.ElectricalConnectionPhaseNameType]float64, len(phases))

	for _, item := range data {
		if item.Value == nil || item.MeasurementId == nil {
			continue
		}

		params, err := electricalConnection.GetParameterDescriptionsForFilter(
			model.ElectricalConnectionParameterDescriptionDataType{MeasurementId: item.MeasurementId})
		if err != nil || len(params) == 0 || params[0].AcMeasuredPhases == nil {
			continue
		}
		param := params[0]

		phase := *param.AcMeasuredPhases
		if _, ok := byPhase[phase]; ok {
			continue // keep the first value per phase
		}
		isPhase := false
		for _, p := range phases {
			if p == phase {
				isPhase = true
				break
			}
		}
		if !isPhase {
			continue
		}
		if toNeutral && param.AcMeasuredInReferenceTo != nil &&
			*param.AcMeasuredInReferenceTo != model.ElectricalConnectionPhaseNameTypeNeutral {
			continue
		}

		if item.ValueState != nil && *item.ValueState != model.MeasurementValueStateTypeNormal {
			return nil, api.ErrDataInvalid
		}

		value := item.Value.GetValue()
		if energyDirection != "" &&
			isOppositeEnergyDirection(localEntity, remoteEntity, *item.MeasurementId, energyDirection) {
			value = -value
		}
		byPhase[phase] = value
	}

	result := make([]float64, 0, len(phases))
	for _, p := range phases {
		if v, ok := byPhase[p]; ok {
			result = append(result, v)
		}
	}
	if len(result) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	return result, nil
}

// isOppositeEnergyDirection reports whether the electrical connection of the
// measurement declares the opposite positive energy direction.
func isOppositeEnergyDirection(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface,
	measurementId model.MeasurementIdType,
	energyDirection model.EnergyDirectionType,
) bool {
	electricalConnection, err := client.NewElectricalConnection(localEntity, remoteEntity)
	if err != nil {
		return false
	}
	desc, err := electricalConnection.GetDescriptionForParameterDescriptionFilter(
		model.ElectricalConnectionParameterDescriptionDataType{MeasurementId: &measurementId})
	if err != nil || desc == nil || desc.PositiveEnergyDirection == nil {
		return false
	}
	return *desc.PositiveEnergyDirection != energyDirection
}
