package tritonhttp

import (
	"bufio"
	"log"
	"os"
	"strings"
)

/**
	Load and parse the mime.types file
**/
func ParseMIME(MIMEPath string) (MIMEMap map[string]string, err error) {

	// We use the util function to create a map for extention: type
	MIMEmap := make(map[string]string)

	// Then we try to open the file
	f, err := os.Open(MIMEPath)
	if err != nil {
		log.Panicln(err)
	}

	//close the file at the end of the program
	defer f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)

	// We scan the file by line, and split space to get extension and type
	for scanner.Scan() {
		myString := scanner.Text()
		entries := strings.Split(myString, " ")
		MIMEmap[entries[0]] = entries[1]
	}

	// If there is any error in scanning, report it.
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return MIMEmap, err
}
