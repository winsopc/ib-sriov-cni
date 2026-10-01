package main

import "testing"

// A DEL the daemon will refuse, because the configuration it carries needs a
// delegate this process cannot exec, must be refused before anything is
// withdrawn: the published member is the controller's record of a binding
// that is still in place, and a refused teardown has not released it.
func TestDaemonDELRefusesADelegateConfigurationBeforeWithdrawing(t *testing.T) {
	cacheDirAt(t)
	vfioPFAt(t)
	f, _ := newAnnotationTest(t)
	mine := annotationEntry("net-a", "0000:03:00.0", 0)
	args := sandboxArgs("sandbox-a")
	publishFrom(t, f, args, mine)
	patches := f.patches
	// The definition changed after ADD: same device, now a VF with IPAM.
	args.StdinData = []byte(`{"cniVersion":"0.3.1","name":"net-a","type":"ib-sriov","deviceID":"0000:03:00.0",` +
		`"ibKubernetesEnabled":true,"publishPfGuid":true,"ipam":{"type":"host-local","subnet":"10.0.0.0/24"}}`)
	err := delNetwork(t.Context(), args, invocation{kube: withdrawer(f)})
	if err == nil {
		t.Fatal("daemon DEL with an IPAM delegate was not refused")
	}
	if f.patches != patches {
		t.Fatal("the refused DEL withdrew the published member first")
	}
	checkEntry(t, f, mine, 1)
}
