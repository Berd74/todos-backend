package main

import (
	"fmt"
	"strings"
)

func main() {

	loop()

}

func start() {

	x := "xd"
	fmt.Printf("type of x : %T \n", x)

	i := 3

	var input string
	fmt.Println("pointer of input : ", &input)
	fmt.Printf("type of input : %T \n", input)
	fmt.Printf("val of input : %v \n", input)

	for i > 0 {
		fmt.Scan(&input)
		fmt.Println("input: " + input)
		i--
	}
}

func array() {

	array := [3]int{-1, 2, 3}

	array[1] = 4

	slice := []string{}

	slice = append(slice, "a")
	slice = append(slice, "b")
	slice = append(slice, "c")
	slice = append(slice, "d")

	var array2 = slice[0:2]

	fmt.Printf("type of array: %T \n", array)
	fmt.Printf("val of array: %v \n", array)

	fmt.Printf("type of slice: %T \n", slice)
	fmt.Printf("val of slice: %v \n", slice)

	fmt.Printf("type of array2: %T \n", array2)

}

func loop() {
	arr := [10]string{"Bob Dun", "Stave Zoomer", "Jim Boomer", "Max Mileniar"}

	fmt.Printf("val of arr: %v \n", arr[6])
	fmt.Printf("val of arr: %v \n", arr)
	fmt.Printf("len of arr: %v \n", len(arr))

	for i, name := range arr {
		fmt.Printf("val of %v: %v \n", i, name)
		if name != "" {
			var f = strings.Fields(name)[0]
			var s = strings.Fields(name)[1]
			fmt.Printf("%v %v \n", f, s)
		}
	}

}
