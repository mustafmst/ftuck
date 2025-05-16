package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func main() {
	var flag1 int
	flag.IntVar(&flag1, "flag1", 0, "flag1 of command")
	var user string
	flag.StringVar(&user, "user", "guest", "user")
	fmt.Print("Welcome to FTUCK. Let me TUCK Your files XD\n")
	slog.Info("arguments before parsing", "args", strings.Join(os.Args[1:], " "), "flag1", flag1, "fla2", user)
	flag.CommandLine.Parse(os.Args[2:])
	slog.Info("parsed", "flag1", flag1, "flag2", user)
}
