import "fmt"

fun reverse(s *[length]int) *[length]int {
	for i:=0;i<length/2;i++{
		(s*)[i],(s*)[length-1-i] = (s*)[length-1-i],(s*)[i];	
	}
	return s;
}