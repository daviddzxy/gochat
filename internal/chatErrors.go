package internal

type ClientAlreadyInRoom struct{}

func (e ClientAlreadyInRoom) Error() string {
	return "client is already in a room"
}

type ClientDoesNotBelongToAnyRoom struct{}

func (e ClientDoesNotBelongToAnyRoom) Error() string {
	return "client does not belong to any room"
}

var ClientAlreadyInRoomError = ClientAlreadyInRoom{}
var ClientDoesNotBelongToAnyRoomError = ClientDoesNotBelongToAnyRoom{}
