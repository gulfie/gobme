// 
//
package main
// - "go1.9"


/*
	gone through several iterations
		
	this is now going to be the tool to wrap it all together and tie it up with a nice bow.


	TODO : 

		tldr; 
		must : update stub startup script to 
			x a) be quite 
			xb) take side channle commands ( e.g. just install, or install over here )  and ignore archiver commands? 
				
			x c) (startup script?) exec the correct thing out of the goos_goarch tree 
				rename to app_launcher_and_multi_exec_ctl   alamec ? ... something.

			x d) does exec from non local dir work.
			x e) switcheru the goos_goarch/thing into the place of the original .shar  ( move the orignial .shar to maybe something.multi )
			(no) f) optional remove everything besides the switcheru bin.... option default to no.
			(not yet) g) (punt on, but) figure out the case of perms not right for this user to unarchvie it ( $HOME/.go_multi/ maybe?) ... maybe just run the thing... 
		
		b) todo: 
			XXX exec from different directores
				paramiterize a "decopress from samedir as archive" option 

			patches to makeself 
				1) quite in the header should be paramarizably quited in the makeself call like --noprogress 
				2) a paraterizable var of takeheaderargsfromenv (MS_HEADER_ARGS) in makeself and the required mechanism to make that work

			roll back debug copying of makeself and makeself-header 

		c) todo : 
			strip out comments / debug / dev stuff 
			/ harden against diff users and diff pwd 
			x import the script as a literal
		
		d) todo : 
			realtive and path based things seem to work.. there are now enough test cases where we'll need to automate up this thing.... so twisty

			Test plan : ( pre checkin / post checkin tests ) 
			x a) does it compile ( given some set of os/architectures ) 
			x b) does it produce an artifact with the right name (^^)
			x c) ( is there other junk around the build dir?)
			/ d) for a set of different conditions what does the thing do when executed first, second and thrid times? 
				i) tmp dir with copied artifact into it 	
				ii) tmpdir with artifact /  different working dir, higher, lower in the fs tree 
				iii) different perms on the files / directories? 
			e) do the above - compilation, on each different os/arch we can get running simply 
			d) above in light of failure cases, e.g. permissions issues, disk space issues, or chatter? 
			f) using a different test artifact than hello_wide_world, maybe one that includes more external libs. ( does dynamic lnking cross compiling work enough?) 
				
			... okay getting that some ways toward done. 
			



	TODO2: gobme side 
		a) check attempt to start removing dependencies from the start script and the makeself stuf. 
			sed / grep 


		b) cli args, 
			x 0) precheck with a build of the local arch, if it fails, fail quick
				test? 
			x 1) paramiterize the build target goos and goarch into lists /sets 
			/ 2) -dev flag to use filse for makeself
				2.a) debug / verbose ? 
			3) readme.md including license incumbering by gnu 
				recomendations to swap out makeself for a go compatable something ( keep thank you to makeself team, cuz' awesome) 
			/ 4) start the testing? 
				/ a) do some dev in shell, then consider migrating to pure go. 
			x 5) see what's up with the $GOROOT adn GOPATH on disk structures?  how can we fit into that reasonably? 
			x 6) move makeself into the gobme tree 
			x 7) github
			8) hostify? 
			9) deltizer? 
			10) telemetry packages? 				
			11) a load balancer? 
			12) style consistency...
			13) different OS testing. 
				13.a) compiling on FreeBSD 
				13.b) running on OpenBSD 
			
	
	TODO3: 
		Interesting idea... security.txt 
		https://securitytxt.org/ 
		
		add flag to GO_MULTI_ALAMEC_ARGS , to allow for 'just show me the security.txt for this thing' 

	TODO4: 
		more testing 
		TEST the texts of the mshtxt.go mstxt.go and   startupscript.go and compares them to the text on disk for the sources. e.g. ( was the devgenmaketxtgo run? ) 

*/


//
// https://stackoverflow.com/questions/6182369/exec-a-shell-command-in-go
//
// harness up the go os/exec to simplify it ?  and may be to add 
// timestamped buffer parts .... 
//
//  maybe take some measurements from other go routines around that process.   telemetry 
// 
//   maybe a popen3 as well... ( with optional result ) 
//  https://gist.github.com/dialogbox/71177892ac6d3879c563
//
//  humm, error codes are more complicated than I thought.
//  
//  http://tldp.org/LDP/abs/html/exitcodes.html


import ( 
	"io/ioutil"
	"os"
	"os/exec"
	"flag"
	"fmt"
	"strings"
	"path/filepath"
	"regexp"
	// "runtime/debug"
	// "go/build/syslist"
)


var debugflg = false
var leavetmp = false 


//
// we'd want to use the non exported goosList and goarchList... they are however unexported from ( go/src/go/build/syslist.go )
//  rooting around the filesystem seems rude at this point. e.g. ( ./runtime/internal/sys/gengoos.go )
//   

