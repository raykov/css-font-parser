package cfp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tt := []struct {
		input  string
		output map[string]any
	}{
		{`12px "Comic`, nil},
		{`12px "Comic, serif`, nil},
		{`12px "Comic, \"serif`, nil},
		{`12px 'Comic`, nil},
		{`12px 'Comic, serif`, nil},

		{`12px serif`, map[string]any{"font-size": "12px", "font-family": []string{"serif"}}},

		{`12px Arial, Verdana, serif`, map[string]any{"font-size": "12px", "font-family": []string{"Arial", "Verdana", "serif"}}},

		{`12px "Times New Roman"`, map[string]any{"font-size": "12px", "font-family": []string{`"Times New Roman"`}}},
		{`12px 'Times New Roman'`, map[string]any{"font-size": "12px", "font-family": []string{`'Times New Roman'`}}},

		{`12px "Times' New Roman"`, map[string]any{"font-size": "12px", "font-family": []string{`"Times' New Roman"`}}},
		{`12px 'Times" New Roman'`, map[string]any{"font-size": "12px", "font-family": []string{`'Times" New Roman'`}}},

		{`12px "Times\" New Roman"`, map[string]any{"font-size": "12px", "font-family": []string{`"Times\" New Roman"`}}},
		{`12px 'Times\' New Roman'`, map[string]any{"font-size": "12px", "font-family": []string{`'Times\' New Roman'`}}},

		{`12px Times New Roman`, map[string]any{"font-size": "12px", "font-family": []string{"Times New Roman"}}},
		{`12px Times New Roman, Comic Sans MS`, map[string]any{"font-size": "12px", "font-family": []string{"Times New Roman", "Comic Sans MS"}}},
		{`12px "Times New Roman", "Comic Sans MS"`, map[string]any{"font-size": "12px", "font-family": []string{`"Times New Roman"`, `"Comic Sans MS"`}}},

		{`12px Red/Black`, nil},
		{`12px "Lucida" Grande`, nil},
		{`12px Ahem!`, nil},
		{`12px Hawaii 5-0`, nil},
		{`12px $42`, nil},

		{`12px Red\/Black`, map[string]any{"font-size": "12px", "font-family": []string{`Red\/Black`}}},
		{`12px Lucida    Grande`, map[string]any{"font-size": "12px", "font-family": []string{"Lucida Grande"}}},
		{`12px Ahem\!`, map[string]any{"font-size": "12px", "font-family": []string{`Ahem\!`}}},
		{`12px \$42`, map[string]any{"font-size": "12px", "font-family": []string{`\$42`}}},
		{`12px €42`, map[string]any{"font-size": "12px", "font-family": []string{`€42`}}},

		{`12px serif`, map[string]any{"font-size": "12px", "font-family": []string{"serif"}}},
		{`xx-small serif`, map[string]any{"font-size": "xx-small", "font-family": []string{"serif"}}},
		{`s-small serif`, map[string]any{"font-size": "s-small", "font-family": []string{"serif"}}},
		{`small serif`, map[string]any{"font-size": "small", "font-family": []string{"serif"}}},
		{`medium serif`, map[string]any{"font-size": "medium", "font-family": []string{"serif"}}},
		{`large serif`, map[string]any{"font-size": "large", "font-family": []string{"serif"}}},
		{`x-large serif`, map[string]any{"font-size": "x-large", "font-family": []string{"serif"}}},
		{`xx-large serif`, map[string]any{"font-size": "xx-large", "font-family": []string{"serif"}}},

		{`larger serif`, map[string]any{"font-size": "larger", "font-family": []string{"serif"}}},
		{`smaller serif`, map[string]any{"font-size": "smaller", "font-family": []string{"serif"}}},

		{`1px serif`, map[string]any{"font-size": "1px", "font-family": []string{"serif"}}},
		{`1em serif`, map[string]any{"font-size": "1em", "font-family": []string{"serif"}}},
		{`1ex serif`, map[string]any{"font-size": "1ex", "font-family": []string{"serif"}}},
		{`1ch serif`, map[string]any{"font-size": "1ch", "font-family": []string{"serif"}}},
		{`1rem serif`, map[string]any{"font-size": "1rem", "font-family": []string{"serif"}}},
		{`1vh serif`, map[string]any{"font-size": "1vh", "font-family": []string{"serif"}}},
		{`1vw serif`, map[string]any{"font-size": "1vw", "font-family": []string{"serif"}}},
		{`1vmin serif`, map[string]any{"font-size": "1vmin", "font-family": []string{"serif"}}},
		{`1vmax serif`, map[string]any{"font-size": "1vmax", "font-family": []string{"serif"}}},
		{`1mm serif`, map[string]any{"font-size": "1mm", "font-family": []string{"serif"}}},
		{`1cm serif`, map[string]any{"font-size": "1cm", "font-family": []string{"serif"}}},
		{`1in serif`, map[string]any{"font-size": "1in", "font-family": []string{"serif"}}},
		{`1pt serif`, map[string]any{"font-size": "1pt", "font-family": []string{"serif"}}},
		{`1pc serif`, map[string]any{"font-size": "1pc", "font-family": []string{"serif"}}},

		{`1 serif`, nil},
		{`xxx-small serif`, nil},
		{`1bs serif`, nil},
		{`100 % serif`, nil},

		{`100% serif`, map[string]any{"font-size": "100%", "font-family": []string{"serif"}}},

		{`1px serif`, map[string]any{"font-size": "1px", "font-family": []string{"serif"}}},
		{`1.1px serif`, map[string]any{"font-size": "1.1px", "font-family": []string{"serif"}}},
		{`-1px serif`, map[string]any{"font-size": "-1px", "font-family": []string{"serif"}}},
		{`-1.1px serif`, map[string]any{"font-size": "-1.1px", "font-family": []string{"serif"}}},
		{`+1px serif`, map[string]any{"font-size": "+1px", "font-family": []string{"serif"}}},
		{`+1.1px serif`, map[string]any{"font-size": "+1.1px", "font-family": []string{"serif"}}},
		{`.1px serif`, map[string]any{"font-size": ".1px", "font-family": []string{"serif"}}},
		{`+.1px serif`, map[string]any{"font-size": "+.1px", "font-family": []string{"serif"}}},
		{`-.1px serif`, map[string]any{"font-size": "-.1px", "font-family": []string{"serif"}}},

		{`12.px serif`, nil},
		{`+---12.2px serif`, nil},
		{`12.1.1px serif`, nil},
		{`10e3px serif`, nil},

		{`12px/16px serif`, map[string]any{"font-size": "12px", "line-height": "16px", "font-family": []string{"serif"}}},
		{`12px/1.5 serif`, map[string]any{"font-size": "12px", "line-height": "1.5", "font-family": []string{"serif"}}},
		{`12px/normal serif`, map[string]any{"font-size": "12px", "font-family": []string{"serif"}}},
		{`12px/105% serif`, map[string]any{"font-size": "12px", "line-height": "105%", "font-family": []string{"serif"}}},

		{`oblique 12px serif`, map[string]any{"font-size": "12px", "font-style": "oblique", "font-family": []string{"serif"}}},
		{`oblique 20deg 12px serif`, map[string]any{"font-size": "12px", "font-style": "oblique 20deg", "font-family": []string{"serif"}}},
		{`oblique 0.02turn 12px serif`, map[string]any{"font-size": "12px", "font-style": "oblique 0.02turn", "font-family": []string{"serif"}}},
		{`oblique .04rad 12px serif`, map[string]any{"font-size": "12px", "font-style": "oblique .04rad", "font-family": []string{"serif"}}},

		{`italic 12px serif`, map[string]any{"font-size": "12px", "font-style": "italic", "font-family": []string{"serif"}}},

		{`small-caps 12px serif`, map[string]any{"font-size": "12px", "font-variant": "small-caps", "font-family": []string{"serif"}}},

		{`bold 12px serif`, map[string]any{"font-size": "12px", "font-weight": "bold", "font-family": []string{"serif"}}},
		{`bolder 12px serif`, map[string]any{"font-size": "12px", "font-weight": "bolder", "font-family": []string{"serif"}}},
		{`lighter 12px serif`, map[string]any{"font-size": "12px", "font-weight": "lighter", "font-family": []string{"serif"}}},

		{`1 12px serif`, map[string]any{"font-size": "12px", "font-weight": "1", "font-family": []string{"serif"}}},
		{`723 12px serif`, map[string]any{"font-size": "12px", "font-weight": "723", "font-family": []string{"serif"}}},
		{`1000 12px serif`, map[string]any{"font-size": "12px", "font-weight": "1000", "font-family": []string{"serif"}}},
		{`1000.00 12px serif`, map[string]any{"font-size": "12px", "font-weight": "1000.00", "font-family": []string{"serif"}}},
		{`1e3 12px serif`, map[string]any{"font-size": "12px", "font-weight": "1e3", "font-family": []string{"serif"}}},
		{`1e+1 12px serif`, map[string]any{"font-size": "12px", "font-weight": "1e+1", "font-family": []string{"serif"}}},
		{`200e-2 12px serif`, map[string]any{"font-size": "12px", "font-weight": "200e-2", "font-family": []string{"serif"}}},
		{`123.456 12px serif`, map[string]any{"font-size": "12px", "font-weight": "123.456", "font-family": []string{"serif"}}},
		{`+123 12px serif`, map[string]any{"font-size": "12px", "font-weight": "+123", "font-family": []string{"serif"}}},

		{`0 12px serif`, map[string]any{"font-size": "12px", "font-family": []string{"serif"}}},
		{`-1 12px serif`, map[string]any{"font-size": "12px", "font-family": []string{"serif"}}},
		{`1000. 12px serif`, map[string]any{"font-size": "12px", "font-family": []string{"serif"}}},
		{`1000.1 12px serif`, map[string]any{"font-size": "12px", "font-family": []string{"serif"}}},
		{`1001 12px serif`, map[string]any{"font-size": "12px", "font-family": []string{"serif"}}},
		{`1.1e3 12px serif`, map[string]any{"font-size": "12px", "font-family": []string{"serif"}}},
		{`1e-2 12px serif`, map[string]any{"font-size": "12px", "font-family": []string{"serif"}}},

		{`ultra-condensed 12px serif`, map[string]any{"font-size": "12px", "font-stretch": "ultra-condensed", "font-family": []string{"serif"}}},
		{`extra-condensed 12px serif`, map[string]any{"font-size": "12px", "font-stretch": "extra-condensed", "font-family": []string{"serif"}}},
		{`condensed 12px serif`, map[string]any{"font-size": "12px", "font-stretch": "condensed", "font-family": []string{"serif"}}},
		{`semi-condensed 12px serif`, map[string]any{"font-size": "12px", "font-stretch": "semi-condensed", "font-family": []string{"serif"}}},
		{`semi-expanded 12px serif`, map[string]any{"font-size": "12px", "font-stretch": "semi-expanded", "font-family": []string{"serif"}}},
		{`expanded 12px serif`, map[string]any{"font-size": "12px", "font-stretch": "expanded", "font-family": []string{"serif"}}},
		{`extra-expanded 12px serif`, map[string]any{"font-size": "12px", "font-stretch": "extra-expanded", "font-family": []string{"serif"}}},
		{`ultra-expanded 12px serif`, map[string]any{"font-size": "12px", "font-stretch": "ultra-expanded", "font-family": []string{"serif"}}},

		{`italic small-caps bold 12px/30px Georgia, serif`, map[string]any{"font-family": []string{"Georgia", "serif"}, "font-size": "12px", "font-style": "italic", "font-variant": "small-caps", "font-weight": "bold", "line-height": "30px"}},
	}

	for _, data := range tt {
		t.Run("with"+data.input, func(t *testing.T) {
			assert.Equal(t, data.output, Parse(data.input))
		})
	}

	for i := 1; i <= 10; i++ {
		input := fmt.Sprintf("%d 12px serif", i*100)
		output := map[string]any{
			"font-size":   "12px",
			"font-weight": fmt.Sprint(i * 100),
			"font-family": []string{"serif"},
		}
		t.Run("with"+input, func(t *testing.T) {
			assert.Equal(t, output, Parse(input))
		})
	}
}

func TestParse2(t *testing.T) {
	data := Parse(`italic small-caps bold 12px/30px Georgia, serif`)
	fmt.Printf("%#v", data)
	assert.True(t, false)
}
