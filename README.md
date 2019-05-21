# foreach-line -- execute a command for each line in a file

I often create lists of files and save them to text files, for later processing. Working on multiple computers with different shells and environments, I often stumble over different for loop syntaxes. `foreach-line` attempts to resolve this issue by behaving in a defined manner no matter what line ending is used -- as one would expect, one execution per line.

I know `xargs` does support this workflow in principle, but I always misremember the flags.

## Installation

```
$ go install github.com/Xjs/foreach-line
```

## Usage

```
$ foreach-line list-of-links.txt wget
```

## Ideas

* Support arbitrary file name positions by a marker like {} (customisable via flag)

## Contributing

Feel free to contact me if you have any ideas where this might be going, open an issue, or a pull request!
