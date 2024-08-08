package enum

type ThemeEnum string

const (
	Light  ThemeEnum = "light"
	Dark   ThemeEnum = "dark"
	System ThemeEnum = "system"
)

func (t ThemeEnum) IsValid() bool {
	return t == Light || t == Dark || t == System
}
