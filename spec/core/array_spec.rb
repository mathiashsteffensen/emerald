EMSpec.describe(Array) do
	describe "#map" do
		arr = [2, 4, 8]

		expect(arr.map { |n| n * 2 }).to(eq([4, 8, 16]))
	end
end
