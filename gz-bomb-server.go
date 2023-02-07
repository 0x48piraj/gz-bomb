package main

import (
"flag"
"fmt"
"io"
"io/ioutil"
"log"
"net/http"
"os"
"os/exec"
"os/signal"
"strings"
)

/* defaults and configs */
const (
	default_file = "bomb.gz"
	default_port = 31337
	default_size = 10240
)

var (
	bomb     *os.File
	verbose  bool
	preserve bool
	file    string
	port    int
	size    int
	// common paths/files for scanners/bots
	handle = []string{"/", "/robots.txt", "/admin", "/.htpasswd"}
	// user agents indicating scanners/bots
	useragents = []string{
		"(Hydra)", ".nasl", "absinthe", "advanced email extractor",
		"arachni/", "autogetcontent", "bilbo", "BFAC",
		"brutus", "brutus/aet", "bsqlbf", "cgichk",
		"cisco-torch", "commix", "core-project/1.0", "crimscanner/",
		"datacha0s", "dirbuster", "domino hunter", "dotdotpwn",
		"email extractor", "fhscan core 1.", "floodgate", "get-minimal",
		"gootkit auto-rooter scanner", "grabber", "grendel-scan", "havij",
		"inspath", "internet ninja", "jaascois", "zmeu", "masscan", "metis", "morfeus fucking scanner",
		"mysqloit", "n-stealth", "nessus", "netsparker",
		"Nikto", "nmap nse", "nmap scripting engine", "nmap-nse",
		"nsauditor", "openvas", "pangolin", "paros",
		"pmafind", "prog.customcrawler", "qualys was", "s.t.a.l.k.e.r.",
		"security scan", "springenwerk", "sql power injector", "sqlmap",
		"sqlninja", "teh forest lobster", "this is an exploit", "toata dragostea",
		"toata dragostea mea pentru diavola", "uil2pn", "user-agent:", "vega/",
		"voideye", "w3af.sf.net", "w3af.sourceforge.net", "w3af.org",
		"webbandit", "webinspect", "webshag", "webtrends security analyzer",
		"webvulnscan", "whatweb", "whcc/", "wordpress hash grabber",
		"xmlrpc exploit", "WPScan", "curl", "index", "wbsearchbot", "proximic", "panscient", "Googlebot", "psbot", "whatsapp", "findxbot", "archive", "Voyager", "pinterest", "ccbot", "bibnum", "Linguee Bot", "libwww", "spider", "wget", "archive.org_bot", "Mediapartners-Google", "Twitterbot", "APIs-Google", "Livelapbot", "msn", "Adidxbot", "europarchive.org", "ip-web-crawler.com", "Facebot", "CrystalSemanticsBot", "bingbot", "fluffy", "Python-urllib", "360Spider", "lipperhey", "binlar", "buzzbot", "postrank", "urlappendbot", "ngbot", "Lipperhey SEO Service", "curl", "grub.org", "Googlebot-Mobile", "lssrocketcrawler", "scaper", "acoonbot", "siteexplorer.info", "yeti", "purebot", "indexer", "Domain Re-Animator Bot", "httrack", "heritrix", "AdsBot", "blexbot", "yacybot", "aihitbot", "FAST-WebCrawler", "msnbot", "sistrix crawler", "xovibot", "ia_archiver", "discobot", "woriobot", "Google favicon", "cXensebot", "DuplexWeb-Google", "baidu", "IOI", "findthatfile", "sogou", "AdvBot", "GingerCrawler", "drupact", "spbot", "wocbot", "facebook", "gnam gnam spider", "googleweblight", "instagram", "y!j-asr", "lb-spider", "careerbot", "antibot", "smtbot", "Aboundex", "bibnum.bnf", "teoma", "changedetection", "gslfbot", "jyxobot", "biglotron", "blekkobot", "tagoobot", "niki-bot", "toplistbot", "google", "yoozBot", "Commons-HttpClient", "citeseerxbot", "UsineNouvelleCrawler", "turnitinbot", "rogerbot", "yandexbot", "wotbox", "bot", "grub", "duckduck", "voilabot", "mediapartners", "TweetmemeBot", "ichiro", "CC Metadata Scaper", "webcrawler", "brainobot", "WeSEE:Search", "ips-agent", "g00g1e", "it2media-domain-crawler", "OrangeBot", "nutch", "netresearchserver", "BUbiNG", "http", "webmon ", "intelium_bot", "web-archive-net.com.bot", "ecosia", "nerdybot", "Applebot", "sitebot", "DuckDuckBot", "integromedb", "FAST Enterprise Crawler", "backlinkcrawler", "openindexspider", "java", "reddit", "seekbot", "Googlebot-Image", "SemanticScholarBot", "linkdex", "webcompanycrawler", "scribdbot", "CyberPatrol", "lssbot", "siteexplorer", "httpunit", "semrush", "ltx71", "exabot", "slurp", "ezooms", "g00g1e.net", "AISearchBot", "ADmantX", "arabot", "Qwantify", "scrape", "mlbot", "fr-crawler", "headless", "Google-Read-Aloud", "speedy", "MegaIndex", "yanga", "SimpleCrawler", "baiduspider", "phpcrawl", "dotbot", "RetrevoPageAnalyzer", "twitter", "ec2linkfinder", "youtube", "bing", "msrbot", "FeedFetcher", "yandex", "domaincrawler", "coccoc", "summify", "GrapeshotCrawler", "gigablast", "webmon", "content crawler spider", "yahoo", "twengabot", "WeSEE", "memorybot", "SemrushBot", "MJ12bot", "seznambot", "googlebot", "bnf.fr_bot", "Google Favicon", "Mail.RU_Bot", "page2rss", "facebookexternalhit", "AddThis", "NerdByNature.Bot", "seokicks-robot", "crawl", "convera", "slack", "elisabot", "A6-Indexer", "crawler4j", "edisterbot", "ahrefsbot", "findlink", "InterfaxScanBot",
	}
)

