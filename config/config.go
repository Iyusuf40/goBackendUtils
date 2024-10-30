package config

var ApiPort = "8081"
var AuthPort = "8082"
var BaseApiUrl = "http://localhost:" + ApiPort + "/api/"
var BaseAuthUrl = "http://localhost:" + AuthPort + "/auth/"

func SetApiPort(port string) {
	ApiPort = port
}

func SetAuthPort(port string) {
	AuthPort = port
}

func SetBaseApiUrl(url string) {
	BaseApiUrl = url
}

func SetBaseAuthUrl(url string) {
	BaseAuthUrl = url
}

// set DBMS to one of "file", "mongodb" or "postgres".
// for testing the lib we recommend to use the file based db
// it is fast and requires no installations. If you set DBMS to
// mongo or postgres, make sure to have them running locally or
// set up the connections with remote instances as your case may be
var DBMS = "file"
var DB_HOST = "localhost"
var DB_USER = "yusuf"
var DB_PASSWORD = "0"

var UsersDatabase = "test"
var UsersRecords = "users"
var UserPassowrdHashCost = 4

func SetDBMS(dbms string) {
	DBMS = dbms
}

func SetDB_HOST(db_host string) {
	DB_HOST = db_host
}

func SetDB_USER(db_user string) {
	DB_USER = db_user
}

func SetDB_PASSWORD(db_password string) {
	DB_PASSWORD = db_password
}

func SetUsersDatabase(usersDatabase string) {
	UsersDatabase = usersDatabase
}

func SetUsersRecords(usersRecords string) {
	UsersRecords = usersRecords
}

func SetUserPasswordHashCost(cost int) {
	UserPassowrdHashCost = cost
}

// set to one of "file" or "redis". If you set it as redis ensure
// to have a running instance setup properly
var TempStoreType = "file"
var TempStoreDb = "test"
var RedisUrl = "localhost:6379"
var RedisPassword = ""

var GmailPassword = "your mail's app password"
var GmailSource = "your_mail@mail.com"

// while creating a user via the POST /api/users router
// if RequireEmailVerification is set to true, the registering
// user will be sent a confirmation link to their email
// that is valid for 24 hours. Only after clicking the link
// would he be registered.
var RequireEmailVerification = false

// you can change this message to your liking
// so far as you keep the `##link##` somewhere in your
// custom message, it will be substituted with the string
// "link" which will be a live link for users to confirm
// their emails and complete signup.
var EmailConfirmationMessage = `
<div>
	<h1>Welcome</h1>
	<p>Complete signup by clicking ##link##.</p>
	<br>
	<br>
	<br>
	<p>powered by Go-Auth https://github.com/iyusuf40/goBackendUtils.</p>
</div>
`

var PasswordResetMessage = `
<div>
	<h1>Reset Password</h1>
	<p>Reset Password by clicking ##link##.</p>
	<br>
	<br>
	<br>
	<p>powered by Go-Auth https://github.com/iyusuf40/goBackendUtils.</p>
</div>
`

// if you change this, make sure to use the new value in
// EmailConfirmationMessage
var LinkSubstitute = "##link##"

func SetTempStoreType(storeType string) {
	TempStoreType = storeType
}

func SetTempStoreDb(tempStoreDb string) {
	TempStoreDb = tempStoreDb
}

func SetRedisUrl(redisUrl string) {
	RedisUrl = redisUrl
}

func SetRedisPassword(redisPassword string) {
	RedisPassword = redisPassword
}

func SetGmailPassword(gmailPassword string) {
	GmailPassword = gmailPassword
}

func SetGmailSource(gmailSource string) {
	GmailSource = gmailSource
}

func SetRequireEmailVerification(requireEmailVerification bool) {
	RequireEmailVerification = requireEmailVerification
}

func SetEmailConfirmationMessage(message string) {
	EmailConfirmationMessage = message
}

func SetPasswordResetMessage(message string) {
	PasswordResetMessage = message
}

func SetLinkSubstitute(linkSubstitute string) {
	LinkSubstitute = linkSubstitute
}

var AllowAllOrigin = false
