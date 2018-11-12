package term

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/admpub/qrencode-go/qrencode"
	colorable "github.com/mattn/go-colorable"
	isatty "github.com/mattn/go-isatty"
	"github.com/mgutz/ansi"
)

func PrintAA(w_in io.Writer, grid *qrencode.BitGrid, inverse bool) {
	// Buffering required for Windows (go-colorable) support
	w := bufio.NewWriterSize(w_in, 1024)

	reset := ansi.ColorCode("reset")
	black := ansi.ColorCode(":black")
	white := ansi.ColorCode(":white")
	if inverse {
		black, white = white, black
	}

	height := grid.Height()
	width := grid.Width()
	line := white + fmt.Sprintf("%*s", width*2+2, "") + reset + "\n"

	fmt.Fprint(w, line)
	for y := 0; y < height; y++ {
		fmt.Fprint(w, white, " ")
		color_prev := white
		for x := 0; x < width; x++ {
			if grid.Get(x, y) {
				if color_prev != black {
					fmt.Fprint(w, black)
					color_prev = black
				}
			} else {
				if color_prev != white {
					fmt.Fprint(w, white)
					color_prev = white
				}
			}
			fmt.Fprint(w, "  ")
		}
		fmt.Fprint(w, white, " ", reset, "\n")
		w.Flush()
	}
	fmt.Fprint(w, line)
	w.Flush()
}

func PrintSixel(w io.Writer, grid *qrencode.BitGrid, inverse bool) {
	black := "0"
	white := "1"

	fmt.Fprint(w,
		"\x1BPq\"1;1",
		"#", black, ";2;0;0;0",
		"#", white, ";2;100;100;100",
	)

	if inverse {
		black, white = white, black
	}

	height := grid.Height()
	width := grid.Width()
	line := "#" + white + "!" + fmt.Sprintf("%d", (width+2)*6) + "~"

	fmt.Fprint(w, line, "-")
	for y := 0; y < height; y++ {
		fmt.Fprint(w, "#", white)
		color := white
		repeat := 6
		var current string
		for x := 0; x < width; x++ {
			if grid.Get(x, y) {
				current = black
			} else {
				current = white
			}
			if current != color {
				fmt.Fprint(w, "#", color, "!", repeat, "~")
				color = current
				repeat = 0
			}
			repeat += 6
		}
		if color == white {
			fmt.Fprintf(w, "#%s!%d~", white, repeat+6)
		} else {
			fmt.Fprintf(w, "#%s!%d~#%s!6~", color, repeat, white)
		}
		fmt.Fprint(w, "-")
	}
	fmt.Fprint(w, line)
	fmt.Fprint(w, "\x1B\\")
}

func Print(text string, level qrencode.ECLevel, inverse ...bool) error {
	var _inverse bool
	if len(inverse) > 0 {
		_inverse = inverse[0]
	}
	grid, err := qrencode.Encode(text, level)
	if err != nil {
		return err
	}

	if !isatty.IsTerminal(os.Stdout.Fd()) {
		PrintSixel(os.Stdout, grid, _inverse)
	} else {
		stdout := colorable.NewColorableStdout()
		PrintAA(stdout, grid, _inverse)
	}
	return nil
}
