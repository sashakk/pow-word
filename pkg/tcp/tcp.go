package tcp

const (
	CodeRequestPuzzle       = 1
	CodeRequestReturnPuzzle = 2
	CodeRequestValidatePoW  = 3
	CodeRequestValidPoW     = 4
	CodeRequestInvalidPoW   = 5

	IndexRequestCode         = 0
	IndexRequestSizeOfPuzzle = 1
	IndexRequestTargetBits   = 2
	IndexSizeOfQuote         = 1
	StartQuote               = 2
	StartPuzzle              = 3
	StartNonce               = 1
	EndNonce                 = 5
	StartHash
	SizeOfHash = 32
)
