# geth

源码

https://github.com/ethereum/go-ethereum/blob/release/1.15/cmd/geth/main.go

https://github.com/ethereum/go-ethereum/blob/release/1.15/cmd/geth/config.go

https://github.com/ethereum/go-ethereum/blob/release/1.15/cmd/utils/flags.go

cmd/geth/main.go

```go
// geth is a command-line client for Ethereum.
package main

import (
	...省略...
    // 依赖第三方包来构建命令行工具
	"github.com/urfave/cli/v2"
)

const (
	clientIdentifier = "geth" // Client identifier to advertise over the network
)

var (
	// flags that configure the node
	...省略...

var app = flags.NewApp("the go-ethereum command line interface")

func init() {
	// Initialize the CLI app and start Geth
    // 调用 geth() 函数，如果用户没有输入其他子命令就调用这个字段指向的函数
	app.Action = geth
    
    // 对应cmd终端输入geth -h，显示按首字母排序
    // 所有支持的子命令
	app.Commands = []*cli.Command{
		// See chaincmd.go:
		initCommand,
		importCommand,
		exportCommand,
		importHistoryCommand,
		exportHistoryCommand,
		importPreimagesCommand,
		removedbCommand,
		dumpCommand,
		dumpGenesisCommand,
		// See accountcmd.go:
		accountCommand,
		walletCommand,
		// See consolecmd.go:
		consoleCommand,
		attachCommand,
		javascriptCommand,
		// See misccmd.go:
		versionCommand,
		versionCheckCommand,
		licenseCommand,
		// See config.go
		dumpConfigCommand,
		// see dbcmd.go
		dbCommand,
		// See cmd/utils/flags_legacy.go
		utils.ShowDeprecated,
		// See snapshot.go
		snapshotCommand,
		// See verkle.go
		verkleCommand,
	}
	if logTestCommand != nil {
		app.Commands = append(app.Commands, logTestCommand)
	}
    // Sort对所有子命令按首字母排序
	sort.Sort(cli.CommandsByName(app.Commands))
	// 所有能够解析的options
	app.Flags = slices.Concat(
		nodeFlags,
		rpcFlags,
		consoleFlags,
		debug.Flags,
		metricsFlags,
	)
	flags.AutoEnvVars(app.Flags, "GETH")

    // 在所有命令执行前调用
	app.Before = func(ctx *cli.Context) error {
        // 设置用于当前程序的CPU上限
		maxprocs.Set() // Automatically set GOMAXPROCS to match Linux container CPU quota.
		flags.MigrateGlobalFlags(ctx)
		if err := debug.Setup(ctx); err != nil {
			return err
		}
		flags.CheckEnvVars(ctx, app.Flags, "GETH")
		return nil
	}
	app.After = func(ctx *cli.Context) error {
		debug.Exit()
		prompt.Stdin.Close() // Resets terminal mode.
		return nil
	}
}

func main() {
    // 程序的入口点，启动一个解析 command line命令的工具
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// prepare manipulates memory cache allowance and setups metric system.
// This function should be called before launching devp2p stack.
// 初始化前的预处理
func prepare(ctx *cli.Context) {
	// If we're running a known preset, log it for convenience.
	switch {
	case ctx.IsSet(utils.SepoliaFlag.Name):
		log.Info("Starting Geth on Sepolia testnet...")

	case ctx.IsSet(utils.HoleskyFlag.Name):
		log.Info("Starting Geth on Holesky testnet...")

	case ctx.IsSet(utils.DeveloperFlag.Name):
		log.Info("Starting Geth in ephemeral dev mode...")
		log.Warn(`You are running Geth in --dev mode. Please note the following:

  1. This mode is only intended for fast, iterative development without assumptions on
     security or persistence.
  2. The database is created in memory unless specified otherwise. Therefore, shutting down
     your computer or losing power will wipe your entire block data and chain state for
     your dev environment.
  3. A random, pre-allocated developer account will be available and unlocked as
     eth.coinbase, which can be used for testing. The random dev account is temporary,
     stored on a ramdisk, and will be lost if your machine is restarted.
  4. Mining is enabled by default. However, the client will only seal blocks if transactions
     are pending in the mempool. The miner's minimum accepted gas price is 1.
  5. Networking is disabled; there is no listen-address, the maximum number of peers is set
     to 0, and discovery is disabled.
