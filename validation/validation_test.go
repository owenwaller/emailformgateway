// Copyright (c) 2024 Owen Waller. All rights reserved.
package validation

import "testing"

func TestRemoveEscapeSequences(t *testing.T) {
	s := "\f\fx\a\a\b"
	expected := "x"
	s = RemoveEscapeSequences(s)
	if expected != s {
		t.Fatalf("Single line case. Could not remove escape sequences. Expected %#v, got %#v\n", expected, s)
	}
	s = "\f\n\fx\n\a\a\n\b"
	expected = "x"
	s = RemoveEscapeSequences(s)
	if expected != s {
		t.Fatalf("Multiline case. Could not remove escape sequences. Expected %#v, got %#v\n", expected, s)
	}
}

func TestRemoveEmailHeaders(t *testing.T) {
	s := "Hello:To:Worldbcc:"
	s = RemoveEmailHeaders(s)
	expected := "Hello:World"
	if expected != s {
		t.Fatalf("Single line case. Did not strip email headers. Expected \"%v\" got\"%v\"\n", expected, s)
	}

	s = "Hello:\nTo:WORLD\nbcc:Subject:"
	expected = "Hello:\nWORLD\nSubject:"
	s = RemoveEmailHeaders(s)
	if expected != s {
		t.Fatalf("Multi-line case. Did not strip email headers. Expected \"%v\" got\"%v\"\n", expected, s)
	}
}

func TestAcceptUnicodeLettersSpacesPunctuation(t *testing.T) {
	var s string
	var expected bool
	var result bool
	// punctuation
	// * is punctuation
	s = "%,.?!\"&(*"
	expected = true
	result = AcceptUnicodeLettersSpacesPunctuation(s)
	if result != expected {
		t.Fatalf("All punctuation case failed.\n")
	}

	s = "ABCdef" // space is a separator we also include tab, CR and LF
	expected = true
	result = AcceptUnicodeLettersSpacesPunctuation(s)
	if result != expected {
		t.Fatalf("All ASCII letter case failed.\n")
	}

	s = "\u65e5\u672c\u8A9E" // space is a separator we also include tab, CR and LF
	expected = true
	result = AcceptUnicodeLettersSpacesPunctuation(s)
	if result != expected {
		t.Fatalf("All Unicode letter case failed.\n")
	}

	s = "\n\t\r " // space is a separator we also include tab, CR and LF
	expected = true
	result = AcceptUnicodeLettersSpacesPunctuation(s)
	if result != expected {
		t.Fatalf("All separator case failed.\n")
	}

	// Number
	s = "123"
	expected = false
	result = AcceptUnicodeLettersSpacesPunctuation(s)
	if result != expected {
		t.Fatalf("All number case passed - it should fail.\n")
	}

	// Control
	s = "\f" // form feed is a control cahracter
	expected = false
	result = AcceptUnicodeLettersSpacesPunctuation(s)
	if result != expected {
		t.Fatalf("All control case passed - it should fail.\n")
	}

	// Math
	s = "+" // + is a math cahracter
	expected = false
	result = AcceptUnicodeLettersSpacesPunctuation(s)
	if result != expected {
		t.Fatalf("All math case passed - it should fail.\n")
	}

	// Symbol
	s = "$" // $ is a symbol
	expected = false
	result = AcceptUnicodeLettersSpacesPunctuation(s)
	if result != expected {
		t.Fatalf("All symbol case passed - it should fail.\n")
	}

	// letters, separation and punctuation case
	s = "ABC !"
	expected = true
	result = AcceptUnicodeLettersSpacesPunctuation(s)
	if result != expected {
		t.Fatalf("Letter separation and punctuation case failed.\n")
	}

	// letters, separation, punctuation, numbers case
	s = "ABC ! 123"
	expected = false
	result = AcceptUnicodeLettersSpacesPunctuation(s)
	if result != expected {
		t.Fatalf("Number, letter separation and punctuation case passed - should fail.\n")
	}

	// letters, separation, punctuation, numbers control case
	s = "ABC ! \n 123"
	expected = false
	result = AcceptUnicodeLettersSpacesPunctuation(s)
	if result != expected {
		t.Fatalf("Number, letter separation, punctuation and control case passed - should fail.\n")
	}

	// letters, separation, punctuation, numbers, control and math case
	s = "ABC ! \n 123 ="
	expected = false
	result = AcceptUnicodeLettersSpacesPunctuation(s)
	if result != expected {
		t.Fatalf("Number, letter separation, punctuation, control and math case passed - should fail.\n")
	}

	// letters, separation, punctuation, numbers, control, math and symbol case
	s = "ABC ! \n 123 = £"
	expected = false
	result = AcceptUnicodeLettersSpacesPunctuation(s)
	if result != expected {
		t.Fatalf("Number, letter separation, punctuation, control, math and symbol case passed - should fail.\n")
	}

}

func TestAdditionalControlCharacters(t *testing.T) {
	var s string
	var expected bool
	var result bool

	s = "\t"
	expected = true
	result = AcceptUnicodeLettersSpacesPunctuation(s)
	if result != expected {
		t.Fatalf("\\t test failed - should pass\n")
	}

	s = "\a"
	expected = false
	result = AcceptUnicodeLettersSpacesPunctuation(s)
	if result != expected {
		t.Fatalf("\\t control charcter test passed - should fail.\n")
	}
}

func TestAcceptAllUnicodeExceptControl(t *testing.T) {
	var s string
	var expected bool
	var result bool
	// punctuation, math, symbol, letters, separators and numbers
	s = "%,.?!\"&(* += $£ ABCdef    1234 \t"
	expected = true
	result = AcceptAllUnicodeExceptControl(s)
	if result != expected {
		t.Fatalf("Everything except control case failed.\n")
	}

	// control and punctuation, math, symbol, letters, separators and numbers
	s = "\f \t %,.?!\"&(* += $£ ABCdef    1234"
	expected = false
	result = AcceptAllUnicodeExceptControl(s)
	if result != expected {
		t.Fatalf("Everything including control case passed - should fail.\n")
	}
}

func TestRemoveScriptTagsAndContents(t *testing.T) {
	var s string
	var expected string
	var result string
	// simple script case
	s = "<script>window.alert(\"some tesxt\");</script>"
	expected = ""
	result = RemoveScriptTagsAndContents(s)
	if expected != result {
		t.Fatalf("Failed to remove all script tags and contents. Expected \"%v\", got \"%v\"\n", expected, result)
	}

	// script in a string case
	s = "hello <SCRIPT>window.alert(\"some tesxt\");</script>world"
	expected = "hello world"
	result = RemoveScriptTagsAndContents(s)
	if expected != result {
		t.Fatalf("Failed to remove all script tags and contents but keep other text. Expected \"%v\", got \"%v\"\n", expected, result)
	}

}
