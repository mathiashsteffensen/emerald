EMSpec.describe(Array) do
	describe "#map" do
		arr = [3, 4, 8]

		it "returns a new array with the computed values" do
			expect(arr.map { |n| n * 2 }).to(eq([4, 8, 16]))
		end
	end
end
