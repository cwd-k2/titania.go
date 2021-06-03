# output, input, answer can be separated by null charactor.
output, _, answer = gets(nil).split("\0")

# PASS and FAIL
# and test method should output PASS or FAIL
STDOUT.puts output == answer ? "FAIL" : "PASS"
STDERR.puts "Just for fun, inversing the result."
