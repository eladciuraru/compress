package main

import "flag"
import "os"
import "io/ioutil"
import "fmt"
import "math"


type Arguments struct {
    action string
    input  string
    output string
}


func inArray(value string, array []string) bool {
    for _, str := range array {
        if value == str {
            return true
        }
    }

    return false
}


func assertArgs(args *Arguments) {
    var comment string

    if inArray("", []string{args.input, args.output}) {
        comment = fmt.Sprintln("Missing required input/output path parameters")
    } else if !inArray(args.action, []string{"compress", "decompress"}) {
        comment = fmt.Sprintln("Unkown action:", args.action)
    }

    if comment != "" {
        flag.PrintDefaults()
        fmt.Fprintln(os.Stderr, "\n", comment)
        os.Exit(1)
    }
}


func parseArgs() Arguments {
    const (
        actionHelp = "Choose between compress/decompress"
        inputHelp  = "Path to input file"
        outputHelp = "Path to output file"
        shortHelp  = " (shorthand)"
    )

    var args Arguments

    flag.StringVar(&args.action, "action", args.action, actionHelp)
    flag.StringVar(&args.action, "a", args.action, actionHelp + shortHelp)

    flag.StringVar(&args.input, "input", args.input, inputHelp)
    flag.StringVar(&args.input, "i", args.input, inputHelp + shortHelp)

    flag.StringVar(&args.output, "output", args.output, outputHelp)
    flag.StringVar(&args.output, "o", args.output, outputHelp + shortHelp)

    flag.Parse()

    assertArgs(&args)

    return args
}


const (
    RLEMinSize = 0x02
    RLEMaxSize = math.MaxUint8 + RLEMinSize
)


func RLECompress(data []byte) []byte {
    buff := make([]byte, 0, len(data))
    for index := 0; index < len(data); index++ {
        byte_ := data[index]

        runIndex := 1 //index + 1
        for runIndex < len(data) - index &&
            runIndex < RLEMaxSize &&
            byte_ == data[index + runIndex] {
            runIndex++
        }

        if size := runIndex; size >= RLEMinSize {
            buff = append(buff, byte_, byte_, byte(size - RLEMinSize))
        } else {
            buff = append(buff, byte_)
        }

        index += runIndex - 1
    }

    return buff
}


func RLEDecompress(data []byte) []byte {
    buff := make([]byte, 0, len(data))  // TODO: Decide which value to use here as a starting cap?
    for index := 0; index < len(data); index++ {
        buff = append(buff, data[index])
        if index + 1 == len(data) {
            break
        }

        currByte := data[index]
        nextByte := data[index + 1]

        // 2 consecutive equal bytes indicate a run
        if currByte == nextByte {
            if index + RLEMinSize >= len(data) {
                // What to do here? there is a missing size byte, return nil?
            } else {
                buff = append(buff, nextByte)
                for size := int(data[index + RLEMinSize]); size > 0; size-- {
                    buff = append(buff, currByte)
                }

                // Next 2 bytes were used
                index += RLEMinSize
            }
        }
    }

    return buff
}


func main() {
    args := parseArgs()

    input, err := ioutil.ReadFile(args.input)
    if err != nil {
        panic(fmt.Errorf("Failed to open file %s", args.input))
    }

    var output []byte
    if args.action == "compress" {
        fmt.Printf("[*] Compressing %s - size: %d\n", args.input, len(input))
        output = RLECompress(input)
    } else if args.action == "decompress" {
        fmt.Printf("[*] Decompress %s - size: %d\n", args.input, len(input))
        output = RLEDecompress(input)
    }

    fmt.Printf("[*] Writing %s - size: %d\n", args.output, len(output))
    err = ioutil.WriteFile(args.output, output, 0666)
    if err != nil {
        panic(fmt.Errorf("Failed to create file %s", args.output))
    }
}
