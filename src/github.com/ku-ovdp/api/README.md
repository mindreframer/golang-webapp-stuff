HTTP REST API for openvoicedata.org
===================================

This project intends to produce a read/write api for [api.openvoicedata.org](http://api.openvoicedata.org/docs).

Building
========
Assuming you have a functional go environment you can simply:

$ go get github.com/ku-ovdp/api

To add swagger (api documentation viewer):

* $ cd $GOPATH/src/github.com/ku-ovdp
* $ git submodule update --init --recursive

Running
=======

$ ./api -baseURL="http://localhost:5000" -persistenceBackend="dummy"

Expected environment variables for production:
- S3_ACCESS_KEY / S3_PRIVATE_KEY - for S3 interoperation
- S3_BUCKET - base S3 url including bucket (e.g. "https://openvoicedata.s3.amazonaws.com/")
- MONGOURI - for mongo persistence (if using mgo backend)