func init() {
	flag.BoolVar(&verbose, "verbose", false, "Provide usage verbosity")
	flag.BoolVar(&preserve, "preserve", false, "Preserve the bomb")
	flag.StringVar(&file, "filename", default_file, "Name of the file for the bomb")
	flag.IntVar(&size, "size", default_size, "Size of bomb to create")
	flag.IntVar(&port, "port", default_port, "HTTPd port to use")
	flag.Parse()
}

func diffuse_bomb(f string) error {
	if preserve {
		return nil
	}
	if verbose {
		fmt.Printf("[-] Diffusing the bomb: %s\n", f)
	}
	err := os.Remove(f)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func fire_in_the_hole(w http.ResponseWriter, r *http.Request) {
	var ua string = r.UserAgent()

	for j := 0; j < len(useragents); j++ {
		if ua == "" || strings.Contains(ua, useragents[j]) {
			fi, _ := bomb.Stat()
			fsize := fi.Size()
			if verbose {
				fmt.Printf("[*] Serving %d bytes gzipped to %s (UA:\"%s\")\n",
					fsize, r.RemoteAddr, useragents[j])
			}

			w.Header().Add("Content-Encoding", "gzip")
			w.Header().Add("Content-type", "text/html; charset=UTF-8")
			w.Header().Add("Content-Length", fmt.Sprintf("%d", fsize))
			w.WriteHeader(http.StatusOK)

			bomb.Seek(0, os.SEEK_SET)
			var o io.Reader = bomb
			_, err := io.CopyN(w, o, int64(fsize))
			if err != nil {
				fmt.Fprintf(os.Stderr, "[-] Error writing bomb (\"%s\")\n", err)
				return
			}

			return
		}
	}
	w.WriteHeader(http.StatusForbidden)
}

func main() {
	var err error
	c_sig := make(chan os.Signal, 1)
	signal.Notify(c_sig, os.Interrupt, os.Kill)
	defer bomb.Close()
	defer diffuse_bomb(file)

	// create the bomb
	bomb, err = os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		log.Fatal(err)
	}

	s, _ := bomb.Stat()
	if s.Size() <= 0 {
		fmt.Printf("[+] Creating the bomb: %s. Wait...\n", file)

		// use /usr/bin/dd and /us/bin/gzip to create content faster, reliably
		// pipe the output to file
		fill_cmd := fmt.Sprintf("(echo '<html>' && dd if=/dev/zero bs=1M count=%d) | gzip -5 -c --", size)
		cmd := exec.Command("sh", "-c", fill_cmd)
		cmd.Stdout = bomb
		cmd_err, _ := cmd.StderrPipe()
		err = cmd.Start()
		if err != nil {
			error_out, _ := ioutil.ReadAll(cmd_err)
			println(error_out)
			log.Fatal(err)
		}
		err = cmd.Wait()
		if err != nil {
			log.Fatal(err)
		}
		println("[+] Done!")
	} else {
		println("[!] Already exists, using that file.")
	}

	for i := 0; i < len(handle); i++ {
		http.HandleFunc(handle[i], fire_in_the_hole)
	}

	go func() {
		if verbose {
			fmt.Printf("[*] Serving the bomb on port %d/TCP\n", port)
		}
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			fmt.Fprintf(os.Stderr, "[-] Error creating the HTTPd server (\"%s\")\n", err)
			os.Exit(1)
		}
		}()
		sig := <-c_sig
		fmt.Printf("\r[!] Caught signal: %s. Exiting...\n", sig.String())

	}