`)

	case !ctx.IsSet(utils.NetworkIdFlag.Name):
		log.Info("Starting Geth on Ethereum mainnet...")
	}
	// If we're a full node on mainnet without --cache specified, bump default cache allowance
	if !ctx.IsSet(utils.CacheFlag.Name) && !ctx.IsSet(utils.NetworkIdFlag.Name) {
		// Make sure we're not on any supported preconfigured testnet either
		if !ctx.IsSet(utils.HoleskyFlag.Name) &&
			!ctx.IsSet(utils.SepoliaFlag.Name) &&
			!ctx.IsSet(utils.DeveloperFlag.Name) {
			// Nope, we're really on mainnet. Bump that cache up!
			log.Info("Bumping default cache on mainnet", "provided", ctx.Int(utils.CacheFlag.Name), "updated", 4096)
			ctx.Set(utils.CacheFlag.Name, strconv.Itoa(4096))
		}
	}
}

// geth is the main entry point into the system if no special subcommand is run.
// It creates a default node based on the command line arguments and runs it in
// blocking mode, waiting for it to be shut down.
// 如果没有指定其他子命令，这个geth就是默认的系统入口
// 以阻塞的方式运行节点，直到节点被关闭
func geth(ctx *cli.Context) error {
    // 检查无效子命令
	if args := ctx.Args().Slice(); len(args) > 0 {
		return fmt.Errorf("invalid command: %q", args[0])
	}
	// 预处理
	prepare(ctx)
    // 创建并配置一个完整的以太坊节点
	stack := makeFullNode(ctx)
	defer stack.Close()
	// 启动节点
	startNode(ctx, stack, false)
    // 等待，直到节点关闭
	stack.Wait()
	return nil
}

// startNode boots up the system node and all registered protocols, after which
// it starts the RPC/IPC interfaces and the miner.
// 启动节点
func startNode(ctx *cli.Context, stack *node.Node, isConsole bool) {
	// Start up the node itself
	utils.StartNode(ctx, stack, isConsole)

	if ctx.IsSet(utils.UnlockedAccountFlag.Name) {
		log.Warn(`The "unlock" flag has been deprecated and has no effect`)
	}

	// Register wallet event handlers to open and auto-derive wallets
	events := make(chan accounts.WalletEvent, 16)
	stack.AccountManager().Subscribe(events)

	// Create a client to interact with local geth node.
	rpcClient := stack.Attach()
	ethClient := ethclient.NewClient(rpcClient)

	go func() {
		// Open any wallets already attached
		for _, wallet := range stack.AccountManager().Wallets() {
			if err := wallet.Open(""); err != nil {
				log.Warn("Failed to open wallet", "url", wallet.URL(), "err", err)
			}
		}
		// Listen for wallet event till termination
		for event := range events {
			switch event.Kind {
			case accounts.WalletArrived:
				if err := event.Wallet.Open(""); err != nil {
					log.Warn("New wallet appeared, failed to open", "url", event.Wallet.URL(), "err", err)
				}
			case accounts.WalletOpened:
				status, _ := event.Wallet.Status()
				log.Info("New wallet appeared", "url", event.Wallet.URL(), "status", status)

				var derivationPaths []accounts.DerivationPath
				if event.Wallet.URL().Scheme == "ledger" {
					derivationPaths = append(derivationPaths, accounts.LegacyLedgerBaseDerivationPath)
				}
				derivationPaths = append(derivationPaths, accounts.DefaultBaseDerivationPath)

				event.Wallet.SelfDerive(derivationPaths, ethClient)

			case accounts.WalletDropped:
				log.Info("Old wallet dropped", "url", event.Wallet.URL())
				event.Wallet.Close()
			}
		}
	}()

	// Spawn a standalone goroutine for status synchronization monitoring,
	// close the node when synchronization is complete if user required.
	if ctx.Bool(utils.ExitWhenSyncedFlag.Name) {
		go func() {
			sub := stack.EventMux().Subscribe(downloader.DoneEvent{})
			defer sub.Unsubscribe()
			for {
				event := <-sub.Chan()
				if event == nil {
					continue
				}
				done, ok := event.Data.(downloader.DoneEvent)
				if !ok {
					continue
				}
				if timestamp := time.Unix(int64(done.Latest.Time), 0); time.Since(timestamp) < 10*time.Minute {
					log.Info("Synchronisation completed", "latestnum", done.Latest.Number, "latesthash", done.Latest.Hash(),
						"age", common.PrettyAge(timestamp))
					stack.Close()
				}
			}
		}()
	}
}
```

