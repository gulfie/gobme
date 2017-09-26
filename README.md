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

1. ~~ initial release ~~
2. some testing, both local on single architectures. 
3. version # 
4. some visual use examples. 


### future possibilties. 

1. Include the build invornment in the gob. In an ideal world someone would be able to crack open the archive and reproduce the artifacts. 
2. worm based autonomic distributed computing. 


  
