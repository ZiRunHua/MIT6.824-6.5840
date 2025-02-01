package indented

import (
	"bufio"
	"os"
	"strings"
)

var (
	Stdout = NewIndentedWriter("  ")
)

type Writer struct {
	writer *bufio.Writer
	indent string
}

func NewIndentedWriter(indent string) *Writer {
	return &Writer{
		writer: bufio.NewWriter(os.Stdout),
		indent: indent,
	}
}

func (iw *Writer) Write(p []byte) (int, error) {
	// 将输入分割成多行并逐行添加缩进
	lines := strings.Split(string(p), "\n")
	for i, line := range lines {
		if len(line) > 0 {
			_, _ = iw.writer.WriteString(iw.indent)
		}
		_, _ = iw.writer.WriteString(line)
		if i < len(lines)-1 {
			_, _ = iw.writer.WriteString("\n")
		}
	}
	return len(p), iw.writer.Flush()
}