```go
 // 创建并配置一个完整的以太坊节点
stack := makeFullNode(ctx)

// config.go
// makeFullNode loads geth configuration and creates the Ethereum backend.
// 根据配置创建全节点
func makeFullNode(ctx *cli.Context) *node.Node {
    // 解析命令行参数和配置文件，生成节点实例
	stack, cfg := makeConfigNode(ctx)
	if ctx.IsSet(utils.OverrideCancun.Name) {
		v := ctx.Uint64(utils.OverrideCancun.Name)
		cfg.Eth.OverrideCancun = &v
	}
	if ctx.IsSet(utils.OverrideVerkle.Name) {
		v := ctx.Uint64(utils.OverrideVerkle.Name)
		cfg.Eth.OverrideVerkle = &v
	}

    // 开启性能监控
	// Start metrics export if enabled
	utils.SetupMetrics(&cfg.Metrics)
    // 将以太坊协议服务（如区块同步、交易池、共识机制）注册到节点中
	backend, eth := utils.RegisterEthService(stack, &cfg.Eth)

	// Create gauge with geth system and build information
	if eth != nil { // The 'eth' backend may be nil in light mode
		var protos []string
		for _, p := range eth.Protocols() {
			protos = append(protos, fmt.Sprintf("%v/%d", p.Name, p.Version))
		}
		metrics.NewRegisteredGaugeInfo("geth/info", nil).Update(metrics.GaugeInfoValue{
			"arch":      runtime.GOARCH,
			"os":        runtime.GOOS,
			"version":   cfg.Node.Version,
			"protocols": strings.Join(protos, ","),
		})
	}
    // 提供基于事件的日志订阅功能
	// Configure log filter RPC API.
	filterSystem := utils.RegisterFilterAPI(stack, backend, &cfg.Eth)

	// Configure GraphQL if requested.
	if ctx.IsSet(utils.GraphQLEnabledFlag.Name) {
         // 添加 GraphQL 接口
		utils.RegisterGraphQLService(stack, backend, filterSystem, &cfg.Node)
	}
	// Add the Ethereum Stats daemon if requested.
	if cfg.Ethstats.URL != "" {
        // 若配置了 Ethstats.URL，将节点状态上报到指定的统计服务器
		utils.RegisterEthStatsService(stack, backend, cfg.Ethstats.URL)
	}
    
    // 指定 SyncTargetFlag 时，注册工具以测试区块同步性能。
	// Configure full-sync tester service if requested
	if ctx.IsSet(utils.SyncTargetFlag.Name) {
		hex := hexutil.MustDecode(ctx.String(utils.SyncTargetFlag.Name))
		if len(hex) != common.HashLength {
			utils.Fatalf("invalid sync target length: have %d, want %d", len(hex), common.HashLength)
		}
		utils.RegisterFullSyncTester(stack, eth, common.BytesToHash(hex))
	}

	if ctx.IsSet(utils.DeveloperFlag.Name) {
        // 用于本地测试
		// Start dev mode.
		simBeacon, err := catalyst.NewSimulatedBeacon(ctx.Uint64(utils.DeveloperPeriodFlag.Name), eth)
		if err != nil {
			utils.Fatalf("failed to register dev mode catalyst service: %v", err)
		}
		catalyst.RegisterSimulatedBeaconAPIs(stack, simBeacon)
		stack.RegisterLifecycle(simBeacon)
	} else if ctx.IsSet(utils.BeaconApiFlag.Name) {
        // 配置与外部共识客户端的交互
		// Start blsync mode.
		srv := rpc.NewServer()
		srv.RegisterName("engine", catalyst.NewConsensusAPI(eth))
		blsyncer := blsync.NewClient(utils.MakeBeaconLightConfig(ctx))
		blsyncer.SetEngineRPC(rpc.DialInProc(srv))
		stack.RegisterLifecycle(blsyncer)
	} else {
        // ‌默认模式‌：注册 catalyst 服务，支持以太坊主网的共识逻辑
		// Launch the engine API for interacting with external consensus client.
		err := catalyst.Register(stack, eth)
		if err != nil {
			utils.Fatalf("failed to register catalyst service: %v", err)
		}
	}
	return stack
}

// makeConfigNode loads geth configuration and creates a blank node instance.
func makeConfigNode(ctx *cli.Context) (*node.Node, gethConfig) {
    // 先获取节点默认配置后，加载设置相应配置
	cfg := loadBaseConfig(ctx)
    // 生成新的节点
	stack, err := node.New(&cfg.Node)
	if err != nil {
		utils.Fatalf("Failed to create the protocol stack: %v", err)
	}
	// Node doesn't by default populate account manager backends
	if err := setAccountManagerBackends(stack.Config(), stack.AccountManager(), stack.KeyStoreDir()); err != nil {
		utils.Fatalf("Failed to set account manager backends: %v", err)
	}
	// 将与以太坊相关的命令行标志应用到配置中
	utils.SetEthConfig(ctx, stack, &cfg.Eth)
	if ctx.IsSet(utils.EthStatsURLFlag.Name) {
		cfg.Ethstats.URL = ctx.String(utils.EthStatsURLFlag.Name)
	}
	applyMetricConfig(ctx, &cfg)

	return stack, cfg
}

// flags.go
// 配置和启用以太坊节点的性能监控系统‌，收集并导出运行时指标（如 CPU、内存、网络、区块处理速度等）
// SetupMetrics configures the metrics system.
func SetupMetrics(cfg *metrics.Config) {
	if !cfg.Enabled {
		return
	}
	log.Info("Enabling metrics collection")
	metrics.Enable()//激活全局指标注册表

	// InfluxDB exporter.
	var (
		enableExport   = cfg.EnableInfluxDB
		enableExportV2 = cfg.EnableInfluxDBV2
	)
    // InfluxDB 导出,禁止同时启用 v1 和 v2，避免配置冲突
	if cfg.EnableInfluxDB && cfg.EnableInfluxDBV2 {
		Fatalf("Flags %v can't be used at the same time", strings.Join([]string{MetricsEnableInfluxDBFlag.Name, MetricsEnableInfluxDBV2Flag.Name}, ", "))
	}
	var (
		endpoint = cfg.InfluxDBEndpoint
		database = cfg.InfluxDBDatabase
		username = cfg.InfluxDBUsername
		password = cfg.InfluxDBPassword

		token        = cfg.InfluxDBToken
		bucket       = cfg.InfluxDBBucket
		organization = cfg.InfluxDBOrganization
		tagsMap      = SplitTagsFlag(cfg.InfluxDBTags)
	)
	if enableExport {
		log.Info("Enabling metrics export to InfluxDB")
		go influxdb.InfluxDBWithTags(metrics.DefaultRegistry, 10*time.Second, endpoint, database, username, password, "geth.", tagsMap)
	} else if enableExportV2 {
		tagsMap := SplitTagsFlag(cfg.InfluxDBTags)
		log.Info("Enabling metrics export to InfluxDB (v2)")
		go influxdb.InfluxDBV2WithTags(metrics.DefaultRegistry, 10*time.Second, endpoint, token, bucket, organization, "geth.", tagsMap)
	}

	// Expvar exporter.
	if cfg.HTTP != "" {
        // 若配置了 cfg.HTTP 和 cfg.Port，启动独立的 HTTP 服务器，通过 /debug/metrics 提供实时指标数据。
		address := net.JoinHostPort(cfg.HTTP, fmt.Sprintf("%d", cfg.Port))
		log.Info("Enabling stand-alone metrics HTTP endpoint", "address", address)
		exp.Setup(address)
	} else if cfg.HTTP == "" && cfg.Port != 0 {
        // 若仅指定端口但未配置 HTTP 地址，记录警告日志。
		log.Warn(fmt.Sprintf("--%s specified without --%s, metrics server will not start.", MetricsPortFlag.Name, MetricsHTTPFlag.Name))
	}

    // 启动后台协程，每 3 秒采集一次进程资源使用情况（CPU、内存、文件描述符等
	// Enable system metrics collection.
	go metrics.CollectProcessMetrics(3 * time.Second)
}
```





