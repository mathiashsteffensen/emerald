# Emerald - A Ruby VM

Emerald is a Ruby compiler & virtual machine written in Go.

### Building Emerald

Run the following from your command line:
```bash
git clone git@github.com:mathiashsteffensen/emerald.git && \
cd emerald && \
make build
```

This will build an `emerald` & an `iem` executable in the current directory.

To run a source file of Ruby code:
```bash
./emerald main.rb
```

To start the Emerald VM in interactive mode:
```bash
./iem
```

### Supported language features
This is still quite far away from being a real implementation.
The below is a list of the features on the roadmap and the ones that have already been implemented.

NOTE: To say they have been implemented does not mean any features are guaranteed to be compatible with the reference Ruby implementation.

- [x] everything is an object
    - [x] allow method calls on everything
    - [x] operators are method calls
- [ ] full UTF8 support
    - [ ] Unicode identifier
    - [ ] Unicode symbols
- [x] method definitions
    - [x] with parens
    - [x] without parens
    - [x] return keyword
    - [ ] default values for parameters
    - [ ] keyword arguments
    - [ ] block arguments
    - [ ] hash as last argument without braces
    - [x] yield
- [x] method calls
    - [x] with parens
    - [x] without parens
    - [x] without parens with args
    - [x] with block arguments
- [ ] conditionals
    - [x] if
    - [x] if modifier
    - [x] if/else
    - [ ] if/elif/else
    - [ ] ternary `? : `
    - [ ] unless
    - [ ] unless modifier
    - [ ] unless/else
    - [ ] case
    - [x] `||`
    - [x] `&&`
- [ ] control flow
    - [ ] for loop
    - [x] while loop
    - [ ] until loop
    - [ ] break
    - [ ] next
    - [ ] redo
    - [ ] flip flop
- [ ] numbers
    - [ ] integers
        - [x] integer arithmetics
        - [x] integers `1234`
        - [x] integers with underscores `1_234`
        - [ ] decimal numbers `0d170`, `0D170`
        - [ ] octal numbers `0252`, `0o252`, `0O252`
        - [ ] hexadecimal numbers `0xaa`, `0xAa`, `0xAA`, `0Xaa`, `0XAa`, `0XaA`
        - [ ] binary numbers `0b10101010`, `0B10101010`
    - [ ] floats
        - [x] float arithmetics
        - [x] `12.34`
        - [ ] `1234e-2`
        - [ ] `1.234E1`
        - [x] floats with underscores `2.2_22`
- [x] booleans
- [ ] strings
    - [x] double quoted
    - [ ] single quoted
    - [ ] character literals (`?\n`, `?a`,...)
    - [ ] `%q{}`
    - [ ] `%Q{}`
    - [ ] heredoc
        - [ ] without indentation (`<<EOF`)
        - [ ] indented (`<<-EOF`)
        - [ ] “squiggly” heredoc `<<~`
        - [ ] quoted heredoc
            - [ ] single quotes `<<-'HEREDOC'`
            - [ ] double quotes `<<-"HEREDOC"`
            - [ ] backticks <<-\`HEREDOC\`"
    - [ ] escaped characters
        - [ ] `\a` bell, ASCII 07h (BEL)
        - [ ] 	`\b` backspace, ASCII 08h (BS)
        - [ ] 	`\t` horizontal tab, ASCII 09h (TAB)
        - [ ] 	`\n` newline (line feed), ASCII 0Ah (LF)
        - [ ] 	`\v` vertical tab, ASCII 0Bh (VT)
        - [ ] 	`\f` form feed, ASCII 0Ch (FF)
        - [ ] 	`\r` carriage return, ASCII 0Dh (CR)
        - [ ] 	`\e` escape, ASCII 1Bh (ESC)
        - [ ] 	`\s` space, ASCII 20h (SPC)
        - [ ] 	`\\` backslash, \
        - [ ] 	`\nnn` octal bit pattern, where nnn is 1-3 octal digits ([0-7])
        - [ ] 	`\xnn` hexadecimal bit pattern, where nn is 1-2 hexadecimal digits ([0-9a-fA-F])
        - [ ] `\unnnn` Unicode character, where nnnn is exactly 4 hexadecimal digits ([0-9a-fA-F])
        - [ ] `\u{nnnn ...}` Unicode character(s), where each nnnn is 1-6 hexadecimal digits ([0-9a-fA-F])
        - [ ] `\cx` or `\C-x` control character, where x is an ASCII printable character
        - [ ] `\M-x` meta character, where x is an ASCII printable character
        - [ ] `\M-\C-x` meta control character, where x is an ASCII printable character
        - [ ] `\M-\cx` same as above
        - [ ] `\c\M-x` same as above
        - [ ] `\c?` or `\C-?` delete, ASCII 7Fh (DEL)
    - [ ] interpolation `#{}`
    - [ ] automatic concatenation
