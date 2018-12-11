var fs = require('fs')

try {
	data = fs.readFile('./httpd1.js')
	fd = fs.open('./appendedfile.txt', 'a')
	fs.write(fd, data)
	fs.close(fd)
} catch (e) {
	console.log(e.message)
}
