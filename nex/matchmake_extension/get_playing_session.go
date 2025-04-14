package matchmake_extension

import (
	// "fmt"

	"github.com/PretendoNetwork/nex-go/v2"
	nex_types "github.com/PretendoNetwork/nex-go/v2/types"
	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-extension"
	"github.com/PretendoNetwork/yo-kai-watch-2/globals"
)

// replace later
func GetPlayingSession(err error, packet nex.PacketInterface, callID uint32, lstPid nex_types.List[nex_types.PID]) (*nex.RMCMessage, *nex.Error) {
	//fmt.Println(lstPid)

	l := nex_types.NewUInt32(0)
	rmcResponseStream := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)

	l.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseBody)
	rmcResponse.ProtocolID = matchmake_extension.ProtocolID
	rmcResponse.MethodID = matchmake_extension.MethodGetPlayingSession
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
