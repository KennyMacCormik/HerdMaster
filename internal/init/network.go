package init

import (
	"github.com/KennyMacCormik/HerdMaster/internal/config"
	"github.com/KennyMacCormik/HerdMaster/internal/network"
	"github.com/KennyMacCormik/HerdMaster/internal/network/http"
	"log/slog"
)

func Endpoint(conf config.Config, lg *slog.Logger) network.Endpoint {
	return http.New(conf, lg)
}
