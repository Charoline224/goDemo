package main

import "fmt"

fun main(){

}

fun rotate(s *[]int,n int) *[]int {
	res = make([]int, len(s));
	for i:=0;i<len(s);i++{
		res[(i+n)%len(s)] = s[i];
	}
}