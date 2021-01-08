// WARNING: Work in progress. Goal: Validate the work done at https://sites.google.com/site/mathematicalmonotheism/
// and other similar sites.
//
// License: This code is released under the God public license (Gpl).
//
// Example of running this program:
//	go run biblecodes.go
//
// Jeremiah 33:3 Call to me and I will answer you and tell you great and unsearchable things you do not know.
// Dan 2:22 He revealeth the deep and secret things.
package main

import (
	"fmt"
	"html/template"
	//	"io/ioutil"
	"math"
	"math/big"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	//hm is short for hebrew map
	hm = make(map[string]int)

	//hom is short for hebrew ordinal map
	hom = make(map[string]int)

	//mispar shemi map
	msm = make(map[string]int)

	//gm is short for greek map
	gm = make(map[string]int)

	english_gematria = make(map[string]int)

	english_alphabet = make(map[string]int)

	first37elements = make([]elements, 37)
)

//Source of this function is: https://rosettacode.org/wiki/Semiprime#Go
func semiprime(n int) bool {
	nf := 0
	for i := 2; i <= n; i++ {
		for n%i == 0 {
			if nf == 2 {
				return false
			}
			nf++
			n /= i
		}
	}
	return nf == 2
}

// return list of primes less than N
// Source of this function is: https://stackoverflow.com/questions/21854191/generating-prime-numbers-in-go
func sieveOfEratosthenes(N int) (primes []int) {
	b := make([]bool, N)
	for i := 2; i < N; i++ {
		if b[i] == true {
			continue
		}
		primes = append(primes, i)
		for k := i * i; k < N; k += i {
			b[k] = true
		}
	}
	return
}

type elements struct {
	atomicmass   float64
	name         string
	symbol       string
	atomicnumber int
}

func init_english_alphabet() {
	english_alphabet["A"] = 1
	english_alphabet["B"] = 2
	english_alphabet["C"] = 3
	english_alphabet["D"] = 4
	english_alphabet["E"] = 5
	english_alphabet["F"] = 6
	english_alphabet["G"] = 7
	english_alphabet["H"] = 8
	english_alphabet["I"] = 9
	english_alphabet["J"] = 10
	english_alphabet["K"] = 11
	english_alphabet["L"] = 12
	english_alphabet["M"] = 13
	english_alphabet["N"] = 14
	english_alphabet["O"] = 15
	english_alphabet["P"] = 16
	english_alphabet["Q"] = 17
	english_alphabet["R"] = 18
	english_alphabet["S"] = 19
	english_alphabet["T"] = 20
	english_alphabet["U"] = 21
	english_alphabet["V"] = 22
	english_alphabet["W"] = 23
	english_alphabet["X"] = 24
	english_alphabet["Y"] = 25
	english_alphabet["Z"] = 26
}

func init_english_gematria() {
	english_gematria["A"] = 6
	english_gematria["B"] = 12
	english_gematria["C"] = 18
	english_gematria["D"] = 24
	english_gematria["E"] = 30
	english_gematria["F"] = 36
	english_gematria["G"] = 42
	english_gematria["H"] = 48
	english_gematria["I"] = 54
	english_gematria["J"] = 60
	english_gematria["K"] = 66
	english_gematria["L"] = 72
	english_gematria["M"] = 78
	english_gematria["N"] = 84
	english_gematria["O"] = 90
	english_gematria["P"] = 96
	english_gematria["Q"] = 102
	english_gematria["R"] = 108
	english_gematria["S"] = 114
	english_gematria["T"] = 120
	english_gematria["U"] = 126
	english_gematria["V"] = 132
	english_gematria["W"] = 138
	english_gematria["X"] = 144
	english_gematria["Y"] = 150
	english_gematria["Z"] = 156
}

//data came from: https://www.lenntech.com/periodic/mass/atomic-mass.htm
func init_periodictable() {

	first37elements = []elements{
		{1.0079, "Hydrogen", "H", 1},
		{4.0026, "Helium", "He", 2},
		{6.941, "Lithium", "Li", 3},
		{9.0122, "Beryllium", "Be", 4},
		{10.811, "Boron", "B", 5},
		{12.0107, "Carbon", "C", 6},
		{14.0067, "Nitrogen", "N", 7},
		{15.9994, "Oxygen", "O", 8},
		{18.9984, "Fluorine", "F", 9},
		{20.1797, "Neon", "Ne", 10},
		{22.9897, "Sodium", "Na", 11},
		{24.305, "Magnesium", "Mg", 12},
		{26.9815, "Aluminum", "Al", 13},
		{28.0855, "Silicon", "Si", 14},
		{30.9738, "Phosphorus", "P", 15},
		{32.065, "Sulfur", "S", 16},
		{35.453, "Chlorine", "Cl", 17},
		{39.948, "Argon", "Ar", 18},
		{39.0983, "Potassium", "K", 19},
		{40.078, "Calcium", "Ca", 20},
		{44.9559, "Scandium", "Sc", 21},
		{47.867, "Titanium", "Ti", 22},
		{50.9415, "Vanadium", "V", 23},
		{51.9961, "Chromium", "Cr", 24},
		{54.938, "Manganese", "Mn", 25},
		{55.845, "Iron", "Fe", 26},
		{58.9332, "Cobalt", "Co", 27},
		{58.6934, "Nickel", "Ni", 28},
		{63.546, "Copper", "Cu", 29},
		{65.39, "Zinc", "Zn", 30},
		{69.723, "Gallium", "Ga", 31},
		{72.64, "Germanium", "Ge", 32},
		{74.9216, "Arsenic", "As", 33},
		{78.96, "Selenium", "Se", 34},
		{79.904, "Bromine", "Br", 35},
		{83.8, "Krypton", "Kr", 36},
		{85.4678, "Rubidium", "Rb", 37},
	}
}

//alphabet was derived from https://en.wikipedia.org/wiki/Hebrew_alphabet
func init_hebrewlanguage() {
	hm["ת"] = 400
	hm["ש"] = 300
	hm["ר"] = 200
	hm["ק"] = 100
	hm["ץ"] = 90
	hm["צ"] = 90
	hm["פ"] = 80
	hm["ע"] = 70
	hm["ס"] = 60
	hm["נ"] = 50
	hm["נ"] = 50
	hm["ן"] = 50
	hm["ם"] = 40
	hm["מ"] = 40
	hm["ל"] = 30
	hm["כ"] = 20
	hm["י"] = 10
	hm["ט"] = 9
	hm["ח"] = 8
	hm["ז"] = 7
	hm["ו"] = 6
	hm["ה"] = 5
	hm["ד"] = 4
	hm["ג"] = 3
	hm["ב"] = 2
	hm["א"] = 1
}

