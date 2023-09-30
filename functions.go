package main

func (p PathExists) Score() bool {
	return true
}

func (f FileContains) Score() bool {
	return true
}
