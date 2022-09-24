module Cache
  class << self
    def conditional_set(key, value)
        db[key] ||= value
    end

    def get(key)
        db[key]
    end

    def db
      @db ||= {}
    end
  end
end

EMSpec.describe "Conditional assignment" do
  describe "||=" do
    context "when it has not been set yet" do
        it "sets and returns the default value" do
            expect(Cache.db).to(eq({}))
            expect(Cache.get("key")).to(eq(nil))
        end
    end
  end
end
