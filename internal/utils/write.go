package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ivankuchin/cisco-config-mapper/internal/cisco"
)

func Write(configs map[string]cisco.Cisco, dst_dir string) error {
	for f, config := range configs {
		err := write_config(config, dst_dir+f)
		if err != nil {
			return err
		}
	}
	return nil
}

func write_config(config cisco.Cisco, dst_file string) error {
	f, err := os.Create(dst_file)
	if err != nil {
		error_meassage := "Error writing to file: " + dst_file
		log.Println(error_meassage)
		return errors.New(error_meassage)
	}

	defer f.Close()

	fmt.Fprintln(f, strings.Join(config.Before_ifaces, "\n"))

	for _, iface := range config.Ifaces {
		fmt.Fprintln(f, "interface "+iface.Name)
		fmt.Fprintln(f, strings.Join(iface.Content, "\n"))
	}

	fmt.Fprintln(f, strings.Join(config.After_ifaces, "\n"))

	return nil
}
