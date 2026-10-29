package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/TwiN/gatus/v5/config"
	"github.com/TwiN/gatus/v5/config/endpoint/ui"
	"github.com/TwiN/gatus/v5/storage/store"
	"github.com/TwiN/gatus/v5/storage/store/common"
	"github.com/TwiN/gatus/v5/storage/store/common/paging"
	"github.com/gofiber/fiber/v2"
)

const (
	badgeColorHexAwesome  = "#40cc11"
	badgeColorHexGreat    = "#94cc11"
	badgeColorHexGood     = "#ccd311"
	badgeColorHexPassable = "#ccb311"
	badgeColorHexBad      = "#cc8111"
	badgeColorHexVeryBad  = "#c7130a"
)

const (
	HealthStatusUp      = "up"
	HealthStatusDown    = "down"
	HealthStatusUnknown = "?"

	badgeFontFamily = "Microsoft YaHei,PingFang SC,Noto Sans SC,sans-serif"
)

var (
	badgeColors = []string{badgeColorHexAwesome, badgeColorHexGreat, badgeColorHexGood, badgeColorHexPassable, badgeColorHexBad}
)

// UptimeBadge handles the automatic generation of badge based on the group name and endpoint name passed.
//
// Valid values for :duration -> 30d, 7d, 24h, 1h
func UptimeBadge(c *fiber.Ctx) error {
	duration := c.Params("duration")
	var from time.Time
	switch duration {
	case "30d":
		from = time.Now().Add(-30 * 24 * time.Hour)
	case "7d":
		from = time.Now().Add(-7 * 24 * time.Hour)
	case "24h":
		from = time.Now().Add(-24 * time.Hour)
	case "1h":
		from = time.Now().Add(-2 * time.Hour) // Because uptime metrics are stored by hour, we have to cheat a little
	default:
		return c.Status(400).SendString("Durations supported: 30d, 7d, 24h, 1h")
	}
	key, err := url.QueryUnescape(c.Params("key"))
	if err != nil {
		return c.Status(400).SendString("invalid key encoding")
	}
	uptime, err := store.Get().GetUptimeByKey(key, from, time.Now())
	if err != nil {
		if errors.Is(err, common.ErrEndpointNotFound) {
			return c.Status(404).SendString(err.Error())
		} else if errors.Is(err, common.ErrInvalidTimeRange) {
			return c.Status(400).SendString(err.Error())
		}
		return c.Status(500).SendString(err.Error())
	}
	c.Set("Content-Type", "image/svg+xml")
	c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Set("Expires", "0")
	return c.Status(200).Send(generateUptimeBadgeSVG(duration, uptime))
}

// ResponseTimeBadge handles the automatic generation of badge based on the group name and endpoint name passed.
//
// Valid values for :duration -> 30d, 7d, 24h, 1h
func ResponseTimeBadge(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		duration := c.Params("duration")
		var from time.Time
		switch duration {
		case "30d":
			from = time.Now().Add(-30 * 24 * time.Hour)
		case "7d":
			from = time.Now().Add(-7 * 24 * time.Hour)
		case "24h":
			from = time.Now().Add(-24 * time.Hour)
		case "1h":
			from = time.Now().Add(-2 * time.Hour) // Because response time metrics are stored by hour, we have to cheat a little
		default:
			return c.Status(400).SendString("Durations supported: 30d, 7d, 24h, 1h")
		}
		key, err := url.QueryUnescape(c.Params("key"))
		if err != nil {
			return c.Status(400).SendString("invalid key encoding")
		}
		averageResponseTime, err := store.Get().GetAverageResponseTimeByKey(key, from, time.Now())
		if err != nil {
			if errors.Is(err, common.ErrEndpointNotFound) {
				return c.Status(404).SendString(err.Error())
			} else if errors.Is(err, common.ErrInvalidTimeRange) {
				return c.Status(400).SendString(err.Error())
			}
			return c.Status(500).SendString(err.Error())
		}
		c.Set("Content-Type", "image/svg+xml")
		c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Set("Expires", "0")
		return c.Status(200).Send(generateResponseTimeBadgeSVG(duration, averageResponseTime, key, cfg))
	}
}

