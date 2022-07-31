package errors

import Errors "errors"

var ResponseHttpGzipDecodeError = Errors.New("http response error gzip decode")
var ResponseReaderCloseError = Errors.New("reader response close error")
var ResponseReadError = Errors.New("response read fatal error. ")
var ResponsePathError = Errors.New("not found resource path. ")
var ResponseEmptyError = Errors.New("empty http body response")
var ResponseNotOkError = Errors.New("response status is not ok")
