package serve

import (
	"github.com/labstack/echo/v4"
	"github.com/warjiang/page-spy-api/config"
	"github.com/warjiang/page-spy-api/container"
	"github.com/warjiang/page-spy-api/util"
)

func Run() {
	err := container.Container().Invoke(func(e *echo.Echo, config *config.Config, staticConfig *config.StaticConfig) {
		if staticConfig != nil {
			hash := staticConfig.GitHash
			version := staticConfig.Version
			if hash == "" {
				hash = "local"
			}
			if version == "" {
				version = "local"
			}
			log.Infof("server info: %s@%s", version, hash)
		}

		for _, ip := range util.GetLocalIPList() {
			log.Infof("LAN address http://%s:%s", ip, config.Port)
		}

		log.Infof("Local address http://localhost:%s", config.Port)
		e.Logger.Fatal(e.Start(":" + config.Port))
	})

	if err != nil {
		log.Fatal(err)
	}
}
