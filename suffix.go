package main

import (
	"bytes"
	"errors"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"os"
	"strings"
	//	"time"
)

func makebooklist() []string {
	booklist := make([]string, 66)

	booklist[0] = "genesis"
	booklist[1] = "exodus"
	booklist[2] = "leviticus"
	booklist[3] = "numbers"
	booklist[4] = "deuteronomy"
	booklist[5] = "joshua"
	booklist[6] = "judges"
	booklist[7] = "ruth"
	booklist[8] = "1 samuel"
	booklist[9] = "2 samuel"
	booklist[10] = "1 kings"
	booklist[11] = "2 kings"
	booklist[12] = "1 chronicles"
	booklist[13] = "2 chronicles"
	booklist[14] = "ezra"
	booklist[15] = "nehemiah"
	booklist[16] = "esther"
	booklist[17] = "job"
	booklist[18] = "psalm"
	booklist[19] = "proverbs"
	booklist[20] = "ecclesiastes"
	booklist[21] = "song of solomon"
	booklist[22] = "isaiah"
	booklist[23] = "jeremiah"
	booklist[24] = "lamentations"
	booklist[25] = "ezekiel"
	booklist[26] = "daniel"
	booklist[27] = "hosea"
	booklist[28] = "joel"
	booklist[29] = "amos"
	booklist[30] = "obadiah"
	booklist[31] = "jonah"
	booklist[32] = "micah"
	booklist[33] = "nahum"
	booklist[34] = "habakkuk"
	booklist[35] = "zephaniah"
	booklist[36] = "haggai"
	booklist[37] = "zechariah"
	booklist[38] = "malachi"
	booklist[39] = "matthew"
	booklist[40] = "mark"
	booklist[41] = "luke"
	booklist[42] = "john"
	booklist[43] = "acts"
	booklist[44] = "romans"
	booklist[45] = "1 corinthians"
	booklist[46] = "2 corinthians"
	booklist[47] = "galatians"
	booklist[48] = "ephesians"
	booklist[49] = "phillippians"
	booklist[50] = "colossians"
	booklist[51] = "1 thessalonians"
	booklist[52] = "2 thessalonians"
	booklist[53] = "1 timothy"
	booklist[54] = "2 timothy"
	booklist[55] = "titus"
	booklist[56] = "philemon"
	booklist[57] = "hebrews"
	booklist[58] = "james"
	booklist[59] = "1 peter"
	booklist[60] = "2 peter"
	booklist[61] = "1 john"
	booklist[62] = "2 john"
	booklist[63] = "3 john"
	booklist[64] = "jude"
	booklist[65] = "revelation"

	return booklist
}

func FindVerse(booklist []string, bible []byte, off int64) (string, int64, error) {
	for i := 0; i < len(booklist); i++ {
		///search backwards 333 characters for verse
		for location := off; location > off-333; location-- {
			slice := string(bible[location : location+int64(len(booklist[i]))])
			if strings.EqualFold(booklist[i], slice) {
				return slice, location, nil
			}
		}
	}

	err := fmt.Sprintf("Could not find: %400s\n", bible[off:off+400])
	return "", -1, errors.New(err)
}

func CleanupData(slice []byte, searchterm string) string {
	output := bytes.Replace(slice, []byte("\x91"), []byte(""), -1)
	output = bytes.Replace(output, []byte("\x92"), []byte(""), -1)
	output = bytes.Replace(output, []byte("\x93"), []byte(""), -1)
	output = bytes.Replace(output, []byte("\x94"), []byte(""), -1)
	output = bytes.Replace(output, []byte("\x97"), []byte(""), -1)
	outputstr := strings.Replace(string(output), searchterm, "\033[0;31m"+searchterm+"\033[0m", -1)
	return outputstr
}

func rewindtoverse(bible []byte, off int64, booklist []string, searchterm string) {
	for x := off; x > off-333; x-- {
		if strings.EqualFold(searchterm, strings.ToLower(string(bible[x:x+int64(len(searchterm))]))) {
			_, location, err := FindVerse(booklist, bible, off)
			if err != nil {
				fmt.Printf("err = %s\n", err)
				return
			}

			outputstr := CleanupData(bible[location:location+400], searchterm)
			//out := strings.Split(outputstr, "\n")
			//fmt.Printf("%s\n\n", out[0])
			fmt.Printf("%400s\n\n", outputstr)
			fmt.Printf("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n")
		}
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("%s <database> <search term>\n", os.Args[0])
		return
	}

	searchterm := os.Args[2]

	booklist := makebooklist()

	b, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("ReadFile error\n")
		return
	}

	//fmt.Printf("time before suffix array creation: %v\n", time.Now())
	index := suffixarray.New(b)
	//fmt.Printf("time after suffix array creation: %v\n", time.Now())

	fd, err := os.Create(os.Args[1] + ".index")
	if err != nil {
		fmt.Printf("os.Create error %s\n", err)
		return
	}

	err = index.Write(fd)
	if err != nil {
		fmt.Printf("index.Write %d\n", err)
		return
	}

	//format the search term to include a space before it
	searchtrm := fmt.Sprintf("%s", searchterm)

	offsets := index.Lookup([]byte(searchtrm), -1)
	for _, off := range offsets {
		rewindtoverse(b, int64(off), booklist, searchtrm)

		//fmt.Printf("off=%d\n", off)
		//_ = booklist
	}
}
