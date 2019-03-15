package cmd

import (
	"github.com/FlameInTheDark/dtbot/api/geoip"
	"github.com/FlameInTheDark/dtbot/bot"
)

// GeoIPCommand handle dice
func GeoIPCommand(ctx bot.Context) {
	ctx.MetricsCommand("geoip")
	ctx.ReplyEmbed("GeoIP", geoip.GetGeoIP(&ctx))
}
