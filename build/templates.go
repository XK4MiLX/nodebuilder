package build
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"text/template"
	"time"
        "net/http"
        "net/url"
	"net"
	"strings"
)

// Backend contains backend specific fields
type Backend struct {
	PackageName                     string             `json:"package_name"`
	BinaryName                      string             `json:"binary_name"`
	Version                         string             `json:"version"`
	BinaryURL                       string             `json:"binary_url"`
	VerificationType                string             `json:"verification_type"`
	VerificationSource              string             `json:"verification_source"`
	ExtractCommand                  string             `json:"extract_command"`
	ExcludeFiles                    []string           `json:"exclude_files"`
	ExecCommandTemplate             string             `json:"exec_command_template"`
	ExecScript                      string             `json:"exec_script"`
	PostinstScriptTemplate          string             `json:"postinst_script_template"`
	ServiceType                     string             `json:"service_type"`
	Mainnet                         bool               `json:"mainnet"`
	NodeType                        string             `json:"node_type"`
	Masternode  			Masternode 	   `json:"masternode"`
	Bootstrap  			Bootstrap 	   `json:"bootstrap"`
	Logs  				Logs	 	   `json:"logs"`
	Corruption			Corruption	   `json:"corruption"`
	Healthcheck			Healthcheck 	   `json:"healthcheck"`
	ServerConfigFile                string             `json:"server_config_file"`
	AdditionalParams                interface{}        `json:"additional_params,omitempty"`
	Platforms                       map[string]Backend `json:"platforms,omitempty"`
}

type Bootstrap struct {
	Enabled  			bool		    `json:"enabled"`
	RemoteSource	  	  	[]string  	    `json:"remote_source"`
	UnpackPathTemplate          	string              `json:"unpack_path_template"`
}

type Logs struct {
	CronExpression  		string		    `json:"cron_expression"`
	SizeLimit	  	  	int 		    `json:"size_limit"`
	FilesTemplate          		string              `json:"files_template"`
	LogsRedirect			bool 		    `json:"logs_redirect"`
}

type Corruption struct {
	Detection   			bool		    `json:"detection"`
	CronExpression			string 		    `json:"cron_expression"`
	LogFilePathTemplate  	  	string 		    `json:"log_file_path_template"`
	CorruptionKeywords		[]string  	    `json:"corruption_keywords"`
	Bootstrap   			bool		    `json:"bootstrap"`
}

type Healthcheck struct {
	ExplorerGetBlockCount   	[]string 	    `json:"explorer_get_block_count"`
	LocalGetBlockCountTemplate    	string 		    `json:"local_get_block_count_template"`
	MaxDifference			int		    `json:"max_difference"`	
}

type Masternode struct {
	KeyValue    	string `json:"-"`
	KeyWord     	string `json:"key_word"`
}

// Config contains the structure of the config
type Config struct {
	Coin struct {
		Name     string `json:"name"`
		Shortcut string `json:"shortcut"`
		Label    string `json:"label"`
		Alias    string `json:"alias"`
	} `json:"coin"`
	Ports struct {
		BackendRPC          int `json:"backend_rpc"`
		BackendMessageQueue int `json:"backend_message_queue"`
		BackendP2P          int `json:"backend_p2p"`
		BackendHttp         int `json:"backend_http"`
		BackendAuthRpc      int `json:"backend_authrpc"`
	} `json:"ports"`
	IPC struct {
		RPCURLTemplate              string `json:"rpc_url_template"`
		RPCUser                     string `json:"rpc_user"`
		RPCPass                     string `json:"rpc_pass"`
		RPCTimeout                  int    `json:"rpc_timeout"`
		MessageQueueBindingTemplate string `json:"message_queue_binding_template"`
		GenerateAuth		    bool   `json:"generate_auth"`
	} `json:"ipc"`
	Backend   Backend `json:"backend"`
	Meta struct {
		BuildDatetime          string `json:"-"` // generated field
		PackageMaintainer      string `json:"package_maintainer"`
		PackageMaintainerEmail string `json:"package_maintainer_email"`
	} `json:"meta"`
	Env struct {
		Version              string `json:"version"`
		BackendInstallPath   string `json:"backend_install_path"`
		BackendDataPath      string `json:"backend_data_path"`
		Architecture         string `json:"architecture"`
	} `json:"-"`
}


