Output "docs/demo.gif"

Set FontSize 14
Set Width 900
Set Height 500
Set Theme "Catppuccin Mocha"

Type "# a small markdown file" Enter
Sleep 500ms
Type `cat <<'EOF' > /tmp/hello.md` Enter
Type "# hello from mend" Enter
Type "" Enter
Type "A tiny example with **bold**, *italic*, and `code`." Enter
Type "" Enter
Type "- one" Enter
Type "- two" Enter
Type "- three" Enter
Type "" Enter
Type '```go' Enter
Type 'package main' Enter
Type 'func main() { println("hi") }' Enter
Type '```' Enter
Type "EOF" Enter
Sleep 300ms
Type `mend -theme monokai -o /tmp/hello.html /tmp/hello.md` Enter
Sleep 1s
Type `ls -la /tmp/hello.html` Enter
Sleep 1s
Type `xdg-open /tmp/hello.html 2>/dev/null || open /tmp/hello.html` Enter
Sleep 3s
