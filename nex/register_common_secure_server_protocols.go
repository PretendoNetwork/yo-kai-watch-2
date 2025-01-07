package nex

import (
	"github.com/PretendoNetwork/nex-go/v2/types"
	commonnattraversal "github.com/PretendoNetwork/nex-protocols-common-go/v2/nat-traversal"
	commonsecure "github.com/PretendoNetwork/nex-protocols-common-go/v2/secure-connection"
	nattraversal "github.com/PretendoNetwork/nex-protocols-go/v2/nat-traversal"
	secure "github.com/PretendoNetwork/nex-protocols-go/v2/secure-connection"
	"github.com/PretendoNetwork/yo-kai-watch-2/globals"

	commonmatchmaking "github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making"
	commonmatchmakingext "github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making-ext"
	commonmatchmakeextension "github.com/PretendoNetwork/nex-protocols-common-go/v2/matchmake-extension"
	matchmaking "github.com/PretendoNetwork/nex-protocols-go/v2/match-making"
	matchmakingext "github.com/PretendoNetwork/nex-protocols-go/v2/match-making-ext"
	matchmakeextension "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-extension"

	"github.com/PretendoNetwork/nex-go/v2"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	matchmakingtypes "github.com/PretendoNetwork/nex-protocols-go/v2/match-making/types"
	notifications_types "github.com/PretendoNetwork/nex-protocols-go/v2/notifications/types"
	ranking "github.com/PretendoNetwork/nex-protocols-go/v2/ranking"
	"strconv"
	"strings"
)
func updateNotificationData(err error, packet nex.PacketInterface, callID uint32, uiType *types.PrimitiveU32, uiParam1 *types.PrimitiveU32, uiParam2 *types.PrimitiveU32, strParam *types.String) (*nex.RMCMessage, *nex.Error) {
    if err != nil {
        common_globals.Logger.Error(err.Error())
        return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
    }
    connection := packet.Sender().(*nex.PRUDPConnection)
    endpoint := connection.Endpoint().(*nex.PRUDPEndPoint)
    
    rmcResponse := nex.NewRMCSuccess(endpoint, nil)
    rmcResponse.ProtocolID = matchmakeextension.ProtocolID
    rmcResponse.MethodID = matchmakeextension.MethodUpdateNotificationData
    rmcResponse.CallID = callID
    return rmcResponse, nil
}
func getFriendNotificationData(err error, packet nex.PacketInterface, callID uint32, uiType *types.PrimitiveS32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	connection := packet.Sender().(*nex.PRUDPConnection)
	endpoint := connection.Endpoint().(*nex.PRUDPEndPoint)

	dataList := types.NewList[*notifications_types.NotificationEvent]()

	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	dataList.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseBody)
	rmcResponse.ProtocolID = matchmakeextension.ProtocolID
	rmcResponse.MethodID = matchmakeextension.MethodGetFriendNotificationData
	rmcResponse.CallID = callID

	return rmcResponse, nil
}

// Is this needed? -Ash
func cleanupSearchMatchmakeSessionHandler(matchmakeSession *matchmakingtypes.MatchmakeSession) {
	_ = matchmakeSession.Attributes.SetIndex(2, types.NewPrimitiveU32(0))
	matchmakeSession.MatchmakeParam = matchmakingtypes.NewMatchmakeParam()
	matchmakeSession.ApplicationBuffer = types.NewBuffer(make([]byte, 0))
	globals.Logger.Info(matchmakeSession.String())
}

func CreateReportDBRecord(_ *types.PID, _ *types.PrimitiveU32, _ *types.QBuffer) error {
	return nil
}

// from nex-protocols-common-go/matchmaking_utils.go
func compareSearchCriteria[T ~uint16 | ~uint32](original T, search string) bool {
	if search == "" { // * Accept any value
		return true
	}

	before, after, found := strings.Cut(search, ",")
	if found {
		min, err := strconv.ParseUint(before, 10, 64)
		if err != nil {
			return false
		}

		max, err := strconv.ParseUint(after, 10, 64)
		if err != nil {
			return false
		}

		return min <= uint64(original) && max >= uint64(original)
	} else {
		searchNum, err := strconv.ParseUint(before, 10, 64)
		if err != nil {
			return false
		}

		return searchNum == uint64(original)
	}
}

func gameSpecificMatchmakeSessionSearchCriteriaChecksHandler(searchCriteria *matchmakingtypes.MatchmakeSessionSearchCriteria, matchmakeSession *matchmakingtypes.MatchmakeSession) bool {
	original := matchmakeSession.Attributes.Slice()
	search := searchCriteria.Attribs.Slice()
	if len(original) != len(search) {
		return false
	}

	
	for index, originalAttribute := range original {
		// ignore dummy criterias for matchmaking
		// everyone ends up in different rooms if you don't skip these
		if index == 1 || index == 4 {
			continue
		}
		searchAttribute := search[index]

		if !compareSearchCriteria(originalAttribute.Value, searchAttribute.Value) {
			return false
		}
	}
		return true
}

func registerCommonSecureServerProtocols() {
	secureProtocol := secure.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(secureProtocol)
	commonSecureProtocol := commonsecure.NewCommonProtocol(secureProtocol)

	commonSecureProtocol.CreateReportDBRecord = CreateReportDBRecord

	natTraversalProtocol := nattraversal.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(natTraversalProtocol)
	commonnattraversal.NewCommonProtocol(natTraversalProtocol)

	matchMakingProtocol := matchmaking.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(matchMakingProtocol)
	commonmatchmaking.NewCommonProtocol(matchMakingProtocol)

	matchMakingExtProtocol := matchmakingext.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(matchMakingExtProtocol)
	commonmatchmakingext.NewCommonProtocol(matchMakingExtProtocol)

	rankingProtocol := ranking.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(rankingProtocol)

	matchmakeExtensionProtocol := matchmakeextension.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(matchmakeExtensionProtocol)
	commonMatchmakeExtensionProtocol := commonmatchmakeextension.NewCommonProtocol(matchmakeExtensionProtocol)
	matchmakeExtensionProtocol.SetHandlerGetFriendNotificationData(getFriendNotificationData)
	matchmakeExtensionProtocol.SetHandlerUpdateNotificationData(updateNotificationData)


	commonMatchmakeExtensionProtocol.CleanupSearchMatchmakeSession = cleanupSearchMatchmakeSessionHandler
	commonMatchmakeExtensionProtocol.GameSpecificMatchmakeSessionSearchCriteriaChecks = gameSpecificMatchmakeSessionSearchCriteriaChecksHandler

}
