Summary:\s\s
There are two programs contained in this directory.  They are:\s\s
biblecodes.go - a program to calculate and verify various Bible codes.\s\s
suffix.go  - a program for searching the bible and other texts using suffix arrays.\s\s
\s\s
Usage:\s\s
1) An exploration of some Bible codes\s\s
go run biblecodes.go\s\s
\s\s
open a browser and go to:\s\s
http://[ip of server]:8080/\s\s
\s\s
2) Command line search engine:\s\s
go run suffix.go <textfile to search> <searchterm>  | less -R\s\s
go run suffix.go bible.txt Jesus\s\s



