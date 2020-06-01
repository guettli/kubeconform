package output

import (
	"bytes"
	"testing"
)

func TestJSONWrite(t *testing.T) {
	type result struct {
		fileName, kind, version string
		err                     error
		skipped                 bool
	}

	for _, testCase := range []struct {
		name        string
		withSummary bool
		verbose     bool

		res    []result
		expect string
	}{
		{
			"a single deployment, no summary, no verbose",
			false,
			false,
			[]result{
				{
					"deployment.yml",
					"Deployment",
					"apps/v1",
					nil,
					false,
				},
			},
			`{
  "resources": []
}
`,
		},
		{
			"a single deployment, summary, no verbose",
			true,
			false,
			[]result{
				{
					"deployment.yml",
					"Deployment",
					"apps/v1",
					nil,
					false,
				},
			},
			`{
  "resources": [],
  "summary": {
    "valid": 1,
    "invalid": 0,
    "errors": 0,
    "skipped": 0
  }
}
`,
		},
		{
			"a single deployment, verbose, with summary",
			true,
			true,
			[]result{
				{
					"deployment.yml",
					"Deployment",
					"apps/v1",
					nil,
					false,
				},
			},
			`{
  "resources": [
    {
      "filename": "deployment.yml",
      "kind": "Deployment",
      "version": "apps/v1",
      "status": "VALID",
      "msg": ""
    }
  ],
  "summary": {
    "valid": 1,
    "invalid": 0,
    "errors": 0,
    "skipped": 0
  }
}
`,
		},
	} {
		w := new(bytes.Buffer)
		o := JSON(w, testCase.withSummary, testCase.verbose)

		for _, res := range testCase.res {
			o.Write(res.fileName, res.kind, res.version, res.err, res.skipped)
		}
		o.Flush()

		if w.String() != testCase.expect {
			t.Fatalf("%s - expected %s, got %s", testCase.name, testCase.expect, w)
		}
	}
}