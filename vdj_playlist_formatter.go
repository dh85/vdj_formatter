package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println()
		fmt.Println("ERROR: Missing playlist input. Please enter a valid format e.g.")
		fmt.Println(".\\vdj_playlist_formatter.exe C:\\Users\\you\\Documents\\my-vdj-playlist.csv")
		fmt.Println()
		os.Exit(1)
	}

	csvPath := os.Args[1]
	if !strings.HasSuffix(csvPath, ".csv") {
		fmt.Println()
		fmt.Println("ERROR: Invalid playlist input. Filepath must end with csv e.g.")
		fmt.Println(".\\vdj_playlist_formatter.exe C:\\Users\\you\\Documents\\my-vdj-playlist.csv")
		fmt.Println()
		os.Exit(1)
	}

	csvFile, err := os.Open(csvPath)
	if err != nil {
		fmt.Println("ERROR:", err)
		fmt.Println()
		os.Exit(1)
	}

	defer func(csvFile *os.File) {
		err := csvFile.Close()
		if err != nil {
			fmt.Println("Error closing CSV:", err)
			fmt.Println()
			os.Exit(1)
		}
	}(csvFile)

	csvReader := csv.NewReader(csvFile)
	csvReader.FieldsPerRecord = -1
	var row []string
	rowTop, err := csvReader.Read()
	if err != nil {
		fmt.Println("ERROR: reading first row of CSV:", err)
		fmt.Println()
		os.Exit(1)
	}

	row = rowTop

	rowSecond, err := csvReader.Read()
	if err != nil {
		fmt.Println("ERROR: reading second row of CSV:", err)
		fmt.Println()
		os.Exit(1)
	}

	row = rowSecond

	artistIndex := Contains(row, "Artist")
	if artistIndex == -1 {
		artistIndex = 1
	}

	titleIndex := Contains(row, "Title")
	if titleIndex == -1 {
		titleIndex = 0
	}

	rows, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	unformatted := welcomeUser()
	fmt.Println()

	output := ""
	for _, row := range rows {
		if unformatted {
			output += fmt.Sprintf("%s - %s\n", row[artistIndex], row[titleIndex])
		} else {
			output += fmt.Sprintf("/me ♬♪♫ now playing... ★★┊ %s - %s ┊★★\n", row[artistIndex], row[titleIndex])
		}
	}

	extractOutput(unformatted, output)
}

func welcomeUser() bool {
	fmt.Println()
	fmt.Println("#####################################################")
	fmt.Println("###   WELCOME TO MULAN'S 3DX PLAYLIST FORMATTER   ###")
	fmt.Println("#####################################################")
	fmt.Println()
	fmt.Println("Would you like your playlist formatted for 3DX or a plain list (e.g. sharing on Discord)?")
	fmt.Println("1 - Formatted for 3DX (default)")
	fmt.Println("2 - Plain list")
	fmt.Println()

	r := bufio.NewReader(os.Stdin)

	for {
		userInput := getUserInput("Your choice", r)
		if userInput == "1" || userInput == "" {
			return false
		} else if userInput != "2" {
			userInput = ""
			fmt.Println("WARNING: Invalid input. Try again with 1 or 2...")
		} else {
			return true
		}
	}
}

func getUserInput(prompt string, reader *bufio.Reader) string {
	fmt.Printf("%s: ", prompt)
	userInput, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Input error:", err)
		os.Exit(1)
		return ""
	}
	userInput = cleanUserInput(userInput)
	return userInput
}

func cleanUserInput(userInput string) string {
	userInput = strings.Replace(userInput, "\n", "", -1)
	userInput = strings.ToLower(userInput)
	userInput = strings.TrimSpace(userInput)
	return userInput
}

func extractOutput(unformatted bool, output string) {
	exitCode := 0
	if unformatted {
		exitCode = printUnformattedOutput(output)
	} else {
		exitCode = createAndOpenFile(output)
	}
	os.Exit(exitCode)
}

func printUnformattedOutput(output string) int {
	fmt.Println("###########################################")
	fmt.Println("###   Unformatted output for playlist   ###")
	fmt.Println("###########################################")
	fmt.Println(output)
	fmt.Println("Successfully listed!!")
	fmt.Println()
	return 0
}

func createAndOpenFile(output string) int {
	filename := "output.txt"
	err := os.WriteFile(filename, []byte(output), 0666)
	if err != nil {
		fmt.Printf("ERROR: writing to %s\n", filename)
		return 1
	}

	openFileErr := open.Run(filename)
	if openFileErr != nil {
		fmt.Printf("Success, but could not open file automatically. Open '%s' in this directory manually.\n", filename)
		return 0
	}

	ex, filepathDirErr := os.Executable()
	if filepathDirErr == nil {
		exPath := filepath.Dir(ex)
		initialFilename := filename
		isWindows := runtime.GOOS == "windows"
		if isWindows {
			filename = fmt.Sprintf("%s\\%s", exPath, initialFilename)
		} else {
			filename = fmt.Sprintf("%s/%s", exPath, initialFilename)
		}
	} else {
		fmt.Println("Error getting filepath:", filepathDirErr)
		return 1
	}

	fmt.Println("Opening...", filename)
	fmt.Println()
	fmt.Println("Successfully formatted!!")
	return 0
}

func Contains(sl []string, name string) int {
	for idx, v := range sl {
		if v == name {
			return idx
		}
	}
	return -1
}
