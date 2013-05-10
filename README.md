Promo Code Server In Golang With MongoDB On Heroku
--------------------------------------------------

Go learning project. Simple implementation of a promo code generator and validator for giving reviewers access to in app purchase protected areas of an iOS app. Minimal query string based authentication and short, easy to type promo codes.

Uses MongoDB hosted via MongoHQ Heroku addon with mgo for persistence and gorilla/mux for routing.

Use
---

To create a new promo code:

`curl http://your-project.herokuapp.com/pc?p=yourpassword`

Outputs code as plain text, sans newline.

To redeem a promo code:

`curl http://your-project.herokuapp.com/promocodes/yourcode`

Outputs `true` or `false` sans newline.

To list valid and redeemed promo codes, visit:

`http://your-project.herokuapp.com/pcall?p=yourpassword`

Run Locally
-----------

Be sure mongodb is installed and running on default port (27017), then:

`PORT=5000 MONGOHQ_URL=mongodb://localhost/promocodes MONGOHQ_DB_NAME=promocodes PROMOPW=password go run web.go`

Deployment
----------

To create a Heroku instance with the Go Buildpack:

`heroku create -b https://github.com/kr/heroku-buildpack-go.git`

Install the MongoHQ Addon:

`heroku addons:add mongohq:sandbox`

Set the following env vars:

`heroku config:set PROMOPW=yourpassword MONGOHQ_DB_NAME=find-this-through-mongohq-web-interface MONGOHQ_URL=also-found-in-mongohq-web-interface`

And push:

`git push heroku master`

Be sure to run `go get` after making any changes to recompile your binary.

