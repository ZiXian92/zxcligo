package zxcligo

import (
	"bytes"
	"fmt"
	"log"
)

func printHelp(ctx *Context) {
	var buf bytes.Buffer

	// Fill up usage format syntax
	buf.WriteString("Usage: ")
	fillCmdUsage(&buf, ctx)
	buf.WriteString(" [arguments]")
	fmt.Println(buf.String())

	printOptions(ctx)
}

func fillCmdUsage(buf *bytes.Buffer, ctx *Context) {
	if ctx == nil {
		return
	}

	fillCmdUsage(buf, ctx.Parent)
	if ctx.Parent != nil {
		buf.WriteString(" ")
	}

	buf.WriteString(ctx.Command)
	buf.WriteString(" [")
	buf.WriteString(ctx.Command)
	buf.WriteString("_options] ")
}

func printOptions(ctx *Context) {
	if ctx == nil {
		return
	}

	printOptions(ctx.Parent)

	fmt.Printf("%s options:\n", ctx.Command)
	log.Printf("%v", ctx.flags)
	ctx.flags.PrintDefaults()
}
