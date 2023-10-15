package pack

import (
	"strings"
	"testing"
)

func TestPackScanner(t *testing.T) {

	inputs := []string{
		"0076ca82a6dff817ec66f44342007202690a93763949 15027957951b64cf874c3557a0f3547bd83b3ff6 refs/heads/master report-status\n",
		"006c0000000000000000000000000000000000000000 cdfdb42577e2506715f8cfeacdbabc092bf63e8d refs/heads/experiment\n",
		"00000009done\n",
	}
	sr := strings.NewReader(strings.Join(inputs, ""))

	var out []string
	scn := NewScanner(sr)
	for scn.Scan() {
		out = append(out, scn.Text())
	}

	if len(out) != 4 {
		t.Fail()
	}
}
