package main 

// a hello world program that does just a little bit more so we can see what we have done.

import (
	"fmt"
	"os"
	"os/exec"
	"encoding/json"
	"sort"

)

func examine_execution_context_return_jsontxt() string {
	//var ret map[string]interface{}

	ret := map[string]interface{} {} 

	wd, err := os.Getwd()
	if nil != err {
		panic(err)
	}

	ret["pwd"] = string(wd)

	ret["args"] = os.Args

	cmd := exec.Command("uname","-a")

	stdouterr , cmderr := cmd.CombinedOutput()

	if nil != cmderr{
		panic(cmderr)
	}

	ret["uname -a"] = string(stdouterr)

	s := os.Environ()

	sort.Strings(s)

	ret["environ"] = s


	byteret , err :=  json.Marshal(ret)

	if nil != err {
		panic(err)
	}

	return string(byteret)
}

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

	fmt.Print("\nJSONBLOB\n")
	fmt.Print(examine_execution_context_return_jsontxt())
	fmt.Print("\n")

	// dump the environment? 

	fmt.Print("Good bye World! Everything ended well enough.\n")
}
