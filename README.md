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

## Flags

`-print`: print each line before executing
`pattern`: replace the given pattern by the line in command line, instead of simply appending it as single argument
  
  Example: `foreachline -pattern {} curl http://example.com/{}`

`-stdin`, `-stdout`, `-stderr`: Attach the given file to the executed commands as stdin/out/err (if pattern is set, it will be replaced in the filename first; output files will be created)

  Example: `foreachline -pattern {} -stdin {} -stdout linecount-{} wc -l`

`-skip-io-fail`: Skip execution of a line if reading the stdin file or creating the stdout/stderr files (if any set) doesn't succeed (if this flag isn't set, the command will be executed anyway with no stdin and stdout and stderr attached to foreach-line's stdout and stderr, respectively)

## Contributing

Feel free to contact me if you have any ideas where this might be going, open an issue, or a pull request!
