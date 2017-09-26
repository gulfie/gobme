#!/bin/sh 

# maybe busybox

# an attempt at a minimalist shell (portable) way to launch the correct go binary. 


# we'll do this closer to the style of makeself for portability reasons (hopefully)

# app launcher and multi exec cli ... alamec 
#
#
#       This is executed possibly with the arguments of the app, so we don't get to use the ARGV
#       instead we use the GO_MULTI_ALAMEC_ARGS ENV var. .. yes.  ick. 
#       we hould be away from collision with that. 
#
#
#       Things we'll want to be doing. 
#       1) figure out what goos/goarch this is executing on and exec that one. 
#       2) switcheru, copy the right goos/goarch thing into the .shar's place ( and move the .shar to .unpacked. 
#       3) 
#
#
#
#
#


# don't trust the ENV versions of GOOS and GOARCH

GOOS="uknown"
GOARCH="unknown"

#set -x 

if false; then 
	pwd
	echo "args : $@"
	export 

	echo "

	#	okay start


	"
	

fi

#	# we'll try with what we are called, we'll temporarily accept an argument for dev.
EXECTARGET=`basename $0` 

RANASSTARTUPSCRIPT="n"


# note we may need to step through the path looking for ourselves. 

if test "$EXECTARGET" = "startupscript";then 
	# we are running from the unpackaging context, we'll need to get the name from another source
	# PWD should be ${EXECTARGET}.unpacked. 

	RANASSTARTUPSCRIPT="y"
	
	# not correct for nomninmal execution !!!! correct for installation only
	RAWDIR=`pwd`
	BASEDIR=`basename "$RAWDIR"`

	# FIX , TODO , remove sed dependency 
	EXECTARGET=`echo "$BASEDIR" | sed 's/\.unpacked$//' `
	#echo "EXECTARGET($EXECTARGET)"
fi

#
#
# these two bits are not normalized well
#


RUNGOAPP="y"
MODE="passthrough"

if test x"$GO_MULTI_ALAMEC_ARGS" != x; then 
	echo "GO_MULTI_ALAMEC_ARGS is defined ($GO_MULTI_ALAMEC_ARGS), NOT executing the underlying go"

	MODE="alamec"

	# parse through the alamec args... like? 
	# probably switcher or JUST_UNPACK_DONT_RUN_APP ( for isntall permissions / users ) 
	SAVEARGS="$@"   # hummm..

	set -- $GO_MULTI_ALAMEC_ARGS


	while true 
	do
		case "$1" in 
		--justunpack | --dontrun) 
			RUNGOAPP="n"
			shift
			;; 
		--help | -h ) 
			echo "there is no help here... fix" 
			exit 2
			;;
		*)
			echo "Unknown flag ($1), quiting in ($0)"
			exit 3
			;;
		esac 
	done
	
	# put the normal argv back
	set -- $SAVEARGS  # ? 

fi





# this'll probably be a huristic thing that'll need a fair amount of testing
# and it won't work well in bare init style containers.

probe_for_goos_goarch()
{
if command uname > /dev/null 
then 
	# probably a unixish
	UNAMES=`uname -s`

	if [ "$?" = "0" ]
	then 
		case "$UNAMES" in 

			# andriod? 
			Linux ) 
				GOOS="linux"
				UNAMEP=`uname -p` 

				if test "$UNAMEP" = "x86_64"
				then 
					GOARCH="amd64"
				elif test "$UNAMEP" = "i386"
				then
					GOARCH="386"
				else 
					# guessing 
					GOARCH=`uname -p | tr 'A-Z' 'a-z' `
				fi
				;; 

			Darwin ) 
				GOOS="darwin"
				
				# huristic and untested
				if uname -a | egrep -i "86_64|amd64" >/dev/null
				then
					GOARCH="adm64"
				else
					GOARCH="386"
				fi	
				
				;; 

			DragonFly ) 
				GOOS="dragonfly"
				GOARCH="amd64" 	# there is only one of them, so it's that. 	
				;;
			FreeBSD ) 
				GOOS="freebsd"
				# guessing 
				GOARCH=`uname -p | tr 'A-Z' 'a-z' ` 
	
				;;

			# nacl, nope  https://github.com/golang/go/wiki/NativeClient
			# maybe 
			
			NetBSD ) 
				GOOS="netbsd"
				UNAMEP=`uname -p`
				
				# internet guess
				if echo "$UNAMEP" | egrep -i "arm" >/dev/null
				then	
					GOARCH="arm"
				elif echo "$UNAMEP" | egrep -i "386" > /dev/null
				then
					GOARCH="386"
	
				elif echo "$UNAMEP" | egrep -i "64|amd" > /dev/null
				then
					GOARCH="amd64"
				fi
				
				;;

			OpenBSD ) 
				GOOS="openbsd"

				# internet guess
				if echo "$UNAMEP" | egrep -i "arm" >/dev/null
				then	
					GOARCH="arm"
				elif echo "$UNAMEP" | egrep -i "386" > /dev/null
				then
					GOARCH="386"
	
				elif echo "$UNAMEP" | egrep -i "64|amd" > /dev/null
				then
					GOARCH="amd64"
				fi
				
				;;
			
			# ?	
			# plan9 
			
			SunOS ) 
				GOOS="solaris"
				GOARCH="amd64"

				;;

			# guessing is there a powershell way?
			CYGWIN* ) 
				GOOS="windows"
				GOARCH="amd64"  # guess. 
				;;

		esac 
	else
		# is this windows or plan9 or something? 
		false	
	fi	
		
