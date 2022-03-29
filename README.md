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

// map[string]interface {}{
//   "font-family":[]string{"Georgia", "serif"},
//   "font-size":"12px",
//   "font-style":"italic",
//   "font-variant":"small-caps",
//   "font-weight":"bold",
//   "line-height":"30px"
// }
```
## License

Licensed under the three-clause BSD license.
