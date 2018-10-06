a <i>Sm</i>ol <i>sh</i>ell
==========================

`smsh` is a small shell based almost entirely on Bash syntax, built in Go and made for ultraportability and simplicity.

`smsh` is designed to take a list of strings as arguments and treat them as a series of statements.
This makes it easy to use the string array of the exec syscall to clearly send a series of statements without more parsing, e.g.
`["bash", "-c", "stmt1; stmt2; stmt3"]` is equivalent to `["smsh", "--", "stmt1", "stmt2", "stmt3"]`.
The utility of this kind of composition probably makes sense to you especially if you've used something like dockerfiles
(it's a lot clearer to append array elements than it is to concatenate strings to the last element of an array and *hope*
that the result is not only syntactically valid but semantically correct).
Each individual statement string does not need escaping with `smsh`; with bash, to be safe, it would (though usually people simply ignore this).

`smsh` will keep a short record of the string of each command it runs, and time each command, and can print this out in a table at the end of the whole `smsh` invocation.
This is sort of like Bash's `set -x` mode plus wrapping each command with `time`, then prettyprinting the final result.
(We print this at the *end* of the whole `smsh` invocation rather than at each command launch because this remains usable even if some of your commands
produce more-than-the-scrollback-limit amounts of output.)
This information can be printed in a human-friendly textual table, or in json, and can be sent to a file descriptor of your choice for mechanical parsing.

`smsh` will default to exit-on-error, unlike bash, which by default ignores exit codes.  This is similar to Bash `set -e` mode.
(See http://redsymbol.net/articles/unofficial-bash-strict-mode/ for a discussion of why exit-on-error is sane.)
If there are commands you explicitly want to ignore the failure of, you can use the `cmd || true` pattern.


Acknowledgements
----------------

`smsh` uses the shell parser https://github.com/mvdan/sh , which is a spectacular piece of work and without which `smsh` could not exist.
