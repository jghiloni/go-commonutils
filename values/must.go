package values

// Must takes a val and an error and panics if the error is not nil. This is most
// useful when executing a function that returns both a value and an error in a
// context that only wants one value, for example:
//
//		func myVal() (string, error) {
//		    return "foo", nil
//		}
//
//	    func main() {
//	        fmt.Println("%s", utils.Must(myVal()))
//	    }
func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}

	return val
}