// HealthBadge handles the automatic generation of badge based on the group name and endpoint name passed.
func HealthBadge(c *fiber.Ctx) error {
	key, err := url.QueryUnescape(c.Params("key"))
	if err != nil {
		return c.Status(400).SendString("invalid key encoding")
	}
	pagingConfig := paging.NewEndpointStatusParams()
	status, err := store.Get().GetEndpointStatusByKey(key, pagingConfig.WithResults(1, 1))
	if err != nil {
		if errors.Is(err, common.ErrEndpointNotFound) {
			return c.Status(404).SendString(err.Error())
		} else if errors.Is(err, common.ErrInvalidTimeRange) {
			return c.Status(400).SendString(err.Error())
		}
		return c.Status(500).SendString(err.Error())
	}
	healthStatus := HealthStatusUnknown
	if len(status.Results) > 0 {
		if status.Results[0].Success {
			healthStatus = HealthStatusUp
		} else {
			healthStatus = HealthStatusDown
		}
	}
	c.Set("Content-Type", "image/svg+xml")
	c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Set("Expires", "0")
	return c.Status(200).Send(generateHealthBadgeSVG(healthStatus))
}

func HealthBadgeShields(c *fiber.Ctx) error {
	key, err := url.QueryUnescape(c.Params("key"))
	if err != nil {
		return c.Status(400).SendString("invalid key encoding")
	}
	pagingConfig := paging.NewEndpointStatusParams()
	status, err := store.Get().GetEndpointStatusByKey(key, pagingConfig.WithResults(1, 1))
	if err != nil {
		if errors.Is(err, common.ErrEndpointNotFound) {
			return c.Status(404).SendString(err.Error())
		} else if errors.Is(err, common.ErrInvalidTimeRange) {
			return c.Status(400).SendString(err.Error())
		}
		return c.Status(500).SendString(err.Error())
	}
	healthStatus := HealthStatusUnknown
	if len(status.Results) > 0 {
		if status.Results[0].Success {
			healthStatus = HealthStatusUp
		} else {
			healthStatus = HealthStatusDown
		}
	}
	c.Set("Content-Type", "application/json")
	c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Set("Expires", "0")
	jsonData, err := generateHealthBadgeShields(healthStatus)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.Status(200).Send(jsonData)
}

func generateUptimeBadgeSVG(duration string, uptime float64) []byte {
	var labelWidth, valueWidth, valueWidthAdjustment int
	switch duration {
	case "30d", "7d", "24h", "1h":
		labelWidth = 110
	default:
	}
	color := getBadgeColorFromUptime(uptime)
	sanitizedValue := strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", uptime*100), "0"), ".") + "%"
	if strings.Contains(sanitizedValue, ".") {
		valueWidthAdjustment = -10
	}
	valueWidth = (len(sanitizedValue) * 11) + valueWidthAdjustment
	width := labelWidth + valueWidth
	labelX := labelWidth / 2
	valueX := labelWidth + (valueWidth / 2)
	svg := []byte(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="20">
  <linearGradient id="b" x2="0" y2="100%%">
    <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
    <stop offset="1" stop-opacity=".1"/>
  </linearGradient>
  <mask id="a">
    <rect width="%d" height="20" rx="3" fill="#fff"/>
  </mask>
  <g mask="url(#a)">
    <path fill="#555" d="M0 0h%dv20H0z"/>
    <path fill="%s" d="M%d 0h%dv20H%dz"/>
    <path fill="url(#b)" d="M0 0h%dv20H0z"/>
  </g>
  <g fill="#fff" text-anchor="middle" font-family="%s" font-size="11">
    <text x="%d" y="15" fill="#010101" fill-opacity=".3">
      可用率 %s
    </text>
    <text x="%d" y="14">
      可用率 %s
    </text>
    <text x="%d" y="15" fill="#010101" fill-opacity=".3">
      %s
    </text>
    <text x="%d" y="14">
      %s
    </text>
  </g>
</svg>`, width, width, labelWidth, color, labelWidth, valueWidth, labelWidth, width, badgeFontFamily, labelX, badgeDurationLabel(duration), labelX, badgeDurationLabel(duration), valueX, sanitizedValue, valueX, sanitizedValue))
	return svg
}

func badgeDurationLabel(duration string) string {
	switch duration {
	case "30d":
		return "30天"
	case "7d":
		return "7天"
	case "24h":
		return "24小时"
	case "1h":
		return "1小时"
	default:
		return duration
	}
}

func healthStatusLabel(healthStatus string) string {
	switch healthStatus {
	case HealthStatusUp:
		return "正常"
	case HealthStatusDown:
		return "异常"
	default:
		return "未知"
	}
}

func getBadgeColorFromUptime(uptime float64) string {
	if uptime >= 0.975 {
		return badgeColorHexAwesome
	} else if uptime >= 0.95 {
		return badgeColorHexGreat
	} else if uptime >= 0.9 {
		return badgeColorHexGood
	} else if uptime >= 0.8 {
		return badgeColorHexPassable
	} else if uptime >= 0.65 {
		return badgeColorHexBad
	}
	return badgeColorHexVeryBad
}

func generateResponseTimeBadgeSVG(duration string, averageResponseTime int, key string, cfg *config.Config) []byte {
	var labelWidth, valueWidth int
	switch duration {
	case "30d":
		labelWidth = 110
	case "7d":
		labelWidth = 105
	case "24h":
		labelWidth = 110
	case "1h":
		labelWidth = 105
	default:
	}
	color := getBadgeColorFromResponseTime(averageResponseTime, key, cfg)
	sanitizedValue := strconv.Itoa(averageResponseTime) + "ms"
	valueWidth = len(sanitizedValue) * 11
	width := labelWidth + valueWidth
	labelX := labelWidth / 2
	valueX := labelWidth + (valueWidth / 2)
	svg := []byte(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="20">
  <linearGradient id="b" x2="0" y2="100%%">
    <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
    <stop offset="1" stop-opacity=".1"/>
  </linearGradient>
  <mask id="a">
    <rect width="%d" height="20" rx="3" fill="#fff"/>
  </mask>
  <g mask="url(#a)">
    <path fill="#555" d="M0 0h%dv20H0z"/>
    <path fill="%s" d="M%d 0h%dv20H%dz"/>
    <path fill="url(#b)" d="M0 0h%dv20H0z"/>
  </g>
  <g fill="#fff" text-anchor="middle" font-family="%s" font-size="11">
    <text x="%d" y="15" fill="#010101" fill-opacity=".3">
      响应 %s
    </text>
    <text x="%d" y="14">
      响应 %s
    </text>
    <text x="%d" y="15" fill="#010101" fill-opacity=".3">
      %s
    </text>
    <text x="%d" y="14">
      %s
    </text>
  </g>
</svg>`, width, width, labelWidth, color, labelWidth, valueWidth, labelWidth, width, badgeFontFamily, labelX, badgeDurationLabel(duration), labelX, badgeDurationLabel(duration), valueX, sanitizedValue, valueX, sanitizedValue))
	return svg
}

