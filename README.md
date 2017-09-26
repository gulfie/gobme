# gobme
A tool to build go mutli execution archives from inside a source dir.


### TLDR
It's cool to compile some go and run it anywhere... now even more so.  Linux to OpenbSD, arm to amd64.  One artifact to run them all.   CD to a directory where 'go build' works and 'gobme'. 

Then take that and run it almost anywhere. 



### How does one use it? 

```

# $user:  go install "github.com/gulfie/gobme"

# $user:  cd ${some_direcotry_containing_a_main_package} 

# $user:  gobme 

The normal go executable has been traded with one that now runs almost everywhere. 

```

So simple, it's easy. 


### how does it work? 

By tieing together go 1.9's simple abitlity to cross compile with a self extracting posix /bin/sh archive

Extracting on first execution at the target site and switcheruing things around so the archive does not extract each time the program is run. 

Some flags are supported, --help for more details.   For speed I'd recomend "gobme --goos linux,openbsd,dragonfly --goarch 386,amd64,mumford"


## NOTES: 

1. This is a rough proof of concept rather than something that should be used in critical sections.
2. It won't run in a bare init only container, sad that. 
3. There are lots of comments in the source that could point out use cases, interfaces, or liablities.  e.g. there is a dependency on sed and a handfull of other POSIX commands.  ( and windows probabyl doesn't work even a little) 


## License:

Some combination of GPL2+ from makeself ( https://github.com/megastep/makeself ) and artistic / oh please be careful with it I haven't even started serious testing yet. 



## TODO


1. initial release 
2. some testing, both local on single architectures. 
3. version # 
4. some visual use examples. 
5. cleanup the visual appearance of the output
6. document the hidden sidechannel commands
7. possibly migrate processing out of the CWD and move it to a temp build dir somewhere
8. stop leaving tailings around , e.g. *.unpacked and the *.sh tools 


### future possibilties. 

1. Include the build invornment in the gob. In an ideal world someone would be able to crack open the archive and reproduce the artifacts. 
2. worm based autonomic distributed computing. 


## Examples


Shorter simple example.  Installation and use ( asusming all your GOPATH and GOROOT is setup correctly ) 

```

# $user : go install github.com/gulfie/gobme

# $user : cd $THE_GOPATH_PLACE_IT_WAS_INSTALLED_INTO/src/github.com/gulfie/gobme/hello_wide_world/

# $user : ls -alrt 
 total 16
-rw-rw-r-- 1 x x  629 Sep 25 18:08 hello_wide_world.go
-rw-rw-r-- 1 x x   17 Sep 25 18:08 .gitignore
drwxrwxr-x 7 x x 4096 Sep 25 18:24 ..
drwxrwxr-x 2 x x 4096 Sep 25 18:24 .

# $user : gobme --goos linux,openbsd,dragonfly --goarch 386,amd64,mumford

Smoke Build worked... proceeding
Working on linux 386
    Build success!

Working on linux amd64
    Build success!

Working on linux mumford
 .. build didn't work moving on
Working on openbsd 386
    Build success!

Working on openbsd amd64
    Build success!

Working on openbsd mumford
 .. build didn't work moving on
Working on dragonfly 386
 .. build didn't work moving on
Working on dragonfly amd64
    Build success!

Working on dragonfly mumford
 .. build didn't work moving on
Thank you to Stephane Peter,  http://makeself.io/ , et. al see source for license and details (GPL2+!)
Header is 581 lines long

About to compress 9912 KB of data...
Adding files to archive named "hello_wide_world"...
./
./goos_goarch/
./goos_goarch/dragonfly_amd64/
./goos_goarch/dragonfly_amd64/hello_wide_world
./goos_goarch/linux_amd64/
./goos_goarch/linux_amd64/hello_wide_world
./goos_goarch/openbsd_386/
./goos_goarch/openbsd_386/hello_wide_world
./goos_goarch/openbsd_amd64/
./goos_goarch/openbsd_amd64/hello_wide_world
./goos_goarch/linux_386/
./goos_goarch/linux_386/hello_wide_world
./startupscript
CRC: 3660133747
MD5: 7fa9732b57111dd70c72f92ff4f9f142

Self-extractable archive "hello_wide_world" successfully created.
Successfull completion maybe

# $user : ls -arlt 
total 3720
-rw-rw-r-- 1 x x     629 Sep 25 18:08 hello_wide_world.go
-rw-rw-r-- 1 x x      17 Sep 25 18:08 .gitignore
drwxrwxr-x 7 x x    4096 Sep 25 18:24 ..
-rwxr-xr-x 1 x x   18143 Sep 25 18:28 makeself.sh
-rwxr-xr-x 1 x x   14425 Sep 25 18:28 makeself-header.sh
drwxr-xr-x 3 x x    4096 Sep 25 18:28 hello_wide_world.unpacked
-rwxrwxr-x 1 x x 3747982 Sep 25 18:28 hello_wide_world
drwxrwxr-x 3 x x    4096 Sep 25 18:28 .
 
# $user : file hello_wide_world

hello_wide_world: POSIX shell script executable (binary data)

# $user : mkdir tmp ; cd tmp ; mv ../hello_wide_world . ; ./hello_wide_world  a b c 
Hello World!
args : 
   0 : (goos_goarch/linux_amd64/hello_wide_world)
   1 : (a)
   2 : (b)
   3 : (c)

wd : ... src/github.com/gulfie/gobme/hello_wide_world/tmp/hello_wide_world.unpacked

Linux desk6 4.4.0-93-generic #116-Ubuntu SMP Fri Aug 11 21:17:51 UTC 2017 x86_64 x86_64 x86_64 GNU/Linux
Good bye World! Everything ended well enough.
x@desk6:


```
(oh nuts... the fist official issue, the WD is wrong on first exec.... nuts ) 
  




