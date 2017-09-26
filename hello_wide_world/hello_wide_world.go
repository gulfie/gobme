package main 

// a hello world program that does just a little bit more so we can see what we have done.

import (
	"fmt"
	"os"
	"os/exec"

)


func main(){


	fmt.Print("Hello World!\n")
	fmt.Print("args : \n")
	for i,arg := range(os.Args) {
		fmt.Printf("   %v : (%v)\n",i , arg)
	} 
	fmt.Print("\n")

	wd, err := os.Getwd()
	if nil != err { 
		panic(err)
	}
	
	fmt.Print("wd : ", wd ,"\n\n")

	cmd := exec.Command("uname","-a")
	stdouterr , err := cmd.CombinedOutput()

	if nil != err {
		panic(err)
	}

	fmt.Printf("%s",stdouterr)

	// dump the environment? 

	fmt.Print("Good bye World! Everything ended well enough.\n")
}
