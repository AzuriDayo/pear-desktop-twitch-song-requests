package songrequests

import "testing"

type behavior struct {
	Query          string
	ExpectedOutput expectedOutput
}

type expectedOutput struct {
	IsNative bool
	Value    string
}

var queryBehaviors = []behavior{
	{
		Query: "PASS_VIDEO_ID",
		ExpectedOutput: expectedOutput{
			IsNative: false,
			Value:    "PASS_VIDEO_ID",
		},
	},
	{
		Query: "https://www.google.com/PASS_VIDEO_ID",
		ExpectedOutput: expectedOutput{
			IsNative: false,
			Value:    "https://www.google.com/PASS_VIDEO_ID",
		},
	},
	{
		Query: "https://youtu.be/PASS_VIDEO_ID",
		ExpectedOutput: expectedOutput{
			IsNative: true,
			Value:    "PASS_VIDEO_ID",
		},
	},
	{
		Query: "https://youtu.be/PASS_VIDEO_ID?shareid=asd",
		ExpectedOutput: expectedOutput{
			IsNative: true,
			Value:    "PASS_VIDEO_ID",
		},
	},
	{
		Query: "https://youtu.be/PASS_VIDEO_ID/wrongurlbutshouldpass&shareid=asd",
		ExpectedOutput: expectedOutput{
			IsNative: true,
			Value:    "PASS_VIDEO_ID",
		},
	},
	{
		Query: "https://www.youtube.com/watch?v=PASS_VIDEO_ID",
		ExpectedOutput: expectedOutput{
			IsNative: true,
			Value:    "PASS_VIDEO_ID",
		},
	},
	{
		Query: "https://music.youtube.com/watch?v=PASS_VIDEO_ID",
		ExpectedOutput: expectedOutput{
			IsNative: true,
			Value:    "PASS_VIDEO_ID",
		},
	},
	{
		Query: "https://music.youtube.com/watch?v=PASS_VIDEO_ID&list=LIST_ID",
		ExpectedOutput: expectedOutput{
			IsNative: true,
			Value:    "PASS_VIDEO_ID",
		},
	},
}

func TestParseSearchQuery(t *testing.T) {
	for _, tt := range queryBehaviors {
		t.Run(tt.Query, func(t *testing.T) {
			out, b := ParseSearchQuery(tt.Query)
			actualOutput := expectedOutput{
				IsNative: b,
				Value:    out,
			}
			if tt.ExpectedOutput != actualOutput {
				t.Errorf("got %s %t, want %s %t", actualOutput.Value, actualOutput.IsNative, tt.ExpectedOutput.Value, tt.ExpectedOutput.IsNative)
			}
		})
	}
}
