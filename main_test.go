package main

import "testing"
import "bytes"
import "strings"


type TestCompress struct {
    input  string
    output string
}


func TestRLECompress(t *testing.T) {
    tests := []TestCompress{
        {input: "AAAABCDAAB", output: "AA\x02BCDAA\x00B"},
        {input: "AAAABCDAA", output: "AA\x02BCDAA\x00"},
        {input: strings.Repeat("A", 260), output: "AA\xffAA\x01"},
        {input: "", output: ""},
    }

    for i, test := range tests {
        input  := []byte(test.input)
        output := []byte(test.output)

        res := RLECompress(input)
        if bytes.Compare(res, output) != 0 {
            t.Errorf("[test:%d] RLE compressing failed on %v, expect %v, got %v",
                     i, input, output, res)
        }
    }
}


func TestRLEDecompress(t *testing.T) {
    tests := []TestCompress{
        {input: "AA\x02BCDAA\x00B", output: "AAAABCDAAB"},
        {input: "AA\x02BCDAA\x00", output: "AAAABCDAA"},
        {input: "\x00\x00\x00\x00A", output: "\x00\x00\x00A"},
        {input: "AA", output: "AA"},  // This supposed to invalid input, but we accept it
        {input: "A", output: "A"},
        {input: "", output: ""},
    }

    for i, test := range tests {
        input  := []byte(test.input)
        output := []byte(test.output)

        res := RLEDecompress(input)
        if bytes.Compare(res, output) != 0{
            t.Errorf("[test:%d] RLE decompressing failed on %v, expect %v, got %v",
                     i, input, output, res)
        }
    }
}
