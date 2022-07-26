package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func problempuller(fileName string) ([]problem, error) {
	//read all the problems from csv file
	//1.open the file
	if fObj, err := os.Open(fileName); err == nil {
		//2.create a new reader
		csvR := csv.NewReader(fObj)
		//3.read the file
		if cLines, err := csvR.ReadAll(); err == nil {
			//4.call the parseproblem function
			return parseproblem(cLines), nil
		} else {
			return nil, fmt.Errorf("error in reading data in csv"+"format from %s file; %s", fileName, err.Error())
		}
	} else {
		return nil, fmt.Errorf("error in opening the %s file;%s", fileName, err.Error())
	}

}

func main() {
	//1.input the name of the file
	fName := flag.String("f", "quiz.csv", "path for csv file")
	//2.set the dduration of the timer
	timer := flag.Int("t", 30, "timer for the quiz")
	flag.Parse()
	//3.pull the problems from the file(calling our problem puller func)
	problems, err := problempuller(*fName)
	//4.hamdle the error
	if err != nil {
		exit(fmt.Sprintf("something went wrong:%s", err.Error()))
	}
	//5.create a variable to count the correct answers
	correctans := 0
	//6.using the duration of the timer, we want to initialise the timer
	tObj := time.NewTimer(time.Duration(*timer) * time.Second)
	ansc := make(chan string)
	//7.loop through the problems, print the questions, we'll accept the answers
problemloop:

	for i, p := range problems {
		var answer string
		fmt.Printf("Problem %d: %s=", i+1, p.q)

		go func() {
			fmt.Scanf("&s", &answer)
			ansc <- answer
		}()
		select {
		case <-tObj.C:
			fmt.Println()
			break problemloop
		case iAns := <-ansc:
			if iAns == p.a {
				correctans++
			}
			if i == len(problems)-1 {
				close(ansc)
			}
		}
	}
	//8.we'll calculate and print the result
	fmt.Printf("Your result is %d out of %d\n", correctans, len(problems))
	fmt.Printf("P		ress enter to Exit")
	<-ansc
}

func parseproblem(lines [][]string) []problem {
	r := make([]problem, len(lines))
	for i := 0; i < len(lines); i++ {
		r[i] = problem{q: lines[i][0], a: lines[i][1]}
	}
	return r
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