func jsonToString(msg json.RawMessage) (string, error) {
	d, err := msg.MarshalJSON()
	if err != nil {
		return "", err
	}
	
	return string(d), nil
}

func getEnv(key string) string {
	ge := os.Getenv(key)
	return string(ge)
}


func generateRandom() (string, error) {
	cmd := exec.Command("/usr/bin/env", "bash", "-c", "pwgen -1 18 -n")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func generateRPCAuth(user, pass string) (string, error) {
	cmd := exec.Command("/usr/bin/env", "bash", "-c", "build/scripts/rpcauth.py \"$0\" \"$1\" | sed -n -e 2p", user, pass)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	
	return out.String(), nil
}

func arrayToString(arr []string) string {
   return strings.Join([]string(arr), "|")
}

// ParseTemplate parses the template
func (c *Config) ParseTemplate() *template.Template {
	templates := map[string]string{
		"IPC.RPCURLTemplate":                      		c.IPC.RPCURLTemplate,
		"IPC.MessageQueueBindingTemplate":         		c.IPC.MessageQueueBindingTemplate,
		"Backend.ExecCommandTemplate":				c.Backend.ExecCommandTemplate,
		"Backend.PostinstScriptTemplate":          		c.Backend.PostinstScriptTemplate,
		"Backend.Healthcheck.LocalGetBlockCountTemplate": 	c.Backend.Healthcheck.LocalGetBlockCountTemplate,
		"Backend.Corruption.LogFilePathTemplate": 		c.Backend.Corruption.LogFilePathTemplate,
		"Backend.Logs.FilesTemplate": 				c.Backend.Logs.FilesTemplate,
		"Backend.Bootstrap.UnpackPathTemplate": 		c.Backend.Bootstrap.UnpackPathTemplate,
	}

	funcMap := template.FuncMap{
		"jsonToString":    	jsonToString,
		"generateRPCAuth": 	generateRPCAuth,
		"containsString":  	strings.Contains,
		"arrayToString":   	arrayToString,
		"generateRandom": 	generateRandom,
		"getPublicIP": 		getPublicIP,
		"getEnv":		getEnv,			
	}

	t := template.New("").Funcs(funcMap)

	for name, def := range templates {
		t = template.Must(t.Parse(fmt.Sprintf(`{{define "%s"}}%s{{end}}`, name, def)))
	}

	return t
}

func copyNonZeroBackendFields(toValue *Backend, fromValue *Backend) {
	from := reflect.ValueOf(*fromValue)
	to := reflect.ValueOf(toValue).Elem()
	for i := 0; i < from.NumField(); i++ {
		if from.Field(i).IsValid() && !from.Field(i).IsZero() {
			to.Field(i).Set(from.Field(i))
		}
	}
}

func checkIPAddress(ip string) bool {
    if net.ParseIP(ip) == nil {
        return false
    }
	
   return true
}

func printResponse(url string) (retErr error, result string) {
    resp, err := http.Get(url)
    if err != nil {
        return err, ""
    }
    defer func() {
        err := resp.Body.Close()
        if err != nil && retErr == nil {
            retErr = err
        }
    }()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err, ""
    }
    ip := strings.Trim(string(body[:]), "\n")
	
    return nil, ip
}

func getPublicIP() string {
        var PublicIP string
        var urls = []string{
                "https://api4.my-ip.io/ip",
                "https://checkip.amazonaws.com",
                "https://api.ipify.org",
        }
        for _, url := range urls {
                err, result := printResponse(url)
                if err == nil {
                  if checkIPAddress(result) == true {
                     PublicIP = result
                     break
                  }
                }
        }
	
	 return PublicIP
}

func isValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}
	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	
	return true
}

