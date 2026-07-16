package connectiondetails

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// logFileMode restricts the trace log to root; it records pod names and
// device identity.
const logFileMode = 0o600

// Logf appends a timestamped line to logFile (multus-style CNI debugging -
// CNI stderr is swallowed into kubelet/multus logs and hard to find) and
// mirrors it to stderr. Best-effort: an empty logFile disables the file sink,
// and file open/write/close failures are surfaced to stderr rather than
// propagated.
func Logf(logFile, format string, a ...interface{}) {
	line := fmt.Sprintf("%s ib-sriov-cni: %s\n",
		time.Now().UTC().Format(time.RFC3339), fmt.Sprintf(format, a...))
	fmt.Fprint(os.Stderr, line)
	if logFile == "" {
		return
	}
	f, err := os.OpenFile(filepath.Clean(logFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, logFileMode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ib-sriov-cni: log open %q failed: %v\n", logFile, err)
		return
	}
	// Check both errors: on an append handle a failed Close can lose
	// buffered data.
	_, werr := f.WriteString(line)
	cerr := f.Close()
	if werr != nil || cerr != nil {
		fmt.Fprintf(os.Stderr, "ib-sriov-cni: log write to %q failed: write=%v close=%v\n",
			logFile, werr, cerr)
	}
}
