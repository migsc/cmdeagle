#!/bin/sh

# Get arguments from environment variables
name=${ARGS_NAME:-World}
age=$ARGS_AGE

# Get flags from environment variables (convert string to boolean/number)
uppercase=$([ "$FLAGS_UPPERCASE" = "true" ] && echo "true" || echo "false")
lowercase=$([ "$FLAGS_LOWERCASE" = "true" ] && echo "true" || echo "false")
repeat=${FLAGS_REPEAT:-1}

# Construct base greeting
greeting="Hello $name!"
if [ -n "$age" ]; then
    greeting="$greeting You are $age years old."
fi

# Apply case transformations
if [ "$uppercase" = "true" ]; then
    greeting=$(echo "$greeting" | tr '[:lower:]' '[:upper:]')
elif [ "$lowercase" = "true" ]; then
    greeting=$(echo "$greeting" | tr '[:upper:]' '[:lower:]')
fi

# Output greeting with repetition
i=0
while [ $i -lt $repeat ]; do
    echo "$greeting"
    i=$((i + 1))
done 