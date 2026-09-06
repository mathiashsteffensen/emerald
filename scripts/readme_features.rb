# Run with: ./emerald scripts/readme_features.rb
# Covers every checked language item in README.md, not unchecked roadmap items.
# Emerald deliberately differs from MRI: integer division may return a Float,
# and Regexp#=~ returns MatchData. Those existing contracts are tested here.
# Embedding, sandbox and build contracts are exercised by Go tests / make.

$checks = 0

def check(label, expected, actual)
  unless expected == actual
    raise "FAIL #{label}: expected #{expected.inspect}, got #{actual.inspect}"
  end
  $checks += 1
end

puts "README: objects, numbers and operators"
check("integer literal", 1234, 1_234)
check("float literal", 12.34, 12.34)
check("negative exponent", 12.34, 1234e-2)
check("uppercase exponent", 12.34, 1.234E1)
check("positive exponent", 123.4, 1.234e+2)
check("float underscores", 2.222, 2.2_22)
check("addition", 7, 3 + 4)
check("subtraction", -1, 3 - 4)
check("multiplication", 12, 3 * 4)
check("division", 3, 12 / 4)
check("fractional integer division", 2.5, 5 / 2)
check("precedence", 14, 2 + 3 * 4)
check("parentheses", 20, (2 + 3) * 4)
check("negation", 5, -(-5))
check("float addition", 3.75, 1.25 + 2.5)
check("float subtraction", -1.25, 1.25 - 2.5)
check("float multiplication", 3.125, 1.25 * 2.5)
check("float division", 0.5, 1.25 / 2.5)
check("float negation", -1.25, -1.25)
check("mixed float arithmetic", 3.5, 1.5 + 2)
check("less", true, 1 < 2)
check("greater", false, 1 > 2)
check("less equal", true, 2 <= 2)
check("greater equal", false, 1 >= 2)
check("equal", true, 2 == 2)
check("not equal", true, 2 != 3)
check("case equal", true, 2 === 2)
check("case unequal", false, 2 === 3)
check("spaceship less", -1, 1 <=> 2)
check("spaceship equal", 0, 2 <=> 2)
check("spaceship greater", 1, 3 <=> 2)
check("float comparison", true, 1.25 < 2.5)
check("float equality", true, 1.25 == 1.25)
check("not true", false, !true)
check("not false", true, !false)
check("not nil", true, !nil)
check("zero truthy", false, !0)
check("empty string truthy", false, !"")
check("empty array truthy", false, ![])
check("integer method", "12", 12.to_s)
check("float method", Float, 1.25.class)
check("true method", TrueClass, true.class)
check("false method", FalseClass, false.class)
check("nil method", NilClass, nil.class)
check("string method", String, "x".class)
check("symbol method", Symbol, :x.class)
check("array method", Array, [].class)
check("hash method", Hash, {}.class)
check("regexp method", Regexp, /x/.class)
check("class is object", Class, String.class)
check("Class class", Class, Class.class)
check("module is object", Module, Enumerable.class)
check("main object", "main", self.to_s)
check("main class", Object, self.class)

puts "README: methods, arguments, return and blocks"
def answer
  42
end

def add(a, b)
  a + b
end

def choose(flag)
  return 7 if flag
  9
end

def keyword_sum(left:, right:)
  left + right
end

def mixed_args(base, extra:)
  base + extra
end

check("bare call", 42, answer)
check("empty parens call", 42, answer())
check("positional call", 7, add(3, 4))
no_parens = add 3, 4
check("call args without parens", 7, no_parens)
check("early return", 7, choose(true))
check("implicit return", 9, choose(false))
check("keyword order", 7, keyword_sum(right: 4, left: 3))
no_parens_keywords = keyword_sum right: 4, left: 3
check("keyword call without parens", 7, no_parens_keywords)
check("mixed keyword args", 7, mixed_args(3, extra: 4))

def empty_method()
end
check("empty definition", nil, empty_method())

def return_nil
  return
  1
end
check("bare return", nil, return_nil)

def twice(n)
  yield(n) + yield(n + 1)
