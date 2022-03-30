## CSS font value parser

Go port of [CSS font value parser](https://github.com/bramstein/css-font-parser)

A simple parser for parsing CSS font values.

```go
import (
	"fmt"

	"github.com/raykov/css-font-parser"
)

data := cfp.Parse(`italic small-caps bold 12px/30px Georgia, serif`)

fmt.Printf("%#v", data)

/*
&cfp.Font{
	Family:[]string{"Georgia", "serif"},
	Size:"12px",
	Style:"italic",
	Variant:"small-caps",
	Weight:"bold",
	Stretch:"",
	LineHeight:"30px",
}
*/
```
## License

Licensed under the three-clause BSD license.
