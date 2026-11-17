package moi

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	internal "github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// handle SPINE events
func (e *MOI) HandleEvent(payload spineapi.EventPayload) {
	// only about events from an inverter entity or device changes for this remote device

	if !e.IsCompatibleEntityType(payload.Entity) {
		return
	}

	if internal.IsEntityAdded(payload) {
		e.deviceConnected(payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.MeasurementDescriptionListDataType:
		e.deviceMeasurementDescriptionDataUpdate(payload.Entity)

	case *model.MeasurementListDataType:
		e.deviceMeasurementDataUpdate(payload)
	}
}

// process required steps when a device is connected
func (e *MOI) deviceConnected(entity spineapi.EntityRemoteInterface) {
	if electricalConnection, err := client.NewElectricalConnection(e.LocalEntity, entity); err == nil {
		if !electricalConnection.HasSubscription() {
			if _, err := electricalConnection.Subscribe(); err != nil {
				logging.Log().Error(err)
			}
		}

		// get electrical connection parameter
		if _, err := electricalConnection.RequestDescriptions(nil, nil); err != nil {
			logging.Log().Error(err)
		}

		if _, err := electricalConnection.RequestParameterDescriptions(nil, nil); err != nil {
			logging.Log().Error(err)
		}
	}

	if measurement, err := client.NewMeasurement(e.LocalEntity, entity); err == nil {
		if !measurement.HasSubscription() {
			if _, err := measurement.Subscribe(); err != nil {
				logging.Log().Error(err)
			}
		}

		// get measurement parameters
		if _, err := measurement.RequestDescriptions(nil, nil); err != nil {
			logging.Log().Error(err)
		}

		if _, err := measurement.RequestConstraints(nil, nil); err != nil {
			logging.Log().Error(err)
		}
	}
}

// the measurement description data of a device was updated
func (e *MOI) deviceMeasurementDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if measurement, err := client.NewMeasurement(e.LocalEntity, entity); err == nil {
		// measurement descriptions received, now get the data
		if _, err := measurement.RequestData(nil, nil); err != nil {
			logging.Log().Error("Error getting measurement list values:", err)
		}
	}
}

// scopeEvents maps measurement scopes to the event they trigger
var scopeEvents = []struct {
	scope model.ScopeTypeType
	event api.EventType
}{
	{model.ScopeTypeTypeACPowerTotal, DataUpdatePower},
	{model.ScopeTypeTypeACPower, DataUpdatePowerPerPhase},
	{model.ScopeTypeTypeACEnergyConsumed, DataUpdateEnergyConsumed},
	{model.ScopeTypeTypeACEnergyProduced, DataUpdateEnergyProduced},
	{model.ScopeTypeTypeACYieldTotal, DataUpdateYieldTotal},
	{model.ScopeTypeTypeACCurrent, DataUpdateCurrentPerPhase},
	{model.ScopeTypeTypeACVoltage, DataUpdateVoltagePerPhase},
	{model.ScopeTypeTypeACFrequency, DataUpdateFrequency},
}

// the measurement data of a device was updated
func (e *MOI) deviceMeasurementDataUpdate(payload spineapi.EventPayload) {
	if e.EventCB == nil {
		return
	}
	measurement, err := client.NewMeasurement(e.LocalEntity, payload.Entity)
	if err != nil {
		return
	}
	for _, se := range scopeEvents {
		filter := model.MeasurementDescriptionDataType{ScopeType: util.Ptr(se.scope)}
		if measurement.CheckEventPayloadDataForFilter(payload.Data, filter) {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, se.event)
		}
	}
}
