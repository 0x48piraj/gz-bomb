# gz-bomb

Web-based "Zip bomb" which eats up all the memory and kills web browsers, scanners and bots. This can also be used to kill systems that download/parse content from a user-provided URL (image-processing servers, AV, websites anything that accept zipped POST data).

Put simply, **gz-bomb-server** creates a big-ass file with repetitive data (in this case zeros), so the compresion is super-duper efficient but when they try to actually look at the content (extract or decompress the gzip'ed data), they'll most likely run out of disk space or RAM.

If you're a risk taker: [Try it yourself](http://127.0.0.1:31337)

## Backstory

If you have ever hosted a website or even administrated a server you'll be very well aware of bad people trying bad things with your stuff.

This is why all server or website admins have to deal with gigabytes of logs full with scanning attempts. So I was wondering...and then remembered the good 'ol ZIP bombs from the golden days. Sadly, web browsers don't understand ZIP, but they do understand GZIP!

## Sooo. Results?

| Client | Result                                             |
|--------|----------------------------------------------------|
| Nikto  | Seems to scan fine but no output is reported       |
| SQLmap | High memory usage until crash                      |
| Curl   | Never stops loading content and CPU usage goes up (with `-o` flag, curl >= [7.55.0](https://github.com/curl/curl/commit/5385450afd61328e7d24b50eeffc2b1571cd9e2f) doesn’t spew binary anymore) |
| Hydra  | Never stops loading content but low CPU usage      |

If you test with other devices/browsers/scripts, please create a [pull](https://github.com/0x48piraj/gz-bomb/pulls).

## Compile

Build it using the Go complier (you don't need to *go get* anything, get it?):

```bash
$ go build -v -o gz-bomb-server gz-bomb-server.go
```

## Usage

```bash
$ ./gz-bomb-server -help
Usage of ./gz-bomb-server:
  -filename string
    	Name of the file for the bomb (default "bomb.gz")
  -port int
    	HTTP port to use (default 31337)
  -preserve
    	Preserve the bomb
  -size int
    	Size of bomb to create (default 10240)
  -verbose
    	Provide usage verbosity
```

Example usage,

```bash
$ ./gz-bomb-server -verbose -size 4096 -port 31337
[*] Creating the bomb: bomb.gz. Wait...
[+] Done!
[*] Serving the bomb on port 31337/TCP
[*] Serving 4168186 bytes gzipped to 127.0.0.1:50346 (UA:"curl")
[!] Caught signal: interrupt. Exiting...
[-] Diffusing the bomb: bomb.gz
```

Ports <= 1024 are privileged ports. You can't use them unless you're root or have the explicit permission to use them.

If you're hell-bent on using these ports, please don't root your way through. Instead, run the server as a non-privileged user and use the setcap utility to grant it the needed permissions.

For example, to allow binding to a low port number (like 80), run setcap once on the executable:

```bash
$ sudo setcap 'cap_net_bind_service=+ep' /path/to/gz-bomb-server
```

You may need to install setcap: `sudo aptitude install libcap2-bin`

## But... Why!?1

Reading about hacker culture and it's rich history is fun but toying with script kiddies is even funnier.