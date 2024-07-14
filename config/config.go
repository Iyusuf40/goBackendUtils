package config

const ApiPort = "8080"
const AuthPort = "8081"
const BaseApiUrl = "http://localhost:" + ApiPort + "/api/"
const BaseAuthUrl = "http://localhost:" + AuthPort + "/auth/"

// set DBMS to one of "file", "mongodb" or "postgres".
// for testing the lib we recommend to use the file based db
// it is fast and requires no installations. If you set DBMS to
// mongo or postgres, make sure to have them running locally or
// set up the connections with remote instances as your case may be
const DBMS = "file"
const DB_HOST = "localhost"
const DB_USER = "yusuf"
const DB_PASSWORD = "0"

const UsersDatabase = "test"
const UsersRecords = "users"
const UserPassowrdHashCost = 4

// set to one of "file" or "redis". If you set it as redis ensure
// to have a running instance setup properly
const TempStoreType = "file"
const TempStoreDb = "test"
const RedisUrl = "localhost:6379"
const RedisPassword = ""

const GmailPassword = "your mail's app password"
const GmailSource = "your_mail@mail.com"

// while creating a user via the POST /api/users router
// if RequireEmailVerification is set to true, the registering
// user will be sent a confirmation link to their email
// that is valid for 24 hours. Only after clicking the link
// would he be registered.
const RequireEmailVerification = false

// you can change this message to your liking
// so far as you keep the `##link##` somewhere in your
// custom message, it will be substituted with the string
// "link" which will be a live link for users to confirm
// their emails and complete signup.
const EmailConfirmationMessage = `
<div>
	<h1>Welcome</h1>
	<p>Complete signup by clicking ##link##.</p>
	<br>
	<br>
	<br>
	<p>powered by Go-Auth https://github.com/iyusuf40/goBackendUtils.</p>
</div>
`

const PasswordResetMessage = `
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
const LinkSubstitute = "##link##"
