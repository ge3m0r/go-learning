package main

func main() {
	DeferClosureLoopV1()

}

func DeferClosureLoopV1() {
	for i := 0; i < 10; i++ {
		defer func() {
			println(i)
		}()
	}
}
