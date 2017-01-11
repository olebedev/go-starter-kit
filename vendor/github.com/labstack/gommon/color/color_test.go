package color

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestText(t *testing.T) {
	fmt.Println("*** colored text ***")
	fmt.Println(Black("black"))
	fmt.Println(Red("red"))
	fmt.Println(Green("green"))
	fmt.Println(Yellow("yellow"))
	fmt.Println(Blue("blue"))
	fmt.Println(Magenta("magenta"))
	fmt.Println(Cyan("cyan"))
	fmt.Println(White("white"))
	fmt.Println(Grey("grey"))
}

func TestBackground(t *testing.T) {
	fmt.Println("*** colored background ***")
	fmt.Println(BlackBg("black background", Wht))
	fmt.Println(RedBg("red background"))
	fmt.Println(GreenBg("green background"))
	fmt.Println(YellowBg("yellow background"))
	fmt.Println(BlueBg("blue background"))
	fmt.Println(MagentaBg("magenta background"))
	fmt.Println(CyanBg("cyan background"))
	fmt.Println(WhiteBg("white background"))
}

func TestEmphasis(t *testing.T) {
	fmt.Println("*** emphasis ***")
	fmt.Println(Reset("reset"))
	fmt.Println(Bold("bold"))
	fmt.Println(Dim("dim"))
	fmt.Println(Italic("italic"))
	fmt.Println(Underline("underline"))
	fmt.Println(Inverse("inverse"))
	fmt.Println(Hidden("hidden"))
	fmt.Println(Strikeout("strikeout"))
}

func TestMixMatch(t *testing.T) {
	fmt.Println("*** mix and match ***")
	fmt.Println(Green("bold green with white background", B, WhtBg))
	fmt.Println(Red("underline red", U))
	fmt.Println(Yellow("dim yellow", D))
	fmt.Println(Cyan("inverse cyan", In))
	fmt.Println(Blue("bold underline dim blue", B, U, D))
}

func TestEnableDisable(t *testing.T) {
	Disable()
	assert.Equal(t, "red", Red("red"))
	Enable()
	assert.NotEqual(t, "green", Green("green"))
}
