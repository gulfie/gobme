package main 


import ( 
	"testing"
	"os"
	"os/exec"
	"io/ioutil"
	"flag"
	"regexp"
	//"errors"
	"fmt"
	"bytes"
	"encoding/json"
	"reflect"
	//"math"
)

// flags! : https://siongui.github.io/2017/04/28/command-line-argument-in-golang-test/
 
func flagboolreturnstub() bool{
	flag.BoolVar(&debugflg, "debug" , false , "Show me what's going on in the tests")
	return true; 
}

var waste = flagboolreturnstub()


 
func TestTesting(t *testing.T){
}
/* // ordering is persurved. 
func TestATesting(t *testing.T){
}
func TestBTesting(t *testing.T){
}
func TestCTesting(t *testing.T){
}

func TestZTesting(t *testing.T){
}
func TestYTesting(t *testing.T){
}
func TestXTesting(t *testing.T){
}
*/

// technically this maybe done already by the go, though we'll want to 
// is there atest ordering?
func Test_simple_building(t *testing.T){


	cmd := exec.Command("go","clean"); // assume clean GOOS GOARCH

	stdouterr, err := cmd.CombinedOutput();

	if nil != err { 
		t.Error(string(stdouterr))
		t.Fail()
	}

	d("it cleaned:",string(stdouterr),"\n")




	cmd = exec.Command("go","build"); // assume clean GOOS GOARCH

	stdouterr, err = cmd.CombinedOutput();

	if nil != err { 
		t.Error(string(stdouterr))
		t.Fail()
	}

	d("it built:",string(stdouterr),"\n")

}


func Test_simple_building_hello(t *testing.T){

	// test hardness may do this.
/*
	origwd , wderr := os.Getwd()

	if nil != wderr {
		t.Error("Can not get working directory");
		t.Fail()
	}
	defer os.Chdir(origwd)
*/	

	os.Chdir("hello_wide_world")

	cmd := exec.Command("go","clean"); // assume clean GOOS GOARCH

	stdouterr, err := cmd.CombinedOutput();

	if nil != err { 
		t.Error(string(stdouterr))
		t.Fail()
	}

	d("it cleaned:",string(stdouterr),"\n")



	cmd = exec.Command("go","build"); // assume Gclean GOOS GOARCH

	stdouterr, err = cmd.CombinedOutput();

	if nil != err { 
		t.Error(string(stdouterr))
		t.Fail()
	}

	d("it built:",string(stdouterr),"\n")

}


// plan, execute the native hello
func Test_simple_use_hello(t *testing.T){

	// test hardness may do this.
/*	origwd , wderr := os.Getwd()

	if nil != wderr {
		t.Error("Can not get working directory");
		t.Fail()
	}
	defer os.Chdir(origwd)
*/
	os.Chdir("hello_wide_world")

	cmd := exec.Command("./hello_wide_world")

	stdouterr,err := cmd.CombinedOutput()


	if nil != err {
		t.Error("can't run native hello_wide_world (" + (fmt.Sprintf("\ntxt:%v\nerr:%v",stdouterr, err))+")\n")
		t.Fail()
	}

	d(string(stdouterr))

}



// plan, execute the native hello
func Test_simple_gobme_of_hww(t *testing.T){

	// test hardness may do this.
/*	origwd , wderr := os.Getwd()

	if nil != wderr {
		t.Error("Can not get working directory");
		t.Fail()
	}
	defer os.Chdir(origwd)
*/
	os.Chdir("hello_wide_world")

	previousDirSlice, _  := ioutil.ReadDir("./")


	// this restricts gobme testing to certan goos/goarch, for speed. 
	//cmd := exec.Command("../gobme" ,"-debug", "-goos","linux,darwin,freebsd", "-goarch","amd64")
	cmd := exec.Command("../gobme" ,"-debug", "-goos","linux", "-goarch","amd64")

	stdouterr,err := cmd.CombinedOutput()


	if nil != err {
		t.Error("can't run native hello_wide_world (" + (fmt.Sprintf("\ntxt:%v\nerr:%v",stdouterr, err))+")\n")
		t.Fail()
	}

	d(string(stdouterr))

	// is there junk in this directory? 

	afterDirSlice, _ := ioutil.ReadDir("./")
	isdiff := false 

	if len(previousDirSlice) != len(afterDirSlice){
		isdiff = true
	}


	if isdiff {
		t.Error("The files before and after are too different\n" + fmt.Sprint("\nbefore :" , previousDirSlice, "\nafter: " , afterDirSlice ,"\n") +"\n")
	}

	for _ , fn := range( []string{"hello_wide_world.unpacked" ,"makeself.sh", "makeself-header.sh" } ) {
	//	d("is fn here? ", fn , "\n");
		if _ , err := os.Stat(fn) ; nil == err  { 
			t.Error("left a file ("+fn+") laying around\n")
		}
	}


	// is there a new executable shellish thing here? 

	filebytes , fberr := ioutil.ReadFile("hello_wide_world")

	if nil != fberr {
		t.Error("Unable to open resulting archive file...("+fmt.Sprint(fberr)+")\n")
	}else{

		findmsmagic := regexp.MustCompile("\\A#!/bin/sh\n# This script was generated using Makeself")

		found := findmsmagic.Find(filebytes)
		if nil == found { 
			t.Error("resulting archive file does not pass type check\n")
			t.Fail()
		}else{
			d("Found the magic in the archive file\n")
		}
	}

}

