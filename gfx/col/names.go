package col

// colours! hardcoded for your pleasure.
const (
	NONE      Colour = 0x00000000
	WHITE     Colour = 0xFFFFFFFF
	BLACK     Colour = 0xFF000000
	RED       Colour = 0xFFFF0000
	BLUE      Colour = 0xFF0000FF
	LIME      Colour = 0xFF00FF00
	LIGHTGREY Colour = 0xFFCCCCCC
	GREY      Colour = 0xFF888888
	DARKGREY  Colour = 0xFF444444
	YELLOW    Colour = 0xFFFFFF00
	FUSCHIA   Colour = 0xFFFF00FF
	CYAN      Colour = 0xFF00FFFF
	MAROON    Colour = 0xFF800000
	OLIVE     Colour = 0xFF808000
	GREEN     Colour = 0xFF008000
	TEAL      Colour = 0xFF008080
	NAVY      Colour = 0xFF000080
	PURPLE    Colour = 0xFF800080
	ORANGE    Colour = 0xFFFFA500
)

var ColourNames = map[Colour]string{
	NONE:      "None",
	WHITE:     "White",
	BLACK:     "Black",
	RED:       "Red",
	BLUE:      "Blue",
	LIME:      "Lime",
	LIGHTGREY: "Light Grey",
	GREY:      "Grey",
	DARKGREY:  "Dark Grey",
	YELLOW:    "Yellow",
	FUSCHIA:   "Fuschia",
	CYAN:      "Cyan",
	MAROON:    "Maroon",
	OLIVE:     "Olive",
	GREEN:     "Green",
	TEAL:      "Teal",
	NAVY:      "Navy",
	PURPLE:    "Purple",
	ORANGE:    "Orange",
}
