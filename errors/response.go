package errors

import Errors "errors"

var ResponseHttpGzipDecodeError = Errors.New("http response error gzip decode")
var ResponseReaderCloseError    = Errors.New("reader response close error")
var ResponseReadError           = Errors.New("response read fatal error. ")

var ResponseEmptyError          = Errors.New("empty http body response")