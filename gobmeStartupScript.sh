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
#export 

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
# given a prefix "" being pwd attempt to find $0 
#
#  
# 
#  getself / getexeccandiate are two distinct ideas, and differnt. 
# 
#	getexeccandiate is the underlying GOOS_GOARCH EXECTARGET that needs to get run. 
#
# ... and we'll need to return the possibly relative path from the users original dir perspective. 
#
#
# 	This is a mess....
#		different cases 
#		a) startup script first time or startup script time X 
#		b) found by relative or anchored path  vs found by PATH
#		c) found by relative PATH vs anchored PATH	
#
#
#
#

getself(){
	TARG="$1" ; shift   # pass in $0 for simpler testing
	PREFIX="$1"

	# clunky
	if test -z "$PREFIX"; then 
		true
	else 
		PREFIX="$PREFIX/"
	fi

	

		# if the path is relative or absolute, we are good and simple 
		if echo "$TARG" | egrep '^(\./|\.\./|/|.*/.*)' >/dev/null; then  # that simplifies... to (/|.*/.*)  # fix test. can we get rid of egrep / grep ?  TODO xxx
			RELATIVEDIR=`dirname "$TARG"`
			target_guess="${PREFIX}${RELATIVEDIR}/${EXECTARGET}.unpacked/goos_goarch/${GOOS}_${GOARCH}/${EXECTARGET}"	
			target_guess_relative="${RELATIVEDIR}/${EXECTARGET}.unpacked/goos_goarch/${GOOS}_${GOARCH}/${EXECTARGET}"
			target_guess_fullpath="$target_guess"
		else 
			# I don't know of a better plan than to walk through the PATH and find the first executable that looks right. 
			PATHDIR=""
			savedPWD=`pwd`  

			if test x"$PREFIX" != x ; then 
				cd "$PREFIX"
			fi 
			
			set -f 
			OLDIFS="$IFS"
			IFS=":"
			for pathdir in $PATH
			do	
				
				if test -x "$pathdir/$TARG"; then 
					PATHDIR="$pathdir"
					
					# if the pathdir is relative, then we'll need to prepend the pwd for the file compares
					if echo "$PATHDIR" | egrep -v '^/' ; then 
						PWD=`pwd`
						FULLPATHDIR="$PWD/$PATHDIR"
					else 
						false	
					fi	
					break	
				fi
			done
			set +f 	
			IFS="$OLDIFS"
			unset OLDIFS

			cd "$savedPWD" ; # ouch.. dd


			if test x"$PATHDIR" != x; then 
				target_guess="${PATHDIR}/${EXECTARGET}.unpacked/goos_goarch/${GOOS}_${GOARCH}/${EXECTARGET}"
				target_guess_relative="$target_guess"
				target_guess_fullpath="${FULLPATHDIR}/${EXECTARGET}.unpacked/goos_goarch/${GOOS}_${GOARCH}/${EXECTARGET}"
			else
				echo "Unable to find ($TARG) in path.... or on the filesystem... where do I come from anyway?($PREFIX)(quiting)"
				exit 7
			fi

			# echo "dev quit, on relative paths for now"
			# exit 2; 
		fi
	

	if test "unknown" != "$GOOS" -a "unkown" != "$GOARCH"; then 
		
		if test -f "$target_guess_fullpath"
		then
	#		echo "yippie, time to go"
			#exec "$target_guess"
			execcanidate="$target_guess_relative" 
		elif test -f "$target_guess_fullpath.exe"
		then
	#		echo "time to windows"
			execcanidate="$target_guess_relative.exe" # if that even works
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


# getexeccanidate
if test x"$RANASSTARTUPSCRIPT" = xy; then 
 	#	target_guess="goos_goarch/${GOOS}_${GOARCH}/${EXECTARGET}" 
	#
	cwd=`pwd`
	# so we need to figure out the $0 of our parent.. that'll take a bit.  or we could just mule in our data with USER_PWD
	# 
	#echo "F IX HERE @" .... NAHH
	getself "$MS_DOLLAR_ZERO" "$MS_USER_PWD" # not as horrible as it may seem
else 
	getself "$0" ""
fi

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

	# FIX, such a commonly possibly variable is less safe.  ( USER_PWD ) 
	if test x"$RUNGOAPP" = xy; then
		cd "$MS_USER_PWD"

		unset MS_USER_PWD 
		unset MS_DOLLAR_ZERO
		# there may be a semi infinitie number of keys like this on different systems. we'll need to weed through them.
		if test x"$MS_OLDPWD" != x; then 
			if test x"OLDPWD" != x ; then 
				OLDPWD="$MS_OLDPWD" 
			fi 
		fi
		unset MS_OLDPWD
		# fixup the execcanidate from the perspective of the USER_PWD
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


