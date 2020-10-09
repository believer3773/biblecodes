Summary:<br>
There are two programs contained in this directory.  They are:<br>
biblecodes.go - a program to calculate and verify various Bible codes.<br>
suffix.go  - a program for searching the bible and other texts using suffix arrays.<br>
<br>
<br>
Usage: An exploration of some Bible codes<br>
go run biblecodes.go<br> <br>

open a browser and go to:<br>
http://[ip of server]:8080/<br>
<br>
Command line search engine:<br>
	go run suffix.go <textfile to search> <searchterm>  | less -R<br>
	go run suffix.go bible.txt Jesus<br>
<br>
<br>
<br>
