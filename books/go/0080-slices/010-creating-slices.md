---
Title: Create a slice
Id: 60
SOId: 4807
---
Slices are the typical way Go programmers store lists of data.

To declare a slice variable use the `[]Type` syntax.

    var a []int

To declare and initialize a slice variable in one line use the `[]Type{values}` syntax.

    var a []int = []int{3, 1, 4, 1, 5, 9}

Another way to initialize a slice is with the `make` function. It takes three arguments: the `Type` of the slice (or [map][ch-732]), the `length`, and the `capacity`.

    a := make([]int, 0, 5)

You can add elements to your new slice using `append`.

    a = append(a, 5)

Check the number of elements in your slice using `len`.

    length := len(a)

Check the capacity of your slice using `cap`. The capacity is the number of elements currently allocated to be in memory for the slice. You can always append to a slice at capacity as Go will automatically create a bigger slice for you.

    capacity := cap(a)

You can access elements in a slice using typical indexing syntax.

    a[0]  // Gets the first member of `a`

You can also use a `for` loop over slices with `range`. The first variable is the index in the specified array, and the second variable is the value for the index.

    for index, value := range a {
        fmt.Printf("Index: %v Value: %v\n", index, value)  // Prints "Index: 0 Value: 5" (and continues until end of slice)
    }

