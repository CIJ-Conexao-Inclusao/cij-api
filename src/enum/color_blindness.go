package enum

type ColorBlindnessEnum string

const (
	Normal       ColorBlindnessEnum = "normal"
	Protanopia   ColorBlindnessEnum = "protanopia"
	Deuteranopia ColorBlindnessEnum = "deuteranopia"
	Tritanopia   ColorBlindnessEnum = "tritanopia"
)

func (c ColorBlindnessEnum) IsValid() bool {
	return c == Normal || c == Protanopia || c == Deuteranopia || c == Tritanopia
}
