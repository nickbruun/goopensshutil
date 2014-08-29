package opensshutil

import (
	"fmt"
	"testing"
)

type configScannerOutput struct {
	Token   configToken
	Literal string
}

type configScannerLineFixture struct {
	Src    string
	Output []configScannerOutput
	Err    bool
}

var (
	// Valid fixtures.
	configScannerLineFixtureVA = configScannerLineFixture{
		"Host *",
		[]configScannerOutput{
			{CT_KEYWORD, "Host"},
			{CT_STRING, "*"},
		},
		false,
	}

	configScannerLineFixtureVB = configScannerLineFixture{
		"Host=*",
		[]configScannerOutput{
			{CT_KEYWORD, "Host"},
			{CT_EQUAL, "="},
			{CT_STRING, "*"},
		},
		false,
	}

	configScannerLineFixtureVC = configScannerLineFixture{
		"Host =*",
		[]configScannerOutput{
			{CT_KEYWORD, "Host"},
			{CT_EQUAL, "="},
			{CT_STRING, "*"},
		},
		false,
	}

	configScannerLineFixtureVD = configScannerLineFixture{
		"Host = *",
		[]configScannerOutput{
			{CT_KEYWORD, "Host"},
			{CT_EQUAL, "="},
			{CT_STRING, "*"},
		},
		false,
	}

	configScannerLineFixtureVE = configScannerLineFixture{
		"\tHost\t = *\t",
		[]configScannerOutput{
			{CT_KEYWORD, "Host"},
			{CT_EQUAL, "="},
			{CT_STRING, "*"},
		},
		false,
	}

	configScannerLineFixtureVF = configScannerLineFixture{
		"     Exec \r   ssh -oSomeOption=123 \t",
		[]configScannerOutput{
			{CT_KEYWORD, "Exec"},
			{CT_KEYWORD, "ssh"},
			{CT_STRING, "-oSomeOption=123"},
		},
		false,
	}

	configScannerLineFixtureVG = configScannerLineFixture{
		"     \"Exec\" \r   \"ssh\"\t-oSomeOption=123 \"\t my command \" \t",
		[]configScannerOutput{
			{CT_STRING, "Exec"},
			{CT_STRING, "ssh"},
			{CT_STRING, "-oSomeOption=123"},
			{CT_STRING, "\t my command "},
		},
		false,
	}

	configScannerLineFixtureVH = configScannerLineFixture{
		"     Exec \r   X=Y ./do.sh \t",
		[]configScannerOutput{
			{CT_KEYWORD, "Exec"},
			{CT_STRING, "X=Y"},
			{CT_STRING, "./do.sh"},
		},
		false,
	}

	configScannerLineFixtureVI = configScannerLineFixture{
		"# Some comment",
		[]configScannerOutput{},
		false,
	}

	configScannerLineFixtureVJ = configScannerLineFixture{
		"# Some identend comment",
		[]configScannerOutput{},
		false,
	}

	configScannerLineFixturesValid = []configScannerLineFixture{
		configScannerLineFixtureVA,
		configScannerLineFixtureVB,
		configScannerLineFixtureVC,
		configScannerLineFixtureVD,
		configScannerLineFixtureVE,
		configScannerLineFixtureVF,
		configScannerLineFixtureVG,
		configScannerLineFixtureVH,
		configScannerLineFixtureVI,
		configScannerLineFixtureVJ,
	}

	// Erroneous fixtures.
	configScannerLineFixtureEA = configScannerLineFixture{
		"Host \"*\"\"I forgot a space\"",
		[]configScannerOutput{
			{CT_KEYWORD, "Host"},
			{CT_STRING, "*"},
			{CT_STRING, "I forgot a space"},
		},
		true,
	}

	configScannerLineFixturesErroneous = []configScannerLineFixture{
		configScannerLineFixtureEA,
	}

	// All fixtures.
	configScannerLineFixtures = append(configScannerLineFixturesValid, configScannerLineFixturesErroneous...)
)

func testConfigScannerOutputFixture(t *testing.T, src string, expectedOutput []configScannerOutput, expectedErr bool) {
	// Scan the output.
	s := newConfigScanner([]byte(src))
	output := make([]configScannerOutput, 0, len(expectedOutput))
	var err error

	for {
		tok, lit := s.Scan()
		output = append(output, configScannerOutput{tok, lit})

		if err = s.Err(); err != nil {
			break
		}

		if tok == CT_EOF || tok == CT_ILLEGAL {
			break
		}
	}

	// Compare the output.
	equal := len(output) == len(expectedOutput)

	if equal {
		for i, a := range output {
			e := expectedOutput[i]

			if a.Token != e.Token || a.Literal != e.Literal {
				equal = false
				break
			}
		}
	}

	if !equal {
		t.Errorf("Scan output mismatch for source `%s`, expected %v but got %v", src, expectedOutput, output)
	} else if expectedErr && err == nil {
		t.Errorf("Expected scan error for source `%s`", src)
	} else if !expectedErr && err != nil {
		t.Errorf("Unepxected scan error source `%s`: %v", src, err)
	}
}

func TestConfigScannerOutput(t *testing.T) {
	// Test line fixtures individually.
	for _, lf := range configScannerLineFixtures {
		// Test without newline.
		{
			o := lf.Output
			if !lf.Err {
				o = append(o, configScannerOutput{CT_EOF, ""})
			}

			testConfigScannerOutputFixture(t, lf.Src, o, lf.Err)
		}

		if lf.Err {
			continue
		}

		// Test with newline.
		testConfigScannerOutputFixture(t, lf.Src+"\n", append(lf.Output, configScannerOutput{CT_EOL, ""}, configScannerOutput{CT_EOF, ""}), lf.Err)

		// Test the statement repeated twice.
		// - this should catch errors in equal sign handling.
		{
			src := fmt.Sprintf("%s\n%s", lf.Src, lf.Src)
			output := make([]configScannerOutput, 0, len(lf.Output)*2+2)
			output = append(output, lf.Output...)
			output = append(output, configScannerOutput{CT_EOL, ""})
			output = append(output, lf.Output...)
			output = append(output, configScannerOutput{CT_EOF, ""})

			testConfigScannerOutputFixture(t, src, output, false)
		}
	}
}
