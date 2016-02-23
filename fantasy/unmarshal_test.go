package fantasy

import (
	"encoding/xml"
	"testing"
	"time"
)

func TestUnmarshalIntAsBool(t *testing.T) {
	type testObject struct {
		XMLName xml.Name  `xml:"content"`
		Expired intAsBool `xml:"expired"`
	}

	var tests = []struct {
		input  []byte
		expect bool
		// are we expecting an error
		err bool
	}{
		{
			[]byte(`<content><expired>0</expired></content>`),
			false,
			false,
		},
		{
			[]byte(`<content><expired>1</expired></content>`),
			true,
			false,
		},
		{
			[]byte(`<content><expired>2</expired></content>`),
			false,
			true,
		},
	}

	for _, test := range tests {
		obj := testObject{}

		err := xml.Unmarshal(test.input, &obj)
		if test.err && err == nil {
			t.Errorf("expecting Unmaarshal intAsBool to return an error for %s", test.input)
		}

		if !test.err && err != nil {
			t.Errorf("Unmarshal intAsBool returned an error: ", err)
		}

		if obj.Expired != intAsBool(test.expect) {
			t.Errorf("intAsBool unmarshalled incorrectly got %s, expected %s", obj.Expired, test.expect)
		}
	}
}

func TestCalendarDate(t *testing.T) {
	type testObject struct {
		XMLName xml.Name     `xml:"content"`
		Date    calendarDate `xml:"date"`
	}

	var tests = []struct {
		input  []byte
		expect string
		// are we expecting an error
		err bool
	}{
		{
			[]byte(`<content><date>2014-05-22</date></content>`),
			"2014-05-22",
			false,
		},
		{
			[]byte(`<content><date>22-05-2014</date></content>`),
			time.Time{}.Format("2006-01-02"),
			true,
		},
	}

	for _, test := range tests {
		obj := testObject{}

		err := xml.Unmarshal(test.input, &obj)
		if test.err && err == nil {
			t.Errorf("expecting Unmaarshal intAsBool to return an error for %s", test.input)
		}

		if !test.err && err != nil {
			t.Errorf("Unmarshal intAsBool returned an error: ", err)
		}
		dateStr := time.Time(obj.Date).Format("2006-01-02")
		if dateStr != test.expect {
			t.Errorf("calendarDate unmarshalled incorrectly got %s, expected %s", dateStr, test.expect)
		}
	}
}
