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

By tieing together 



## NOTES: 

1. This is a rough proof of concept rather than something that should be used in critical sections.
2. 




### future possibilties. 

1. Include the build invornment in the gob. In an ideal world someone would be able to crack open the archive and reproduce the artifacts. 
2. worm based distributed computing. 
  
