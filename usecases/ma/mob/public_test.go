package mob

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

type measurementFixture struct {
	id    uint
	mType model.MeasurementTypeType
	scope model.ScopeTypeType
	value float64
}

func batteryFixture(powerScope model.ScopeTypeType) []measurementFixture {
	return []measurementFixture{
		{0, model.MeasurementTypeTypePower, powerScope, 1500},
		{1, model.MeasurementTypeTypeEnergy, model.ScopeTypeTypeCharge, 1000},
		{2, model.MeasurementTypeTypeEnergy, model.ScopeTypeTypeDischarge, 800},
		{3, model.MeasurementTypeTypePercentage, model.ScopeTypeTypeStateOfCharge, 63.5},
		{4, model.MeasurementTypeTypePercentage, model.ScopeTypeTypeStateOfHealth, 98},
		{5, model.MeasurementTypeTypeTemperature, model.ScopeTypeTypeComponentTemperature, 21.5},
	}
}

func (s *MaMOBSuite) loadFixture(fixture []measurementFixture, direction model.EnergyDirectionType) {
	descs := []model.MeasurementDescriptionDataType{}
	values := []model.MeasurementDataType{}
	params := []model.ElectricalConnectionParameterDescriptionDataType{}
	for i, f := range fixture {
		descs = append(descs, model.MeasurementDescriptionDataType{
			MeasurementId:   util.Ptr(model.MeasurementIdType(f.id)),
			MeasurementType: util.Ptr(f.mType),
			ScopeType:       util.Ptr(f.scope),
		})
		values = append(values, model.MeasurementDataType{
			MeasurementId: util.Ptr(model.MeasurementIdType(f.id)),
			Value:         model.NewScaledNumberType(f.value),
		})
		params = append(params, model.ElectricalConnectionParameterDescriptionDataType{
			ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
			ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(i)),
			MeasurementId:          util.Ptr(model.MeasurementIdType(f.id)),
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
	_, err = rEl.UpdateData(true, model.FunctionTypeElectricalConnectionDescriptionListData,
		&model.ElectricalConnectionDescriptionListDataType{
			ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{
				{
					ElectricalConnectionId:  util.Ptr(model.ElectricalConnectionIdType(0)),
					PositiveEnergyDirection: util.Ptr(direction),
				},
			},
		}, nil, nil)
	assert.Nil(s.T(), err)
	_, err = rEl.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData,
		&model.ElectricalConnectionParameterDescriptionListDataType{ElectricalConnectionParameterDescriptionData: params}, nil, nil)
	assert.Nil(s.T(), err)
}

func (s *MaMOBSuite) Test_IncompatibleEntity() {
	_, err := s.sut.Power(s.mockRemoteEntity)
	assert.Equal(s.T(), api.ErrNoCompatibleEntity, err)
	for _, fn := range []func() (float64, error){
		func() (float64, error) { return s.sut.EnergyCharged(s.mockRemoteEntity) },
		func() (float64, error) { return s.sut.EnergyDischarged(s.mockRemoteEntity) },
		func() (float64, error) { return s.sut.StateOfCharge(s.mockRemoteEntity) },
		func() (float64, error) { return s.sut.StateOfHealth(s.mockRemoteEntity) },
		func() (float64, error) { return s.sut.Temperature(s.mockRemoteEntity) },
	} {
		_, err := fn()
		assert.Equal(s.T(), api.ErrNoCompatibleEntity, err)
	}
}

func (s *MaMOBSuite) Test_NoData() {
	_, err := s.sut.Power(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	_, err = s.sut.StateOfCharge(s.monitoredEntity)
	assert.NotNil(s.T(), err)
}

func (s *MaMOBSuite) Test_Values() {
	s.loadFixture(batteryFixture(model.ScopeTypeTypeACPowerTotal), model.EnergyDirectionTypeConsume)

	power, err := s.sut.Power(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1500.0, power)

	cases := map[string]struct {
		fn   func() (float64, error)
		want float64
	}{
		"charged":     {func() (float64, error) { return s.sut.EnergyCharged(s.monitoredEntity) }, 1000},
		"discharged":  {func() (float64, error) { return s.sut.EnergyDischarged(s.monitoredEntity) }, 800},
		"soc":         {func() (float64, error) { return s.sut.StateOfCharge(s.monitoredEntity) }, 63.5},
		"soh":         {func() (float64, error) { return s.sut.StateOfHealth(s.monitoredEntity) }, 98},
		"temperature": {func() (float64, error) { return s.sut.Temperature(s.monitoredEntity) }, 21.5},
	}
	for name, c := range cases {
		got, err := c.fn()
		assert.Nil(s.T(), err, name)
		assert.Equal(s.T(), c.want, got, name)
	}
}

func (s *MaMOBSuite) Test_ProduceDirectionIsConverted() {
	// a battery declaring "produce" as positive direction that reports +1500 W
	// is discharging -> -1500 W in the charge-positive convention
	s.loadFixture(batteryFixture(model.ScopeTypeTypeACPowerTotal), model.EnergyDirectionTypeProduce)

	power, err := s.sut.Power(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), -1500.0, power)

	// energies are counters and are not sign converted
	charged, err := s.sut.EnergyCharged(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1000.0, charged)
}

func (s *MaMOBSuite) Test_DCPowerFallback() {
	s.loadFixture(batteryFixture(model.ScopeTypeTypeDCPower), model.EnergyDirectionTypeConsume)

	power, err := s.sut.Power(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1500.0, power)
}

func (s *MaMOBSuite) Test_InvalidPowerDoesNotFallBack() {
	s.loadFixture(batteryFixture(model.ScopeTypeTypeACPowerTotal), model.EnergyDirectionTypeConsume)

	rMeas := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	_, fErr := rMeas.UpdateData(true, model.FunctionTypeMeasurementListData, &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(1500),
				ValueState:    util.Ptr(model.MeasurementValueStateTypeError),
			},
		},
	}, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err := s.sut.Power(s.monitoredEntity)
	assert.Equal(s.T(), api.ErrDataInvalid, err)
}

