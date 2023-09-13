package timecode_test

import (
	"fmt"

	"github.com/abema/go-timecode/timecode"
)

func ExampleNewTimecode() {
	tc, err := timecode.NewTimecode(1800, 30000, 1001)
	if err != nil {
		panic(1)
	}
	fmt.Println(tc)
	// Output: 00:01:00:02
}

func ExampleParseTimecode() {
	tc, err := timecode.ParseTimecode("00:09:00:00", 30000, 1001)
	if err != nil {
		panic(1)
	}
	fmt.Println(tc)
	// Output: 00:09:00:02
}

func ExampleReset() {
	tc, err := timecode.NewTimecode(1798, 30000, 1001)
	if err != nil {
		panic(1)
	}
	tcc, _ := timecode.Reset(tc, 1800)
	fmt.Println(tcc)
	// Output: 00:01:00:02
}

func ExampleTimecode_Add() {
	tc1, err := timecode.NewTimecode(1798, 30000, 1001)
	if err != nil {
		panic(1)
	}
	tc2, err := timecode.NewTimecode(2, 30000, 1001)
	if err != nil {
		panic(1)
	}
	tc3, _ := tc1.Add(tc2)
	fmt.Println(tc3)
	// Output:
	// 00:01:00:02
}

func ExampleTimecode_AddFrames() {
	tc1, err := timecode.NewTimecode(1798, 30000, 1001)
	if err != nil {
		panic(1)
	}

	tc2, _ := tc1.AddFrames(2)
	fmt.Println(tc2)
	// Output: 00:01:00:02
}

func ExampleTimecode_Sub() {
	tc1, err := timecode.NewTimecode(1800, 30000, 1001)
	if err != nil {
		panic(1)
	}
	tc2, err := timecode.NewTimecode(2, 30000, 1001)
	if err != nil {
		panic(1)
	}
	tc3, _ := tc1.Sub(tc2)
	fmt.Println(tc3)
	// Output: 00:00:59:28
}

func ExampleTimecode_SubFrames() {
	tc1, err := timecode.NewTimecode(1800, 30000, 1001)
	if err != nil {
		panic(1)
	}
	tc2, _ := tc1.SubFrames(2)
	fmt.Println(tc2)
	// Output: 00:00:59:28
}

func ExampleTimecode_String() {
	tc, err := timecode.NewTimecode(3600, 60000, 1001)
	if err != nil {
		panic(1)
	}
	fmt.Println(tc.String())
	// Output: 00:01:00:04
}
