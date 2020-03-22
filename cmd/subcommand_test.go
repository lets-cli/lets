package cmd

import (
	"testing"

	"github.com/lets-cli/lets/test"
)

func TestSubCommandCmd(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		expect string
	}{
		{
			name:   "should eval global eval_env",
			args:   []string{"test-global-eval_env"},
			expect: "computed_env",
		}, {
			name:   "should parse global env",
			args:   []string{"test-global-env"},
			expect: "static_env",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			rootCmd, bufOut := newTestRootCmd(tt.args)
			test.MockArgs(tt.args)
			err := rootCmd.Execute()

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			outStr := bufOut.String()

			same, exp, got := test.CompareCmdOutput(tt.expect, outStr)
			if !same {
				t.Errorf("wrong command output. \nexpect: %s \ngot:    %s", exp, got)
			}
		})
	}
}
