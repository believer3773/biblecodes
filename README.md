Summary:
There are two programs contained in this directory.  They are:
biblecodes.go - a program to calculate and verify various Bible codes.
suffix.go  - a program for searching the bible and other texts using suffix arrays.

Usage:
1) An exploration of some Bible codes
go run biblecodes.go

open a browser and go to:
http://[ip of server]:8080/

2) Command line search engine:
go run suffix.go <textfile to search> <searchterm>  | less -R
go run suffix.go bible.txt Jesus




