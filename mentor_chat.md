### The Mental Model for Brackets in Go

If you are struggling with brackets, here is the cheat sheet to memorize. Think of them as physical objects in our factory:

*   **`()` Parentheses = Worker Action / Hand-off.** 
    Every time you see `()`, you are either *calling a worker* to do a job, or handing them raw materials. `fmt.Println("hello")` means "Hand the string hello to the Println worker."
*   **`[]` Square Brackets = Shelves / Collections.** 
    Square brackets are used exclusively for storage shelves (Arrays, Slices, and Maps). When you write `dict["info"]`, you are saying "Go to the `dict` shelf, and look in the bin labeled `info`." When you write `[]byte`, you are saying "This is a shelf that holds bytes."
*   **`{}` Curly Braces = Rooms / Boxes.** 
    Curly braces define a physical space. An `if` statement uses `{}` to say "Only the code *inside this room* happens if the condition is true." A `struct` uses `{}` to box its fields together.

---

### 🚨 Code Rescue: Fixing the Switch Statement

You accidentally missed my instruction from the last step! You pasted both the List block and the Dictionary block **outside** of the `switch` statement's room (its curly braces). You also accidentally forgot to pack the actual dictionary *value* inside your dictionary loop!

Because curly braces can get very confusing when they are misaligned, I am going to give you the complete, final version of the `marshalValue` function from top to bottom. 

Please **highlight your entire `marshalValue` function** (from line 28 all the way down to the bottom of the file) and overwrite it with this:

```go
func marshalValue(w *bytes.Buffer, v interface{}) error {
	switch val := v.(type) {
	case string:
		w.WriteString(fmt.Sprintf("%d:%s", len(val), val))
	case int:
		w.WriteString(fmt.Sprintf("i%de", val))
	case []interface{}:
		w.WriteString("l")
		for _, item := range val {
			err := marshalValue(w, item)
			if err != nil {
				return err
			}
		}
		w.WriteString("e")
	case map[string]interface{}:
		w.WriteString("d")
		var keys []string
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		
		for _, k := range keys {
			err := marshalValue(w, k) // Pack the Key
			if err != nil {
				return err
			}
			err = marshalValue(w, val[k]) // Pack the Value
			if err != nil {
				return err
			}
		}
		w.WriteString("e")
	default:
		return fmt.Errorf("encoder doesn't support type: %T", v)
	}
	return nil
}
```

Once you have replaced that, run `go build -v ./...` in the terminal. If it's clean, we can finally calculate the InfoHash!
