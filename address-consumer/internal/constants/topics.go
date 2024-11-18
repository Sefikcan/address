package constants

type KafkaTopicsStruct struct {
	AddressCreated string
	AddressUpdated string
	AddressDeleted string
}

var KafkaTopics = KafkaTopicsStruct{
	AddressUpdated: "address-updated",
	AddressCreated: "address-created",
	AddressDeleted: "address-deleted",
}
