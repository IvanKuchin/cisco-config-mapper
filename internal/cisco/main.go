package cisco

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

func New(fname string) (Cisco, error) {

	content, err := os.ReadFile(fname)
	if err != nil {
		error_message := "Error reading device config file: " + fname
		log.Println(error_message)
		return Cisco{}, errors.New(error_message)
	}

	sh_run := Cisco{}

	var before_flag = true
	var iface_flag = false
	var iface_config []string
	var after_flag = false

	content_cleaned := strings.ReplaceAll(string(content), "\r", "")

	for _, line := range strings.Split(content_cleaned, "\n") {

		switch {
		case before_flag && !strings.HasPrefix(line, "interface"):
			sh_run.Before_ifaces = append(sh_run.Before_ifaces, line)
		case (before_flag || after_flag || iface_flag) && strings.HasPrefix(line, "interface"):
			{
				before_flag = false
				after_flag = false
				iface_flag = true
				iface_config = append(iface_config, line)
			}
		case iface_flag && strings.HasPrefix(line, " "):
			iface_config = append(iface_config, line)
		case after_flag || (iface_flag && !strings.HasPrefix(line, " ")):
			{
				iface_flag = false
				after_flag = true
				sh_run.After_ifaces = append(sh_run.After_ifaces, line)
			}
		default:
			error_message := fmt.Sprintf("Error parsing device config file %s at line %s", fname, line)
			log.Println(error_message)
			return Cisco{}, errors.New(error_message)
		}
	}

	ifaces, err := parse_ifaces(iface_config)
	if err != nil {
		return Cisco{}, err
	}

	sh_run.Ifaces = ifaces

	return sh_run, nil
}

func parse_ifaces(iface_config []string) ([]Iface, error) {

	ifaces := []Iface{}

	for _, line := range iface_config {

		if strings.HasPrefix(line, "interface") {

			iface_name := strings.TrimPrefix(line, "interface ")
			iface := Iface{Name: iface_name}
			ifaces = append(ifaces, iface)

		} else if len(ifaces) > 0 {

			iface := &ifaces[len(ifaces)-1]
			iface.Content = append(iface.Content, line)

		} else {

			error_message := fmt.Sprintf("Error parsing interface config: %s", line)
			log.Println(error_message)
			return []Iface{}, errors.New(error_message)

		}
	}

	return ifaces, nil
}
