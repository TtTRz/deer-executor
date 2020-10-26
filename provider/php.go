/* PHP Compiler Provider
 * (C) 2019 LanceLRQ
 *
 * This code is licenced under the GPLv3.
 */
package provider

import (
	"fmt"
)

type PHPCompileProvider struct {
	CodeCompileProvider
}


func (prov *PHPCompileProvider) Init(code string, workDir string) error {
	prov.isReady = false
	prov.realTime = true
	prov.codeContent = code
	prov.workDir = workDir
	prov.Name = "php"

	err := prov.checkWorkDir()
	if err != nil {
		return err
	}

	err = prov.initFiles(".php", "")
	return err
}

func (prov *PHPCompileProvider) Compile() (result bool, errmsg string) {
	result, errmsg = prov.shell(fmt.Sprintf(CompileCommands.PHP, prov.codeFilePath))
	if result {
		prov.isReady = true
	}
	return
}

func (prov *PHPCompileProvider) GetRunArgs() (args []string) {
	args = []string{ "/usr/bin/php", "-f", prov.codeFilePath }
	return
}

func (prov *PHPCompileProvider) IsCompileError(remsg string) bool {
	return false
}