func init_mispar_shemi() {
	msm["ט"] = 419
	msm["ח"] = 418
	msm["ז"] = 67
	msm["ו"] = 12
	msm["ה"] = 6
	msm["ד"] = 434
	msm["ג"] = 83
	msm["ב"] = 412
	msm["א"] = 111
	msm["ס"] = 120
	msm["ן"] = 106
	msm["נ"] = 106
	msm["ם"] = 80
	msm["מ"] = 80
	msm["ל"] = 74
	msm["ך"] = 100
	msm["כ"] = 100
	msm["י"] = 20
	//msm["ת"] = 406
	msm["ת"] = 416
	msm["ש"] = 350
	msm["ר"] = 510
	msm["ק"] = 186
	msm["ץ"] = 104
	msm["צ"] = 104
	msm["ף"] = 81
	msm["פ"] = 81
	msm["ע"] = 130
}

//hebrew ordinal map (hom) https://sites.google.com/site/mathematicalmonotheism/phi-star-of-genesis-1-1
func init_hebrewlanguage_ordinal() {
	hom["ת"] = 22
	hom["ש"] = 21
	hom["ר"] = 20
	hom["ק"] = 19
	hom["צ"] = 18
	hom["פ"] = 17
	hom["ע"] = 16
	hom["ס"] = 15
	hom["נ"] = 14
	hom["ן"] = 14
	hom["ם"] = 13
	hom["מ"] = 13
	hom["ל"] = 12
	hom["כ"] = 11
	hom["י"] = 10
	hom["ט"] = 9
	hom["ח"] = 8
	hom["ז"] = 7
	hom["ו"] = 6
	hom["ה"] = 5
	hom["ד"] = 4
	hom["ג"] = 3
	hom["ב"] = 2
	hom["א"] = 1
}

//https://en.wikipedia.org/wiki/Greek_alphabet
func init_greeklanguage() {
	gm["Α"] = 1
	gm["α"] = 1
	gm["ά"] = 1
	gm["ᾶ"] = 1 //???
	gm["Β"] = 2
	gm["β"] = 2
	gm["Γ"] = 3
	gm["γ"] = 3
	gm["Δ"] = 4
	gm["δ"] = 4
	gm["Ε"] = 5
	gm["ε"] = 5
	gm["Ζ"] = 7
	gm["ζ"] = 7
	gm["Η"] = 8
	gm["η"] = 8
	gm["Θ"] = 9
	gm["θ"] = 9
	gm["Ι"] = 10
	gm["ι"] = 10
	gm["ὶ"] = 10 //???
	gm["Κ"] = 20
	gm["κ"] = 20
	gm["Λ"] = 30
	gm["λ"] = 30
	gm["Μ"] = 40
	gm["μ"] = 40
	gm["Ν"] = 50
	gm["ν"] = 50
	gm["Ξ"] = 60
	gm["ξ"] = 60
	gm["Ο"] = 70
	gm["ο"] = 70
	gm["ό"] = 70
	gm["Ὁ"] = 70
	gm["ὁ"] = 70

	gm["Π"] = 80
	gm["π"] = 80
	gm["Ρ"] = 100
	gm["ρ"] = 100
	gm["Σ"] = 200
	gm["σ"] = 200
	gm["ς"] = 200
	gm["Τ"] = 300
	gm["τ"] = 300
	gm["Υ"] = 400
	gm["υ"] = 400

	gm["ύ"] = 400

	gm["Φ"] = 500
	gm["φ"] = 500
	gm["Χ"] = 600
	gm["χ"] = 600
	gm["Ψ"] = 700
	gm["ψ"] = 700
	gm["Ω"] = 800
	gm["ω"] = 800
}

//These two values are united through Pi. Observe that a CIRCLE with a circumference of 2368 has a diameter of 754 (Pi = The ratio between the circumference/diameter of a CIRCLE):
func ProofInThePi() []AwesomeMathInformation {
	str := ""
	html := make([]AwesomeMathInformation, 0)

	pi := Pi(10000)      //request 1000 digits of Pi. Yes, this is more than needed for this function
	piafterdec := pi[2:] //skip past "3."

	str = fmt.Sprintf("The first 36 digits of pi after decimal point are {%v}\n\tthen the following digits are: %v", piafterdec[0:36], piafterdec[36:37+8])
	html = append(html, AwesomeMathInformation{str})

	str = fmt.Sprintf("At offset 37 in Pi is 197.  The 6 digits after 197 (169399) add up to 37.")
	html = append(html, AwesomeMathInformation{str})

	if (2 + 3 + 5 + 7 + 11 + 13 + 17 + 19 + 23 + 29 + 31 + 37) == 197 {
		str = fmt.Sprintf("The sum of all Primes up to 37 (2 + 3 + 5 + 7 + 11 + 13 + 17 + 19 + 23 + 29 + 31 + 37) = 197")
		html = append(html, AwesomeMathInformation{str})
	}

	const prec = 370 //this asks for more precision than needed
	a, _ := new(big.Float).SetPrec(prec).SetString("1.0")
	b, _ := new(big.Float).SetPrec(prec).SetString("754.0")
	result := new(big.Float).SetPrec(prec).Quo(a, b) //Quo == Divide

	//compute the digital root or digital sum of the first 360 digits
	floatstr := result.Text('g', 83)[2:]
	//fmt.Println(str)
	jesusgreek_circle_total := 0
	for _, digit := range floatstr {
		dig, _ := strconv.Atoi(string(digit))
		jesusgreek_circle_total = jesusgreek_circle_total + dig
	}

	str = fmt.Sprintf("'Jesus Christ' when represented in Greek characters, the characters add up to 754.  The continuously repeating 83 digits found by dividing 1/754 adds up to: %d or the number of degrees in a circle.", jesusgreek_circle_total)
	html = append(html, AwesomeMathInformation{str})

	///The Pi code of 666
	for i, _ := range piafterdec {
		if strings.Compare(piafterdec[i:i+4], "2701") == 0 {
			str = fmt.Sprintf("At offset %d after the Pi decimal point is:'%v'. By adding the positions of 2701 (i.e. 165+166+167+168) you get 666.", i+1, piafterdec[i:i+4])
			html = append(html, AwesomeMathInformation{str})
			break
		}
	}

	total := float64(0.0)
	s := strings.Split(pi, "37")

	for _, digit := range s[1] {
		digit, _ := strconv.Atoi(string(digit))
		total = total + math.Pow(float64(digit), 2.0)
		//fmt.Println("i=", i, "digit=", digit, "total=", total)
		if total >= 703 {
			break
		}
	}

	str = fmt.Sprintf("The sum of the SQUARES of all the Pi digits immediately following the FIRST appearance of 37 up to the 73rd decimal digit of Pi = T37 =%v.", total)
	html = append(html, AwesomeMathInformation{str})

	str = fmt.Sprintf("Epistle in Greek is ἐπιστολή which has a value of 703.  Epistle means letter and is divisible by 37.")
	html = append(html, AwesomeMathInformation{str})

	//There are 144 hours in the 6 days of creation and the sum of the first 144 decimal digits of Pi (after the decimal point) = 666
	total = float64(0.0)
	for i, _ := range piafterdec {
		if i == 144 {
			break
		}
		digit := piafterdec[i : i+1]
		d, _ := strconv.Atoi(digit)
		total = total + float64(d)
	}
	str = fmt.Sprintf("There are 144 hours in the 6 days of creation and the sum of the first 144 decimal digits of Pi (after the decimal point) =%v.", total)
	html = append(html, AwesomeMathInformation{str})

	//144/666 is an infinitely repeating cycle of 6x6x6  NOTE: float64 cannot repeat infinitely.
	str = fmt.Sprintf("144/666 yields the infinitely repeating sequence of %.15F...", float64(144.0/666.0))
	html = append(html, AwesomeMathInformation{str})

	str = fmt.Sprintf("216 is equal to 6x6x6.")
	html = append(html, AwesomeMathInformation{str})

	str = fmt.Sprintf("666 is also evenly divisible by 37.")
	html = append(html, AwesomeMathInformation{str})

	return html

}

