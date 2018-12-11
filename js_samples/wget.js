var http = require('http')
var options = {
    url: 'http://www.baidu.com',
    method: 'get',
};

http.request(options, function(resp) {
	console.log(resp.statusCode)
	for (var key in resp.headers) {
		console.log(key, resp.headers[key])
	}
	console.log(resp.data)
})
