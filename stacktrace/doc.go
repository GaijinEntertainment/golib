// Package stacktrace provides efficient stack trace capture and formatting for Go programs.
//
// This package captures the call stack of the current goroutine with configurable depth
// and skip parameters, allowing precise control over which frames are recorded.
//
// # Basic Usage
//
// To capture the current stack trace:
//
//	stack := stacktrace.Capture(0, stacktrace.DefaultDepth)
//	fmt.Println(stack.String())
//
// # Controlling Depth
//
// Capture a limited number of frames:
//
//	stack := stacktrace.Capture(0, 10) // Capture up to 10 frames
//
// Capture the complete stack trace:
//
//	stack := stacktrace.Capture(0, math.MaxInt) // Unlimited depth
//
// # Skipping Frames
//
// Skip frames to exclude wrapper functions from the trace:
//
//	func myLogger(msg string) {
//	    // Skip 1 frame to exclude myLogger from the trace
//	    stack := stacktrace.Capture(1, stacktrace.DefaultDepth)
//	    log.Printf("%s\n%s", msg, stack.String())
//	}
//
// # Output Format
//
// The Stack.String() method formats stack traces as:
//
//	function.name
//		/path/to/file.go:123
//	another.function
//		/path/to/other.go:456
//
// This format is compact and easy to parse, with each frame showing the fully-qualified
// function name followed by the source location on the next line, indented with a tab.
//
// # Iterating Over Frames
//
// Access individual frames using the iterator for custom processing or formatting:
//
//	stack := stacktrace.Capture(0, stacktrace.DefaultDepth)
//	for idx, frame := range stack.Frames() {
//	    fmt.Printf("%d: %s at %s:%d\n", idx, frame.Function, frame.File, frame.Line)
//	}
package stacktrace