//This function borrowed from https://github.com/JJ/pigo/blob/master/pigo.go
func Pi(ndigits int64) string {
	if ndigits <= 7 {
		return "3.141595"
	} else {
		digits := big.NewInt(ndigits + 10)
		unity := big.NewInt(0)                 // crea un entero tocho
		unity.Exp(big.NewInt(10), digits, nil) // Le asigna valor
		pi := big.NewInt(0)
		four := big.NewInt(4) // Todos deben ser enteros tocho

		// Serie de McLaurin
		pi.Mul(four, pi.Sub(pi.Mul(four, arccot(5, unity)), arccot(239, unity)))
		output := fmt.Sprintf("%s.%s", pi.String()[0:1], pi.String()[1:ndigits])
		return output
	}
}

//This function borrowed from https://github.com/JJ/pigo/blob/master/pigo.go
func arccot(x int64, unity *big.Int) *big.Int {
	bigx := big.NewInt(x)
	xsquared := big.NewInt(x * x)
	sum := big.NewInt(0)
	sum.Div(unity, bigx)
	xpower := big.NewInt(0)
	xpower.Set(sum)
	n := int64(3)
	zero := big.NewInt(0)
	sign := false

	term := big.NewInt(0)
	for {
		xpower.Div(xpower, xsquared)
		term.Div(xpower, big.NewInt(n))
		if term.Cmp(zero) == 0 {
			break
		}
		if sign {
			sum.Add(sum, term)
		} else {
			sum.Sub(sum, term)
		}
		sign = !sign
		n += 2
	}
	return sum
}

//reverse a string using space as the separation point
func ReverseVerse(input string) []string {
	s := strings.Split(input, " ")

	a := make([]string, 0)
	for x := len(s[:]) - 1; x >= 0; x-- {
		a = append(a, s[x])
	}

	return a
}

//ms = mispar shemi (full name) gematria
func compute_verse_msm(verse string, mapprovided map[string]int) int {
	sentencesum := 0

	rv := ReverseVerse(verse)

	//fmt.Println("MSMDEBUG words=", verse)

	for _, word := range rv {
		lettersum := 0
		letters := strings.Split(word, "")
		for _, letter := range letters {
			val, ok := mapprovided[letter]
			if !ok {
				fmt.Println("In Hebrew word:", strings.Join(letters, ""), "character missing from mapprovided table:'", letter, "'")
				os.Exit(0)
			} else {
				//fmt.Println("MSMDEBUG:", letter, val)
				lettersum += val

				//				if letter == "ת" {
				//					mapprovided[letter] = 406
				//				}

			}
		}

		sentencesum += lettersum
		//fmt.Println("sum for", word, " is ", lettersum)
	}

	//fmt.Println("MSMDEBUG total=", sentencesum, "\n")
	//mapprovided["ת"] += 416

	return sentencesum
}

func compute_verse_hebrew(verse string, mapprovided map[string]int) int {
	sentencesum := 0

	rv := ReverseVerse(verse)

	for _, word := range rv {
		lettersum := 0
		letters := strings.Split(word, "")
		for _, letter := range letters {
			val, ok := mapprovided[letter]
			if !ok {
				fmt.Println("In Hebrew word:", strings.Join(letters, ""), "character missing from mapprovided table:'", letter, "'")
				os.Exit(0)
			} else {
				lettersum += val
			}
		}

		sentencesum += lettersum
		//fmt.Println("sum for", word, " is ", lettersum)
	}

	return sentencesum
}

func compute_verse_greek(verse string) int {
	sentencesum := 0

	rv := ReverseVerse(verse)

	for _, word := range rv {
		lettersum := 0
		letters := strings.Split(word, "")
		for _, letter := range letters {
			val, ok := gm[letter]
			if !ok {
				fmt.Println("In Greek word:", strings.Join(letters, ""), "character missing from gm table:", letter)
				os.Exit(0)
			} else {
				lettersum += val
			}
		}

		sentencesum += lettersum
		//fmt.Println("sum for", word, " is ", lettersum)
	}

	return sentencesum
}

