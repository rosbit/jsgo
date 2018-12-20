# jsgo (A JavaScript interpreter written in Go)

**jsgo** is implemented in Go, and it is an sample application of
[Duktape Bridge for Go](https://github.com/rosbit/duktape-bridge). It is intended to
provide a method of making Go structures to be JavaScript modules. So one can create
framework related jobs with Go, and fulfill the changable logic in JS.

**jsgo** embeds the [Duktape](https://duktape.org) JavaScript, so it
can be used as an Ecmascript E5/E5.1 interpreter. Right now, it contains some modules
such as `http`, `fs`, `url`, `db`, etc. What I have done is just providing a sample of
implementing a JavaScript module in Go. One can produce more modules or change modules
if needed.  Enjoy **jsgo**.

### Binary and Download
   If you don't to want to build `jsgo`, go to [jsgo binary](https://github.com/rosbit/jsgo/releases)
   to download the proper version. Save it as jsgo, `chmod +x jsgo`, it can be used as
   a js interpreter, or as http server, etc..

### Build

**jsgo** only depends on [Duktape Bridge for Go](https://github.com/rosbit/duktape-bridge)
and [go-flags](https://github.com/jessevdk/go-flags). To build an executable `jsgo`, at any
directory, perform the command:

   ```bash
   GOPATH=`pwd` go get -u github.com/rosbit/jsgo
   ```
Now you get a standalone executable `bin/jsgo`, copy it to anywhere you want. It's small
and has no runtime dependency, it is a full JavaScript interpreter.

If you want a `jsgo` with `db` module, which can be acted as a mysql/sqlite3 client. Run this command:

   ```bash
   GOPATH=`pwd` go get -u -tags=db github.com/rosbit/jsgo   # with -tags=db
   ```
Run `bin/jsgo -m`, `db` will appear in the list. Also this `jsgo` has no runtime dependency.

### Usage

#### List built-in modules

  - run `bin/jsgo -m`

#### Run simple codes

Just run `bin/jsgo -e 'js-code-in-quote-marks'`, e.g.:

  - `bin/jsgo -e 'print("hello jsgo")'`
  - `bin/jsgo -e '1 + 1'`

#### Run JavaScript script file

Just run `bin/jsgo <js-file>`

Suppose there's a js file `a.js`:

   ```js
   console.log('hello jsgo');
   ```

Run `bin/jsgo a.js`, That's all.

There are sample JavaScript files under `src/github.com/rosbit/jsgo/js_samples`, which will show how
to use builtin jsgo modules such as `http`, `fs`, `url`, `db`:

   - httpd1.js (Node.js version)

     ```js
     var http = require('http')
       
     var server = http.createServer(function (request, response) {
          response.writeHead(200, {'Content-Type': 'text/plain'})
          response.end('Hello World\n')
     }).listen(8888)

     console.log('Server running at http://127.0.0.1:8888/')
     ```
   - httpd2.js (response body returned directly)

     ```js
     var http = require('http')
     
     var server = http.createServer(function (request, response) {
        return {
            desc: 'json sample',
            ival: 1,
            iaval: [1, 2, 3],
            saval: ['this', 'is', 'a', 'test'],
            mval: {a: 'map', b: 'val', c: 'here'}
        }
     }).listen(8888)

     console.log('Server running at http://127.0.0.1:8888/')
     ```

### Status

The package is not fully tested, so be careful.

### Contribution

Pull requests are welcome! Also, if you want to discuss something send a pull request with proposal and changes.
__Convention:__ fork the repository and make changes on your fork in a feature branch.
