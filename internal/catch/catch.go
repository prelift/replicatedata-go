package catch

func Panic(f func()) (caught any) {
	defer func() { caught = recover() }()
	f()
	return caught
}
