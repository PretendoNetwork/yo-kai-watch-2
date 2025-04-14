package matchmaking

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	match_making "github.com/PretendoNetwork/nex-protocols-go/v2/match-making"
	"github.com/PretendoNetwork/yo-kai-watch-2/globals"
	nex_match_making_database "github.com/PretendoNetwork/yo-kai-watch-2/nex/match_making/database"
)

func FindByOwner(err error, packet nex.PacketInterface, callID uint32, id types.PID, resultRange types.ResultRange) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	connection := packet.Sender().(*nex.PRUDPConnection)
	endpoint := connection.Endpoint().(*nex.PRUDPEndPoint)

	globals.MatchmakingManager.Mutex.RLock()

	gatheringHolders, nexError := nex_match_making_database.FindMatchmakeSessionsByOwner(globals.MatchmakingManager, packet.Sender().(*nex.PRUDPConnection), id, resultRange)
	if nexError != nil {
		globals.MatchmakingManager.Mutex.RUnlock()
		return nil, nexError
	}

	globals.MatchmakingManager.Mutex.RUnlock()

	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	gatheringHolders.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseBody)
	rmcResponse.ProtocolID = match_making.ProtocolID
	rmcResponse.MethodID = match_making.MethodFindByOwner
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