type execrecord struct {
	stderr,stdout []byte
	runerr error
	reportedwd string
	jsonblob map[string]interface{}
}


// exec style env strings. nil == what you have already 
// convert cmdstr to  a []string
func runAndCapture(cmdstr []string, env []string) (er execrecord){

	cmd := exec.Command(cmdstr[0], cmdstr[1:]...)

	var stderr ,stdout bytes.Buffer 
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if nil != env {

		cmd.Env = make([]string, len(env))

		for _ , v := range(env){
			d("copying over env v (",v,")\n")
			cmd.Env = append(cmd.Env,v)
		}
	}

	err  := cmd.Run()
	if nil != err {
		d("cmd.run found an error:(",err,")\n")
	}

	er.runerr = err 

	er.stdout = stdout.Bytes()
	er.stderr = stderr.Bytes()


	findwd := regexp.MustCompile("\nwd : (\\S+)")
	ret := findwd.Find(stdout.Bytes())
	if ret != nil {
		er.reportedwd = string(ret)
	}

	findjson := regexp.MustCompile("\nJSONBLOB\n(.*?)\n")

	jret := findjson.FindSubmatch(stdout.Bytes())


	if nil != jret { 
		jsonbytes := jret[1]
	//	d("json bytes (string) : " , string(jsonbytes) ,"\n")
		err := json.Unmarshal(jsonbytes, &er.jsonblob) 

		if nil != err { 
			panic(err)
		}
	//	d("json unstringed : " , er.jsonblob ,"\n")
	}

	d("er out : " , er , "\n\n")
	return er
}


// env0 := envslicetomapss(er0.jsonblob["environ"])

// stringslice0 :=  interfaceslicetostringslice(er0.jsonblob["environ"].([]interface{})) 
func interfaceslicetostringslice( inslice []interface{} ) ( outslice []string ){ 

	// there must be a better way. FIX TODO, research Go idiums
	//outslice = make([]string , len(inslice.([]interface{})))

//	d(" inslice ", inslice , "\n")


	for i := range(inslice){
		outslice = append(outslice, inslice[i].(string))
	}

//	d(" outslice ", outslice ,"\n")

	return outslice
}


var equalslicer = regexp.MustCompile("\\A([^=]+)=(.*)$")

func envslicetomapss( env []string ) (ret map[string]string ) {

	// d("env env : " , env ,"\n\n\n")
	ret = make(map[string]string)

	// does the map show up magicly?
	for _ , v := range(env) {
		d(".(",v,")")
		parts := equalslicer.FindSubmatch([]byte(v))
		if nil != parts {
			// note, some input avalidation.. some.	
			ret[string(parts[1])] = string(parts[2])
		}
	}

	// d("\nenv slice ret " , ret,"\n")

	return ret
}


func mapsstoenvslice(envmap  map[string]string ) ( environ []string ){

	for k := range(envmap) { 
		environ = append(environ , k + "=" + envmap[k] )
	}

	return environ
}


