/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/
package main

import (
	"github.com/IBAX-io/go-ibax/packages/crypto"
	"runtime"
func main() {
	crypto.InitHash("SHA256")
	crypto.InitCurve("ECDSA")
	runtime.LockOSThread()
	cmd.Execute()
}
