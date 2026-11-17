package internal

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

type mvFixture struct {
	id          uint
	scope       model.ScopeTypeType
	phase       *model.ElectricalConnectionPhaseNameType
	referenceTo *model.ElectricalConnectionPhaseNameType
	value       *float64
	state       *model.MeasurementValueStateType
}

func mvPhase(p model.ElectricalConnectionPhaseNameType) *model.ElectricalConnectionPhaseNameType {
	return &p
}

func mvValue(v float64) *float64 { return &v }

// loadMeasurements writes descriptions, values and electrical connection data
// into the monitored remote entity. withParams=false omits the parameter
// descriptions (electrical connection information missing).
func (s *InternalSuite) loadMeasurements(fixture []mvFixture, direction model.EnergyDirectionType, withParams bool) {
	descs := []model.MeasurementDescriptionDataType{}
	values := []model.MeasurementDataType{}
	params := []model.ElectricalConnectionParameterDescriptionDataType{}
	for i, f := range fixture {
		descs = append(descs, model.MeasurementDescriptionDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(f.id)),
			ScopeType:     util.Ptr(f.scope),
		})
		item := model.MeasurementDataType{MeasurementId: util.Ptr(model.MeasurementIdType(f.id)), ValueState: f.state}
		if f.value != nil {
			item.Value = model.NewScaledNumberType(*f.value)
		}
		values = append(values, item)
		params = append(params, model.ElectricalConnectionParameterDescriptionDataType{
			ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
			ParameterId:             util.Ptr(model.ElectricalConnectionParameterIdType(i)),
			MeasurementId:           util.Ptr(model.MeasurementIdType(f.id)),
			AcMeasuredPhases:        f.phase,
			AcMeasuredInReferenceTo: f.referenceTo,
		})
	}

	rMeas := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, err := rMeas.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData,
		&model.MeasurementDescriptionListDataType{MeasurementDescriptionData: descs}, nil, nil)
	assert.Nil(s.T(), err)
	_, err = rMeas.UpdateData(true, model.FunctionTypeMeasurementListData,
		&model.MeasurementListDataType{MeasurementData: values}, nil, nil)
	assert.Nil(s.T(), err)

	rEl := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	desc := model.ElectricalConnectionDescriptionDataType{ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0))}
	if direction != "" {
		desc.PositiveEnergyDirection = util.Ptr(direction)
	}
	_, err = rEl.UpdateData(true, model.FunctionTypeElectricalConnectionDescriptionListData,
		&model.ElectricalConnectionDescriptionListDataType{
			ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{desc},
		}, nil, nil)
	assert.Nil(s.T(), err)
	if withParams {
		_, err = rEl.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData,
			&model.ElectricalConnectionParameterDescriptionListDataType{ElectricalConnectionParameterDescriptionData: params}, nil, nil)
		assert.Nil(s.T(), err)
	}
}

