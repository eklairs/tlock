package boundedinteger

// Bounded integer is an easy way to cycle around a bounds while increasing or decreasing a number
type BoundedInteger struct {
    // Current value
    Value int

    // Maximum value
    Max int
}

// New instance of BoundedInteger
func New(value, max int) BoundedInteger {
    return BoundedInteger{
        Max: max,
        Value: value,
    }
}

// Increases the value by 1
func (integer *BoundedInteger) Increase() {
    integer.Value = (integer.Value + 1) % integer.Max
}

// Decreases the value by 1
func (integer *BoundedInteger) Decrease() {
    if integer.Value == 0 {
        integer.Value = integer.Max - 1
    } else {
        integer.Value -= 1
    }
}

