package nex

import (
	"fmt"
	"os"
	"strconv"

	"github.com/PretendoNetwork/nex-go/v2"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	"github.com/PretendoNetwork/yo-kai-watch-2/globals"
)

func StartSecureServer() {
	globals.SecureServer = nex.NewPRUDPServer()
	globals.SecureServer.ByteStreamSettings.UseStructureHeader = true

	globals.SecureEndpoint = nex.NewPRUDPEndPoint(1)
	globals.SecureEndpoint.IsSecureEndPoint = true
	globals.SecureEndpoint.ServerAccount = globals.SecureServerAccount
	globals.SecureEndpoint.AccountDetailsByPID = globals.AccountDetailsByPID
	globals.SecureEndpoint.AccountDetailsByUsername = globals.AccountDetailsByUsername
	globals.SecureServer.BindPRUDPEndPoint(globals.SecureEndpoint)

	globals.SecureServer.LibraryVersions.SetDefault(nex.NewLibraryVersion(3, 6, 1))
	globals.SecureServer.AccessKey = "7ab183bb"

	globals.SecureEndpoint.OnData(func(packet nex.PacketInterface) {
		request := packet.RMCMessage()

		//trace is never getting what they want
		fmt.Println("==Yo-kai Watch 2- Secure==")
		fmt.Printf("Protocol ID: %d\n", request.ProtocolID)
		fmt.Printf("Method ID: %d\n", request.MethodID)
		fmt.Println("===============")
	})

	globals.SecureEndpoint.OnError(func(err *nex.Error) {
		globals.Logger.Errorf("Secure: %v", err)
	})

	globals.MatchmakingManager = common_globals.NewMatchmakingManager(globals.SecureEndpoint, globals.Postgres)

	registerCommonSecureServerProtocols()

	port, _ := strconv.Atoi(os.Getenv("PN_YKW2_SECURE_SERVER_PORT"))

	globals.SecureServer.Listen(port)
}