func (s *InternalSuite) Test_MeasurementValueForFilter() {
	filter := model.MeasurementDescriptionDataType{ScopeType: util.Ptr(model.ScopeTypeTypeACPowerTotal)}

	_, err := MeasurementValueForFilter(nil, nil, filter, model.EnergyDirectionTypeConsume)
	assert.Equal(s.T(), api.ErrMetadataNotAvailable, err)

	_, err = MeasurementValueForFilter(s.localEntity, s.monitoredEntity, filter, model.EnergyDirectionTypeConsume)
	assert.Equal(s.T(), api.ErrDataNotAvailable, err)

	// a description without a value is no data
	s.loadMeasurements([]mvFixture{{id: 0, scope: model.ScopeTypeTypeACPowerTotal}}, model.EnergyDirectionTypeConsume, true)
	_, err = MeasurementValueForFilter(s.localEntity, s.monitoredEntity, filter, model.EnergyDirectionTypeConsume)
	assert.Equal(s.T(), api.ErrDataNotAvailable, err)

	// same direction: value as reported
	s.loadMeasurements([]mvFixture{{id: 0, scope: model.ScopeTypeTypeACPowerTotal, value: mvValue(-93.5)}}, model.EnergyDirectionTypeConsume, true)
	v, err := MeasurementValueForFilter(s.localEntity, s.monitoredEntity, filter, model.EnergyDirectionTypeConsume)
	assert.Nil(s.T(), err)
	assert.InDelta(s.T(), -93.5, v, 1e-9)

	// no direction requested: value as reported
	v, err = MeasurementValueForFilter(s.localEntity, s.monitoredEntity, filter, "")
	assert.Nil(s.T(), err)
	assert.InDelta(s.T(), -93.5, v, 1e-9)

	// opposite direction: negated
	s.loadMeasurements([]mvFixture{{id: 0, scope: model.ScopeTypeTypeACPowerTotal, value: mvValue(-93.5)}}, model.EnergyDirectionTypeProduce, true)
	v, err = MeasurementValueForFilter(s.localEntity, s.monitoredEntity, filter, model.EnergyDirectionTypeConsume)
	assert.Nil(s.T(), err)
	assert.InDelta(s.T(), 93.5, v, 1e-9)

	// direction unknown: value as reported
	s.loadMeasurements([]mvFixture{{id: 0, scope: model.ScopeTypeTypeACPowerTotal, value: mvValue(-93.5)}}, "", true)
	v, err = MeasurementValueForFilter(s.localEntity, s.monitoredEntity, filter, model.EnergyDirectionTypeConsume)
	assert.Nil(s.T(), err)
	assert.InDelta(s.T(), -93.5, v, 1e-9)

	// ambiguous filter
	s.loadMeasurements([]mvFixture{
		{id: 0, scope: model.ScopeTypeTypeACPowerTotal, value: mvValue(1)},
		{id: 1, scope: model.ScopeTypeTypeACPowerTotal, value: mvValue(2)},
	}, model.EnergyDirectionTypeConsume, true)
	_, err = MeasurementValueForFilter(s.localEntity, s.monitoredEntity, filter, model.EnergyDirectionTypeConsume)
	assert.Equal(s.T(), api.ErrDataNotAvailable, err)

	// invalid value state
	s.loadMeasurements([]mvFixture{{
		id: 0, scope: model.ScopeTypeTypeACPowerTotal, value: mvValue(1),
		state: util.Ptr(model.MeasurementValueStateTypeError),
	}}, model.EnergyDirectionTypeConsume, true)
	_, err = MeasurementValueForFilter(s.localEntity, s.monitoredEntity, filter, model.EnergyDirectionTypeConsume)
	assert.Equal(s.T(), api.ErrDataInvalid, err)
}

