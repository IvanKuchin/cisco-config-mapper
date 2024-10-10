package textmanipulations

import (
	"reflect"
	"testing"
)

func TestRemove_prefixes(t *testing.T) {
	tests := []struct {
		name     string
		content  []string
		prefixes []string
		expected []string
	}{
		{
			name: "No prefixes",
			content: []string{
				"line1",
				"line2",
				"line3",
			},
			prefixes: []string{},
			expected: []string{
				"line1",
				"line2",
				"line3",
			},
		},
		{
			name: "Single prefix",
			content: []string{
				"hsrp 500",
				" ip add 5.5.5.5/24",
				" timers 1 3",
				" preempt",
				"line2",
				"hsrp 600",
				" ip add 6.6.6.6/24",
				" timers 1 3",
				" preempt",
			},
			prefixes: []string{"hsrp"},
			expected: []string{
				"line2",
			},
		},
		{
			name: "Multiple prefixes",
			content: []string{
				"hsrp 500",
				" ip add 5.5.5.5/24",
				" timers 1 3",
				" preempt",
				"prefix2_line2",
				"line3",
				"aaa authentication login default local",
				"aaa authorization exec default local",
				"aaa accounting exec default start-stop group radius",
			},
			prefixes: []string{"aaa", "prefix2_", "hsrp"},
			expected: []string{
				"line3",
			},
		},
		{
			name: "No matching prefixes",
			content: []string{
				"line1",
				"line2",
				"line3",
			},
			prefixes: []string{"prefix_"},
			expected: []string{
				"line1",
				"line2",
				"line3",
			},
		},
		{
			name:     "Empty content",
			content:  []string{},
			prefixes: []string{"prefix_"},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Remove_prefixes(tt.content, tt.prefixes)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Remove_prefixes() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRemove_lines(t *testing.T) {
	tests := []struct {
		name     string
		content  []string
		lines    []string
		expected []string
	}{
		{
			name: "No lines",
			content: []string{
				"line1",
				"line2",
				"line3",
			},
			lines: []string{},
			expected: []string{
				"line1",
				"line2",
				"line3",
			},
		},
		{
			name: "Single line",
			content: []string{
				"hsrp 500",
				" ip add 5.5.5.5/24",
				" timers 1 3",
				" preempt",
				"line2",
				"hsrp 600",
				" ip add 6.6.6.6/24",
				" timers 1 3",
				" preempt",
			},
			lines: []string{"hsrp 500"},
			expected: []string{
				"line2",
				"hsrp 600",
				" ip add 6.6.6.6/24",
				" timers 1 3",
				" preempt",
			},
		},
		{
			name: "Multiple lines",
			content: []string{
				"hsrp 500",
				" ip add 5.5.5.5/24",
				" timers 1 3",
				" preempt",
				"prefix2_line2",
				"line3",
				"aaa authentication login default local",
				"aaa authorization exec default local",
				"aaa accounting exec default start-stop group radius",
			},
			lines: []string{"prefix2_", "hsrp 500", "aaa accounting exec default start-stop group radius"},
			expected: []string{
				"prefix2_line2",
				"line3",
				"aaa authentication login default local",
				"aaa authorization exec default local",
			},
		},
		{
			name: "No matching lines",
			content: []string{
				"line1",
				"line2",
				"line3",
			},
			lines: []string{"prefix_"},
			expected: []string{
				"line1",
				"line2",
				"line3",
			},
		},
		{
			name:     "Empty content",
			content:  []string{},
			lines:    []string{"prefix_"},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Remove_lines(tt.content, tt.lines)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Remove_prefixes() = %v, want %v", result, tt.expected)
			}
		})
	}
}
