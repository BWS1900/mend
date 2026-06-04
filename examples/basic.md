# mend

A small example showing common markdown features.

## Inline

Text with **bold**, *italic*, ~~strikethrough~~, `inline code`, and a [link](https://example.com).

## Lists

- one
- two
- three

1. first
2. second
3. third

## Code

```go
package main

import "fmt"

func main() {
    fmt.Println("hello, mend")
}
```

```python
def fib(n):
    a, b = 0, 1
    for _ in range(n):
        yield a
        a, b = b, a + b
```

## Table

| feature | supported |
| --- | --- |
| GFM tables | yes |
| footnotes | yes |
| code highlight | yes |

## Quote

> mend turns markdown into a single HTML file with
> embedded CSS. No JavaScript required.

---

_eof_
