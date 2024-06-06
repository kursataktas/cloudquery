package main

import (
	"context"
	"log"

	"github.com/cloudquery/cloudquery/plugins/destination/snowflake/client"
	internalPlugin "github.com/cloudquery/cloudquery/plugins/destination/snowflake/resources/plugin"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/serve"
)

func main() {
	p := plugin.NewPlugin(
		internalPlugin.Name,
		internalPlugin.Version,
		client.New,
		plugin.WithBuildTargets([]plugin.BuildTarget{
			{OS: plugin.GoOSLinux, Arch: plugin.GoArchAmd64},
			{OS: plugin.GoOSLinux, Arch: plugin.GoArchArm64},
			{OS: plugin.GoOSWindows, Arch: plugin.GoArchAmd64, CGO: true},
			{OS: plugin.GoOSDarwin, Arch: plugin.GoArchAmd64},
			{OS: plugin.GoOSDarwin, Arch: plugin.GoArchArm64},
		}),
		plugin.WithKind(internalPlugin.Kind),
		plugin.WithTeam(internalPlugin.Team),
		plugin.WithJSONSchema(client.JSONSchema),
	)

	if err := serve.Plugin(p, serve.WithDestinationV0V1Server()).Serve(context.Background()); err != nil {
		log.Fatalf("failed to serve plugin: %v", err)
	}
}
