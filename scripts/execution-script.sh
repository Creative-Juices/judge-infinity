#!/bin/bash
# The first argument is the time limit ( uses default unit seconds )
# The second argument is the execution command
# The third argument is the input path of the testcase
# The fourth argument is the output path of the testcase

# Check for missing shell arguments
if [ $# -lt 4 ]
then
    exit 3
fi

timelimit=$1
runcommand=$2
inputpath=$3
outputpath=$4

exitcode=4

if timeout $timelimit $runcommand < $inputpath > output 2> runerr
then
    if timeout 12 ./cmp $outputpath output
    then
        exitcode=$?
    else
        exitcode=$?
        if [ "$exitcode" == "143" ]
        then
            exitcode=5
        fi
    fi
else
    if [ "$?" == "143" ]
    then
        exitcode=5
    else
        exitcode=6
    fi
fi

exit $exitcode