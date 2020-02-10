package cmd

import (
	"github.com/kindritskyiMax/lets/test"
	"strings"
	"testing"
)

func TestSubCommandCmd(t *testing.T) {
	t.Run("run test-options", func(t *testing.T) {
		args := []string{"test-options"}
		rootCmd, bufOut := newTestRootCmd(args, test.GetTestConfig())
		test.MockArgs(args)
		err := rootCmd.Execute()

		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		outStr := bufOut.String()
		expect := `Flags command
            LETSOPT_KV_OPT=
            LETSOPT_BOOL_OPT=false
            LETSOPT_ARGS=
			LETSCLI_KV_OPT=
            LETSCLI_BOOL_OPT=
            LETSCLI_ARGS=`

		same, exp, got := test.CompareCmdOutput(expect, outStr)
		if !same {
			t.Errorf("wrong output. \nexpect %s, \ngot:   %s", exp, got)
		}
	})

	t.Run("run test-options", func(t *testing.T) {
		args := []string{"test-options", "--kv-opt"}
		rootCmd, _ := newTestRootCmd(args, test.GetTestConfig())
		test.MockArgs(args)
		err := rootCmd.Execute()

		if err == nil {
			t.Fatalf("must fail in docopts parsing")
		}

		expectErr := "failed to parse docopt options for cmd test-options: --kv-opt requires argument"
		if !strings.Contains(err.Error(), expectErr) {
			t.Errorf("must fail with error, \nexect: %s, \ngot:   %s", expectErr, err)
		}
		if !strings.Contains(err.Error(), "Usage:") {
			t.Errorf("must show Usage, \ngot:   %s", err)
		}
		if !strings.Contains(err.Error(), "Options:") {
			t.Errorf("must show Options, \ngot:   %s", err)
		}
	})

	tests := []struct {
		name   string
		args   []string
		expect string
		//envvars map[string]string
	}{
		{
			name: "should parse --kv-opt value",
			args: []string{"test-options", "--kv-opt=hello"},
			expect: `Flags command
            LETSOPT_KV_OPT=hello
            LETSOPT_BOOL_OPT=false
            LETSOPT_ARGS=
			LETSCLI_KV_OPT=--kv-opt hello
            LETSCLI_BOOL_OPT=
            LETSCLI_ARGS=`,
		}, {
			name: "should parse --kv-opt and --bool-opt",
			args: []string{"test-options", "--kv-opt=hello", "--bool-opt"},
			expect: `Flags command
            LETSOPT_KV_OPT=hello
            LETSOPT_BOOL_OPT=true
            LETSOPT_ARGS=
			LETSCLI_KV_OPT=--kv-opt hello
            LETSCLI_BOOL_OPT=--bool-opt
            LETSCLI_ARGS=`,
		}, {
			name: "should parse --kv-opt, --bool-opt and --args",
			args: []string{"test-options", "--kv-opt=hello", "--bool-opt", "myarg1", "myarg2"},
			expect: `Flags command
            LETSOPT_KV_OPT=hello
            LETSOPT_BOOL_OPT=true
            LETSOPT_ARGS=myarg1 myarg2
			LETSCLI_KV_OPT=--kv-opt hello
            LETSCLI_BOOL_OPT=--bool-opt
            LETSCLI_ARGS=myarg1 myarg2`, // maybe prepend lets subcommand
		}, {
			name: "should parse only --args",
			args: []string{"test-options", "myarg1", "myarg2"},
			expect: `Flags command
            LETSOPT_KV_OPT=
            LETSOPT_BOOL_OPT=false
            LETSOPT_ARGS=myarg1 myarg2
			LETSCLI_KV_OPT=
            LETSCLI_BOOL_OPT=
            LETSCLI_ARGS=myarg1 myarg2`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootCmd, bufOut := newTestRootCmd(tt.args, test.GetTestConfig())
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
