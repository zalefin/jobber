package main

import (
    "os"
    "fmt"
    "bufio"
    "os/exec"
    "sync"
    "github.com/google/shlex"
    "strings"
    "strconv"
)

var sem chan int

func ParseCPUCount() int {
    // run the lscpu command and parse the number of CPUs on the system

    out, _ := exec.Command("lscpu").Output()
    outstring := strings.TrimSpace(string(out))
    lines := strings.Split(outstring, "\n")

    for _, line := range lines {
        fields := strings.Split(line, ":")
        if len(fields) < 2 {
            continue
        }
        key := strings.TrimSpace(fields[0])
        value := strings.TrimSpace(fields[1])

        switch key {
        case "CPU(s)":
            t, _ := strconv.Atoi(value)
            return int(t)
        }
    }
    // no match found, ERROR
    return -1
}

func main() {
    // check if second argument exists
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "Usage: ./jobber [TARGETS PATH]")
        os.Exit(-1)
    }
    fpath := os.Args[1]

    n_cpu := ParseCPUCount()
    if (n_cpu == -1) {
        // parse CPU count failed
        fmt.Fprintln(os.Stderr, "Error: Parse CPU count failed!")
        os.Exit(-1)
    }

    // init channel semaphore
    sem = make(chan int, n_cpu)

    // create scanner pointer
    var scanner *bufio.Scanner

    // check if reading from stdin
    if fpath == "-" {
        fmt.Println("READING STDIN")
        scanner = bufio.NewScanner(os.Stdin)
    } else {
        file, err := os.Open(fpath)
        if err != nil {
            fmt.Fprintln(os.Stderr, "bad path!")
            os.Exit(-1)
        }
        defer file.Close()
        scanner = bufio.NewScanner(file)
    }

    var commands [][]string

    for scanner.Scan() {
        args, err := shlex.Split(scanner.Text());
        if err != nil {
            fmt.Fprintln(os.Stderr, "parse error!")
            os.Exit(-1)
        }
        commands = append(commands, args)
    }

    var wg sync.WaitGroup

    for _, command := range commands {
        sem <- 1
        fmt.Println("ran job")
        wg.Add(1)
        go func(command []string) {
            defer wg.Done()
            cmd := exec.Command(command[0], command[1:]...)
            out, _ := cmd.Output()
            fmt.Print(string(out))
            <- sem
        }(command)
    }
    // wait for all jobs to complete
    wg.Wait()
}
