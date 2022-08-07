package cmd

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/saintmalik/dca-tool/model"
	"github.com/spf13/cobra"
)

//go:embed all:templates/*
var TempFs embed.FS

//go:embed all:config.yaml
var b string

//go:embed all:config.json
var c string

var configYaml = b
var configJSON = c

const (
	foldername = "cmd"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "This command allows you to set DCA Configurations",
	Long:  `This command allows you to set your Binance API and Secret and other Configurations`,
	Run: func(cmd *cobra.Command, args []string) {
		setApi, _ := cmd.Flags().GetString("credapi")
		if setApi == "reset" {
			openbrowser("http://localhost:4046/api")
			http.HandleFunc("/api", creds)
			fmt.Println("Starting Server to set Binance API and Secret")
			panic(http.ListenAndServe("localhost:4046", nil))
			} else {
			 	main()
		}
	},
}

var tmpl *template.Template

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.PersistentFlags().String("credapi", "reset", "Set your Binance API and SecretKey credentials")
	tmpl = template.Must(template.ParseFS(TempFs, "templates/*.html"))
}

func main() {
	openbrowser("http://localhost:4046/")
	http.HandleFunc("/", form)
	fmt.Println("Starting Server to set Configurations")
	panic(http.ListenAndServe("localhost:4046", nil))
}

func creds(w http.ResponseWriter, r *http.Request) {
	model.Bapi = r.FormValue("bapi")
	model.Bsecret = r.FormValue("bsecret") //picking up the value from the form

	f, err := os.Create(filepath.Join(foldername, configYaml))
	if err != nil {
		fmt.Println(err)
		return
	}

	l, err := f.WriteString(`{"api":"` + model.Bapi + `","secretkey":"` + model.Bsecret + `"}`)
	if err != nil {
		fmt.Println("Error Writing to the Config.yaml file", err)
		f.Close()
		return
	}

	fmt.Println(l, "Your Binance API and Secret are set, dont let anyone access this file on your computer")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "creds.html", nil)
	if err != nil {
		log.Fatal("Error loading index template: ", err)
	}
}

func form(w http.ResponseWriter, r *http.Request) {
	model.Coinid = r.FormValue("coinid")
	model.Amount = r.FormValue("amount")
	model.Alloted = r.FormValue("allottedpercent")
	model.Buyinterval = r.FormValue("buyingintervals")
	model.Fee = r.FormValue("feepercent")
	model.Testvalue = r.FormValue("testing") //picking up the value from the form

	f, err := os.Create(filepath.Join(foldername, configJSON))
	if err != nil {
		fmt.Println(err)
		return
	}

	l, err := f.WriteString(`{"coins":"` + model.Coinid + `","amount":"` + model.Amount + `","percent":"` + model.Alloted + `","fee":"` + model.Fee + `","testing":"` + model.Testvalue + `","buyintervals":"` + model.Buyinterval + `"}`)
	if err != nil {
		fmt.Println("Your Binance API and Secret are set, dont let anyone access this file on you", err)
		f.Close()
		return
	}

	fmt.Println(l, "Your DCA Configurations are set, quit the server and and run dca run")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = tmpl.ExecuteTemplate(w, "creds.html", nil)
	if err != nil {
		log.Fatal("Error loading index template: ", err)
	}
}

func openbrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}