func Fibonacci() []AwesomeMathInformation {

	html := make([]AwesomeMathInformation, 0)
	str := ""

	str = fmt.Sprintf("Every 37th element of Fibonacci sequence taken mod 73 equals zero")
	html = append(html, AwesomeMathInformation{str})

	str = fmt.Sprintf("Every 19th element of Fibonacci sequence taken mod 37 equals zero")
	html = append(html, AwesomeMathInformation{str})

	str = fmt.Sprintf("The number 19 is mid-way between 1 and 37")
	html = append(html, AwesomeMathInformation{str})

	return html //comment this line out to see many numbers where this rule holds true.

	modulo := new(big.Int)
	var counter int64 = 0
	seventythree := new(big.Int)
	seventythree = big.NewInt(73)

	// Initialize two big ints with the first two numbers in the sequence.
	var limit big.Int
	a := big.NewInt(0)
	b := big.NewInt(1)

	// Initialize limit as 10^2020, 2020 picked for the year this program was written
	limit.Exp(big.NewInt(10), big.NewInt(2020), nil)
	// Loop while a is smaller than 10^2020
	for a.Cmp(&limit) < 0 {

		// Compute the next Fibonacci number, storing it in a.
		a.Add(a, b)

		// Swap a and b so that b is the next number in the sequence.
		a, b = b, a

		counter = counter + 1
		if (counter % 37) == 0 { //mod 37
			result := modulo.Mod(a, seventythree) //mod 73
			if result.Cmp(big.NewInt(0)) == 0 {
				a_str := a.String()
				result_str := result.String()
				fmt.Printf("%s mod 73 == %s\n", a_str, result_str)
			}
		}

		if (counter % 19) == 0 {
			//mod 19.  19 is at the center between 1 and 37.
			//the 19th Fibonacci number is 4181.  The following
			//equation yields Phi (1.618...). (41x81)/2053.
			//https://sites.google.com/site/mathematicalmonotheism/phi-star-of-genesis-1-1

			result := modulo.Mod(a, big.NewInt(37))
			if result.Cmp(big.NewInt(0)) == 0 {
				a_str := a.String()
				result_str := result.String()
				fmt.Printf("%s mod 37 == %s\n", a_str, result_str)
			}
		}
	}

	return nil
}

type importantnumbers struct {
	//hebrewgenesiscount int
	jesushebrewcount int
	jesusgreekcount  int
}

//validating claims of atomic proof of Christianity with web site: https://www.lenntech.com/periodic/mass/atomic-mass.htm
/*
type elements struct {
	atomicmass float64
	name	   string
	symbol	   string
	atomicnumber int
}
*/
func doawesomechemistrymath(inums importantnumbers) []AwesomeChemistryInformation {
	jgc := inums.jesusgreekcount

	composites := []int{0, 4, 6, 8, 9, 10, 12, 14, 15, 16, 18, 20, 21, 22, 24, 25, 26, 27, 28, 30, 32, 33, 34, 35, 36, 38, 39, 40, 42, 44, 45, 46, 48, 49, 50, 51, 52, 54, 55, 56, 57, 58, 60, 62, 63, 64, 65, 66, 68, 69, 70, 72, 74, 75, 76, 77, 78, 80, 81, 82, 84, 85, 86, 87, 88, 90, 91, 92, 93, 94, 95, 96, 98, 99, 100, 102, 104, 105, 106, 108, 110, 111, 112, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 128, 129, 130, 132, 133, 134, 135, 136, 138, 140, 141, 142, 143, 144, 145, 146, 147, 148, 150, 152, 153, 154, 155, 156, 158, 159, 160, 161, 162, 164, 165, 166, 168, 169, 170, 171, 172, 174, 175, 176, 177, 178, 180, 182, 183, 184, 185, 186, 187, 188, 189, 190, 192, 194, 195, 196, 198, 200}

	///math fact #1
	total := 0
	for _, e := range first37elements {
		neutrons := int(math.Round(float64(e.atomicmass - float64(e.atomicnumber))))
		total += (composites[e.atomicnumber] + composites[neutrons])
		//fmt.Println("element", e.name, "protons:", e.atomicnumber, composites[e.atomicnumber], "neutrons:", neutrons, composites[neutrons])
	}

	html := make([]AwesomeChemistryInformation, 0)
	str := ""

	if total == jgc {
		str = fmt.Sprintf("The letters in 'Jesus Christ' when written in Greek add up to 2368 and the first 37 elements (proton and neutron composite numbers) add up to: %d", total)
	} else {
		str = fmt.Sprintf("fail - chemistry test")
	}

	html = append(html, AwesomeChemistryInformation{str})

	return html
}

func doawesomednabiomath() {
	/*
		• 754 + (712) = The number of PROTONS in the universal amino acids (1466)
		• 754 + (712 - 197) = The number of NEUTRONS in the universal amino acids (1269)
	*/
}

//https://sites.google.com/site/mathematicalmonotheism/
func doawesomemath(primes, semiprimes []int, inums importantnumbers) []AwesomeMathInformation {
	//compute prime factors of Genesis 1:1 in Hebrew
	//First verse in Hebrew downloaded from: https://www.ccel.org/a/anonymous/hebrewot/Genesis.html
	//296+407+395+401+86+203+913 == 2701 == 37x73
	genesis11 := "בראשית ברא אלהים את השמים ואת הארץ" //1 In the beginning God created the heaven and the earth. KJV
	hebrewgenesiscount := compute_verse_hebrew(genesis11, hm)

	html := make([]AwesomeMathInformation, 0)
	str := ""

	str = fmt.Sprintf("Genesis 1:1 In the beginning God created the heavens and the earth.")
	html = append(html, AwesomeMathInformation{str})

	str = fmt.Sprintf("Genesis 1:1 in Hebrew is %v has a standard Hebrew count of characters equal to %v.", genesis11, hebrewgenesiscount)
	html = append(html, AwesomeMathInformation{str})

	//2701 == (37 * 73)
	if hebrewgenesiscount == (37 * 73) {
		str = fmt.Sprintf("The prime factors of 2701 are 37 and 73, in other words, 37x73 is equal to 2701.")
		html = append(html, AwesomeMathInformation{str})

		str := fmt.Sprintf("Chokmah (חכמה) is the Hebrew Greek Word for wisdom.  Adding the Hebrew ordinal values of each Hebrew character in 'חכמה' equals %d.", hom["ח"]+hom["כ"]+hom["מ"]+hom["ה"])
		html = append(html, AwesomeMathInformation{str})

		str = fmt.Sprintf("By taking each character for the Hebrew word for wisdom which is 'הםכח' and add up the numerical value for each character you get %d.", hm["ח"]+hm["כ"]+hm["מ"]+hm["ה"])
		html = append(html, AwesomeMathInformation{str})

		str = fmt.Sprintf("37 is the 12th prime number and 73 is the 21st prime number. Note 12 and 21 are mirror images.")
		html = append(html, AwesomeMathInformation{str})

		//now do the mirror reflection of 2701 trick
		if (2701 + 1072) == 3773 {
			str = fmt.Sprintf("The mirror value of 2701 is 1072.  When you add (2701+1072) you get 3773.")
			html = append(html, AwesomeMathInformation{str})
		}
	}

	hebrew_mispar_shemi_count := compute_verse_msm(genesis11, msm)
	str = fmt.Sprintf("Genesis 1:1 in Hebrew using Mispar Shemi gematria (e.g. letter counting) equals %d, which is equal to 74 squared. e.g. 74x74.", hebrew_mispar_shemi_count)
	html = append(html, AwesomeMathInformation{str})

	cnt := compute_english_alphabet("JESUS")
	str = fmt.Sprintf("Adding the positions of each letter of 'JESUS' equals %d.  J=10, E=5, S=19, U=21, S=19.", cnt)
	html = append(html, AwesomeMathInformation{str})

	str = fmt.Sprintf("74 is the 37th even number.")
	html = append(html, AwesomeMathInformation{str})

	/*
		a := "בראשית" //1 In the beginning God created the heaven and the earth. KJV
		b := "ברא"    //1 In the beginning God created the heaven and the earth. KJV
		c := "אלהים"  //1 In the beginning God created the heaven and the earth. KJV
		d := "את"     //1 In the beginning God created the heaven and the earth. KJV
		e := "השמים"  //1 In the beginning God created the heaven and the earth. KJV
		f := "ואת"    //1 In the beginning God created the heaven and the earth. KJV
		g := "הארץ"   //1 In the beginning God created the heaven and the earth. KJV

		acnt := compute_verse_msm(a, msm)
		bcnt := compute_verse_msm(b, msm)
		ccnt := compute_verse_msm(c, msm)
		dcnt := compute_verse_msm(d, msm)
		ecnt := compute_verse_msm(e, msm)
		fcnt := compute_verse_msm(f, msm)
		gcnt := compute_verse_msm(g, msm)
		fmt.Println("INFO2", d, dcnt, e, ecnt)
		fmt.Println("INFO2", acnt, bcnt, ccnt, dcnt, ecnt, fcnt, gcnt)
		//INFO2 1809 1033 291 517 536 529 731
	*/

	/*
		//TODO
		754 - (37th + 73rd Triprimes) =
		ORDINAL value of Genesis 1:1 (298)
	*/

	return html

}

