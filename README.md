# CVE-2021-4034
> January 25, 2022 | An00bRektn

This is a golang implementation of CVE-2021-4034 based on (*read as:* blatantly stolen from) [clubby789's](https://github.com/clubby789/CVE-2021-4034) implementation of the vulnerability discovered by [Qualys](https://blog.qualys.com/vulnerabilities-threat-research/2022/01/25/pwnkit-local-privilege-escalation-vulnerability-discovered-in-polkits-pkexec-cve-2021-4034).

I made a blog post going into further detail about the vulnerability, mainly for myself, but for anyone who wants a more drawn out explanation, which you can find [here](https://an00brektn.github.io/pwnkit/).

## Mitigation
Canonical has already released a patched version in the APT package manager, so all Ubuntu distros that are not EOL that are up-to-date (`sudo apt update && sudo apt upgrade`), should be good. If your distro does not have a patch yet, the current "band-aid" solution is to remove the SUID bit like so: `sudo chmod 0755 $(which pkexec)`.

## FAQ

#### How are you doing?
Good, thanks for asking 😊

#### Did you discover this vulnerability?
No. This was from Qualys. Did you not read the first part?

#### Why clubby789?
Saw them create a new repository in my GitHub feed and decided I'd try to recreate it in another language (yes I have nothing better to do).

#### Why Golang? Shouldn't you just use C?
I just wanted to get some more Golang practice in. Also, currently not smart enough to figure out the exact details of the vuln without looking at a PoC. I get it now, but I needed to copy the code to understand the code. Also, golang compiles statically, which is neat :)

#### Who asked?
idk

## Advisory
This code should be used for authorized penetration testing and/or educational purposes only. Any misuse of the program will not be the responsibility of the author. Only run this exploit against machines that you have express permission to run it against.
