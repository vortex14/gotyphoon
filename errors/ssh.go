package errors

import Errors "errors"

var ErrorSshCloseClient = Errors.New("error ssh close client")
var ErrorSshCloseSession = Errors.New("failed to create session")