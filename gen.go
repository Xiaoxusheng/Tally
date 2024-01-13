package main

/*
此时使用自动化脚本
*/
//go:generate firewall-cmd --list-all
//go:generate rm -rf tally
//go:generate go mod tidy
//go:generate go build -o tally -ldflags "-s -w"
