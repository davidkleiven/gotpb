package gotpb

func panicOnErr(e error) {
	if e != nil {
		panic(e)
	}
}