// android nacl , both taken out because it looks like we can't run them.
// atman is in the same boat ( https://github.com/atmanos/atmanos ) 
var fullgooscommalist = "android,darwin,dragonfly,freebsd,linux,nacl,netbsd,openbsd,plan9,solaris,windows,zos"
	
var arggooscommalist string
// [...]string { "darwin" , "dragonfly" , "freebsd" , "linux" , "netbsd" , "openbsd" , "plan9" , "solaris" , "windows" , "zos"  } 
//var goosSlice = [...]string { "dragonfly" , "freebsd" , "linux" }
var goosSlice []string //  = [...]string { "dragonfly" , "linux" }

var fullgoarchcommalist = "386,amd64,amd64p32,arm,armbe,arm64,arm64be,ppc64,ppc64le,mips,mipsle,mips64,mips64le,mips64p32,mips64p32le,ppc,s390,s390x,sparc,sparc64"
var arggoarchcommalist string
// = [...]string { "386" , "amd64" , "amd64p32" , "arm" , "armbe" , "arm64" , "arm64be" , "ppc64" , "ppc64le" , "mips" , "mipsle" , "mips64" , "mips64le" , "mips64p32" , "mips64p32le" , "ppc" , "s390" , "s390x" , "sparc" , "sparc64" }

var  goarchSlice []string //  = [...]string { "386" , "amd64" } 


func getBuildTarget() string{

	fullcwd , err := os.Getwd()
	if nil != err {
		panic(err)
	}

	splitfullcwd := strings.Split(fullcwd,"/")
	//fmt.Print(splitfullcwd)

	buildtarget := splitfullcwd[(len(splitfullcwd) -1)]

	return buildtarget
}



func cleanAndBuild(goos, goarch string) (outerr string,  err error) {

			cleancmd := exec.Command("go" , "clean")
			cleanouterr, cleanerr := cleancmd.CombinedOutput()

			if (nil != cleanerr){
				fmt.Print(cleanouterr)
				return string(cleanouterr), cleanerr
			}

			// fmt.Print("Cleaned\n");

			buildcmd := exec.Command("go", "build")

			//build nil versions in whatever the user thinks is wise, e.g. whatever the current GOOS and GOARCH are in. 
			newenv  := os.Environ()

			if "" != goos{
				newenv = append(newenv , "GOOS="+goos)
			}

			if "" !=  goarch{
				newenv = append(newenv , "GOARCH="+goarch) 
			}

			buildcmd.Env = newenv
			buildouterr, builderr := buildcmd.CombinedOutput()

			return  string(buildouterr) , builderr

}



func cleanAndBuildGoosGoarch(){
	var err error

	// does the build work at all? 
	smokeouterr, smokeerr := cleanAndBuild("","")

	if nil != smokeerr{
		fmt.Print(smokeouterr)
		fmt.Print(smokeerr)
		panic("Unable to do a smoke check simple build, quiting")
	}else{
		fmt.Print("Smoke Build worked... proceeding\n")
	}




	// what is the binname going to be? 

	buildtarget := getBuildTarget()
	unpackeddir := buildtarget + ".unpacked/"

	// proactive cleanup
	if _,filerr := os.Stat(unpackeddir) ; nil == filerr  {
		d("removing unpackeddir\n")
		os.RemoveAll(unpackeddir)
	}


	check(os.MkdirAll(unpackeddir+"goos_goarch",0755))



	for  _ ,  goos := range goosSlice {
		for _ , goarch := range goarchSlice {

			fmt.Printf("Working on %s %s\n", goos , goarch)

			buildouterr, builderr := cleanAndBuild(goos,goarch)

			if  nil == builderr {
				fmt.Print("    Build success!\n"+string(buildouterr)+"\n");

				targetdir := fmt.Sprintf("%s/goos_goarch/%s_%s",unpackeddir, goos,goarch)
				err = os.MkdirAll(targetdir, 0755)

				if nil != err {
					panic(err)
				}

				tmptarget := buildtarget 

				// windows appends .exe to the end. which breaks the simplicity
				_, staterr := os.Stat(tmptarget)

				// generalize? 
				if nil != staterr {
					tmptarget = tmptarget + ".exe"
					_ , staterr := os.Stat(tmptarget)
					if  nil != staterr { 
						fmt.Print("Can not find build result... sorry\n");
						panic(staterr)
					}
				}

				err = os.Rename(tmptarget, fmt.Sprintf("%s/%s",targetdir,tmptarget) )

				if nil != err {
					panic(err)
				}
			}else{
	//			fmt.Print(string(buildouterr))
	//			fmt.Print(builderr)
				fmt.Print(" .. build didn't work moving on\n");
			}
		}
	}



}


// maybe the stack trace is enought to make this worth it. 
// see : https://github.com/go-errors/errors ? 
// 	this whole github based build system... humm.... 
func check(e error){
	if nil != e { 
		fmt.Print("check detected a falure\n",e)
	//	fmt.Print("Stack trace : ", Stack())
		panic(e)
	}
}


