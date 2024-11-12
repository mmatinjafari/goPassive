package main

import (
    "fmt"
    "os"
    "os/exec"
    "strings"
    "io/ioutil"
    "time"
)

func runCommandInZsh(command string) (string, error) {
    cmd := exec.Command("bash", "-c", command)
    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Printf("Error occurred: %s\n", string(output)) // print full error output
        return "", err
    }
    return strings.TrimSpace(string(output)), nil
}

func gatherUrls(domain string) error {
    tempFile := "/tmp/tempfile-" + fmt.Sprintf("%d", time.Now().Unix())
    
    commands := []string{
        fmt.Sprintf("echo https://%s | tee %s", domain, tempFile),
        fmt.Sprintf("echo %s | waybackurls | sort -u | tee -a %s", domain, tempFile),
        fmt.Sprintf("echo %s | gau --subs -u | tee -a %s", domain, tempFile),
    }

    for _, command := range commands {
        fmt.Println("Executing command:", command)
        if _, err := runCommandInZsh(command); err != nil {
            return err
        }
    }

    // Process the temp file and output results
    data, err := ioutil.ReadFile(tempFile)
    if err != nil {
        return err
    }

    uniqueUrls := make(map[string]bool)
    lines := strings.Split(string(data), "\n")
    for _, line := range lines {
        if line != "" && !strings.HasPrefix(line, "#") { // filter out empty or comment lines
            uniqueUrls[line] = true
        }
    }

    // Write the unique URLs to a new file
    outputFile := domain + ".urls"
    file, err := os.Create(outputFile)
    if err != nil {
        return err
    }
    defer file.Close()

    for url := range uniqueUrls {
        _, err := file.WriteString(url + "\n")
        if err != nil {
            return err
        }
    }

    fmt.Printf("URLs gathered and saved to: %s\n", outputFile)
    return nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: nice_passive <domain>")
        os.Exit(1)
    }

    domain := os.Args[1]
    if err := gatherUrls(domain); err != nil {
        fmt.Printf("Error gathering URLs for %s: %v\n", domain, err)
        os.Exit(1)
    }
}

