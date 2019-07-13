package util

import (
	"errors"
	"io/ioutil"
	"os/exec"
)

// CheckErr if err is not nill, call panic to dump error message and die
func CheckErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// WgToolIsExisted return true if found wg, othewise return false
func WgToolIsExisted() bool {
	whichCmd := exec.Command("which", "wg")

	err := whichCmd.Run()

	return err == nil
}

// WgGenPrivateKey run command wg genkey to generates a new private key
func WgGenPrivateKey() (string, error) {
	wgCmd := exec.Command("wg", "genkey")

	output, err := wgCmd.Output()

	return string(output), err
}

// WgGenPublicKey run command wg pubkey to generates a new public key
func WgGenPublicKey(privatekey string) (string, error) {
	wgCmd := exec.Command("wg", "genpsk")

	in, err := wgCmd.StdinPipe()
	if err != nil {
		return "", errors.New("can't create pipe that will be connected to the command's standard input")
	}
	out, err := wgCmd.StdoutPipe()
	if err != nil {
		return "", errors.New("can't create pipe that will be connected to the command's standard output")
	}
	wgCmd.Start()

	in.Write([]byte(privatekey))
	in.Close()
	output, err := ioutil.ReadAll(out)
	if err != nil {
		return "", errors.New("can't read output from wg genpsk")
	}
	wgCmd.Wait()

	return string(output), nil
}

// if util.WgToolIsExisted() {
// 	fmt.Println("wg is ready for use")
// } else {
// 	panic("cannot found wg")
// }
//
// prikey, err := util.WgGenPrivateKey()
//
// if err != nil {
// 	panic(err.Error())
// }
// fmt.Println("new private key", prikey)
//
// pubkey, err := util.WgGenPublicKey(prikey)
//
// if err != nil {
// 	panic(err.Error())
// }
//
// fmt.Println("new public key", pubkey)
