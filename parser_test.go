package cfp

import (
	"encoding/json"
	"errors"
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

		{`12px serif`, map[string]any{"font-size": "12px", "font-family": []any{"serif"}}},

		{`12px Arial, Verdana, serif`, map[string]any{"font-size": "12px", "font-family": []any{"Arial", "Verdana", "serif"}}},

		{`12px "Times New Roman"`, map[string]any{"font-size": "12px", "font-family": []any{`"Times New Roman"`}}},
		{`12px 'Times New Roman'`, map[string]any{"font-size": "12px", "font-family": []any{`'Times New Roman'`}}},

		{`12px "Times' New Roman"`, map[string]any{"font-size": "12px", "font-family": []any{`"Times' New Roman"`}}},
		{`12px 'Times" New Roman'`, map[string]any{"font-size": "12px", "font-family": []any{`'Times" New Roman'`}}},

		{`12px "Times\" New Roman"`, map[string]any{"font-size": "12px", "font-family": []any{`"Times\" New Roman"`}}},
		{`12px 'Times\' New Roman'`, map[string]any{"font-size": "12px", "font-family": []any{`'Times\' New Roman'`}}},

		{`12px Times New Roman`, map[string]any{"font-size": "12px", "font-family": []any{"Times New Roman"}}},
		{`12px Times New Roman, Comic Sans MS`, map[string]any{"font-size": "12px", "font-family": []any{"Times New Roman", "Comic Sans MS"}}},
		{`12px "Times New Roman", "Comic Sans MS"`, map[string]any{"font-size": "12px", "font-family": []any{`"Times New Roman"`, `"Comic Sans MS"`}}},

		{`12px Red/Black`, nil},
		{`12px "Lucida" Grande`, nil},
		{`12px Ahem!`, nil},
		{`12px Hawaii 5-0`, nil},
		{`12px $42`, nil},

		{`12px Red\/Black`, map[string]any{"font-size": "12px", "font-family": []any{`Red\/Black`}}},
		{`12px Lucida    Grande`, map[string]any{"font-size": "12px", "font-family": []any{"Lucida Grande"}}},
		{`12px Ahem\!`, map[string]any{"font-size": "12px", "font-family": []any{`Ahem\!`}}},
		{`12px \$42`, map[string]any{"font-size": "12px", "font-family": []any{`\$42`}}},
		{`12px €42`, map[string]any{"font-size": "12px", "font-family": []any{`€42`}}},

		{`12px serif`, map[string]any{"font-size": "12px", "font-family": []any{"serif"}}},
		{`xx-small serif`, map[string]any{"font-size": "xx-small", "font-family": []any{"serif"}}},
		{`s-small serif`, map[string]any{"font-size": "s-small", "font-family": []any{"serif"}}},
		{`small serif`, map[string]any{"font-size": "small", "font-family": []any{"serif"}}},
		{`medium serif`, map[string]any{"font-size": "medium", "font-family": []any{"serif"}}},
		{`large serif`, map[string]any{"font-size": "large", "font-family": []any{"serif"}}},
		{`x-large serif`, map[string]any{"font-size": "x-large", "font-family": []any{"serif"}}},
		{`xx-large serif`, map[string]any{"font-size": "xx-large", "font-family": []any{"serif"}}},

		{`larger serif`, map[string]any{"font-size": "larger", "font-family": []any{"serif"}}},
		{`smaller serif`, map[string]any{"font-size": "smaller", "font-family": []any{"serif"}}},

		{`1px serif`, map[string]any{"font-size": "1px", "font-family": []any{"serif"}}},
		{`1em serif`, map[string]any{"font-size": "1em", "font-family": []any{"serif"}}},
		{`1ex serif`, map[string]any{"font-size": "1ex", "font-family": []any{"serif"}}},
		{`1ch serif`, map[string]any{"font-size": "1ch", "font-family": []any{"serif"}}},
		{`1rem serif`, map[string]any{"font-size": "1rem", "font-family": []any{"serif"}}},
		{`1vh serif`, map[string]any{"font-size": "1vh", "font-family": []any{"serif"}}},
		{`1vw serif`, map[string]any{"font-size": "1vw", "font-family": []any{"serif"}}},
		{`1vmin serif`, map[string]any{"font-size": "1vmin", "font-family": []any{"serif"}}},
		{`1vmax serif`, map[string]any{"font-size": "1vmax", "font-family": []any{"serif"}}},
		{`1mm serif`, map[string]any{"font-size": "1mm", "font-family": []any{"serif"}}},
		{`1cm serif`, map[string]any{"font-size": "1cm", "font-family": []any{"serif"}}},
		{`1in serif`, map[string]any{"font-size": "1in", "font-family": []any{"serif"}}},
		{`1pt serif`, map[string]any{"font-size": "1pt", "font-family": []any{"serif"}}},
		{`1pc serif`, map[string]any{"font-size": "1pc", "font-family": []any{"serif"}}},

		{`1 serif`, nil},
		{`xxx-small serif`, nil},
		{`1bs serif`, nil},
		{`100 % serif`, nil},

		{`100% serif`, map[string]any{"font-size": "100%", "font-family": []any{"serif"}}},

		{`1px serif`, map[string]any{"font-size": "1px", "font-family": []any{"serif"}}},
		{`1.1px serif`, map[string]any{"font-size": "1.1px", "font-family": []any{"serif"}}},
		{`-1px serif`, map[string]any{"font-size": "-1px", "font-family": []any{"serif"}}},
		{`-1.1px serif`, map[string]any{"font-size": "-1.1px", "font-family": []any{"serif"}}},
		{`+1px serif`, map[string]any{"font-size": "+1px", "font-family": []any{"serif"}}},
		{`+1.1px serif`, map[string]any{"font-size": "+1.1px", "font-family": []any{"serif"}}},
		{`.1px serif`, map[string]any{"font-size": ".1px", "font-family": []any{"serif"}}},
		{`+.1px serif`, map[string]any{"font-size": "+.1px", "font-family": []any{"serif"}}},
		{`-.1px serif`, map[string]any{"font-size": "-.1px", "font-family": []any{"serif"}}},

		{`12.px serif`, nil},
		{`+---12.2px serif`, nil},
		{`12.1.1px serif`, nil},
		{`10e3px serif`, nil},

		{`12px/16px serif`, map[string]any{"font-size": "12px", "line-height": "16px", "font-family": []any{"serif"}}},
		{`12px/1.5 serif`, map[string]any{"font-size": "12px", "line-height": "1.5", "font-family": []any{"serif"}}},
		{`12px/normal serif`, map[string]any{"font-size": "12px", "font-family": []any{"serif"}}},
		{`12px/105% serif`, map[string]any{"font-size": "12px", "line-height": "105%", "font-family": []any{"serif"}}},

		{`oblique 12px serif`, map[string]any{"font-size": "12px", "font-style": "oblique", "font-family": []any{"serif"}}},
		{`oblique 20deg 12px serif`, map[string]any{"font-size": "12px", "font-style": "oblique 20deg", "font-family": []any{"serif"}}},
		{`oblique 0.02turn 12px serif`, map[string]any{"font-size": "12px", "font-style": "oblique 0.02turn", "font-family": []any{"serif"}}},
		{`oblique .04rad 12px serif`, map[string]any{"font-size": "12px", "font-style": "oblique .04rad", "font-family": []any{"serif"}}},

		{`italic 12px serif`, map[string]any{"font-size": "12px", "font-style": "italic", "font-family": []any{"serif"}}},

		{`small-caps 12px serif`, map[string]any{"font-size": "12px", "font-variant": "small-caps", "font-family": []any{"serif"}}},

		{`bold 12px serif`, map[string]any{"font-size": "12px", "font-weight": "bold", "font-family": []any{"serif"}}},
		{`bolder 12px serif`, map[string]any{"font-size": "12px", "font-weight": "bolder", "font-family": []any{"serif"}}},
		{`lighter 12px serif`, map[string]any{"font-size": "12px", "font-weight": "lighter", "font-family": []any{"serif"}}},

		{`1 12px serif`, map[string]any{"font-size": "12px", "font-weight": "1", "font-family": []any{"serif"}}},
		{`723 12px serif`, map[string]any{"font-size": "12px", "font-weight": "723", "font-family": []any{"serif"}}},
		{`1000 12px serif`, map[string]any{"font-size": "12px", "font-weight": "1000", "font-family": []any{"serif"}}},
		{`1000.00 12px serif`, map[string]any{"font-size": "12px", "font-weight": "1000.00", "font-family": []any{"serif"}}},
		{`1e3 12px serif`, map[string]any{"font-size": "12px", "font-weight": "1e3", "font-family": []any{"serif"}}},
		{`1e+1 12px serif`, map[string]any{"font-size": "12px", "font-weight": "1e+1", "font-family": []any{"serif"}}},
		{`200e-2 12px serif`, map[string]any{"font-size": "12px", "font-weight": "200e-2", "font-family": []any{"serif"}}},
		{`123.456 12px serif`, map[string]any{"font-size": "12px", "font-weight": "123.456", "font-family": []any{"serif"}}},
		{`+123 12px serif`, map[string]any{"font-size": "12px", "font-weight": "+123", "font-family": []any{"serif"}}},

		{`0 12px serif`, map[string]any{"font-size": "12px", "font-family": []any{"serif"}}},
		{`-1 12px serif`, map[string]any{"font-size": "12px", "font-family": []any{"serif"}}},
		{`1000. 12px serif`, map[string]any{"font-size": "12px", "font-family": []any{"serif"}}},
		{`1000.1 12px serif`, map[string]any{"font-size": "12px", "font-family": []any{"serif"}}},
		{`1001 12px serif`, map[string]any{"font-size": "12px", "font-family": []any{"serif"}}},
		{`1.1e3 12px serif`, map[string]any{"font-size": "12px", "font-family": []any{"serif"}}},
		{`1e-2 12px serif`, map[string]any{"font-size": "12px", "font-family": []any{"serif"}}},

		{`ultra-condensed 12px serif`, map[string]any{"font-size": "12px", "font-stretch": "ultra-condensed", "font-family": []any{"serif"}}},
		{`extra-condensed 12px serif`, map[string]any{"font-size": "12px", "font-stretch": "extra-condensed", "font-family": []any{"serif"}}},
		{`condensed 12px serif`, map[string]any{"font-size": "12px", "font-stretch": "condensed", "font-family": []any{"serif"}}},
		{`semi-condensed 12px serif`, map[string]any{"font-size": "12px", "font-stretch": "semi-condensed", "font-family": []any{"serif"}}},
		{`semi-expanded 12px serif`, map[string]any{"font-size": "12px", "font-stretch": "semi-expanded", "font-family": []any{"serif"}}},
		{`expanded 12px serif`, map[string]any{"font-size": "12px", "font-stretch": "expanded", "font-family": []any{"serif"}}},
		{`extra-expanded 12px serif`, map[string]any{"font-size": "12px", "font-stretch": "extra-expanded", "font-family": []any{"serif"}}},
		{`ultra-expanded 12px serif`, map[string]any{"font-size": "12px", "font-stretch": "ultra-expanded", "font-family": []any{"serif"}}},

		{`italic small-caps bold 12px/30px Georgia, serif`, map[string]any{"font-family": []any{"Georgia", "serif"}, "font-size": "12px", "font-style": "italic", "font-variant": "small-caps", "font-weight": "bold", "line-height": "30px"}},

		{`100 12px serif`, map[string]any{"font-size": "12px", "font-weight": "100", "font-family": []any{"serif"}}},
		{`200 12px serif`, map[string]any{"font-size": "12px", "font-weight": "200", "font-family": []any{"serif"}}},
		{`300 12px serif`, map[string]any{"font-size": "12px", "font-weight": "300", "font-family": []any{"serif"}}},
		{`400 12px serif`, map[string]any{"font-size": "12px", "font-weight": "400", "font-family": []any{"serif"}}},
		{`500 12px serif`, map[string]any{"font-size": "12px", "font-weight": "500", "font-family": []any{"serif"}}},
		{`600 12px serif`, map[string]any{"font-size": "12px", "font-weight": "600", "font-family": []any{"serif"}}},
		{`700 12px serif`, map[string]any{"font-size": "12px", "font-weight": "700", "font-family": []any{"serif"}}},
		{`800 12px serif`, map[string]any{"font-size": "12px", "font-weight": "800", "font-family": []any{"serif"}}},
		{`900 12px serif`, map[string]any{"font-size": "12px", "font-weight": "900", "font-family": []any{"serif"}}},
		{`1000 12px serif`, map[string]any{"font-size": "12px", "font-weight": "1000", "font-family": []any{"serif"}}},
	}

	for _, data := range tt {
		t.Run("with"+data.input, func(t *testing.T) {
			font := Parse(data.input)
			if font == nil {
				assert.Equal(t, data.output, map[string]any(nil))
			} else {
				r, err := json.Marshal(font)
				if err != nil {
					fmt.Println(err)
				}

				result := make(map[string]any)
				err = json.Unmarshal(r, &result)
				if err != nil {
					fmt.Println(err)
				}

				assert.Equal(t, data.output, result)
			}
		})
	}
}

func TestFont_Error(t *testing.T) {
	f := &Font{error: errors.New("error")}
	assert.EqualError(t, f.Error(), "error")
}
