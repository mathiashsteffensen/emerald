class SizeMatters
  include Comparable

  attr_reader :str

  def <=>(other)
    str.size <=> other.str.size
  end

  def initialize(str)
    @str = str
  end

  def inspect
    @str
  end
end

EMSpec.describe("Comparable") do
	obj = SizeMatters.new("Z")

	describe "#==" do
		it "returns true when strings have same size" do
			other = SizeMatters.new("Z")
			expect(obj == other).to(eq(true))
		end

		it "returns false when strings don't have same size" do
			other = SizeMatters.new("ZZ")
			expect(obj == other).to(eq(false))
		end
	end
end
