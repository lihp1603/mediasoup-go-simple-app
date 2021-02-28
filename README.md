# mediasoup-go-simple-app

A demo application of [mediasoup-go](https://github.com/jiyeyuran/mediasoup-go). 

ref: https://github.com/mkhahani/mediasoup-sample-app

## Installation

Clone the project:

```tex
$ git clone https://github.com/lihp1603/mediasoup-go-simple-app.git
$ cd mediasoup-go-simple-app
```

Set up the mediasoup-go-sample-app server:

```tex
$ cd server
$ go build
```

Make sure TLS certificates reside in `/etc/certs` directory with names `fullchain.pem` and `privkey.pem`.

Set up the mediasoup-go-simple-app browser app:

```
$ cd app
$ npm install
```

## Run

- Run the  server application in a terminal:

```tex
$ cd server
$ ./server -ip xxx -port xxx 
```

- In a different terminal build and run the browser application:

```tex
$ cd app
$ npm start
```

- Enjoy.

## License

MIT