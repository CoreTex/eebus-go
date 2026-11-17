package moi

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

// measurement ids as announced by a KOSTAL PLENTICORE G3 inverter entity
const (
	idPowerTotal     = 14
	idPowerL1        = 21
	idPowerL2        = 22
	idPowerL3        = 23
	idYieldTotal     = 9
	idEnergyConsumed = 18
	idEnergyProduced = 19
	idCurrentL1      = 15
	idCurrentL2      = 16
	idCurrentL3      = 17
	idVoltageL1N     = 25
	idVoltageL2N     = 27
	idVoltageL3N     = 29
	idVoltageL1L2    = 24
	idVoltageL2L3    = 26
	idVoltageL3L1    = 28
	idFrequency      = 20
)

type measurementFixture struct {
	id          uint
	mType       model.MeasurementTypeType
	scope       model.ScopeTypeType
	phase       *model.ElectricalConnectionPhaseNameType
	referenceTo *model.ElectricalConnectionPhaseNameType
	value       *float64
}

func phase(p model.ElectricalConnectionPhaseNameType) *model.ElectricalConnectionPhaseNameType {
	return &p
}

// kostalFixture mirrors the descriptions and values a KOSTAL PLENTICORE G3
// reported (taken from a real SPINE trace, values scaled). Notable details:
//   - energies/yields are reported negative (positive energy direction is consume)
//   - phase-to-neutral and phase-to-phase voltages share the acVoltage scope and
//     use the same AcMeasuredPhases values, distinguished only by
//     AcMeasuredInReferenceTo
func kostalFixture() []measurementFixture {
	n := phase(model.ElectricalConnectionPhaseNameTypeNeutral)
	a := phase(model.ElectricalConnectionPhaseNameTypeA)
	b := phase(model.ElectricalConnectionPhaseNameTypeB)
	c := phase(model.ElectricalConnectionPhaseNameTypeC)
	v := func(f float64) *float64 { return &f }
	return []measurementFixture{
		{idPowerTotal, model.MeasurementTypeTypePower, model.ScopeTypeTypeACPowerTotal, nil, nil, v(-93.5)},
		{idPowerL1, model.MeasurementTypeTypePower, model.ScopeTypeTypeACPower, a, n, v(-26.8)},
		{idPowerL2, model.MeasurementTypeTypePower, model.ScopeTypeTypeACPower, b, n, v(-32.5)},
		{idPowerL3, model.MeasurementTypeTypePower, model.ScopeTypeTypeACPower, c, n, v(-34.3)},
		{idYieldTotal, model.MeasurementTypeTypeEnergy, model.ScopeTypeTypeACYieldTotal, nil, nil, v(-27990163.9)},
		{idEnergyConsumed, model.MeasurementTypeTypeEnergy, model.ScopeTypeTypeACEnergyConsumed, nil, nil, v(152253.3)},
		{idEnergyProduced, model.MeasurementTypeTypeEnergy, model.ScopeTypeTypeACEnergyProduced, nil, nil, v(307869.8)},
		{idCurrentL1, model.MeasurementTypeTypeCurrent, model.ScopeTypeTypeACCurrent, a, nil, v(-0.4)},
		{idCurrentL2, model.MeasurementTypeTypeCurrent, model.ScopeTypeTypeACCurrent, b, nil, v(-0.5)},
		{idCurrentL3, model.MeasurementTypeTypeCurrent, model.ScopeTypeTypeACCurrent, c, nil, v(-0.6)},
		// phase-to-phase first, to prove the order of the data does not matter
		{idVoltageL1L2, model.MeasurementTypeTypeVoltage, model.ScopeTypeTypeACVoltage, a, b, v(398.1)},
		{idVoltageL2L3, model.MeasurementTypeTypeVoltage, model.ScopeTypeTypeACVoltage, b, c, v(399.2)},
		{idVoltageL3L1, model.MeasurementTypeTypeVoltage, model.ScopeTypeTypeACVoltage, c, a, v(400.3)},
		{idVoltageL1N, model.MeasurementTypeTypeVoltage, model.ScopeTypeTypeACVoltage, a, n, v(228.7)},
		{idVoltageL2N, model.MeasurementTypeTypeVoltage, model.ScopeTypeTypeACVoltage, b, n, v(229.6)},
		{idVoltageL3N, model.MeasurementTypeTypeVoltage, model.ScopeTypeTypeACVoltage, c, n, v(230.6)},
		{idFrequency, model.MeasurementTypeTypeFrequency, model.ScopeTypeTypeACFrequency, nil, nil, v(50.0)},
	}
}

