module DiffHelpers
  class Differ
    # rubocop:disable Metrics/AbcSize
    # rubocop:disable Metrics/MethodLength
    def diff(actual, expected, path = [])
      return if actual == expected

      types = [actual, expected].map(&:class).uniq

      if types == [Hash]
        compare_arrays(actual.keys, expected.keys, path, actual)
        expected.each do |k, expected_val|
          diff(actual[k], expected_val, path + [k])
        end

      elsif types == [Array]
        if compare_arrays(actual, expected, path)
          expected.each_with_index do |expected_el, i|
            diff(actual[i], expected_el, path + ["[#{i}]"])
          end
        end

      elsif types.size == 1 # same class
        puts "Mismatched values in #{path}:"
        puts "\t  actual=#{actual}"
        puts "\texpected=#{expected}"

      else
        puts "Mismatched types in #{path}:"
        puts "\t  actual=#{types[0]} value=#{actual.inspect}"
        puts "\texpected=#{types[1]} value=#{expected.inspect}"
      end
    end
    # rubocop:enable Metrics/MethodLength
    # rubocop:enable Metrics/AbcSize

    private

    def compare_arrays(actual, expected, path, context = actual)
      if expected.size != actual.size
        puts "Extra/missing elements in #{path}:"
        puts "\tactual=#{actual.size} expected=#{expected.size}"
        puts "\textra=#{actual - expected}"
        puts "\tmissing=#{expected - actual}"
        puts "\tcontext=#{context}"
        false
      else
        true
      end
    end
  end

  def diff_assert(actual, expected)
    expect(actual).to eq(expected)
  rescue RSpec::Expectations::ExpectationNotMetError
    Differ.new.diff(actual, expected)
    raise
  end
end