end

def yield_nothing
  yield
end

def yield_pair
  yield 2, 3
end
check("yield repeatedly", 14, twice(3) { |n| n * 2 })
check("yield no args", 42, yield_nothing { answer })
check("yield multiple args", 5, yield_pair { |a, b| a + b })
check("brace block", [2, 4, 6], [1, 2, 3].map { |n| n * 2 })
sum = 0
[1, 2, 3].each do |n|
  sum += n
end
check("do block and captured local", 6, sum)
check("nested blocks", [[2, 3], [3, 4]], [1, 2].map { |a| [1, 2].map { |b| a + b } })

puts "README: conditionals and loops"
check("if true", 1, if true; 1; end)
check("if false", nil, if false; 1; end)
check("if else", 2, if false; 1; else; 2; end)
check("elsif", 2, if false; 1; elsif true; 2; else; 3; end)
check("elsif fallback", 3, if false; 1; elsif false; 2; else; 3; end)
check("ternary true", 1, true ? 1 : 2)
check("ternary false", 2, false ? 1 : 2)
check("unless false", 1, unless false; 1; end)
check("unless true", nil, unless true; 1; end)
check("unless else", 2, unless true; 1; else; 2; end)
modified = 0
modified = 1 if true
modified = 2 if false
modified = 3 unless true
check("modifiers skipped", 1, modified)
modified = 4 unless false
check("unless modifier", 4, modified)
check("or returns left", 3, 3 || 4)
check("or returns right", 4, nil || 4)
check("and returns left", false, false && 4)
check("and returns right", 4, 3 && 4)
check("and nil", nil, nil && 4)
check("boolean precedence", true, true || false && false)
$effects = 0
def effect
  $effects += 1
  9
end
true || effect
false && effect
check("short circuit", 0, $effects)
false || effect
true && effect
check("evaluated boolean operands", 2, $effects)
check("lazy ternary", 1, true ? 1 : effect)
check("ternary skips branch", 2, $effects)
check("case match", "two", case 2; when 1; "one"; when 2, 3; "two"; else; "other"; end)
check("case else", "other", case 4; when 1; "one"; else; "other"; end)
check("case miss", nil, case 4; when 1; "one"; end)
check("case class", true, case "x"; when String; true; else; false; end)
check("case regexp", true, case "abc"; when /b/; true; else; false; end)
i = 0
sum = 0
loop_value = while i < 5
  sum += i
  i += 1
end
check("while iterations", 10, sum)
check("while result", nil, loop_value)
while false
  sum = 100
end
check("while zero iterations", 10, sum)

