# gojs (A JavaScript interpreter written in Go)

**gojs** is implemented in Go, and it is an sample application of
[Duktape Bridge for Go](https://github.com/rosbit/duktape-bridge). It is intended to
provide a method of making Go structures to be JavaScript modules. So one can create
framework related jobs with Go, and fulfill the changable logic in JS.

**gojs** embeds the [Duktape](https://duktape.org) JavaScript, so it
can be used as an Ecmascript E5/E5.1 interpreter. Right now, it contains some modules
such as `http`, `fs`, `url`, etc. What I have done is just providing a sample of
implementing a JavaScript module in Go. One can produce more modules if needed.
Enjoy **gojs**.

### Build

**gojs** only depends on [Duktape Bridge for Go](https://github.com/rosbit/duktape-bridge)
and [go-flags](https://github.com/jessevdk/go-flags),
the steps to build **gojs** are just pulling the related codes with go tool:

   1. at any directory, `mkdir src`
   2. and run the following command:

       ```bash
       GOPATH=`pwd` go get https://github.com/jessevdk/go-flags
       GOPATH=`pwd` go get github.com/rosbit/duktape-bridge/duk-bridge-go
       GOPATH=`pwd` go get github.com/rosbit/gojs
       ```

   3. now you get a standalone executable `bin/gojs`, copy it to anywhere you want. It's small
      and has no dependency, it is a full JavaScript interpreter.

### Usage

#### Run simple codes

Just run `bin/gojs -e 'js-code-in-quote-marks'`, e.g.:

  - `bin/gojs -e 'print("hello gojs")'`
  - `bin/gojs -e '1 + 1'`

#### Run JavaScript script file

Just run `bin/gojs <js-file>`

Suppose there's a js file `a.js`:

   ```js
   console.log('hello gojs');
   ```

Run `bin/gojs a.js`, That's all.

There are sample JavaScript files under `src/github.com/rosbit/gojs/js_samples`, which will show how
to use builtin gojs modules such as `http`, `fs`, `url`:

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
