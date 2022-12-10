package constants

//  Shield Types
const (
	Invulnerable string = "INVULNERABLE"
	Active       string = "Active"
	Broken       string = "BROKEN"
	Disabled     string = "Disabled"
	Unavailable  string = "UNAVAILABLE"
)

//  Planet Types
const (
	Primitive string = "PRIMITIVE"
	Resource  string = "RESOURCE"
	Bot       string = "BOT"
	User      string = "USER"
	Base      string = "BASE"
)

// Resources
const (
	Water    string = "WATER"
	Graphene string = "GRAPHENE"
	Shelio   string = "SHELIO"
)

// Experience Constant Types
const (
	UserExperiences string = "USER"
	ClanExperiences string = "CLAN"
)

// Paltan Roles
const (
	PaltanLeader    string = "LEADER"
	PaltanSubLeader string = "SUB_LEADER"
	PaltanMember    string = "MEMBER"
)

//  Building States
const (
	WorkingState   string = "WORKING"
	UpgradingState string = "UPGRADING"
)

//  Building Types
const (
	WaterMiningPlant        string = "WATER_MINING_PLANT"
	GrapheneMiningPlant     string = "GRAPHENE_MINING_PLANT"
	Shield                  string = "SHIELD"
	PopulationControlCenter string = "POPULATION_CONTROL_CENTER"
	AttackProductionCenter  string = "ATTACK_PRODUCTION_CENTER"
	DefenceProductionCenter string = "DEFENCE_PRODUCTION_CENTER"
	DiamondStorage          string = "DIAMOND_STORAGE"
	WaterPressureTank       string = "WATER_PRESSURE_TANK"
	WarRoom                 string = "WAR_ROOM"
	CommunicationCenter     string = "COMMUNICATION_CENTER"
	ResearchLab             string = "RESEARCH_LAB"
	Library                 string = "LIBRARY"
)

//  Unit Types
const (
	Scout              string = "SCOUT"
	ResourceCarrier    string = "RESOURCE CARRIER"
	AircraftCarrier    string = "AIRCRAFT CARRIER"
	Attacker           string = "ATTACKER"
	DefenceShipCarrier string = "DEFENCE SHIP CARRIER"
	Defender           string = "DEFENDER"
)

//  Ship Types
const (
	SystemDistanceBlocks float64 = 150
	SectorDistanceBlocks float64 = 20
)

//  Mission Constants
const (
	AttackMission   string = "ATTACK"
	SpyMission      string = "SPY"
	DepartureState  string = "DEPARTURE"
	InProgressState string = "IN_PROGRESS"
	ReturningState  string = "RETURNING"
	CompletedState  string = "COMPLETED"
)

// Authenticators
const (
	GoogleAuthenticator   string = "GOOGLE"
	FacebookAuthenticator string = "FACEBOOK"
)