func (s *InternalSuite) Test_MeasurementValuesPerPhase() {
	filter := model.MeasurementDescriptionDataType{ScopeType: util.Ptr(model.ScopeTypeTypeACVoltage)}
	a := mvPhase(model.ElectricalConnectionPhaseNameTypeA)
	b := mvPhase(model.ElectricalConnectionPhaseNameTypeB)
	c := mvPhase(model.ElectricalConnectionPhaseNameTypeC)
	n := mvPhase(model.ElectricalConnectionPhaseNameTypeNeutral)
	abc := mvPhase(model.ElectricalConnectionPhaseNameTypeAbc)

	data, err := MeasurementValuesPerPhase(nil, nil, filter, "", true)
	assert.Equal(s.T(), api.ErrMetadataNotAvailable, err)
	assert.Nil(s.T(), data)

	data, err = MeasurementValuesPerPhase(s.localEntity, s.monitoredEntity, filter, "", true)
	assert.Equal(s.T(), api.ErrDataNotAvailable, err)
	assert.Nil(s.T(), data)

	fixture := []mvFixture{
		// phase-to-phase first; reported with the same AcMeasuredPhases
		{id: 24, scope: model.ScopeTypeTypeACVoltage, phase: a, referenceTo: b, value: mvValue(398.1)},
		{id: 26, scope: model.ScopeTypeTypeACVoltage, phase: b, referenceTo: c, value: mvValue(399.2)},
		{id: 28, scope: model.ScopeTypeTypeACVoltage, phase: c, referenceTo: a, value: mvValue(400.3)},
		// out of order on purpose
		{id: 29, scope: model.ScopeTypeTypeACVoltage, phase: c, referenceTo: n, value: mvValue(230.6)},
		{id: 25, scope: model.ScopeTypeTypeACVoltage, phase: a, referenceTo: n, value: mvValue(228.7)},
		{id: 27, scope: model.ScopeTypeTypeACVoltage, phase: b, referenceTo: n, value: mvValue(229.6)},
		// not a single phase: ignored
		{id: 30, scope: model.ScopeTypeTypeACVoltage, phase: abc, value: mvValue(1)},
		// no phase information: ignored
		{id: 31, scope: model.ScopeTypeTypeACVoltage, value: mvValue(2)},
		// no value: ignored
		{id: 32, scope: model.ScopeTypeTypeACVoltage, phase: a, referenceTo: n},
	}
	s.loadMeasurements(fixture, model.EnergyDirectionTypeConsume, true)

	data, err = MeasurementValuesPerPhase(s.localEntity, s.monitoredEntity, filter, "", true)
	assert.Nil(s.T(), err)
	assert.InDeltaSlice(s.T(), []float64{228.7, 229.6, 230.6}, data, 1e-9)

	// without the neutral filter the first value per phase wins (phase-to-phase here)
	data, err = MeasurementValuesPerPhase(s.localEntity, s.monitoredEntity, filter, "", false)
	assert.Nil(s.T(), err)
	assert.InDeltaSlice(s.T(), []float64{398.1, 399.2, 400.3}, data, 1e-9)

	// energy direction conversion
	power := []mvFixture{
		{id: 21, scope: model.ScopeTypeTypeACPower, phase: a, referenceTo: n, value: mvValue(-26.8)},
		{id: 22, scope: model.ScopeTypeTypeACPower, phase: b, referenceTo: n, value: mvValue(-32.5)},
	}
	powerFilter := model.MeasurementDescriptionDataType{ScopeType: util.Ptr(model.ScopeTypeTypeACPower)}
	s.loadMeasurements(power, model.EnergyDirectionTypeProduce, true)
	data, err = MeasurementValuesPerPhase(s.localEntity, s.monitoredEntity, powerFilter, model.EnergyDirectionTypeConsume, true)
	assert.Nil(s.T(), err)
	assert.InDeltaSlice(s.T(), []float64{26.8, 32.5}, data, 1e-9)

	// missing parameter descriptions: no phase can be determined
	s.loadMeasurements(power, model.EnergyDirectionTypeConsume, false)
	rEl := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	_, fErr := rEl.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData,
		&model.ElectricalConnectionParameterDescriptionListDataType{}, nil, nil)
	assert.Nil(s.T(), fErr)
	data, err = MeasurementValuesPerPhase(s.localEntity, s.monitoredEntity, powerFilter, "", true)
	assert.Equal(s.T(), api.ErrDataNotAvailable, err)
	assert.Nil(s.T(), data)

	// invalid value state
	s.loadMeasurements([]mvFixture{{
		id: 21, scope: model.ScopeTypeTypeACPower, phase: a, referenceTo: n, value: mvValue(1),
		state: util.Ptr(model.MeasurementValueStateTypeError),
	}}, model.EnergyDirectionTypeConsume, true)
	_, err = MeasurementValuesPerPhase(s.localEntity, s.monitoredEntity, powerFilter, "", true)
	assert.Equal(s.T(), api.ErrDataInvalid, err)
}
