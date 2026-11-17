package mob

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
func (e *MOB) HandleEvent(payload spineapi.EventPayload) {
	// only about events from a battery entity or device changes for this remote device

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
func (e *MOB) deviceConnected(entity spineapi.EntityRemoteInterface) {
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
func (e *MOB) deviceMeasurementDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
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
	{model.ScopeTypeTypeDCPower, DataUpdatePower},
	{model.ScopeTypeTypeCharge, DataUpdateEnergyCharged},
	{model.ScopeTypeTypeDischarge, DataUpdateEnergyDischarged},
	{model.ScopeTypeTypeStateOfCharge, DataUpdateStateOfCharge},
	{model.ScopeTypeTypeStateOfHealth, DataUpdateStateOfHealth},
	{model.ScopeTypeTypeComponentTemperature, DataUpdateTemperature},
}

// the measurement data of a device was updated
func (e *MOB) deviceMeasurementDataUpdate(payload spineapi.EventPayload) {
	if e.EventCB == nil {
		return
	}
	measurement, err := client.NewMeasurement(e.LocalEntity, payload.Entity)
	if err != nil {
		return
	}
	sent := map[api.EventType]bool{}
	for _, se := range scopeEvents {
		if sent[se.event] {
			continue
		}
		filter := model.MeasurementDescriptionDataType{ScopeType: util.Ptr(se.scope)}
		if measurement.CheckEventPayloadDataForFilter(payload.Data, filter) {
			sent[se.event] = true
			e.EventCB(payload.Ski, payload.Device, payload.Entity, se.event)
		}
	}
}