puts "README: strings, escapes, symbols, arrays and hashes"
check("empty string", 0, "".length)
check("string addition", "abcd", "ab" + "cd")
check("interpolation", "n=3, yes=true, nil=", "n=#{1 + 2}, yes=#{true}, nil=#{nil}")
check("nested interpolation", "outer inner 3", "outer #{"inner #{1 + 2}"}")
check("interpolation order", "a1b2", "a#{1}b#{2}")
# Regex hex escapes provide an oracle independent of string escape decoding.
check("bell escape", true, /^\x07$/ === "\a")
check("backspace escape", true, /^\x08$/ === "\b")
check("tab escape", true, /^\x09$/ === "\t")
check("newline escape", true, /^\x0a$/ === "\n")
check("vertical tab escape", true, /^\x0b$/ === "\v")
check("form feed escape", true, /^\x0c$/ === "\f")
check("carriage return escape", true, /^\x0d$/ === "\r")
check("space escape", true, /^\x20$/ === "\s")
check("symbol literal", :word, :"word")
check("quoted symbol", "two words", :"two words".to_s)
check("singleton symbols", :word, "word".to_sym)
check("empty array", 0, [].length)
array = [1, "two", nil, [4]]
check("array indexing", "two", array[1])
check("nested indexing", 4, array[3][0])
check("negative indexing", [4], array[-1])
check("index out of bounds", nil, array[4])
check("negative out of bounds", nil, array[-5])
check("empty indexing", nil, [][0])
appended = []
append_result = appended << 1
append_result << 2
appended << 3
check("append aliases receiver", [1, 2, 3], append_result)
check("append contents", [1, 2, 3], appended)
check("hash rocket", 1, {"a" => 1}["a"])
check("hash labels", 2, {b: 2}[:b])
check("hash mixed", 2, {"a" => 1, b: 2}[:b])
check("missing hash key", nil, {}[:missing])
check("duplicate hash key", 2, {a: 1, a: 2}[:a])
check("symbol string distinction", [1, 2], [{a: 1, "a" => 2}[:a], {a: 1, "a" => 2}["a"]])
keys = [nil, true, false, 1, 1.25, "str", :sym, [], {}, /x/, Object, Enumerable, Object.new]
key_hash = {}
index = 0
keys.each do |key|
  key_hash[key] = index
  index += 1
end
index = 0
keys.each do |key|
  check("object hash key #{index}", index, key_hash[key])
  index += 1
end
check("regexp literal", "/a(b)/", /a(b)/.inspect)
match = /a(b)/ =~ "zabz"
check("pattern match", MatchData, match.class)
check("match capture", "b", $1)
check("match global", "ab", $&)
check("match miss", nil, /a/ =~ "zzz")

puts "README: variables and assignment operators"
local = 1
local = local + 1
check("local reassignment", 2, local)
$shared = 3
def read_shared
  $shared
end
check("global across method", 3, read_shared)
check("unset global", nil, $readme_unset)
a = 8
a += 4
check("plus assign", 12, a)
a -= 2
check("minus assign", 10, a)
a *= 3
check("multiply assign", 30, a)
a /= 5
check("divide assign", 6, a)
a = nil
a ||= 7
check("or assign nil", 7, a)
a ||= effect
check("or assign preserves", 7, a)
a &&= 8
check("and assign truthy", 8, a)
a = false
a &&= effect
check("and assign false", false, a)
a ||= 9
check("or assign false", 9, a)
check("assign short circuit effects", 2, $effects)
$shared += 4
check("global op assignment", 7, $shared)

puts "README: modules, classes, inheritance, self and assignment methods"
module ReadmeSpace
  module Labels
    def label
      "label"
    end
  end
  class Parent
    include Labels
    def initialize(value)
      @value = value
    end
    def value
      @value
    end
    def value=(value)
      @value = value
      -1
    end
    def identity
      self
    end
    def description
      "parent"
    end
    def increment
      @value += 1
    end
    def unset
      @unset
    end
    class << self
      def kind
        "parent class"
      end
    end
  end
  class Child < Parent
    def description
      "child"
    end
  end
end
parent = ReadmeSpace::Parent.new(2)
child = ReadmeSpace::Child.new(3)
check("constant scope access", Class, ReadmeSpace::Parent.class)
check("constructor", 2, parent.value)
check("inherited constructor", 3, child.value)
check("inherited module", "label", child.label)
check("parent method", "parent", parent.description)
check("method override", "child", child.description)
check("self", child, child.identity)
check("uninitialized ivar", nil, child.unset)
check("ivar op assign", 4, child.increment)
check("instance isolation", 2, parent.value)
check("class method", "parent class", ReadmeSpace::Parent.kind)
check("inherited class method", "parent class", ReadmeSpace::Child.kind)
check("setter expression", 10, (child.value = 10))
check("setter mutation", 10, child.value)
child.value += 2
check("setter op assignment", 12, child.value)
class << child
  def description
    "singleton"
  end
end
check("singleton override", "singleton", child.description)
check("singleton isolation", "child", ReadmeSpace::Child.new(1).description)
class ReadmeOperator
  def +(other)
    other + 10
  end
  def <<(other)
    other + 20
  end
end
check("operator dispatch", 12, ReadmeOperator.new + 2)
check("append operator dispatch", 22, ReadmeOperator.new << 2)

puts "README: rescue and feature interactions"
def rescued_argument_error
  add(1)
rescue ArgumentError
  "rescued"
end
check("rescue matching class", "rescued", rescued_argument_error)
def rescued_standard_error
  raise "boom"
rescue StandardError
  "rescued"
end
check("rescue superclass", "rescued", rescued_standard_error)
def rescued_yield
  yield
rescue StandardError
  "rescued"
end
check("rescue yielded error", "rescued", rescued_yield { raise "boom" })
check("rescue block", [7], [1].map do |n|
  raise "boom"
rescue StandardError
  7
end)
def no_error
  5
rescue StandardError
  6
end
check("rescue not taken", 5, no_error)
# Mix constructors, nested calls, block capture, branches, setters and hash keys.
results = {}
[1, 2, 3].each do |n|
  item = ReadmeSpace::Child.new(n)
  item.value += 1
  results[n] = if item.value > 2
    "#{item.description}:#{item.value}"
  else
    :small
  end
end
check("combined branches", [:small, "child:3", "child:4"], [results[1], results[2], results[3]])
puts "README: edge cases and evaluation order"
# Comments also work after code. The hash inside this string is not a comment.
check("hash character in string", "#", "#") # trailing comment
check("empty if body", nil, if true; end)
check("empty else body", nil, if false; 1; else; end)
check("empty unless body", nil, unless false; end)
$effects = 0
check("side-effecting or", 9, effect || 10)
check("or evaluates left once", 1, $effects)
check("side-effecting and", 10, effect && 10)
check("and evaluates left once", 2, $effects)
def false_effect
  $effects += 1
  false
end
check("false left and", false, false_effect && 10)
check("false and evaluates left once", 3, $effects)

def captured_sum(values)
  total = 0
  values.each do |n|
    total += n
  end
  total
end
check("mutating captured method local", 6, captured_sum([1, 2, 3]))
check("fresh closure invocation", 9, captured_sum([4, 5]))

class ReadmeIndex
  def initialize
    @values = {}
  end
  def [](key)
    @values[key]
  end
  def []=(key, value)
    @values[key] = value
    -1
  end
end
indexed = ReadmeIndex.new
check("custom index assignment", 10, (indexed[:key] = 10))
check("custom index getter", 10, indexed[:key])
indexed[:key] += 2
indexed[:key] -= 1
indexed[:key] *= 3
indexed[:key] /= 3
check("custom compound index assignment", 11, indexed[:key])
$receiver_calls = 0
$key_calls = 0
$indexed = indexed
def index_receiver
  $receiver_calls += 1
  $indexed
end
def index_key
  $key_calls += 1
  :key
end
index_receiver[index_key] += 1
check("compound receiver once", 1, $receiver_calls)
check("compound index once", 1, $key_calls)
$receiver_calls = 0
$key_calls = 0
index_receiver[index_key] ||= effect
check("conditional receiver once", 1, $receiver_calls)
check("conditional index once", 1, $key_calls)
check("conditional value preserved", 12, indexed[:key])
indexed[:missing] ||= 7
indexed[:key] &&= 8
check("conditional index writes", [7, 8], [indexed[:missing], indexed[:key]])

def missing_keyword
  keyword_sum(left: 1)
rescue ArgumentError
  :missing
end
check("missing required keyword", :missing, missing_keyword)
def unknown_keyword
  keyword_sum(left: 1, right: 2, extra: 3)
rescue ArgumentError
  :unknown
end
check("unknown keyword", :unknown, unknown_keyword)
def multiple_rescues
  raise "boom"
rescue ArgumentError
  :wrong
rescue RuntimeError, TypeError => error
  error.message
end
check("multiple rescues and binding", "boom", multiple_rescues)

def plain_global_reader
  readme_plain_global
end
readme_plain_global = 42
check("Emerald plain globals", 42, plain_global_reader)
# Separate instances of the same class must remain separate hash keys.
first_key = Object.new
second_key = Object.new
object_keys = {first_key => 1, second_key => 2}
check("distinct object hash keys", [1, 2], [object_keys[first_key], object_keys[second_key]])
check("hash assignment result", 3, (object_keys[first_key] = 3))
puts "PASS: #{$checks} README feature checks"
$checks