// loadFixture writes the fixture into the remote features. direction is the
// positive energy direction of the electrical connection ("" = none).
func (s *MaMOISuite) loadFixture(fixture []measurementFixture, direction model.EnergyDirectionType) {
	descs := []model.MeasurementDescriptionDataType{}
	values := []model.MeasurementDataType{}
	params := []model.ElectricalConnectionParameterDescriptionDataType{}
	for i, f := range fixture {
		descs = append(descs, model.MeasurementDescriptionDataType{
			MeasurementId:   util.Ptr(model.MeasurementIdType(f.id)),
			MeasurementType: util.Ptr(f.mType),
			CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
			ScopeType:       util.Ptr(f.scope),
		})
		item := model.MeasurementDataType{MeasurementId: util.Ptr(model.MeasurementIdType(f.id))}
		if f.value != nil {
			item.Value = model.NewScaledNumberType(*f.value)
			item.ValueState = util.Ptr(model.MeasurementValueStateTypeNormal)
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
	desc := model.ElectricalConnectionDescriptionDataType{
		ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
		PowerSupplyType:        util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
	}
	if direction != "" {
		desc.PositiveEnergyDirection = util.Ptr(direction)
	}
	_, err = rEl.UpdateData(true, model.FunctionTypeElectricalConnectionDescriptionListData,
		&model.ElectricalConnectionDescriptionListDataType{
			ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{desc},
		}, nil, nil)
	assert.Nil(s.T(), err)
	_, err = rEl.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData,
		&model.ElectricalConnectionParameterDescriptionListDataType{ElectricalConnectionParameterDescriptionData: params}, nil, nil)
	assert.Nil(s.T(), err)
}

func (s *MaMOISuite) Test_IncompatibleEntity() {
	_, err := s.sut.Power(s.mockRemoteEntity)
	assert.Equal(s.T(), api.ErrNoCompatibleEntity, err)
	_, err = s.sut.PowerPerPhase(s.mockRemoteEntity)
	assert.Equal(s.T(), api.ErrNoCompatibleEntity, err)
	_, err = s.sut.EnergyConsumed(s.mockRemoteEntity)
	assert.Equal(s.T(), api.ErrNoCompatibleEntity, err)
	_, err = s.sut.EnergyProduced(s.mockRemoteEntity)
	assert.Equal(s.T(), api.ErrNoCompatibleEntity, err)
	_, err = s.sut.YieldTotal(s.mockRemoteEntity)
	assert.Equal(s.T(), api.ErrNoCompatibleEntity, err)
	_, err = s.sut.CurrentPerPhase(s.mockRemoteEntity)
	assert.Equal(s.T(), api.ErrNoCompatibleEntity, err)
	_, err = s.sut.VoltagePerPhase(s.mockRemoteEntity)
	assert.Equal(s.T(), api.ErrNoCompatibleEntity, err)
	_, err = s.sut.Frequency(s.mockRemoteEntity)
	assert.Equal(s.T(), api.ErrNoCompatibleEntity, err)
}

func (s *MaMOISuite) Test_NoData() {
	_, err := s.sut.Power(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	_, err = s.sut.YieldTotal(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	data, err := s.sut.VoltagePerPhase(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
}

func (s *MaMOISuite) Test_KostalValues() {
	s.loadFixture(kostalFixture(), model.EnergyDirectionTypeConsume)

	power, err := s.sut.Power(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.InDelta(s.T(), -93.5, power, 1e-9)

	perPhase, err := s.sut.PowerPerPhase(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.InDeltaSlice(s.T(), []float64{-26.8, -32.5, -34.3}, perPhase, 1e-9)

	consumed, err := s.sut.EnergyConsumed(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.InDelta(s.T(), 152253.3, consumed, 1e-6)

	produced, err := s.sut.EnergyProduced(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.InDelta(s.T(), 307869.8, produced, 1e-6)

	// reported negative, returned as magnitude
	yield, err := s.sut.YieldTotal(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.InDelta(s.T(), 27990163.9, yield, 1e-3)

	currents, err := s.sut.CurrentPerPhase(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.InDeltaSlice(s.T(), []float64{-0.4, -0.5, -0.6}, currents, 1e-9)

	// only phase-to-neutral voltages, ordered L1..L3
	voltages, err := s.sut.VoltagePerPhase(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.InDeltaSlice(s.T(), []float64{228.7, 229.6, 230.6}, voltages, 1e-9)

	freq, err := s.sut.Frequency(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.InDelta(s.T(), 50.0, freq, 1e-9)
}

func (s *MaMOISuite) Test_ProduceDirectionIsConverted() {
	s.loadFixture(kostalFixture(), model.EnergyDirectionTypeProduce)

	power, err := s.sut.Power(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.InDelta(s.T(), 93.5, power, 1e-9)

	perPhase, err := s.sut.PowerPerPhase(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.InDeltaSlice(s.T(), []float64{26.8, 32.5, 34.3}, perPhase, 1e-9)

	// voltages carry no energy direction
	voltages, err := s.sut.VoltagePerPhase(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.InDeltaSlice(s.T(), []float64{228.7, 229.6, 230.6}, voltages, 1e-9)
}

func (s *MaMOISuite) Test_InvalidValueState() {
	s.loadFixture(kostalFixture(), model.EnergyDirectionTypeConsume)

	rMeas := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rMeas.UpdateData(true, model.FunctionTypeMeasurementListData, &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(idPowerTotal)),
				Value:         model.NewScaledNumberType(1),
				ValueState:    util.Ptr(model.MeasurementValueStateTypeError),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(idVoltageL1N)),
				Value:         model.NewScaledNumberType(1),
				ValueState:    util.Ptr(model.MeasurementValueStateTypeError),
			},
		},
	}, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err := s.sut.Power(s.monitoredEntity)
	assert.Equal(s.T(), api.ErrDataInvalid, err)
	_, err = s.sut.VoltagePerPhase(s.monitoredEntity)
	assert.Equal(s.T(), api.ErrDataInvalid, err)
}
