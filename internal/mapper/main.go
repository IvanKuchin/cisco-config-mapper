package mapper

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ivankuchin/cisco-config-mapper/internal/cisco"
	"github.com/ivankuchin/cisco-config-mapper/internal/config"
)

func Map(mappings map[string]config.Config, src_dir string, dst_dir string) (map[string]cisco.Cisco, error) {

	var configs = make(map[string]cisco.Cisco)
	var final_configs = make(map[string]cisco.Cisco)

	for f, _ := range mappings {
		sh_run, err := cisco.New(src_dir + f)
		if err != nil {
			return final_configs, err
		}

		configs[f] = sh_run
	}

	for f, _ := range mappings {
		fmt.Print("\t" + f)
		final_config, err := convert_config(configs[f], mappings[f])
		if err != nil {
			return final_configs, err
		}
		final_configs[f] = final_config
		fmt.Println("\tsuccess")
	}

	return final_configs, nil
}

func convert_config(sh_run cisco.Cisco, mapping config.Config) (cisco.Cisco, error) {
	var final_config cisco.Cisco

	// remove "end" to avoid adding it to the final config
	mapping.Remove_lines = append(mapping.Remove_lines, "end")

	final_config.Ifaces = map_iface(sh_run.Ifaces, mapping.Iface_mappings)
	final_config.Before_ifaces = remove_prefixes(sh_run.Before_ifaces, mapping.Remove_prefixes)
	final_config.After_ifaces = remove_prefixes(sh_run.After_ifaces, mapping.Remove_prefixes)
	final_config.Before_ifaces = remove_lines(final_config.Before_ifaces, mapping.Remove_lines)
	final_config.After_ifaces = remove_lines(final_config.After_ifaces, mapping.Remove_lines)
	final_config.After_ifaces = append_text(final_config.After_ifaces, mapping.Append)

	return final_config, nil
}

func map_iface(ifaces []cisco.Iface, mappings []config.Iface_mapping) []cisco.Iface {
	result := make([]cisco.Iface, 0)

	for _, iface := range ifaces {
		for _, mapping := range mappings {

			re := regexp.MustCompile(mapping.From + "\\b")
			if re.MatchString("interface " + iface.Name) {
				iface.Name = re.ReplaceAllString("interface "+iface.Name, mapping.To)
				result = append(result, iface)
				break
			}
		}
	}

	return result
}

func remove_prefixes(content []string, prefixes []string) []string {
	var result []string

	for i := 0; i < len(content); i++ {
		line := content[i]
		found := false
		for _, prefix := range prefixes {
			if strings.HasPrefix(line, prefix) {
				found = true

				// remove block after line if idented
				curr_identation := get_identation(line)
				block_identation := strings.Repeat(" ", curr_identation+1)
				for i < len(content)-1 && strings.HasPrefix(content[i+1], block_identation) {
					i++
				}

				break
			}
		}
		if !found {
			result = append(result, line)
		}
	}

	return result
}

func remove_lines(content []string, prefixes []string) []string {
	var result []string

	for i := 0; i < len(content); i++ {
		line := content[i]
		found := false
		for _, prefix := range prefixes {
			if line == prefix {
				found = true

				// remove block after line if idented
				curr_identation := get_identation(line)
				block_identation := strings.Repeat(" ", curr_identation+1)
				for i < len(content)-1 && strings.HasPrefix(content[i+1], block_identation) {
					i++
				}

				break
			}
		}
		if !found {
			result = append(result, line)
		}
	}

	return result
}

func append_text(content []string, text string) []string {
	return append(content, text)
}

func get_identation(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}
