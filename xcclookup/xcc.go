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
	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/logger"
	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/services"
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
              ...::----------::::..            .::-======+****+===-:.          
                 ...::::::::::..                  ...:::::::::::..              
                     ........                         .........
`
	cacheMgr = types.NewCacheManager()
)

func main() {
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
		return
	}

	// Initialize logger with config
	log := logger.New(cfg.Logging.GetLogLevel())
	log.Info("Starting Quicksilver Cross Chain Claims", "version", Version, "commit", GitCommit)
	log.Debug("Logo", "logo", Logo)

	// Create context with logger
	ctx := context.Background()
	ctx = logger.WithLogger(ctx, log)

	// Initialize cache manager
	err = cacheMgr.Add(ctx, &types.Cache[prewards.ConnectionProtocolData]{}, cfg.SourceLcd+"/quicksilver/participationrewards/v1/protocoldata/ProtocolDataTypeConnection/", types.DataTypeProtocolData, time.Minute*5)
	if err != nil {
		log.Error("Failed to add connection cache", "error", err)
		return
	}
	err = cacheMgr.Add(ctx, &types.Cache[prewards.OsmosisParamsProtocolData]{}, cfg.SourceLcd+"/quicksilver/participationrewards/v1/protocoldata/ProtocolDataTypeOsmosisParams/", types.DataTypeProtocolData, time.Hour*24)
	if err != nil {
		log.Error("Failed to add osmosis params cache", "error", err)
		return
	}
	err = cacheMgr.Add(ctx, &types.Cache[prewards.UmeeParamsProtocolData]{}, cfg.SourceLcd+"/quicksilver/participationrewards/v1/protocoldata/ProtocolDataTypeUmeeParams/", types.DataTypeProtocolData, time.Hour*24)
	if err != nil {
		log.Error("Failed to add umee params cache", "error", err)
		return
	}
	err = cacheMgr.Add(ctx, &types.Cache[prewards.MembraneProtocolData]{}, cfg.SourceLcd+"/quicksilver/participationrewards/v1/protocoldata/ProtocolDataTypeMembraneParams/", types.DataTypeProtocolData, time.Hour*24)
	if err != nil {
		log.Error("Failed to add membrane params cache", "error", err)
		return
	}
	err = cacheMgr.Add(ctx, &types.Cache[prewards.OsmosisPoolProtocolData]{}, cfg.SourceLcd+"/quicksilver/participationrewards/v1/protocoldata/ProtocolDataTypeOsmosisPool/", types.DataTypeProtocolData, time.Minute*5)
	if err != nil {
		log.Error("Failed to add osmosis pool cache", "error", err)
		return
	}
	err = cacheMgr.Add(ctx, &types.Cache[prewards.OsmosisClPoolProtocolData]{}, cfg.SourceLcd+"/quicksilver/participationrewards/v1/protocoldata/ProtocolDataTypeOsmosisCLPool/", types.DataTypeProtocolData, time.Minute*5)
	if err != nil {
		log.Error("Failed to add osmosis CL pool cache", "error", err)
		return
	}
	err = cacheMgr.Add(ctx, &types.Cache[prewards.LiquidAllowedDenomProtocolData]{}, cfg.SourceLcd+"/quicksilver/participationrewards/v1/protocoldata/ProtocolDataTypeLiquidToken/", types.DataTypeProtocolData, time.Minute*5)
	if err != nil {
		log.Error("Failed to add liquid denom cache", "error", err)
		return
	}
	err = cacheMgr.Add(ctx, &types.Cache[icstypes.Zone]{}, cfg.SourceLcd+"/quicksilver/interchainstaking/v1/zones", types.DataTypeZone, time.Hour*24)
	if err != nil {
		log.Error("Failed to add zone cache", "error", err)
		return
	}

	log.Debug("Cache manager initialized", "source_lcd", cfg.SourceLcd)

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)

	types.AddMocks(ctx, &cacheMgr, cfg.Mocks.OsmosisPools)
	types.AddMocks(ctx, &cacheMgr, cfg.Mocks.Connections)
	types.AddMocks(ctx, &cacheMgr, cfg.Mocks.UmeeParams)
	types.AddMocks(ctx, &cacheMgr, cfg.Mocks.MembraneParams)

	r := mux.NewRouter()
	connections, err := types.GetCache[prewards.ConnectionProtocolData](ctx, &cacheMgr)
	if err != nil {
		log.Error("Failed to get connections cache", "error", err)
		return
	}

	// Create services
	versionService := services.NewVersionService(&types.VersionService{})

	// Create claims service
	claimsService := services.NewClaimsService(cfg, &cacheMgr)

	r.HandleFunc("/cache", handlers.GetCacheHandler(ctx, cfg, &cacheMgr))
	r.HandleFunc("/{address}/epoch", handlers.GetAssetsHandler(ctx, cfg, &cacheMgr, claimsService, types.GetHeights(connections), types.OutputEpoch))
	r.HandleFunc("/{address}/current", handlers.GetAssetsHandler(ctx, cfg, &cacheMgr, claimsService, types.GetZeroHeights(connections), types.OutputCurrent))
	r.HandleFunc("/version", handlers.GetVersionHandler(versionService))
	http.Handle("/", r)

	bindPort := 8090
	if cfg.BindPort != 0 {
		bindPort = cfg.BindPort
	}

	log.Info("Starting HTTP server", "port", bindPort)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", bindPort),
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("Server error", "error", err.Error())
		return
	}
}
