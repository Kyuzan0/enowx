package daemon

import (
	"encoding/json"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/enowx/enowx/core/agent"
)

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)
	args := os.Args[3:]
	_ = json.NewEncoder(os.Stdout).Encode(args)
}

func TestPsQuotePreservesArguments(t *testing.T) {
	testCases := [][]string{
		{"--append-system-prompt", agent.Guide},
		{"simple", "words"},
		{"don't", "it's 'great'"},
		{`say "hello"`, `nested "quotes" inside`},
		{`'hello' and "world"`, `"don't" do that`},
		{`eval "$("$ENOWX_BIN" display)"`},
		{"$VAR", "$1", "$@", "${VAR}"},
		{"line 1\nline 2\r\nline 3\ttab"},
		{`path\to\"file"`, `C:\Program Files\app`},
		{"", "next"},
	}

	for _, tc := range testCases {
		cmdStr := "& " + psQuote(append([]string{os.Args[0], "-test.run=TestHelperProcess", "--"}, tc...))
		cmd := exec.Command("powershell.exe", "-NoLogo", "-Command", cmdStr)
		cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("powershell execution failed: %v: %s", err, out)
		}
		var got []string
		if err := json.Unmarshal(out, &got); err != nil {
			t.Fatalf("failed to parse helper output %q: %v", string(out), err)
		}
		if !reflect.DeepEqual(got, tc) {
			t.Errorf("mismatch:\n got: %#v\nwant: %#v", got, tc)
		}
	}
}
