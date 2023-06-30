package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	prewards "github.com/ingenuity-build/quicksilver/x/participationrewards/types"
	"github.com/ingenuity-build/xcclookup/pkgs/handlers"
	"github.com/ingenuity-build/xcclookup/pkgs/types"
	"gopkg.in/yaml.v3"
)

const (
	// Bech32Prefix defines the Bech32 prefix used for EthAccounts
	Bech32Prefix = "quick"

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = Bech32Prefix
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = Bech32Prefix + sdk.PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)

var (
	GitCommit string
	Version   string = "development"
	Logo             = `
                               .........                                        
                       ..::-----------------::..                                
                   ..::---------------------------.                             
                ..:---------:::::::::::::::-----=-==:.                          
             ..::-------:::::::::::::::::::::::---====-:                        
           ..::------:::::::::::::::::::::::::::::--=-==-:                      
         ..::-----::::::::::::::::::::::::::::::::::---===-.                    
        ..:------:::::::::::::::::::::::::::::::::::::-====-                    
       .::------:::::::::::::::::::::::::::::::::::::::--===                    
     ..::------:::::::::::::::-----:::::::::::::::::::::----                    
     .::-:----::::::::-=++**###%%%%%%#*=:::::::::::::::::--:                    
    .::-::---::::::=+******#%%%%%%#*+=+##+-::::::::::::::::                     
   ..:--:::::::::=+++****#%@%%%%*:     -%=-=++++++===-::::.                     
   .::--:::::::-=++++**#%@@@@%%#.     .*%....+***++++==:..                      
   .:--:::::::-=+++++*%@@@@@@@%%-..:-+#%%:...-***+++=-:.                        
  ..:--::::::-==++++*%@@@@@@@@%%%%%%%%%%%*:.:+**+=-:..                          
  ..:--:::::-===++++%@@@@@@@@@@%%%%%%%%%%%%#*+-::..                             
  ..:---:::::===+++*@@@@@@@@@@@@%%%%%%%#*+-::...              ....              
   .::--:::::-==+++*@@@@@@@@@@@@@%%%*+=--:...             ..::------.           
   ..:-=-:::::-=++++@@@@@@@@@@@%#*===----.              ..::----=-:-=-          
    .:--=::::::-++++*@@@@@@@@#+====-----:              ..:----=--:::-=-         
    ..:--=::::::-++++%@@@@%*====--------.              .:--=+=+-:::::-=:        
     ..:---::::::-+++*%@%+====---------:             ..:+#%*.:*::+++=-==        
      ..:--=-:::::-=+*++====----------:.           .:=#@@@%%%%#*+*++++==        
       ..:--=-::::::-====-----------::.          ..:+**@@@@%%%%###***==:        
         .::--=--::::-=----------::::.          ..:==+*#@@@@%%%%####+=-         
          ..::--=--::----------:::::.           .:--===+#%@@@%%%%#+==-.         
            ..::---=--------::::::..            .::-======+****+===-:.          
              ...::----------::::..              ..:----======---:..            
                 ...::::::::::..                  ...:::::::::::..              
                     ........                         .........
`

	connectionManager    types.CacheManager[prewards.ConnectionProtocolData]
	poolsManager         types.CacheManager[prewards.OsmosisPoolProtocolData]
	osmosisParamsManager types.CacheManager[prewards.OsmosisParamsProtocolData]
	umeeParamsManager    types.CacheManager[prewards.UmeeParamsProtocolData]
	tokenManager         types.CacheManager[prewards.LiquidAllowedDenomProtocolData]
)

func main() {
	fmt.Println(Logo)
	fmt.Printf("Quicksilver - Cross Chain Claims %s (%s)\n", Version, GitCommit)

	var fileName, action string
	flag.StringVar(&fileName, "f", "", "YAML file to parse.")
	flag.StringVar(&action, "a", "serve", "Action - either 'serve' or 'backend'.")
	flag.Parse()

	if fileName == "" {
		fmt.Println("Please provide config file by using -f option")
		return
	}

	yamlFile, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		return
	}

	var cfg types.Config
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		fmt.Printf("Error parsing config file: %s\n", err)
	}

	connectionManager.Init(cfg.SourceLcd+"/quicksilver/participationrewards/v1/protocoldata/ProtocolDataTypeConnection/", time.Minute*5)
	osmosisParamsManager.Init(cfg.SourceLcd+"/quicksilver/participationrewards/v1/protocoldata/ProtocolDataTypeOsmosisParams/", time.Hour*24)
	umeeParamsManager.Init(cfg.SourceLcd+"/quicksilver/participationrewards/v1/protocoldata/ProtocolDataTypeUmeeParams/", time.Hour*24)
	poolsManager.Init(cfg.SourceLcd+"/quicksilver/participationrewards/v1/protocoldata/ProtocolDataTypeOsmosisPool/", time.Minute*5)
	tokenManager.Init(cfg.SourceLcd+"/quicksilver/participationrewards/v1/protocoldata/ProtocolDataTypeLiquidToken/", time.Minute*5)
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)

	switch action {
	case "serve":
		r := mux.NewRouter()
		r.HandleFunc("/cache", handlers.GetCacheHandler(cfg, &connectionManager, &poolsManager, &osmosisParamsManager, &tokenManager))
		r.HandleFunc("/{address}/epoch", handlers.GetEpochHandler(cfg, &connectionManager, &poolsManager, &osmosisParamsManager, &umeeParamsManager, &tokenManager))
		r.HandleFunc("/{address}/current", handlers.GetCurrentHandler(cfg, &connectionManager, &poolsManager, &osmosisParamsManager, &tokenManager))
		// r.HandleFunc("/{address}/airdrop/{claimId}", handlers.AirdropHandler)
		r.HandleFunc("/version", handlers.GetVersionHandler(Version))
		http.Handle("/", r)

		if err := http.ListenAndServe(":8090", nil); err != nil {
			fmt.Printf("Error: %v\n", err.Error())
			return
		}
	case "autoxcc":
		// connect to DB to fetch addresses.
		// periodically requery this list.
		// poll epoch
		// iterate addreses, make query, and submit txs.
		// easy :)
	default:
		fmt.Println("Please specify '-a serve' or '-a backend' command line args")
		return
	}
}
