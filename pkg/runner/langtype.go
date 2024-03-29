package runner

import "path/filepath"

// Detect language type from filename's extension.
func LangType(filename string) string {
	switch filepath.Ext(filename) {
	case ".c":
		return "c"
	case ".cc", ".cpp":
		return "cpp"
	case ".m":
		return "objective-c"
	case ".java":
		return "java"
	case ".kt":
		return "kotlin"
	case ".scala":
		return "scala"
	case ".swift":
		return "swift"
	case ".cs":
		return "csharp"
	case ".go":
		return "go"
	case ".hs":
		return "haskell"
	case ".erl":
		return "erlang"
	case ".pl":
		return "perl"
	case ".py": // python2 は考えません
		return "python3"
	case ".rb":
		return "ruby"
	case ".php":
		return "php"
	case ".sh":
		return "bash"
	case ".r":
		return "r"
	case ".js":
		return "javascript"
	case ".ts":
		return "typescript"
	case ".coffee":
		return "coffeescript"
	case ".vb":
		return "vb"
	case ".cbl", ".cob":
		return "cobol"
	case ".fs":
		return "fsharp"
	case ".d":
		return "d"
	case ".clj":
		return "clojure"
	case ".exs":
		return "elixir"
	case ".sql":
		return "mysql"
	case ".rs":
		return "rust"
	case ".scm":
		return "scheme"
	case ".lisp":
		return "commonlisp"
	case ".nako3":
		return "nadesiko"
	default:
		return "plain"
	}
}