fi
}



execcanidate=""


#
# fall back strategy of ... try them all and see what works. 
#	which is horrible because some of them core dump and such ( or launch qeum and have that fail ) 
 

getexeccanidate(){

	#
	## @$!#@$ windows...
	#


	if test x"$RANASSTARTUPSCRIPT" = xy; then 
		target_guess="goos_goarch/${GOOS}_${GOARCH}/${EXECTARGET}" 
	else
		# we need to recover the relative directory to the right place. 
	
		# if the path is relative or absolute, we are good and simple 
		if echo "$0" | egrep '^(\./|\.\./|/|.*/.*)' >/dev/null; then  # that simplifies... to (/|.*/.*)  # fix test. can we get rid of egrep / grep ?  TODO xxx
			RELATIVEDIR=`dirname "$0"`
			target_guess="${RELATIVEDIR}/${EXECTARGET}.unpacked/goos_goarch/${GOOS}_${GOARCH}/${EXECTARGET}"	
		else 
			# I don't know of a better plan than to walk through the PATH and find the first executable that looks right. 
			PATHDIR=""
			for pathdir in $PATH
			do	
				if test -x "$pathdir/$0"; then 
					PATHDIR="$pathdir"
					break	
				fi
			done
	
			if test x"$PATHDIR" != x; then 
				target_guess="${PATHDIR}/${EXECTARGET}.unpacked/goos_goarch/${GOOS}_${GOARCH}/${EXECTARGET}"
			else
				echo "Unable to find ($0) in path.... or on the filesystem... where do I come from anyway?(quiting)"
				exit 7
			fi

			# echo "dev quit, on relative paths for now"
			# exit 2; 
		fi
	fi

	if test "unknown" != "$GOOS" -a "unkown" != "$GOARCH"; then 
		
		if test -f "$target_guess"
		then
	#		echo "yippie, time to go"
			#exec "$target_guess"
			execcanidate="$target_guess" 
		elif test -f "$target_guess.exe"
		then
	#		echo "time to windows"
			execcanidate="$target_guess.exe" # if that even works
		else
			echo "there was no matching execuable around ($target_guess), quiting."
			exit 3
		fi

	else
		echo "there is no binary matching the probed GOOS ($GOOS) GOARCH ($GOARCH) as ($target_guess) "
		echo "switching to guns"
		echo "okay not yet, quiting"
		exit 4
	fi
}


probe_for_goos_goarch
getexeccanidate

if false; then  
	echo "(${GOOS}_${GOARCH})"
	echo "$execcanidate"
fi



switcheru(){

	# catch the first time switcheru
	if test x"$RANASSTARTUPSCRIPT" = xy; then 
		# do the switcheru. mv the .shar to the unpacked directory 
		#  FIX TODO error  check XXX
		# this may explode if more complicated args go to makeself-header

		if test -L "$RAWDIR/../$EXECTARGET"; then 
			NOTE="looks like archive was already switcheru'd" 
		else 
			mv "$RAWDIR/../$EXECTARGET" "$RAWDIR/$EXECTARGET"
			# consider allowing a relative directory rather than an absolute.  FIX XXX TODO
			ln -s "$RAWDIR/startupscript" "$RAWDIR/../$EXECTARGET"	
		fi
	fi

}

# now execution depends on mode. 
#
# in passthrough, startup the prog with $@
# e.g. no alamec args
if test "$MODE" = "passthrough"; then 
# 	echo "running in passthrough mode"
	# may be from startup or from after switcheru

	switcheru

	if test x"$RUNGOAPP" = xy; then
		exec "$execcanidate" $@
	else
		echo "not running go app for some reason";
	fi

elif test "$MODE" = "alamec"; then 
	echo "running in mode alamec... what are we doing again?"
	# we could show up here for two reasons... someone opened up the package then manually ran the startup script...

	switcheru

	if test x"$RUNGOAPP" = xy; then
		exec "$execcanidate" $@
	else
		echo "not running go app for some reason";
	fi

	
else
	echo "Unknown MODE ($MODE), quiting"
	exit 6;
fi

echo "devquit or failure... not sure which"

exit 8


