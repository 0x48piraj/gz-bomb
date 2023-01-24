# gz-bomb

Web-based "Zip bomb" which eats up all the memory and kills web browsers, scanners and bots. This can also be used to kill systems that download/parse content from a user-provided URL (image-processing servers, AV, websites anything that accept zipped POST data).

Put simply, **gz-bomb** creates a big-ass file with repetitive data (in this case zeros), so the compresion is super-duper efficient but when they try to actually look at the content (extract or decompress the gzip'ed data), they'll most likely run out of disk space or RAM.

If you're a risk taker: [Try it yourself](DEMO)

## Backstory

If you have ever hosted a website or even administrated a server you'll be very well aware of bad people trying bad things with your stuff.

This is why all server or website admins have to deal with gigabytes of logs full with scanning attempts. So I was wondering...and then remembered the good 'ol ZIP bombs from the golden days. Sadly, web browsers don't understand ZIP, but they do understand GZIP!

## Sooo. Results?

| Client | Result                                             |
|--------|----------------------------------------------------|
| Nikto  | Seems to scan fine but no output is reported       |
| SQLmap | High memory usage until crash                      |
| Curl   | Never stops loading content and CPU usage goes up  |
| Hydra  | Never stops loading content but low CPU usage      |

If you test with other devices/browsers/scripts, please create a [pull](https://github.com/0x48piraj/gz-bomb/pulls).

## Compile

Build it using the Go complier (you don't need to *go get* anything):

```bash
$ go build -v -o gz-bomb gz-bomb.go
```

## Usage

```bash
$ ./gz-bomb -help
Usage of ./gz-bomb:
  -filename string
    	Name of the file for the bomb (default "bomb.gz")
  -port int
    	HTTP port to use (default 80)
  -preserve
    	Preserve the bomb
  -size int
    	Size of bomb to create (default 10240)
  -verbose
    	Provide usage verbosity
```

Example usage,

```bash
$ ./gz-bomb -verbose -size 4096 -port 81 -preserve
[*] Creating the bomb: bomb.gz. Wait...
[+] Done!
[*] Serving the bomb on port 81/TCP
[*] Serving 1042079 bytes gzipped to [::1]:34062 (UA:"curl")
[!] Caught signal: interrupt. Exiting...
[-] Diffusing the bomb: bomb.gz
```

## But... Why!?1

Reading about hacker culture and it's rich history is fun and toying with script kiddies is even funnier.