// http://reason.landmarkbiblebaptist.net/888.html
func compute_english_gematria(s string) int {
	total := 0

	for _, char := range s {
		total += english_gematria[string(char)]
	}

	return total
}

func compute_english_alphabet(s string) int {
	total := 0

	for _, char := range s {
		total += english_alphabet[string(char)]
	}

	return total
}

// http://reason.landmarkbiblebaptist.net/888.html
func english_gematria_jesus() []Gematria {
	list888 := []string{"THE TRINITY", "MESSIAH JESUS", "THE LION OF JUDAH", "MORNING STAR", "THE KING JESUS", "FINISHED CROSS",
		"DIVINE PRESENCE", "OMNIPRESENT", "THE LORD'S TIME", "COMING TRUTH", "KING OF THE SABBATH", "BIBLICAL PROPHET",
		"THE HOUSE OF GOD", "LAW OF LIBERTY", "SCRIPTURES", "THE PASSOVER", "LAMB OF GOD SACRIFICE", "LIVE IN CHRIST",
		"SAVED IN JESUS", "JESUS FORGAVE", "CHRIST RECEIVED", "SPIRIT BIRTH"}

	g := make([]Gematria, 0)

	for _, entry := range list888 {
		val := compute_english_gematria(entry)
		//str = fmt.Sprintf("'%v',english gematria value is %d\n", entry, val)
		//g = append(g, Gematria{str, val})
		g = append(g, Gematria{Name: entry, Count: val})
	}

	return g
}

// http://reason.landmarkbiblebaptist.net/888.html
func english_gematria_666() []Gematria {
	list666 := []string{"A SATANIC MARK", "SATAN'S SEAL", "A SATANIC PLAN", "RECEIVE A MARK", "CHOICE OF DOOM", "FOREHEAD SIGN", "THE HAND OR HEAD",
		"EDEN TEMPTED", "PEOPLE SIN", "DEVIL'S HEIR", "HUMANITY", "SON OF SIN", "WICKED WILL", "STUBBORN", "SCORNERS", "LUSTFUL",
		"TREACHERIES", "CORRUPT", "FLOOD OF NOAH", "COMPUTER", "TERMINALS", "BIO-IMPLANT", "LUCIFER HELL", "HELL: TWICE DEAD",
		"HELL: A REAL PLACE", "HELL BURNS", "SET ON FIRE", "DEAL FROM HELL", "STUPID DEAL", "DAVID ALLENDER",
	}

	g := make([]Gematria, 0)

	for _, entry := range list666 {
		val := compute_english_gematria(entry)
		//str = fmt.Sprintf("%v\n", entry, val)
		g = append(g, Gematria{Name: entry, Count: val})
	}

	return g
}

//from revelation's Mark of the Beast Exposed - Revelation's Ancient Discoveries - youtube.com
func roman_gematria_666() {
	popetitle := "VICARIUS FILII DEI"

	charmap := make(map[string]int)
	charmap["I"] = 1
	charmap["V"] = 5
	charmap["C"] = 100
	charmap["A"] = 0
	charmap["R"] = 0
	charmap["U"] = 5
	charmap["S"] = 0
	charmap["F"] = 0
	charmap["L"] = 50
	charmap["D"] = 500

	total := 0
	for _, v := range popetitle {
		val, ok := charmap[string(v)]
		if ok {
			total += val
		}
	}

	fmt.Println(popetitle, "adds up to", total)
}

type VerseInfo struct {
	Verse string
	Done  bool
}

type LanguageInformation struct {
	Character string
	Integer   string
}

type AwesomeMathInformation struct {
	Html string
}

type AwesomeChemistryInformation struct {
	Html string
}

type Gematria struct {
	Name  string
	Count int
}

type BibleData struct {
	PageTitle        string
	VerseList        []VerseInfo
	JesusHebrew      string
	JesusHebrewCount int
	HebrewInfo       []LanguageInformation
	JesusGreek       string
	JesusGreekCount  int
	GreekInfo        []LanguageInformation
	AwesomeMath      []AwesomeMathInformation
	AwesomeChemistry []AwesomeChemistryInformation
	GematriaJesus    []Gematria
	GematriaEvil     []Gematria
}