- [ ] arrays
    - [x] array literal `[1,2]`
    - [x] array indexing `arr[2]`
    - [ ] splat
    - [ ] array decomposition
    - [ ] implicit array assignment
    - [ ] array of strings `%w{}`
    - [ ] array of symbols `%i{}`
- [x] nil
- [x] hashes
    - [x] literal with `=>` notation
    - [x] literal with `key:` notation
    - [x] indexing `hash[:foo]`
    - [x] every Ruby Object can be a hash key
- [ ] symbols
    - [x] `:symbol`
    - [ ] `:"symbol"`
    - [ ] `:"symbol"` with interpolation
    - [ ] `:'symbol'`
    - [ ] `%s{symbol}`
    - [x] singleton symbols
- [ ] regexp
    - [x] `/regex/`
    - [ ] `%r{regex}`
- [ ] ranges
    - [ ] `..` inclusive
    - [ ] `...` exclusive
- [ ] procs `->`
- [x] variables
    - [x] variable assignments
    - [x] globals
    - [x] Ruby globals ($ notation)
- [ ] operators
    - [x] `+`
    - [x] `-`
    - [x] `/`
    - [x] `*`
    - [x] `!`
    - [x] `<`
    - [x] `>`
    - [ ] `**` (pow)
    - [ ] `%` (modulus)
    - [ ] `&` (AND)
    - [ ] `^` (XOR)
    - [ ] `>>` (right shift)
    - [ ] `<<` (left shift, append)
    - [x] `==` (equal)
    - [x] `!=` (not equal)
    - [ ] `===` (case equality)
    - [x] `=~` (pattern match)
    - [ ] `!~` (does not match)
    - [x] `<=>` (comparison or spaceship operator)
    - [x] `<=` (less or equal)
    - [x] `>=` (greater or equal)
    - [ ] assignment operators
        - [ ] `+=`
        - [ ] `-=`
        - [ ] `/=`
        - [ ] `*=`
        - [ ] `%=`
        - [ ] `**=`
        - [ ] `&=`
        - [ ] `|=`
        - [ ] `^=`
        - [ ] `<<=`
        - [ ] `>>=`
        - [x] `||=`
        - [x] `&&=`
- [ ] error handling
    - [ ] begin
    - [x] rescue
    - [ ] ensure
    - [ ] retry
- [ ] constants
- [ ] scope operator `::`
  - [x] Constant access `MyMod::MyClass`
  - [ ] Method access `String::new`
- [ ] classes
    - [x] class objects
    - [x] class Class
    - [x] instance variables
    - [ ] class variables
    - [x] class methods
    - [x] instance methods
    - [x] method overrides
    - [ ] private
    - [ ] protected
    - [ ] public
    - [x] inheritance
    - [x] constructors
    - [x] new
    - [x] `self`
    - [x] singleton classes (also known as the metaclass or eigenclass) `class << self`
    - [ ] singleton methods `def self.method`
    - [x] assigment methods
    - [ ] super in methods
- [x] modules
- [x] object main
- [x] comments '#'
- [ ] C extension compatability