// 100% batteries included, thus we portage in our own tools.
// this takes a while to do, should the production of this function be automated

func dropCopyOfMakeself(){
	fmt.Print("Thank you to Stephane Peter,  http://makeself.io/ , et. al see source for license and details (GPL2+!)\n")


	check(ioutil.WriteFile("makeself.sh",[]byte(Mstxt), 0755))
	check(ioutil.WriteFile("makeself-header.sh",[]byte(Mshtxt), 0755))

}

func dropCopyOfStartupScript(prefix string){

	check(ioutil.WriteFile(prefix + "/" + "startupscript",[]byte(Sutxt), 0755))

}



// fmt.Print only when we are logging 
func d (a ...interface{})(){
	if debugflg {
		fmt.Print(a... )
	}	
}


func df (format string , a ...interface{}){
	if debugflg { 
		fmt.Printf(format , a...)
	}
}



// find our bin, then look for the forked_makeself/makeself{,-header}.sh files, 
// note this won't run with 'go run ...'
func genmakeselfgo(){

	// find ourselves : https://stackoverflow.com/questions/18537257/how-to-get-the-directory-of-the-currently-running-file
	ex , err := os.Executable()
	if nil != err { 
		panic(err)
	}

	lookforbackticks := regexp.MustCompile("`")

	// valid for dev, not for nominal runs.
	execdir := filepath.Dir(ex)
	df("found ourselves at (%v)\n",execdir)

	//  srcfilename, destfilename (.go file), variablename 
	relativefn := [][3]string{
		{ "forked_makeself/makeself.sh" , "mstxt.go", "Mstxt" },
		{ "forked_makeself/makeself-header.sh" , "mshtxt.go" , "Mshtxt" },
		{ "gobmeStartupScript.sh" , "startupscript.go", "Sutxt" }   }


	for _ , t := range relativefn { 
		srcfn,dstfn,varname := t[0] , t[1], t[2] 
		// 
		bytetxt , err := ioutil.ReadFile(srcfn)
		if nil != err { 
			panic(err)
		}

		strtxt := string(bytetxt)

		strtxt = lookforbackticks.ReplaceAllString(strtxt,"` + \"`\" + `")

		strtxt = `package main 
// this is generated by go_build_multi_exec from (` +  srcfn + `) do not manually edit
// see -help for more details

// ` + varname +  ` is the exported text of an internal part of the shell archive subsystem
var ` + varname + " = `" + strtxt + "`\n\n"


		err = ioutil.WriteFile(dstfn,[]byte(strtxt), 0755)
		if  nil != err{
			panic (err)
		}
	}



	//for _ , fn :
	d(relativefn)
	

}



func main(){

	// FIX , TODO , make a manual list of mandatory GOOS_GOARCH that must pass? 
	flag.StringVar(&arggooscommalist, "goos", fullgooscommalist, "comma seperated list of GOOS to attempt to compile for, successes are include d in the multi")
	flag.StringVar(&arggoarchcommalist, "goarch", fullgoarchcommalist, "comma seperated list of GOARCH to attempt to compile for, success are include in the multi")
		
	var devgentxtgo = false 

	flag.BoolVar(&devgentxtgo,"devgenmaketxtgo" , false , "an internal thing to allow us to bring the makeself and start scripts into the binary as strings rather than having multiple files laying around, dev use only")

	flag.BoolVar(&debugflg , "debug", false , "more verbose debuggin is enabled, more junk to the stdout/stderr")
	flag.BoolVar(&leavetmp, "leavetmp", false , "a debugging flag to leave open makeself stuff and archive")

	flag.Parse()

	if devgentxtgo { 
		d("generating texts\n")
		genmakeselfgo()	
		d("successful quiting\n");
		return 
	}


	// parse out the comma lists to the slices
	// more error checking? 
	goosSlice = strings.Split(arggooscommalist,",")
	goarchSlice = strings.Split(arggoarchcommalist,",")


	// bahhh, just make a tempdir somewhere..
	buildtarget := getBuildTarget()
	unpackeddir := buildtarget + ".unpacked"


	cleanAndBuildGoosGoarch()
	dropCopyOfMakeself()


	check(os.MkdirAll(unpackeddir,0755))

	dropCopyOfStartupScript(unpackeddir)

	// "./makeself.sh --target 

	// fido is a place holder
	cmd := exec.Command("./makeself.sh" ,"--nox11", "--noprogress" , "--target" , unpackeddir , unpackeddir , buildtarget , "multi architecture self installing go build of " + buildtarget , "./startupscript")
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr 

	err := cmd.Run()

	check(err)

	if ! leavetmp {
		d("removing working files\n")
		os.RemoveAll(unpackeddir)
		os.Remove("makeself.sh")
		os.Remove("makeself-header.sh")
	}else{
		d("leaving working files\n")
	}

	fmt.Print("Successfull completion maybe\n")

}



// for the builder we don't care exquisitely about the failures, only that they are. 

/*

*/

