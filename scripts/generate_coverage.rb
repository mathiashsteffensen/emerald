#!/usr/bin/env ruby

require 'json'
require 'fileutils'
require 'open3'

# Ensure tmp directory exists
FileUtils.mkdir_p('tmp')

# Run tests and get coverage
if ARGV.empty?
  puts "Usage: #{$0} <packages>"
  exit 1
end

test_packages = ARGV.join(' ')
puts "Running tests for #{test_packages}..."
test_command = "EM_TEST=1 go test #{test_packages} -coverprofile=tmp/coverage.out"

test_output = ""
exit_status = nil
Open3.popen2e(test_command) do |stdin, stdout_err, wait_thr|
  while line = stdout_err.gets
    puts line
    test_output << line
  end
  exit_status = wait_thr.value
end

unless exit_status.success?
  puts "\nTests failed with exit code #{exit_status.exitstatus}. Aborting coverage generation."
  exit exit_status.exitstatus
end

# Run go tool cover -func
puts "\nProcessing coverage data..."
cover_command = "go tool cover -func=tmp/coverage.out"
cover_output = `#{cover_command}`

unless $?.success?
  puts "Failed to run 'go tool cover -func'. Make sure tmp/coverage.out exists and is valid."
  exit $?.exitstatus
end

coverage_data = {
  packages: {},
  functions: {},
  total: 0.0
}

# Parse package coverage from test_output
test_output.each_line do |line|
  if line =~ /ok\s+(\S+)\s+.*coverage:\s+([\d.]+)%/
    package = $1
    percentage = $2.to_f
    coverage_data[:packages][package] = percentage
  end
end

# Parse function coverage from cover_output
cover_output.each_line do |line|
  if line =~ /^(\S+)\/([^\/:]+):(\d+):\s+(\S+)\s+([\d.]+)%/
    package = $1
    file_name = $2
    line_num = $3
    func_name = $4
    percentage = $5.to_f
    
    coverage_data[:functions][package] ||= {}
    key = "#{file_name}:#{func_name}"
    coverage_data[:functions][package][key] = percentage
  elsif line =~ /^total:\s+\(statements\)\s+([\d.]+)%/
    coverage_data[:total] = $1.to_f
  end
end

# Sort packages by percentage descending
coverage_data[:packages] = coverage_data[:packages].sort_by { |_k, v| -v }.to_h

# Sort functions by package (matching package order) and then by percentage descending
sorted_functions = {}
coverage_data[:packages].each_key do |pkg|
  if coverage_data[:functions][pkg]
    sorted_functions[pkg] = coverage_data[:functions][pkg].sort_by { |_k, v| -v }.to_h
  end
end
coverage_data[:functions] = sorted_functions

File.write('coverage.json', JSON.pretty_generate(coverage_data))
puts "coverage.json generated."

# Generate HTML report
puts "Generating HTML report..."
`go tool cover -html=tmp/coverage.out -o tmp/coverage.html`

if $?.success?
  puts "tmp/coverage.html generated."
else
  puts "Failed to generate HTML report."
  exit $?.exitstatus
end