//assumes a hash map of type map[string]int
//returns a list of hash keys and values in sorted order
func retrieve_hash_map(m map[string]int) ([]string, []string) {
	keys := make([]string, 0)
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	ints := make([]string, 0)
	listentry := ""
	for _, k := range keys {
		listentry = fmt.Sprintf("%d", hm[k])
		ints = append(ints, listentry)
	}

	return keys, ints
}

func sendEarlyHtml(w http.ResponseWriter) {
	w.Write([]byte("<!DOCTYPE html>\n"))
	w.Write([]byte("<html>\n"))
	w.Write([]byte("<head>\n"))
	w.Write([]byte("<meta charset=\"utf-8\"\">\n"))
	w.Write([]byte("<script src=\"https://ajax.googleapis.com/ajax/libs/jquery/3.5.1/jquery.min.js\">\n"))
	w.Write([]byte("</script>\n"))
	w.Write([]byte("<style>\n"))

	w.Write([]byte(".block {\ndisplay: block;\n width: 100%;\n border: none;\n background-color: #4CAF50;\n"))
	w.Write([]byte("color: white;\npadding: 14px 28px;\nfont-size: 16px;\ncursor: pointer;\ntext-align: center;\n}"))
	w.Write([]byte("block:hover {\nbackground-color: #ddd;color: black;\n"))
	w.Write([]byte(" } \n\n"))

	w.Write([]byte(".arrow {\n"))

	w.Write([]byte("  border: solid black;\n"))
	w.Write([]byte("  border-width: 0 3px 3px 0;\n"))
	w.Write([]byte("  display: inline-block;\n"))
	w.Write([]byte("  padding: 3px;\n"))
	w.Write([]byte("}\n\n"))
	w.Write([]byte(".right {\n"))
	w.Write([]byte("  transform: rotate(-45deg);\n"))
	w.Write([]byte("  -webkit-transform: rotate(-45deg);\n"))
	w.Write([]byte("}\n\n"))
	w.Write([]byte(".left {\n"))
	w.Write([]byte("  transform: rotate(135deg);\n"))
	w.Write([]byte("-webkit-transform: rotate(135deg);\n"))
	w.Write([]byte("}\n\n"))
	w.Write([]byte(".up {\n"))
	w.Write([]byte("  transform: rotate(-135deg);\n"))
	w.Write([]byte("  -webkit-transform: rotate(-135deg);\n"))
	w.Write([]byte("}\n\n"))
	w.Write([]byte(".down {\n"))
	w.Write([]byte("  transform: rotate(45deg);\n"))
	w.Write([]byte("  -webkit-transform: rotate(45deg);\n"))
	w.Write([]byte("}\n\n"))

	w.Write([]byte(".triangle {\n"))
	w.Write([]byte("width: 74;\n"))
	w.Write([]byte("height: 64;\n"))
	w.Write([]byte("border-left:   725px solid transparent;\n"))
	w.Write([]byte("border-right:  725px solid transparent;\n"))
	w.Write([]byte("border-bottom: 725px solid #555;\n"))
	w.Write([]byte("}\n"))
	w.Write([]byte("</style>\n"))
	w.Write([]byte("</head>\n"))
	w.Write([]byte("<body>\n\n"))

	/*
		w.Write([]byte(".triangle {\n"))
		w.Write([]byte("width: 74;\n"))
		w.Write([]byte("height: 64;\n"))
		w.Write([]byte("border: solid 750px;\n"))
		w.Write([]byte("border-color: transparent transparent gray transparent; \n"))
		w.Write([]byte("}\n"))
		w.Write([]byte("</style>\n"))
		w.Write([]byte("</head>\n"))
		w.Write([]byte("<body>\n\n"))
	*/

}

func sendEarlyHtml_fd(w *os.File) {
	w.Write([]byte("<!DOCTYPE html>\n"))
	w.Write([]byte("<html>\n"))
	w.Write([]byte("<head>\n"))
	w.Write([]byte("<!-- Global site tag (gtag.js) - Google Analytics -->\n<script async src=\"https://www.googletagmanager.com/gtag/js?id=G-9BHPJN08TF\"></script>\n<script>window.dataLayer = window.dataLayer || []; function gtag(){dataLayer.push(arguments);} gtag('js', new Date()); gtag('config', 'G-9BHPJN08TF'); </script>\n"))
	w.Write([]byte("<script data-ad-client=\"ca-pub-5646233493379677\" async src=\"https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js\"></script>\n"))
	w.Write([]byte("<meta charset=\"utf-8\"\">\n"))
	w.Write([]byte("<script src=\"https://ajax.googleapis.com/ajax/libs/jquery/3.5.1/jquery.min.js\">\n"))
	w.Write([]byte("</script>\n"))
	w.Write([]byte("<style>\n"))

	w.Write([]byte(".block {\ndisplay: block;\n width: 100%;\n border: none;\n background-color: #4CAF50;\n"))
	w.Write([]byte("color: white;\npadding: 14px 28px;\nfont-size: 16px;\ncursor: pointer;\ntext-align: center;\n}"))
	w.Write([]byte("block:hover {\nbackground-color: #ddd;color: black;\n"))
	w.Write([]byte(" } \n\n"))

	w.Write([]byte(".arrow {\n"))

	w.Write([]byte("  border: solid black;\n"))
	w.Write([]byte("  border-width: 0 3px 3px 0;\n"))
	w.Write([]byte("  display: inline-block;\n"))
	w.Write([]byte("  padding: 3px;\n"))
	w.Write([]byte("}\n\n"))
	w.Write([]byte(".right {\n"))
	w.Write([]byte("  transform: rotate(-45deg);\n"))
	w.Write([]byte("  -webkit-transform: rotate(-45deg);\n"))
	w.Write([]byte("}\n\n"))
	w.Write([]byte(".left {\n"))
	w.Write([]byte("  transform: rotate(135deg);\n"))
	w.Write([]byte("-webkit-transform: rotate(135deg);\n"))
	w.Write([]byte("}\n\n"))
	w.Write([]byte(".up {\n"))
	w.Write([]byte("  transform: rotate(-135deg);\n"))
	w.Write([]byte("  -webkit-transform: rotate(-135deg);\n"))
	w.Write([]byte("}\n\n"))
	w.Write([]byte(".down {\n"))
	w.Write([]byte("  transform: rotate(45deg);\n"))
	w.Write([]byte("  -webkit-transform: rotate(45deg);\n"))
	w.Write([]byte("}\n\n"))

	w.Write([]byte(".triangle {\n"))
	w.Write([]byte("width: 74;\n"))
	w.Write([]byte("height: 64;\n"))
	w.Write([]byte("border-left:   725px solid transparent;\n"))
	w.Write([]byte("border-right:  725px solid transparent;\n"))
	w.Write([]byte("border-bottom: 725px solid #555;\n"))
	w.Write([]byte("}\n"))
	w.Write([]byte("</style>\n"))
	w.Write([]byte("</head>\n"))
	w.Write([]byte("<body>\n\n"))

	/*
		w.Write([]byte(".triangle {\n"))
		w.Write([]byte("width: 74;\n"))
		w.Write([]byte("height: 64;\n"))
		w.Write([]byte("border: solid 750px;\n"))
		w.Write([]byte("border-color: transparent transparent gray transparent; \n"))
		w.Write([]byte("}\n"))
		w.Write([]byte("</style>\n"))
		w.Write([]byte("</head>\n"))
		w.Write([]byte("<body>\n\n"))
	*/

}

