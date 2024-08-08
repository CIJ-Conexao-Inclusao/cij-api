package model

import (
	"cij_api/src/enum"
)

type Config struct {
	FontSize       int                     `json:"font_size"`
	ScreenReader   bool                    `json:"screen_reader"`
	VoiceCapture   bool                    `json:"voice_capture"`
	Theme          enum.ThemeEnum          `json:"theme"`
	ColorBlindness enum.ColorBlindnessEnum `json:"color_blindness"`
	SystemColors   SystemColors            `json:"system_colors"`
}

type SystemColors struct {
	PrimaryColors SystemPrimaryColors `json:"primary_colors"`
	ChartColors   SystemChartColors   `json:"chart_colors"`
}

type SystemPrimaryColors struct {
	PrimaryColor       string `json:"primary_color"`
	SecondaryColor     string `json:"secondary_color"`
	FontColor          string `json:"font_color"`
	SecondaryFontColor string `json:"secondary_font_color"`
	InputColor         string `json:"input_color"`
	BackgroundColor    string `json:"background_color"`
}

type SystemChartColors map[enum.DisabilityCategoryEnum]string

var DefaultConfig = Config{
	FontSize:       16,
	ScreenReader:   false,
	VoiceCapture:   false,
	Theme:          enum.Light,
	ColorBlindness: enum.Normal,
	SystemColors: SystemColors{
		PrimaryColors: SystemPrimaryColors{
			PrimaryColor:       "#004AAD",
			SecondaryColor:     "#003379",
			FontColor:          "#000000",
			SecondaryFontColor: "#999999",
			InputColor:         "#EEEEEE",
			BackgroundColor:    "#FFFFFF",
		},
		ChartColors: SystemChartColors{
			enum.Motor:        "#003379",
			enum.Hearing:      "#004AAD",
			enum.Intellectual: "#309ACD",
			enum.Psychosocial: "#086F6F",
			enum.Visual:       "#09BCB6",
		},
	},
}
