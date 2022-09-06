EMSpec.describe(Array) do
	arr = [2, 4, 8]

	describe "#map" do
		it "returns a new array with the computed values" do
			expect(arr.map { |n| n * 2 }).to(eq([4, 8, 16]))
		end
	end

	describe "#first" do
		context "with no arguments" do
			it "returns the first element of the array" do
				expect(arr.first).to(eq(2))
			end
		end

		context "with an argument" do
			it "returns the first n elements in the array" do
				expect(arr.first(2)).to(eq([2, 4]))
			end
		end
	end
end
