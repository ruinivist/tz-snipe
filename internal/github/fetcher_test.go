package github

import "testing"

func TestGetPushPatches(t *testing.T) {
	// this IS a BAD test
	// I write these to test out individual functions as I write them, so they not not
	// in any good sense a "test" test rather more like if what I wrote even works at all
	// is there a better way?
	patches, err := getPushPatches("torvalds")
	if err != nil {
		t.Fatal(err)
	}

	if len(patches) == 0 {
		t.Fatalf("no patches found")
	}
}

func TestGetUniqTzFromPatches(t *testing.T) {
	// another test in the same vein
	patches, _ := getPushPatches("torvalds")
	tzs, err := getUniqueTzFromPatches(patches)
	if err != nil {
		t.Fatal(err)
	}

	if len(tzs) == 0 {
		t.Fatalf("no timezones found")
	}
}