func getBadgeColorFromResponseTime(responseTime int, key string, cfg *config.Config) string {
	thresholds := ui.GetDefaultConfig().Badge.ResponseTime.Thresholds
	if endpoint := cfg.GetEndpointByKey(key); endpoint != nil {
		thresholds = endpoint.UIConfig.Badge.ResponseTime.Thresholds
	}
	// the threshold config requires 5 values, so we can be sure it's set here
	for i := range 5 {
		if responseTime <= thresholds[i] {
			return badgeColors[i]
		}
	}
	return badgeColorHexVeryBad
}

func generateHealthBadgeSVG(healthStatus string) []byte {
	color := getBadgeColorFromHealth(healthStatus)
	labelWidth := 48
	valueWidth := 40
	statusLabel := healthStatusLabel(healthStatus)

	width := labelWidth + valueWidth
	labelX := labelWidth / 2
	valueX := labelWidth + (valueWidth / 2)
	svg := []byte(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="20">
  <linearGradient id="b" x2="0" y2="100%%">
    <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
    <stop offset="1" stop-opacity=".1"/>
  </linearGradient>
  <mask id="a">
    <rect width="%d" height="20" rx="3" fill="#fff"/>
  </mask>
  <g mask="url(#a)">
    <path fill="#555" d="M0 0h%dv20H0z"/>
    <path fill="%s" d="M%d 0h%dv20H%dz"/>
    <path fill="url(#b)" d="M0 0h%dv20H0z"/>
  </g>
  <g fill="#fff" text-anchor="middle" font-family="%s" font-size="11">
    <text x="%d" y="15" fill="#010101" fill-opacity=".3">
      健康
    </text>
    <text x="%d" y="14">
      健康
    </text>
    <text x="%d" y="15" fill="#010101" fill-opacity=".3">
      %s
    </text>
    <text x="%d" y="14">
      %s
    </text>
  </g>
</svg>`, width, width, labelWidth, color, labelWidth, valueWidth, labelWidth, width, badgeFontFamily, labelX, labelX, valueX, statusLabel, valueX, statusLabel))

	return svg
}

func generateHealthBadgeShields(healthStatus string) ([]byte, error) {
	color := getBadgeShieldsColorFromHealth(healthStatus)
	data := map[string]interface{}{
		"schemaVersion": 1,
		"label":         "状态",
		"message":       healthStatusLabel(healthStatus),
		"color":         color,
	}
	return json.Marshal(data)
}

func getBadgeColorFromHealth(healthStatus string) string {
	if healthStatus == HealthStatusUp {
		return badgeColorHexAwesome
	} else if healthStatus == HealthStatusDown {
		return badgeColorHexVeryBad
	}
	return badgeColorHexPassable
}

func getBadgeShieldsColorFromHealth(healthStatus string) string {
	if healthStatus == HealthStatusUp {
		return "brightgreen"
	} else if healthStatus == HealthStatusDown {
		return "red"
	}
	return "yellow"
}
