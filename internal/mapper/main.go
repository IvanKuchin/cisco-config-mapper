package mapper

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/ivankuchin/cisco-config-mapper/internal/cisco"
	"github.com/ivankuchin/cisco-config-mapper/internal/config"
	textmanipulations "github.com/ivankuchin/cisco-config-mapper/internal/text-manipulations"
)

func Map(mappings map[string]config.Config, src_dir string, dst_dir string) (map[string]cisco.Cisco, error) {

	var configs = make(map[string]cisco.Cisco)
	var final_configs = make(map[string]cisco.Cisco)
	var sorted_mappings = make([]string, 0)

	for k, _ := range mappings {
		sorted_mappings = append(sorted_mappings, k)
	}

	sort.Strings(sorted_mappings)

	for _, f := range sorted_mappings {
		sh_run, err := cisco.New(src_dir + f)
		if err != nil {
			return final_configs, err
		}

		configs[f] = sh_run
	}

	for _, f := range sorted_mappings {
		fmt.Printf("conversion of %v started\n", f)
		final_config, err := convert_config(configs[f], mappings[f])
		if err != nil {
			fmt.Printf("conversion of %v failed\n", f)
			return final_configs, err
		}
		final_configs[f] = final_config
		fmt.Printf("conversion of %v successful\n", f)
	}

	return final_configs, nil
}

func convert_config(sh_run cisco.Cisco, mapping config.Config) (cisco.Cisco, error) {
	var final_config cisco.Cisco

	// remove "end" to avoid adding it to the final config
	mapping.Remove_lines = append(mapping.Remove_lines, "end")

	converted_ifaces, err := map_ifaces(sh_run.Ifaces, mapping.Iface_mappings)
	if err != nil {
		return final_config, err
	}
	final_config.Ifaces = converted_ifaces

	final_config.Before_ifaces = textmanipulations.Remove_prefixes(sh_run.Before_ifaces, mapping.Remove_prefixes)
	final_config.After_ifaces = textmanipulations.Remove_prefixes(sh_run.After_ifaces, mapping.Remove_prefixes)
	final_config.Before_ifaces = textmanipulations.Remove_lines(final_config.Before_ifaces, mapping.Remove_lines)
	final_config.After_ifaces = textmanipulations.Remove_lines(final_config.After_ifaces, mapping.Remove_lines)
	final_config.Before_ifaces = replace_iface_names(final_config.Before_ifaces, mapping.Iface_mappings)
	final_config.After_ifaces = replace_iface_names(final_config.After_ifaces, mapping.Iface_mappings)
	final_config.Before_ifaces = prepend_text(final_config.Before_ifaces, mapping.Prepend)
	final_config.After_ifaces = append_text(final_config.After_ifaces, mapping.Append)

	return final_config, nil
}

func replace_iface_names_one_line(src string, ifaces []config.Iface_mapping) string {
	for _, iface_mappings := range ifaces {
		re := regexp.MustCompile("\\b" + iface_mappings.From + "\\b")
		if re.MatchString(src) {
			return re.ReplaceAllString(src, iface_mappings.To)
		}
	}

	return src
}

func replace_iface_names(content []string, ifaces []config.Iface_mapping) []string {
	var result []string

	for _, line := range content {
		result = append(result, replace_iface_names_one_line(line, ifaces))
	}

	return result
}

func prepend_text(content []string, text string) []string {
	return append([]string{text}, content...)
}

func append_text(content []string, text string) []string {
	return append(content, text)
}