func main() {
	//initialize a bunch of hash tables and string lists
	init_hebrewlanguage()
	init_hebrewlanguage_ordinal()
	init_mispar_shemi()
	init_greeklanguage()
	//init_greeklanguage_ordinal() //TODO
	init_periodictable()
	init_english_gematria()
	init_english_alphabet()

	//roman_gematria_666()

	//print primes between 2 and 400 into a list for use later in the program
	primes := sieveOfEratosthenes(400)
	_ = primes

	//initialize semiprimes array to -1
	semiprimes := make([]int, 500)
	for i := 0; i < 500; i++ {
		semiprimes[i] = -1
	}

	//find semiprimes up to 500 for use later in the program
	cnt := 0
	for i := 0; i < 500; i++ {
		if semiprime(i) == true {
			semiprimes[cnt] = i
			cnt = cnt + 1
		}
	}

	jesushebrew := "חישםה עשוהי" //Jesus Christ in Hebrew
	jesushebrewcount := compute_verse_hebrew(jesushebrew, hm)

	jesusgreek := "Ιησούς Χριστός" //Jesus Christ in Greek
	jesusgreekcount := compute_verse_greek(jesusgreek)

	//satanhebrew := "שטן" //Satan in Hebrew = 359, 359th day of year is Christmas
	//satanhebrewcount := compute_verse_hebrew(satanhebrew, hm)

	inums := importantnumbers{jesushebrewcount, jesusgreekcount}

	//get the sorted list of hash keys
	hebrewkeys, _ := retrieve_hash_map(hm)
	greekkeys, _ := retrieve_hash_map(gm)

	hi := make([]LanguageInformation, 0)
	for _, key := range hebrewkeys {
		integerval := fmt.Sprintf("%d", hm[key])
		hi = append(hi, LanguageInformation{key, integerval})
	}

	gi := make([]LanguageInformation, 0)
	for _, key := range greekkeys {
		integerval := fmt.Sprintf("%d", gm[key])
		gi = append(gi, LanguageInformation{key, integerval})
	}

	data := BibleData{
		PageTitle: "Biblical Numerics",
		//VerseList: v,
		JesusHebrew:      jesushebrew,
		JesusHebrewCount: jesushebrewcount,
		HebrewInfo:       hi,
		JesusGreek:       jesusgreek,
		JesusGreekCount:  jesusgreekcount,
		GreekInfo:        gi,
		AwesomeMath:      nil,
		AwesomeChemistry: nil,
		GematriaJesus:    nil,
		GematriaEvil:     nil,
	}

	DOCS := "docs"

	err := os.MkdirAll(DOCS, 0755)
	if err != nil {
		fmt.Println("error creating" + DOCS + "directory")
		return
	}

	/////////////////////////////////////////////////////
	jesushebrewfd, err := os.OpenFile(DOCS+"/jesushebrew.html", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("error opening jesushebrew.html")
		return
	}

	jesushebrewtmpl := template.Must(template.ParseFiles("jesushebrew.html"))
	data.PageTitle = "Jesus Christ in Hebrew"
	sendEarlyHtml_fd(jesushebrewfd) //Write any early html like <head> and <style>
	err = jesushebrewtmpl.Execute(jesushebrewfd, data)
	if err != nil {
		fmt.Println("jesushebrew.html template failed.")
		return
	}

	/////////////////////////////////////////////////////

	jesusgreekfd, err := os.OpenFile(DOCS+"/jesusgreek.html", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("error opening jesusgreek.html")
		return
	}

	jesusgreektmpl := template.Must(template.ParseFiles("jesusgreek.html"))
	data.PageTitle = "Jesus Christ in Greek"
	sendEarlyHtml_fd(jesusgreekfd) //Write any early html like <head> and <style>
	err = jesusgreektmpl.Execute(jesusgreekfd, data)
	if err != nil {
		fmt.Println("jesusgreek.html template failed.")
		return
	}

	/////////////////////////////////////////////////////

	awesomemath0fd, err := os.OpenFile(DOCS+"/awesomemath0.html", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("error opening jesusgreek.html")
		return
	}

	awesomemath0tmpl := template.Must(template.ParseFiles("awesomemath0.html"))
	data.PageTitle = "Some of the Awesome Math in Genesis 1:1"
	awesomehtmllist := doawesomemath(primes, semiprimes, inums)
	data.AwesomeMath = awesomehtmllist
	sendEarlyHtml_fd(awesomemath0fd) //Write any early html like <head> and <style>
	err = awesomemath0tmpl.Execute(awesomemath0fd, data)
	if err != nil {
		fmt.Println("jesusgreek.html template failed.")
		return
	}

	/////////////////////////////////////////////////////

	awesomemath1fd, err := os.OpenFile(DOCS+"/awesomemath1.html", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("error opening awesomemath1.html")
		return
	}

	awesomemath1tmpl := template.Must(template.ParseFiles("awesomemath1.html"))
	data.PageTitle = "Awesome Math with Jesus."
	html := make([]AwesomeMathInformation, 0)

	awesomechemistrylist := doawesomechemistrymath(inums)
	for _, entry := range awesomechemistrylist {
		html = append(html, AwesomeMathInformation{entry.Html})
	}

	jesusthenazarene := "ιησους ο ναζωραιος"
	jesusnazarenecnt := compute_verse_greek(jesusthenazarene)
	satangreek := "ὁ καλούμενος Διάβολος καὶ Ὁ Σατανᾶς"
	satangreekcnt := compute_verse_greek(satangreek)

	str := ""
	const prec = 370 //this asks for more precision than needed
	a, _ := new(big.Float).SetPrec(prec).SetString("1.0")
	b, _ := new(big.Float).SetPrec(prec).SetString("754.0")
	result := new(big.Float).SetPrec(prec).Quo(a, b) //Quo == Divide

	//compute the digital root or digital sum of the first 360 digits
	floatstr := result.Text('g', 83)[2:]
	//fmt.Println(str)
	jesusgreek_circle_total := 0
	for _, digit := range floatstr {
		dig, _ := strconv.Atoi(string(digit))
		jesusgreek_circle_total = jesusgreek_circle_total + dig
	}

	str = fmt.Sprintf("The letters in 'Jesus Christ' when written in Hebrew add up to 754.  The continuously repeating 83 digits found by dividing 1/754 adds up to: %d or the number of degrees in a circle.", jesusgreek_circle_total)
	html = append(html, AwesomeMathInformation{str})

	str = fmt.Sprintf("Observe that a CIRCLE with a circumference of 2368 has a diameter of 754, the ratio of which is %.2F or Pi.\n", float64(inums.jesusgreekcount)/float64(inums.jesushebrewcount))
	html = append(html, AwesomeMathInformation{str})
	str = fmt.Sprintf("The New Testament was originally written in Koine Greek.")
	html = append(html, AwesomeMathInformation{str})
	str = fmt.Sprintf("The word 'JESUS' when represented in Koine Greek add up to 888. JESUS is a perfect multiple of 37. E.g. 24x37=888")
	html = append(html, AwesomeMathInformation{str})
	str = fmt.Sprintf("The word 'CHRIST' when represented in Koine Greek add up to 1480.  1480 is also a perfect multiple of 37.  E.g. 40x37=1480")
	html = append(html, AwesomeMathInformation{str})
	str = fmt.Sprintf("8 + 8 + 8 + (1 + 4 + 8 + 0) = 37")
	html = append(html, AwesomeMathInformation{str})

	part1 := semiprimes[37-1] + semiprimes[73-1]
	if inums.jesusgreekcount+part1 == 37*73 {
		str = fmt.Sprintf("%d + (37th SemiPrime: %d + 73rd SemiPrime: %d) = (37x73) = %d.", inums.jesusgreekcount, semiprimes[37-1], semiprimes[73-1], 37*73)
		html = append(html, AwesomeMathInformation{str})
	}

	//subtract 1 from index on primes because it is zero based
	part1 = inums.jesusgreekcount - (primes[37-1] + primes[73-1])
	part2 := math.Pow(3, 7) - math.Pow(7, 3)
	if part1 == int(part2) {
		str = fmt.Sprintf("2368 - (37th prime: %d 73rd prime: %d) = 3^7 - 7^3 = %d. 1844 is a perfect multiple of 37.  E.g. 37x49=1844.", primes[37-1], primes[73-1], 1844)
		html = append(html, AwesomeMathInformation{str})
		//fmt.Println("fact -", "3^7 - 7^3 == ", "2368 - ((37th Prime:", primes[37-1], ") + (73rd Prime:", primes[73-1], ")) == ", part2)
	}

	str = fmt.Sprintf("'Jesus Of Nazareth The King Of The Jews' from John 19:19 in Greek is %s. Greek character count is %d. 2197=13x13x13", jesusthenazarene, jesusnazarenecnt)
	html = append(html, AwesomeMathInformation{str})
	str = fmt.Sprintf("'who is called the devil and Satan' from Rev 12:9 in Greek is '%s'.  Greek character count is %d.  2197=13x13x13", satangreek, satangreekcnt)
	html = append(html, AwesomeMathInformation{str})
	data.AwesomeMath = html

	sendEarlyHtml_fd(awesomemath1fd) //Write any early html like <head> and <style>
	err = awesomemath1tmpl.Execute(awesomemath1fd, data)
	if err != nil {
		fmt.Println("awesomemath1.html template failed.")
		return
	}

	/////////////////////////////////////////////////////

	awesomemath2fd, err := os.OpenFile(DOCS+"/awesomemath2.html", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("error opening awesomemath2.html")
		return
	}

	awesomemath2tmpl := template.Must(template.ParseFiles("awesomemath2.html"))
	data.PageTitle = "Awesome Math with Fibonacci, Pi, and the non-evil version of 666."
	sendEarlyHtml_fd(awesomemath2fd) //Write any early html like <head> and <style>
	html = make([]AwesomeMathInformation, 0)

	fiblist := Fibonacci()
	pilist := ProofInThePi()
	for _, fib := range fiblist {
		html = append(html, fib)
	}

	for _, pientry := range pilist {
		html = append(html, pientry)
	}

	data.AwesomeMath = html
	awesomemath2tmpl.Execute(awesomemath2fd, data)

	/////////////////////////////////////////////////////

	gematriafd, err := os.OpenFile(DOCS+"/gematria.html", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("error opening gematria.html")
		return
	}

	gematriatmpl := template.Must(template.ParseFiles("gematria.html"))
	data.PageTitle = "English Gematria of Words - based on the Sumerian code"
	sendEarlyHtml_fd(gematriafd) //Write any early html like <head> and <style>
	data.GematriaJesus = english_gematria_jesus()
	data.GematriaEvil = english_gematria_666()
	gematriatmpl.Execute(gematriafd, data)

	/////////////////////////////////////////////////////

	referencesfd, err := os.OpenFile(DOCS+"/refs.html", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("error opening references.html")
		return
	}

	referencestmpl := template.Must(template.ParseFiles("refs.html"))
	data.PageTitle = "References:"
	sendEarlyHtml_fd(referencesfd) //Write any early html like <head> and <style>
	referencestmpl.Execute(referencesfd, data)

	/////////////////////////////////////////////////////

	indexfd, err := os.OpenFile(DOCS+"/index.html", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("error opening index.html")
		return
	}

	indextmpl := template.Must(template.ParseFiles("index.html"))
	data.PageTitle = "A Survey of Biblical Numerics and Mathematical Monotheism (Bible Codes)"
	sendEarlyHtml_fd(indexfd) //Write any early html like <head> and <style>
	indextmpl.Execute(indexfd, data)

	/////////////////////////////////////////////////////

	notesfd, err := os.OpenFile(DOCS+"/notes.html", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("error opening notes.html")
		return
	}

	notestmpl := template.Must(template.ParseFiles("notes.html"))
	data.PageTitle = "Unverified (by biblecodes.go) Bible Numerics Notes."
	sendEarlyHtml_fd(notesfd) //Write any early html like <head> and <style>
	notestmpl.Execute(notesfd, data)

	/////////////////////////////////////////////////////

}