//
// 	we only care ( or don't care) about some keys, so pay atention to them, or igrnore them.
//
//
//	isdiff = comparemapsonkeys( []string { }, false , env0 , evn1 )
func comparemapsonkeys( keyset []string, onlycheck bool, left map[string]string , right map[string]string ) (ismateriallydiff bool ){ 

	// d("comparemapsonkeys\n",left, right ,"\n\n")
	var keystolookat []string 

	if onlycheck{
		d(" onlycheck\n")
		keystolookat = keyset
	}else{
		d("everything but..\n")
		// union of keys from each, minus the keyset 
		// https://stackoverflow.com/questions/21362950/golang-getting-a-slice-of-keys-from-a-map
		kmap := make(map[string]bool)
		ksmap := make(map[string]bool)

		for _ , k := range  keyset {
	//		d("k1 (",k,")\n")
			ksmap[k] = true
		}

		for k:= range left {
	//		d("k2 (",k,")\n")
			kmap[k] = true
		}
		for k:= range right {
	//		d("k3 (",k,")\n")
			kmap[k] = true
		}

		for k:= range kmap {
	//		d("k4 (",k,")\n")
			if ! ksmap[k] { 
				keystolookat = append(keystolookat, k)
			}
		}
	}

	// d("keys chosen\n", keystolookat ,"\n\n")

	for _ , k := range keystolookat {
		//l , r := left[k] , right[k] 
		l, lok := left[k]
		r, rok := right[k]
		d("checking : k (",k,") => (",l ," , " ,r,")\n")
		if lok  { 
			if rok { 
				if l != r { 
					ismateriallydiff = true	
					d("diff at k(",k,") are different (",l," vs " ,r,")\n") 
				}
			}else{
				ismateriallydiff = true
				d("key (",k,") not present on the right side\n")
				// break

			}
		}else{
			ismateriallydiff = true
			d("key (",k,") not present on the left side\n")
			// break
		}
	}


	return ismateriallydiff
}




func testcompareTwoJSONEr( t *testing.T, testprefix string,  left , right execrecord)  ( isgood bool) { 

	// check args
	// next check the environments. 

	// Do the same thing again, 
	// but for executions that use relative PATH variables. 

	//	d("json unstringed : " , er.jsonblob ,"\n")
	isdiff := false 

	// should check for minimum keyset on er 
	// reading json is a pita. 

	d("working on (",testprefix,")\n")


	if _ , exists := left.jsonblob["environ"] ; ! exists  { 
		t.Error("environ is not present in left\n")
		isgood = false;
	}

	if _ , exists := right.jsonblob["environ"] ; ! exists  { 
		t.Error("environ is not present in right\n")
		isgood = false;
	}

	// there must be a better way. FIX TODO, research Go idiums

	stringslice0 := make([]string , len(left.jsonblob["environ"].([]interface{})))

	stringslice0 =  interfaceslicetostringslice(left.jsonblob["environ"].([]interface{})) 

	env0 := envslicetomapss(stringslice0)

	stringslice1 := make([]string , len(right.jsonblob["environ"].([]interface{})))

	stringslice1 =  interfaceslicetostringslice(right.jsonblob["environ"].([]interface{})) 

	env1 := envslicetomapss(stringslice1)

	isdiff = comparemapsonkeys( []string { }, false , env0 , env1 )

	d("isdiff (" , isdiff , ")\n")

	if isdiff { 
		t.Error("environment sets seem materially different, run with -args -debug to find out")
		isgood = false;
	}

	// check argv really quick.
	

	if ! reflect.DeepEqual( left.jsonblob["args"].([]interface{})  , right.jsonblob["args"].([]interface{}) ){
		t.Error("argument strings to the two invocations are different , see -args -debug \n")
		isgood = false;
	}

	return isgood;
}





