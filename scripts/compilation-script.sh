#!/bin/bash
# The first argument is the compilation command

# Check for missing compilation arguments
if [ $# -lt 1 ]
then
	exit 3
fi

# Check if compilation is needed, for example not needed for Python
if [ "$*" = "compilation-not-needed" ]
then
	exit 0
fi

bash -c "timeout 12 $* 2> cmperr" || exit 7

# If we've reached here then compilation was successful
# Or there were some warnings in compilation
# For example, division by zero is a warning in C++ but compilation is successful
exit 0
