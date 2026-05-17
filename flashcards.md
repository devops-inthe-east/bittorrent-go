# 🧠 Golang & Systems Flashcards: Phase 1 (Bencode Parser)

**Q: What is an `interface{}` in Go?**
A: A "Mystery Box". It represents a value of ANY type. When a function returns an `interface{}`, we don't know if it's an int, string, or map until we check it.

**Q: How do we check what's inside the Mystery Box?**
A: With a **Type Assertion**. 
Syntax: `value, ok := mysteryBox.(expectedType)`. If `ok` is true, the box contained what we expected!

**Q: What is the purpose of the `_` (underscore) in Go?**
A: It is the **Blank Identifier** (The Trash Can). We use it when a function returns multiple values (like a result and an error), but we only want to keep one of them. It prevents "unused variable" compiler errors.

**Q: Why do we declare `err` and check `if err != nil` immediately after calling a function?**
A: If a function (the worker) trips and fails, the item it was supposed to bring back might be broken or empty. `nil` means "nothing". If `err != nil` (error is NOT nothing), we must handle it immediately before our program crashes trying to use a broken item.

**Q: What is the hierarchy of Struct, Object, Variable, and Function?**
A: 
- **Struct:** The Blueprint (defines what a machine looks like).
- **Object:** The physical Machine taking up space in memory.
- **Variable:** The Name Tag we slap on the machine to refer to it.
- **Function:** The Worker that takes Inputs (raw materials) and returns Outputs (finished goods).

**Q: Why is our `decodeList` function recursive?**
A: Lists are like Russian Nesting Dolls; a list can contain another list inside it. To open the inner list, the function must call `Decode()` again (looping back into the parsing logic).
