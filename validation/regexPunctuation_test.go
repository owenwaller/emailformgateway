// Copyright (c) 2024 Owen Waller. All rights reserved.
package validation

import (
	"regexp"
	"testing"
)

func TestUnicodePuncutation(t *testing.T) {
	s := "("
	re := regexp.MustCompile(`\pP+`)
	expected := re.MatchString(s)
	if expected != true {
		t.Errorf("( is not in punctuation class.\n")
	}

	re = regexp.MustCompile(`\pL+`)
	expected = re.MatchString(s)
	if expected != false {
		t.Errorf("( is in letter class.\n")
	}

	re = regexp.MustCompile(`\pN+`)
	expected = re.MatchString(s)
	if expected != false {
		t.Errorf("( is in number class.\n")
	}

	re = regexp.MustCompile(`\pS+`)
	expected = re.MatchString(s)
	if expected != false {
		t.Errorf("( is in Symbol class.\n")
	}

	re = regexp.MustCompile(`\pM+`)
	expected = re.MatchString(s)
	if expected != false {
		t.Errorf("( is in mark class.\n")
	}

	re = regexp.MustCompile(`\pC+`)
	expected = re.MatchString(s)
	if expected != false {
		t.Errorf("( is in control class.\n")
	}

	re = regexp.MustCompile(`\pZ+`)
	expected = re.MatchString(s)
	if expected != false {
		t.Errorf("( is in separator class.\n")
	}

}
