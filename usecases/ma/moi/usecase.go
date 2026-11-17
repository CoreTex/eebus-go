package moi

import (
	"errors"

	"github.com/enbility/eebus-go/api"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	usecase "github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type MOI struct {
	*usecase.UseCaseBase
}

var _ ucapi.MaMOIInterface = (*MOI)(nil)

// Add support for the Monitoring of Inverter (MOI) use case
// as a Monitoring Appliance actor
//
// Parameters:
//   - localEntity: The local entity which should support the use case
//   - eventCB: The callback to be called when an event is triggered (optional, can be nil)
func NewMOI(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *MOI {
	validActorTypes := []model.UseCaseActorType{model.UseCaseActorTypeInverter}
	validEntityTypes := []model.EntityTypeType{model.EntityTypeTypeInverter}

	// Scenario 1 (momentary power) is mandatory, all others are optional.
	// All scenarios are served by the Measurement and ElectricalConnection
	// features of the inverter entity.
	//
	// The scenario numbering follows what devices announce for
	// "monitoringOfInverter" v1.0.0 (e.g. KOSTAL PLENTICORE G3: 1-7).
	scenarioFeatures := []model.FeatureTypeType{
		model.FeatureTypeTypeElectricalConnection,
		model.FeatureTypeTypeMeasurement,
	}
	useCaseScenarios := make([]api.UseCaseScenario, 0, 7)
	for scenario := model.UseCaseScenarioSupportType(1); scenario <= 7; scenario++ {
		useCaseScenarios = append(useCaseScenarios, api.UseCaseScenario{
			Scenario:       scenario,
			Mandatory:      scenario == 1,
			ServerFeatures: scenarioFeatures,
		})
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeMonitoringAppliance,
		model.UseCaseNameTypeMonitoringOfInverter,
		"1.0.0",
		"release",
		useCaseScenarios,
		eventCB,
		UseCaseSupportUpdate,
		validActorTypes,
		validEntityTypes,
		false,
	)

	uc := &MOI{
		UseCaseBase: usecase,
	}

	_ = localEntity.Device().Events().Subscribe(uc)

	return uc
}

func (e *MOI) AddFeatures() error {
	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeElectricalConnection,
		model.FeatureTypeTypeMeasurement,
	}
	for _, feature := range clientFeatures {
		if f := e.LocalEntity.GetOrAddFeature(feature, model.RoleTypeClient); f == nil {
			return errors.New("failed to add feature: " + string(feature))
		}
	}

	return nil
}