// 
func Test_simple_gobme_of_hww_local_usage(t *testing.T){

	os.Chdir("hello_wide_world")

	// should have our hello_wide_world archive built. 

	e := os.MkdirAll("tmp/tmp2/tmp3/tmp4/tmp5/tmp6", 0755 )
	if nil != e {
		t.Error("failed to make tmpdir stack")
		t.Fail()
	}
	if ! leavetmp {
	//	defer os.RemoveAll("tmp")
	}


	os.Chdir("tmp/tmp2/tmp3")
	defer  os.Chdir("../../../") // so the Remove all works.

	cmd := exec.Command("cp" , "../../../hello_wide_world", "tmp4/")
	stdouterr , err := cmd.CombinedOutput()

	if nil != err{
		t.Error("Unable to copy archive file into test space", stdouterr, err)
		t.Fail()
	}

	//  pause and think for a  bit. 

	// oh neat... makeself helpfully pops up a window with startup script output (presumably when there is no terminal ) 
	er0 := run_and_capture([]string{"tmp4/hello_wide_world"},nil)

	if nil != er0.runerr {
		t.Error("failure to run hello_wide_world from a differnet direcotry\n", er0.runerr )
		t.Fail()
	}

	er1 := run_and_capture([]string{"tmp4/hello_wide_world"},nil)

	// we'll only want to compare some of the outputs. 	
		d("er0.wd", string(er0.reportedwd))
		d("\n")
		d(string(er0.stdout))
		d("\n")
		d(string(er0.stderr))
		d("\ner1.wd", string(er1.reportedwd))
		d("\n")
		d(string(er1.stdout))
		d("\n")
		d(string(er1.stderr))

	// possibly redundant due to shell env PWD 
	if er0.reportedwd != er1.reportedwd { 
		d("er0.wd", string(er0.reportedwd))
		d("\n")
		d(string(er0.stdout))
		d("\ner1.wd", string(er1.reportedwd))
		d("\n")
		d(string(er1.stdout))
		t.Error("first and second call to hww give different wd\n" )
	}


	testcompare_two_json_er( t , "calling up a directory",  er0 , er1  ) 


	os.Chdir("tmp4/tmp5")
	defer os.Chdir ("../../")

	cmd = exec.Command("cp" ,"../../../../../hello_wide_world" ,".")
	stdouterr , err = cmd.CombinedOutput()

	if nil != err{
		t.Error(err)
	}

	d("doing er2\n")
	er2 := run_and_capture([]string{"./hello_wide_world"},nil)

	d("doing er3\n")
	er3 := run_and_capture([]string{"./hello_wide_world"},nil)

	testcompare_two_json_er( t, "relative path same dir", er2 ,er3 ) 


	os.Chdir("tmp6")
	defer os.Chdir("../")

	cmd = exec.Command("cp" ,"../../../../../../hello_wide_world" ,".")
	stdouterr , err = cmd.CombinedOutput()

	if nil != err{
		t.Error(err)
		t.Fail()
	}

/*

	d("\n\n\n\n\nWHAT\n\n\n\n\n")

	cmd = exec.Command("pwd")
	stdouterr, err = cmd.CombinedOutput()
	fmt.Print("pwd : " , string(stdouterr))



	cmd = exec.Command("ls","-alrt")
	stdouterr, err = cmd.CombinedOutput()
	fmt.Print("ls -alrt : " , string(stdouterr))


	cmd = exec.Command("bash","-c" ,"export",)
	stdouterr, err = cmd.CombinedOutput()
	fmt.Print("export : " , string(stdouterr))
*/
	//
	// add . as the first thing in the PATH to test non anchorned calls. 
	// 

	// my path , not it's path. 

	origpath , _  := os.LookupEnv("PATH")
	defer os.Setenv("PATH",origpath)
	
	// just one dot, without the / is still valid..
	os.Setenv("PATH","./:"+origpath)

/*
	saved := os.Enviorn() 
	envmap := envslicetomapss(os.Environ())
	d("\nenvmap ",envmap,"\n")

	PATH, isok := envmap["PATH"]

	if isok { 
		envmap["PATH"] = "./:" + PATH +":./"
	}

	//delete(envmap,"PWD")

	env := mapsstoenvslice(envmap)
	d("\nenv : " , env )
	d("\n\n")
*/

/*
	d("what is it?\n")
	ert1 := run_and_capture([]string{"ls","-alrt"}, nil)
	fmt.Print("stdout : " ,string(ert1.stdout),"\n")
	fmt.Print("stderr : " ,string(ert1.stderr),"\n")

	ert2 := run_and_capture([]string{"bash","-c","export"}, nil)
	fmt.Print("stdout : " ,string(ert2.stdout),"\n")
	fmt.Print("stderr : " ,string(ert2.stderr),"\n")

	ert3 := run_and_capture([]string{"head hello_wide_world"}, nil)
	fmt.Print("stdout : " ,string(ert3.stdout),"\n")
	fmt.Print("stderr : " ,string(ert3.stderr),"\n")
*/

	// . fails 
	d("doing er4\n")

	//er4 := run_and_capture("hello_wide_world" , append(os.Environ() , "PATH=./:/bin/:/usr/bin/:/usr/local/bin" ))
	er4 := runAndCapture([]string{"hello_wide_world"}, nil)
//	fmt.Print("stdout : " ,string(er4.stdout),"\n")
//	fmt.Print("stderr : " ,string(er4.stderr),"\n")


	d("doing er5\n")
	//er5 := run_and_capture("hello_wide_world" , env)
	er5 := runAndCapture([]string{"hello_wide_world"} , nil)
//	fmt.Print("stdout : " ,string(er5.stdout),"\n")
//	fmt.Print("stderr : " ,string(er5.stderr),"\n")

	testcompareTwoJSONEr( t, "same pwd unancored via PATH", er4 ,er5 ) 


//  	more along this line would be a good idea. 

// 	



	d("completed\n\n")
}




func Test_over(t *testing.T){
	d("\n\nlast test done\n\n");
}
