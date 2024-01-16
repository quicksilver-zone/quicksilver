package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
	"gopkg.in/yaml.v3"

	"github.com/ingenuity-build/xcclookup/pkgs/handlers"
	"github.com/ingenuity-build/xcclookup/pkgs/types"
)

const (
	// Bech32Prefix defines the Bech32 prefix used for EthAccounts.
	Bech32Prefix = "quick"

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address.
	Bech32PrefixAccAddr = Bech32Prefix
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key.
	Bech32PrefixAccPub = Bech32Prefix + sdk.PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address.
	Bech32PrefixValAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key.
	Bech32PrefixValPub = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address.
	Bech32PrefixConsAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key.
	Bech32PrefixConsPub = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)

var (
	GitCommit string
	Version   = "development"
	Logo      = `
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

	connectionManager     types.CacheManager[prewards.ConnectionProtocolData]
	osmosisPoolsManager   types.CacheManager[prewards.OsmosisPoolProtocolData]
	crescentPoolsManager  types.CacheManager[prewards.CrescentPoolProtocolData]
	osmosisParamsManager  types.CacheManager[prewards.OsmosisParamsProtocolData]
	umeeParamsManager     types.CacheManager[prewards.UmeeParamsProtocolData]
	crescentParamsManager types.CacheManager[prewards.CrescentParamsProtocolData]
	tokenManager          types.CacheManager[prewards.LiquidAllowedDenomProtocolData]
	zonesManager          types.CacheManager[icstypes.Zone]
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

	ctx := context.Background()
	connectionManager.Init(ctx, cfg.SourceLcd+"/quicksilver/participationrewards/v1/protocoldata/ProtocolDataTypeConnection/", types.DataTypeProtocolData, time.Minute*5)
	osmosisParamsManager.Init(ctx, cfg.SourceLcd+"/quicksilver/participationrewards/v1/protocoldata/ProtocolDataTypeOsmosisParams/", types.DataTypeProtocolData, time.Hour*24)
	umeeParamsManager.Init(ctx, cfg.SourceLcd+"/quicksilver/participationrewards/v1/protocoldata/ProtocolDataTypeUmeeParams/", types.DataTypeProtocolData, time.Hour*24)
	crescentParamsManager.Init(ctx, cfg.SourceLcd+"/quicksilver/participationrewards/v1/protocoldata/ProtocolDataTypeCrescentParams/", types.DataTypeProtocolData, time.Hour*24)
	osmosisPoolsManager.Init(ctx, cfg.SourceLcd+"/quicksilver/participationrewards/v1/protocoldata/ProtocolDataTypeOsmosisPool/", types.DataTypeProtocolData, time.Minute*5)
	crescentPoolsManager.Init(ctx, cfg.SourceLcd+"/quicksilver/participationrewards/v1/protocoldata/ProtocolDataTypeCrescentPool/", types.DataTypeProtocolData, time.Minute*5)
	tokenManager.Init(ctx, cfg.SourceLcd+"/quicksilver/participationrewards/v1/protocoldata/ProtocolDataTypeLiquidToken/", types.DataTypeProtocolData, time.Minute*5)
	zonesManager.Init(ctx, cfg.SourceLcd+"/quicksilver/interchainstaking/v1/zones", types.DataTypeZone, time.Hour*24)
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)

	switch action {
	case "serve":
		r := mux.NewRouter()
		r.HandleFunc("/cache", handlers.GetCacheHandler(ctx, cfg, &connectionManager, &osmosisPoolsManager, &osmosisParamsManager, &tokenManager))
		r.HandleFunc("/{address}/epoch", handlers.GetEpochHandler(ctx, cfg, &connectionManager, &osmosisPoolsManager, &crescentPoolsManager, &osmosisParamsManager, &umeeParamsManager, &crescentParamsManager, &tokenManager, &zonesManager))
		r.HandleFunc("/{address}/current", handlers.GetCurrentHandler(ctx, cfg, &connectionManager, &osmosisPoolsManager, &crescentPoolsManager, &osmosisParamsManager, &umeeParamsManager, &crescentParamsManager, &tokenManager, &zonesManager))
		// r.HandleFunc("/{address}/airdrop/{claimId}", handlers.AirdropHandler)
		r.HandleFunc("/version", handlers.GetVersionHandler(Version))
		http.Handle("/", r)

		server := &http.Server{
			Addr:              ":8090",
			Handler:           r,
			ReadHeaderTimeout: 10 * time.Second,
		}

		if err := server.ListenAndServe(); err != nil {
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
