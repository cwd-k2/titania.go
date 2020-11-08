#!/usr/bin/env ruby
require 'yaml'
require 'json'
puts YAML.dump(JSON.parse STDIN.read).gsub(/\n\n'/, "\\n'")
