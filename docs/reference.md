
###### name

The name of your argument also serves as the name of the environment variable cmdeagle will set for you when it runs the  `start` script within a shell.


###### type

The `type` key denotes a format that cmdeagle will attempt to parse into an environment variable for you when the command is executed. The default type is `string` and effectively does nothing, taking the argument as is.

It's worth noting that environment variables are strings in most systems, and all this key does is help you verify and format the into those string values.

- `string`: The argument is taken as is. This is the default type.
- `number`: The argument is parsed as a number. If it's not a valid number, cmdeagle will fail and print an error.
- `integer`: The argument is parsed as an integer. If it's not a valid integer, or  cmdeagle will fail and print an error. If you pass a number that's not an integer, it will be rounded down to the nearest integer.
- `boolean`: The argument is parsed as a boolean. If it's not a valid boolean, cmdeagle will fail and print an error.
- `url`: The argument is parsed as a URL. If it's not a valid URL, cmdeagle will fail and print an error.
- `filepath`: The argument is parsed as a file path. If it's not a valid file path, cmdeagle will fail and print an error.

###### description

A description of your argument. This is used to generate help text for your CLI.


###### required

Whether the argument is required. If it is, cmdeagle will fail to build your CLI if it's not provided.


###### default

A default value for your argument. This is used if the argument is not provided.
