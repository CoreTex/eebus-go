package moi

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *MaMOISuite) Test_Events() {
	payload := spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}
	s.sut.HandleEvent(payload)

	payload.Entity = s.monitoredEntity
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeEntityChange
	payload.ChangeType = spineapi.ElementChangeAdd
	s.sut.HandleEvent(payload)

	payload.ChangeType = spineapi.ElementChangeRemove
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeAdd
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeUpdate
	payload.Data = util.Ptr(model.MeasurementDescriptionListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.MeasurementListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.NodeManagementUseCaseDataType{})
	s.sut.HandleEvent(payload)
}

func (s *MaMOISuite) Test_Failures() {
	s.sut.deviceConnected(s.mockRemoteEntity)

	s.sut.deviceMeasurementDescriptionDataUpdate(s.mockRemoteEntity)
}

func (s *MaMOISuite) Test_deviceMeasurementDataUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.monitoredEntity,
	}
	s.sut.deviceMeasurementDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	s.loadFixture(kostalFixture(), model.EnergyDirectionTypeConsume)

	// every scope of the fixture triggers its event
	for _, f := range kostalFixture() {
		s.eventCalled = false
		payload.Data = &model.MeasurementListDataType{
			MeasurementData: []model.MeasurementDataType{
				{
					MeasurementId: util.Ptr(model.MeasurementIdType(f.id)),
					Value:         model.NewScaledNumberType(1),
				},
			},
		}
		s.sut.deviceMeasurementDataUpdate(payload)
		assert.True(s.T(), s.eventCalled, "no event for scope %s", f.scope)
	}

	// no callback configured: nothing happens
	s.eventCalled = false
	s.sut.EventCB = nil
	s.sut.deviceMeasurementDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)
}
