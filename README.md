# Reinvented W̶h̶e̶e̶l̶ CLI Parser

CLI parser that's strongly typed using generics, and doesn't force you to conform to a model of "commands" or "actions." This library only parses, it does not run your code for you. This lets you inspect the results of the argument parsing and work with them however you please before running the rest of your program. Other frameworks tie you into a model of "command, subcommand, and function that gets called for the commands." No "pre/post exec" hooks because _you're the one running your code_, `go-wheel` does not dispatch your commands for you. This is a common but unnecessary anti-pattern in CLI frameworks.

<more details to come>
