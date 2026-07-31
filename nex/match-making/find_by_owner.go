package nex_match_making

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	match_making "github.com/PretendoNetwork/nex-protocols-go/v2/match-making"
	match_making_types "github.com/PretendoNetwork/nex-protocols-go/v2/match-making/types"
	"github.com/PretendoNetwork/yo-kai-watch-2/globals"
	nex_match_making_database "github.com/PretendoNetwork/yo-kai-watch-2/nex/match-making/database"
)

func FindByOwner(err error, packet nex.PacketInterface, callID uint32, id types.PID, resultRange types.ResultRange) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	connection := packet.Sender().(*nex.PRUDPConnection)
	endpoint := connection.Endpoint().(*nex.PRUDPEndPoint)

	globals.MatchmakingManager.Mutex.RLock()

	gatheringIDs, nexError := nex_match_making_database.FindGatheringIDsByOwner(globals.MatchmakingManager, connection, id, resultRange)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		globals.MatchmakingManager.Mutex.RUnlock()
		return nil, nexError
	}

	gatheringHolders := types.NewList[match_making_types.GatheringHolder]()
	if len(gatheringIDs) > 0 {
		gatheringHolders, nexError = nex_match_making_database.GetDetailedGatheringsByID(globals.MatchmakingManager, uint64(connection.PID()), gatheringIDs)
		if nexError != nil {
			common_globals.Logger.Error(nexError.Error())
			globals.MatchmakingManager.Mutex.RUnlock()
			return nil, nexError
		}
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
