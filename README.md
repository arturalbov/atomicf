# atomicf

[![version](https://img.shields.io/github/release/arturalbov/atomicf.svg)](https://github.com/arturalbov/atomicf/releases/latest)
[![Coverage Status](https://coveralls.io/repos/github/arturalbov/atomicf/badge.svg?branch=master)](https://coveralls.io/github/arturalbov/atomicf?branch=master)
[![LoC](https://tokei.rs/b1/github/arturalbov/atomicf)](https://github.com/arturalbov/atomicf)
[![license](https://img.shields.io/github/license/arturalbov/atomicf.svg)](https://github.com/arturalbov/atomicf/blob/master/LICENSE)
[![contributors](https://img.shields.io/github/contributors/arturalbov/atomicf.svg)](https://github.com/arturalbov/atomicf/graphs/contributors)
[![contribute](https://img.shields.io/badge/contributions-welcome-orange.svg)](https://github.com/arturalbov/atomicf/graphs/contributors)
[![telegram](https://img.shields.io/badge/Telegram-%40arturalbov-blue.svg?style=social&logo=telegram)](https://t.me/username)

```$xslt
import "github.com/arturalbov/atomicf"
```

atomicf is a go package for ensure files with atomic writing by using logging.
Implemented according to https://danluu.com/file-consistency/. Should prevent file corruption on any Linux filesystem.


It contains `type AtomicFile struct` decoration over standard `os.File` where `func (file *AtomicFile) Write(data []bytes) (n int, err error)`  and `func (file *AtomicFile) WriteAt(b []byte, off int64) (n int, err error)` are atomic.

File could be recovered after corruption using `func (file *AtomicFile) Recover() (err error)` function.

## Issues

If you have any problems with or questions about this package, please contact us
through a [GitHub issue](https://github.com/arturalbov/atomicf/issues).

## Contributing

You are invited to contribute new features, fixes, or updates, large or small;
I am always thrilled to receive pull requests, and do my best to process them
as fast as I can.