// kostalBatteryFixture is the battery entity [2,1,1] of a KOSTAL PLENTICORE G3
// as recorded via "Download diagnostics" (scale -1 applied): DC values, energy
// counters under dcChargeEnergy/dcDischargeEnergy, no value state on most
// values, no value for stateOfHealth.
func kostalBatteryFixture() []measurementFixture {
	return []measurementFixture{
		{7, model.MeasurementTypeTypePercentage, model.ScopeTypeTypeStateOfCharge, 69.0},
		{8, model.MeasurementTypeTypeEnergy, model.ScopeTypeTypeStateOfEnergy, 304704.0},
		{10, model.MeasurementTypeTypeCapacity, model.ScopeTypeTypeUseableCapacity, 4416.0},
		{4, model.MeasurementTypeTypePower, model.ScopeTypeTypeDCPower, -1250.0},
		{2, model.MeasurementTypeTypeCurrent, model.ScopeTypeTypeDCCurrent, -4.7},
		{5, model.MeasurementTypeTypeVoltage, model.ScopeTypeTypeDCVoltage, 263.7},
		{1, model.MeasurementTypeTypeEnergy, model.ScopeTypeTypeDCChargeEnergy, 0},
		{3, model.MeasurementTypeTypeEnergy, model.ScopeTypeTypeDCDischargeEnergy, 0},
		{0, model.MeasurementTypeTypeCount, model.ScopeTypeTypeLoadCycleCount, 1357},
		{6, model.MeasurementTypeTypeTemperature, model.ScopeTypeTypeComponentTemperature, 27.5},
	}
}

func (s *MaMOBSuite) Test_KostalBattery() {
	s.loadFixture(kostalBatteryFixture(), model.EnergyDirectionTypeConsume)

	cases := map[string]struct {
		fn   func() (float64, error)
		want float64
	}{
		"power (dc fallback, discharging)": {func() (float64, error) { return s.sut.Power(s.monitoredEntity) }, -1250},
		"current (discharging)":            {func() (float64, error) { return s.sut.Current(s.monitoredEntity) }, -4.7},
		"voltage":                          {func() (float64, error) { return s.sut.Voltage(s.monitoredEntity) }, 263.7},
		"soc":                              {func() (float64, error) { return s.sut.StateOfCharge(s.monitoredEntity) }, 69},
		// reported as is (the KOSTAL value is off by a factor of 100: 69 % x 4416 Wh)
		"state of energy":          {func() (float64, error) { return s.sut.StateOfEnergy(s.monitoredEntity) }, 304704},
		"useable capacity":         {func() (float64, error) { return s.sut.UseableCapacity(s.monitoredEntity) }, 4416},
		"load cycles":              {func() (float64, error) { return s.sut.LoadCycleCount(s.monitoredEntity) }, 1357},
		"temperature":              {func() (float64, error) { return s.sut.Temperature(s.monitoredEntity) }, 27.5},
		"charged (dc fallback)":    {func() (float64, error) { return s.sut.EnergyCharged(s.monitoredEntity) }, 0},
		"discharged (dc fallback)": {func() (float64, error) { return s.sut.EnergyDischarged(s.monitoredEntity) }, 0},
	}
	for name, c := range cases {
		got, err := c.fn()
		assert.Nil(s.T(), err, name)
		assert.InDelta(s.T(), c.want, got, 1e-9, name)
	}

	// no value reported for stateOfHealth
	_, err := s.sut.StateOfHealth(s.monitoredEntity)
	assert.NotNil(s.T(), err)
}

func (s *MaMOBSuite) Test_NewValuesIncompatibleEntity() {
	for _, fn := range []func() (float64, error){
		func() (float64, error) { return s.sut.StateOfEnergy(s.mockRemoteEntity) },
		func() (float64, error) { return s.sut.UseableCapacity(s.mockRemoteEntity) },
		func() (float64, error) { return s.sut.LoadCycleCount(s.mockRemoteEntity) },
		func() (float64, error) { return s.sut.Voltage(s.mockRemoteEntity) },
		func() (float64, error) { return s.sut.Current(s.mockRemoteEntity) },
	} {
		_, err := fn()
		assert.Equal(s.T(), api.ErrNoCompatibleEntity, err)
	}
}

func (s *MaMOBSuite) Test_ChargeScopePreferredOverDC() {
	fixture := append(batteryFixture(model.ScopeTypeTypeACPowerTotal),
		measurementFixture{20, model.MeasurementTypeTypeEnergy, model.ScopeTypeTypeDCChargeEnergy, 5})
	s.loadFixture(fixture, model.EnergyDirectionTypeConsume)
	got, err := s.sut.EnergyCharged(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1000.0, got) // "charge" wins over "dcChargeEnergy"
}
