package libfunc

func Stage0() func([]int, int) []int {

	multiply := func(values []int, multiplier int) []int {
		multipliedValues := make([]int, len(values))
		for i, v := range values {
			multipliedValues[i] = v * multiplier
		}
		return multipliedValues
	}

	return multiply
}

func StringThing() {

	generator := func(done <-chan interface{}, strings ...string) <-chan string {
		sStream := make(chan string)
		go func() {
			defer close(sStream)
			for _, i := range strings {
				select {
				case <-done:
					return
				case sStream <- i:
				}
			}
		}()
		return sStream
	}

	addString := func(
		done <-chan interface{},
		stringStream <-chan string,
		tempString string,
	) <-chan string {

		bucketStream := make(chan string)
		go func() {
			defer close(bucketStream)
			for i := range stringStream {
				select {
				case <-done:
					return
				case bucketStream <- tempString + " " + i:
				}
			}
		}()
		return bucketStream
	}

	done := make(chan interface{})
	defer close(done)

	ss := generator(done, "one", "two")
	addString(done, ss, ",")

	//stringStream :=
	//
	//addString(done, )

}
