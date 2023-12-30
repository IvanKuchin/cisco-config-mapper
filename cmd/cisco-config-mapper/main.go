package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ivankuchin/cisco-config-mapper/internal/cli"
	"github.com/ivankuchin/cisco-config-mapper/internal/config"
	"github.com/ivankuchin/cisco-config-mapper/internal/mapper"
	"github.com/ivankuchin/cisco-config-mapper/internal/utils"
)

func SetLogFlags() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}

func main() {
	SetLogFlags()

	cli.Execute()

	config_dir := cli.ConfigDir
	src_dir := cli.SrcDir
	dst_dir := cli.DstDir

	fmt.Println("Reading mapping files from " + config_dir)
	err := config.Read(config_dir)
	if err != nil {
		os.Exit(1)
	}

	fmt.Println("Checking source configs " + src_dir)
	err = config.CheckSrcFolder(src_dir)
	if err != nil {
		os.Exit(1)
	}

	fmt.Println("Converting configs...")
	final_configs, err := mapper.Map(config.Get(), src_dir, dst_dir)
	if err != nil {
		os.Exit(1)
	}

	fmt.Println("Writing configs to " + dst_dir)
	err = utils.Write(final_configs, dst_dir)
	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
