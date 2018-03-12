---
Title: Compare time and date
Id: 170
SOId: 32577
---
Sometime you will need to know, with 2 dates objects, if there are corresponding to the same date, or find which date is after the other.

In **Go**, there is 4 way to compare dates:
- `date1 == date2`, returns `true` when the 2 are the same moment
- `date1 != date2`, returns `true` when the 2 are different moment
- `date1.Before(date2)`, returns `true` when the first is strictly before the second
- `date1.After(date2)`, returns `true` when the first is strictly after the second

> WARNING: When the 2 Time to compare are the same (or correspond to the exact same date), functions `After` and `Before` will return `false`, as a date is neither before nor after itself
> - `date1 == date1`, returns `true`
> - `date1 != date1`, returns `false`
> - `date1.After(date1)`, returns `false`
> - `date1.Before(date1)`, returns `false`

<!-- break -->

> TIPS: If you need to know if a date is before or equal another one, just need to combine the 4 operators
> - `date1 == date2 || date1.After(date2)`, returns `true` when date1 is after or equal date2
>or using `! (date1.Before(date2))`
> - `date1 == date2 || date1.Before(date2)`, returns `true` when date1 is before or equal date2
>or using `!(date1.After(date2))`

Some examples to see how to use:

```
    // Init 2 dates for example
    var date1 = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
    var date2 = time.Date(2017, time.July, 25, 16, 22, 42, 123, time.UTC)
    var date3 = time.Date(2017, time.July, 25, 16, 22, 42, 123, time.UTC)

    bool1 := date1.Before(date2) // true, because date1 is before date2
    bool2 := date1.After(date2) // false, because date1 is not after date2

    bool3 := date2.Before(date1) // false, because date2 is not before date1
    bool4 := date2.After(date1) // true, because date2 is after date1

    bool5 := date1 == date2 // false, not the same moment
    bool6 := date1 == date3 // true, different objects but representing the exact same time

    bool7 := date1 != date2 // true, different moments
    bool8 := date1 != date3 // false, not different moments

    bool9 := date1.After(date3) // false, because date1 is not after date3 (that are the same)
    bool10:= date1.Before(date3) // false, because date1 is not before date3 (that are the same)

    bool11 := !(date1.Before(date3)) // true, because date1 is not before date3
    bool12 := !(date1.After(date3)) // true, because date1 is not after date3
```