// LoadConfig loads the config files
func LoadConfig(configsDir, coin string, url string) (*Config, error) {
	config := new(Config)

	var d *json.Decoder

        if isValidUrl(url) == false {
        	f, err := os.Open(filepath.Join(configsDir, "coins", coin + ".json"))
		if err != nil {
		    return  nil, fmt.Errorf("Config file for %s not found, create config or use URL", coin)
	       	 }
		d = json.NewDecoder(f)
	} else {
    		resp, bas := http.Get(url)
    		if bas != nil {
        		 return nil, bas
    		}
    		defer resp.Body.Close()
                d = json.NewDecoder(resp.Body)
	}
	
	err := d.Decode(config)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(filepath.Join(configsDir, "environ.json"))
	if err != nil {
		return nil, err
	}
	d = json.NewDecoder(f)
	err = d.Decode(&config.Env)
	if err != nil {
		return nil, err
	}

	config.Meta.BuildDatetime = time.Now().Format("Mon, 02 Jan 2006 15:04:05 -0700")
	config.Env.Architecture = runtime.GOARCH

	if !isEmpty(config, "backend") {
		// set platform specific fields to config
		platform, found := config.Backend.Platforms[runtime.GOARCH]
		if found {
			copyNonZeroBackendFields(&config.Backend, &platform)
		}

		//switch config.Backend.ServiceType {
		//case "forking":
		//case "simple":
		//default:
		//	return nil, fmt.Errorf("Invalid service type: %s", config.Backend.ServiceType)
		//}

		switch config.Backend.VerificationType {
		case "":
		case "gpg":
		case "sha256":
		case "gpg-sha256":
		default:
			return nil, fmt.Errorf("Invalid verification type: %s", config.Backend.VerificationType)
		}
	}
	
	return config, nil
}

func isEmpty(config *Config, target string) bool {
	switch target {
	case "backend":
		return config.Backend.PackageName == ""
	default:
		panic("Invalid target name: " + target)
	}
}

func RemoveContents(dir string) error {
    files, err := filepath.Glob(filepath.Join(dir, "*"))
    if err != nil {
        return err
    }
    for _, file := range files {
        err = os.RemoveAll(file)
        if err != nil {
            return err
        }
    }
	
    return nil
}

// GeneratePackageDefinitions generate the package definitions from the config
func GeneratePackageDefinitions(config *Config, templateDir, outputDir , coin string) error {
	templ := config.ParseTemplate()
	err := makeOutputDir(outputDir)
	if err != nil {
		return err
	}

	for _, subdir := range []string{"backend"} {
		if isEmpty(config, subdir) {
			continue
		}

		root := filepath.Join(templateDir, subdir)

		err = os.Mkdir(filepath.Join(outputDir, subdir), 0755)
		if err != nil {
			return err
		}

		err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("%s: %s", path, err)
			}

			if path == root {
				return nil
			}
			if filepath.Base(path)[0] == '.' {
				return nil
			}

			subpath := path[len(root)-len(subdir):]

			if info.IsDir() {
				err = os.Mkdir(filepath.Join(outputDir, subpath), info.Mode())
				if err != nil {
					return fmt.Errorf("%s: %s", path, err)
				}
				return nil
			}

			t := template.Must(templ.Clone())
			t = template.Must(t.ParseFiles(path))
			err = writeTemplate(filepath.Join(outputDir, subpath), info, t, config)
			if err != nil {
				return fmt.Errorf("%s: %s", path, err)
			}

			return nil
		})
		if err != nil {
			return err
		}
	}

	if !isEmpty(config, "backend") {
		if err := writeBackendServerConfigFile(config, outputDir, coin); err != nil {
			return err
		}
		if err := writeBackendExecScript(config, outputDir); err != nil {
			return err
		}
	}
	
	return nil
}

func makeOutputDir(path string) error {
	err := os.RemoveAll(path)
	if err == nil {
		err = os.Mkdir(path, 0755)
	}
	return err
}

func writeTemplate(path string, info os.FileInfo, templ *template.Template, config *Config) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer f.Close()

	return templ.ExecuteTemplate(f, "main", config)
}

func writeBackendServerConfigFile(config *Config, outputDir, coin string) error {
	out, err := os.OpenFile(filepath.Join(outputDir, "backend/server.conf"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	if config.Backend.ServerConfigFile == "" {
		return nil
	} else {
		in, err := os.Open(filepath.Join(outputDir, "backend/config", config.Backend.ServerConfigFile))
		if err != nil {
			return err
		}
		defer in.Close()

		_, err = io.Copy(out, in)
		if err != nil {
			return err
		}
	}
	
	return nil
}


func writeBackendExecScript(config *Config, outputDir string) error {

        if config.Backend.ExecScript == "" {
		return nil
	}

	out, err := os.OpenFile(filepath.Join(outputDir, "backend/exec.sh"), os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer out.Close()

	in, err := os.Open(filepath.Join(outputDir, "backend/scripts", config.Backend.ExecScript))
	if err != nil {
		return err
	}
	defer in.Close()

	_, err = io.Copy(out, in)
        if err != nil {
                return err
        }

        return nil

}
