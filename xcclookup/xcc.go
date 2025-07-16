package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"

	sdk "github.com/cosmos/cosmos-sdk/types"

	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"

	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/handlers"
	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/types"
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
	Logo = `
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
	cacheMgr = types.NewCacheManager()
)

func main() {
	fmt.Println(Logo)
	version, err := types.GetVersion()
	if err != nil {
		fmt.Printf("Error getting version: %s\n", err)
		return
	}
	fmt.Printf("xcclookup (Cross Chain Claims) %s\n", string(version))

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
	cacheMgr.Add(ctx, &types.Cache[prewards.ConnectionProtocolData]{}, cfg.SourceLcd+"/quicksilver/participationrewards/v1/protocoldata/ProtocolDataTypeConnection/", types.DataTypeProtocolData, time.Minute*5)
	cacheMgr.Add(ctx, &types.Cache[prewards.OsmosisParamsProtocolData]{}, cfg.SourceLcd+"/quicksilver/participationrewards/v1/protocoldata/ProtocolDataTypeOsmosisParams/", types.DataTypeProtocolData, time.Hour*24)
	cacheMgr.Add(ctx, &types.Cache[prewards.UmeeParamsProtocolData]{}, cfg.SourceLcd+"/quicksilver/participationrewards/v1/protocoldata/ProtocolDataTypeUmeeParams/", types.DataTypeProtocolData, time.Hour*24)
	cacheMgr.Add(ctx, &types.Cache[prewards.OsmosisPoolProtocolData]{}, cfg.SourceLcd+"/quicksilver/participationrewards/v1/protocoldata/ProtocolDataTypeOsmosisPool/", types.DataTypeProtocolData, time.Minute*5)
	cacheMgr.Add(ctx, &types.Cache[prewards.OsmosisClPoolProtocolData]{}, cfg.SourceLcd+"/quicksilver/participationrewards/v1/protocoldata/ProtocolDataTypeOsmosisCLPool/", types.DataTypeProtocolData, time.Minute*5)
	cacheMgr.Add(ctx, &types.Cache[prewards.LiquidAllowedDenomProtocolData]{}, cfg.SourceLcd+"/quicksilver/participationrewards/v1/protocoldata/ProtocolDataTypeLiquidToken/", types.DataTypeProtocolData, time.Minute*5)
	cacheMgr.Add(ctx, &types.Cache[icstypes.Zone]{}, cfg.SourceLcd+"/quicksilver/interchainstaking/v1/zones", types.DataTypeZone, time.Hour*24)
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)

	types.AddMocks(ctx, &cacheMgr, cfg.Mocks.OsmosisPools)
	types.AddMocks(ctx, &cacheMgr, cfg.Mocks.Connections)
	types.AddMocks(ctx, &cacheMgr, cfg.Mocks.UmeeParams)

	r := mux.NewRouter()
	connections := types.GetCache[prewards.ConnectionProtocolData](ctx, &cacheMgr)

	r.HandleFunc("/cache", handlers.GetCacheHandler(ctx, cfg, &cacheMgr))
	r.HandleFunc("/{address}/epoch", handlers.GetAssetsHandler(ctx, cfg, &cacheMgr, types.GetHeights(connections), types.OutputEpoch))
	r.HandleFunc("/{address}/current", handlers.GetAssetsHandler(ctx, cfg, &cacheMgr, types.GetZeroHeights(connections), types.OutputCurrent))
	r.HandleFunc("/version", handlers.GetVersionHandler())
	http.Handle("/", r)

	bindPort := 8090
	if cfg.BindPort != 0 {
		bindPort = cfg.BindPort
	}
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", bindPort),
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return
	}
}
