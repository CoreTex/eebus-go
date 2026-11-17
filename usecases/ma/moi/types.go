package moi

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntitiesScenarios` to get the current data
	UseCaseSupportUpdate api.EventType = "ma-moi-UseCaseSupportUpdate"

	// Inverter total active AC power updated
	//
	// Use `Power` to get the current data
	DataUpdatePower api.EventType = "ma-moi-DataUpdatePower"

	// Inverter phase specific active AC power updated
	//
	// Use `PowerPerPhase` to get the current data
	DataUpdatePowerPerPhase api.EventType = "ma-moi-DataUpdatePowerPerPhase"

	// Inverter consumed AC energy updated
	//
	// Use `EnergyConsumed` to get the current data
	DataUpdateEnergyConsumed api.EventType = "ma-moi-DataUpdateEnergyConsumed"

	// Inverter produced AC energy updated
	//
	// Use `EnergyProduced` to get the current data
	DataUpdateEnergyProduced api.EventType = "ma-moi-DataUpdateEnergyProduced"

	// Inverter total AC yield updated
	//
	// Use `YieldTotal` to get the current data
	DataUpdateYieldTotal api.EventType = "ma-moi-DataUpdateYieldTotal"

	// Inverter phase specific AC current updated
	//
	// Use `CurrentPerPhase` to get the current data
	DataUpdateCurrentPerPhase api.EventType = "ma-moi-DataUpdateCurrentPerPhase"

	// Inverter phase specific AC voltage updated
	//
	// Use `VoltagePerPhase` to get the current data
	DataUpdateVoltagePerPhase api.EventType = "ma-moi-DataUpdateVoltagePerPhase"

	// Inverter AC frequency updated
	//
	// Use `Frequency` to get the current data
	DataUpdateFrequency api.EventType = "ma-moi-DataUpdateFrequency"
)
