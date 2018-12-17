var db = require('db')

var options = {
	host: 'your_mysqld_host', // a hostname or IP, or an unix domain socket file path `/tmp/mysql.sock`
	port: 3306,               // default mysqld listening port. for unix domain socket, port will be neglected.
	user: 'saas',             // user name
	password: 'saas1234',     // password
	database: 'saas'          // database is not necessary
}
/* the following is for sqlite3 options.
var options = {
	type: 'sqlite3',
	db: './test.db'
}*/
var conn = db.createConnection(options)
try {
	conn.query('select * from admin', function(err, resultSet) {
		if (err) {
			console.log(err)
			return
		}

		var columns = resultSet.fields()
		console.log(columns)
		while (true) {
			var row = resultSet.next() // or resultSet.nextRow() to get row array.
			if (row) {
				console.log(row)
			} else {
				break
			}
		}
	})
} catch (ex) {
	console.log(ex.message)
} finally {
	conn.close()
}
