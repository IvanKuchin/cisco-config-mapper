package mapper

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/ivankuchin/cisco-config-mapper/internal/cisco"
	"github.com/ivankuchin/cisco-config-mapper/internal/config"
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

	converted_ifaces, err := map_iface(sh_run.Ifaces, mapping.Iface_mappings)
	if err != nil {
		return final_config, err
	}
	final_config.Ifaces = converted_ifaces

	final_config.Before_ifaces = remove_prefixes(sh_run.Before_ifaces, mapping.Remove_prefixes)
	final_config.After_ifaces = remove_prefixes(sh_run.After_ifaces, mapping.Remove_prefixes)
	final_config.Before_ifaces = remove_lines(final_config.Before_ifaces, mapping.Remove_lines)
	final_config.After_ifaces = remove_lines(final_config.After_ifaces, mapping.Remove_lines)
	final_config.Before_ifaces = replace_iface_names(final_config.Before_ifaces, mapping.Iface_mappings)
	final_config.After_ifaces = replace_iface_names(final_config.After_ifaces, mapping.Iface_mappings)
	final_config.Before_ifaces = prepend_text(final_config.Before_ifaces, mapping.Prepend)
	final_config.After_ifaces = append_text(final_config.After_ifaces, mapping.Append)

	return final_config, nil
}

func map_iface(ifaces []cisco.Iface, iface_mappings []config.Iface_mapping) ([]cisco.Iface, error) {
	var err error

	result := make([]cisco.Iface, 0)

	for _, iface := range ifaces {
		for j, _ := range iface_mappings {

			re := regexp.MustCompile("\\b" + iface_mappings[j].From + "\\b")
			if re.MatchString(iface.Name) {

				// actual iface name replacement / prepending / appending
				iface.Name = re.ReplaceAllString(iface.Name, iface_mappings[j].To)
				if iface_mappings[j].Prepend != "" {
					iface.Content = append([]string{iface_mappings[j].Prepend}, iface.Content...)
				}
				if iface_mappings[j].Append != "" {
					iface.Content = append(iface.Content, iface_mappings[j].Append)
				}

				iface_mappings[j].Visited = true
				result = append(result, iface)
				break
			}
		}
	}

	for _, iface_mapping := range iface_mappings {
		if !iface_mapping.Visited {
			fmt.Println("ERROR: interface " + iface_mapping.From + " not found")
			err = fmt.Errorf("ERROR: required interfaces not found")
		}
	}

	return result, err
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

func remove_prefixes(content []string, prefixes []string) []string {
	var result []string

	for i := 0; i < len(content); i++ {
		line := content[i]
		found := false
		for _, prefix := range prefixes {
			if strings.HasPrefix(line, prefix) {
				found = true
				i = remove_following_block(get_identation(line), i, content)
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
				i = remove_following_block(get_identation(line), i, content)
				break
			}
		}
		if !found {
			result = append(result, line)
		}
	}

	return result
}

func remove_following_block(identation int, i int, content []string) int {
	block_identation := strings.Repeat(" ", identation+1)
	for i < len(content)-1 && strings.HasPrefix(content[i+1], block_identation) {
		i++
	}

	return i
}

func prepend_text(content []string, text string) []string {
	return append([]string{text}, content...)
}

func append_text(content []string, text string) []string {
	return append(content, text)
}

func get_identation(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}
