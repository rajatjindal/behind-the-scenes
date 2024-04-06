module github.com/rajatjindal/behind-the-scenes/api

go 1.20

require (
	github.com/fermyon/spin-go-sdk v0.0.0-20240220234050-48ddef7a2617
	github.com/google/uuid v1.3.1
	github.com/gorilla/mux v1.8.0
	github.com/sirupsen/logrus v1.9.4-0.20230606125235-dd1b4c2e81af
	github.com/slack-go/slack v0.12.3
)

require (
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	github.com/stretchr/testify v1.8.3 // indirect
	golang.org/x/sys v0.19.0 // indirect
)

replace github.com/sirupsen/logrus v1.9.3 => github.com/sirupsen/logrus v1.9.4-0.20230606125235-dd1b4c2e81af
