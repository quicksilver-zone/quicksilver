package containers

// ImageConfig contains all images and their respective tags
// needed for running e2e tests.
type ImageConfig struct {
	InitRepository string
	InitTag        string

	QuicksilversRepository string
	QuicksilverTag         string

	HermesRepository string
	HermesTag        string

	ICQRepository string
	ICQTag        string

	XCCLookupRepository string
	XCCLookupTag        string
}

//nolint:deadcode
const (
	// Current Git branch quicksilver repo/version. It is meant to be built locally.
	// It is used when skipping upgrade by setting QUICKSILVER_E2E_SKIP_UPGRADE to true).
	// This image should be pre-built with `make docker-build-debug` either in CI or locally.

	CurrentBranchQuickSilverRepository = "quicksilver"
	CurrentBranchQuicksilverTag        = "debug"
	// Pre-upgrade quicksilver repo/tag to pull.
	// It should be uploaded to Docker Hub. QUICKSILVER_E2E_SKIP_UPGRADE should be unset
	// for this functionality to be used.
	previousVersionQuicksilverRepository = "quicksilverzone/quicksilver"
	previousVersionQuicksilverTag        = "v1.2.4"
	// Pre-upgrade repo/tag for quicksilver initialization (this should be one version below upgradeVersion)
	previousVersionInitRepository = "quicksilverzone/quicksilver"
	previousVersionInitTag        = "v1.2.7"
	// Hermes repo/version for relayer
	hermesRepository = "informalsystems/hermes"
	hermesTag        = "1.3.0"
	// ICQ repo/version for relayer
	icqRepository = "quicksilverzone/interchain-queries"
	icqTag        = "v0.8.8-beta.1"
	// xccLookup repo/version
	xccLookupRepository = "quicksilverzone/xcclookup"
	xccLookupTag        = "v0.4.2"
)

// NewImageConfig returns ImageConfig needed for running e2e test.
// If isUpgrade is true, returns images for running the upgrade
// If isFork is true, utilizes provided fork height to initiate fork logic
func NewImageConfig(isUpgrade, isFork bool) ImageConfig {
	config := ImageConfig{
		HermesRepository:    hermesRepository,
		HermesTag:           hermesTag,
		ICQRepository:       icqRepository,
		ICQTag:              icqTag,
		XCCLookupRepository: xccLookupRepository,
		XCCLookupTag:        xccLookupTag,
	}

	if !isUpgrade {
		// If upgrade is not tested, we do not need InitRepository and InitTag
		// because we directly call the initialization logic without
		// the need for Docker.
		config.QuicksilversRepository = CurrentBranchQuickSilverRepository
		config.QuicksilverTag = CurrentBranchQuicksilverTag
		return config
	}

	// If upgrade is tested, we need to utilize InitRepository and InitTag
	// to initialize older state with Docker
	config.InitRepository = previousVersionInitRepository
	config.InitTag = previousVersionInitTag

	if isFork {
		// Forks are state compatible with earlier versions before fork height.
		// Normally, validators switch the binaries pre-fork height
		// Then, once the fork height is reached, the state breaking-logic
		// is run.
		config.QuicksilversRepository = CurrentBranchQuickSilverRepository
		config.QuicksilverTag = CurrentBranchQuicksilverTag
	} else {
		// Upgrades are run at the time when upgrade height is reached
		// and are submitted via a governance proposal. Therefore, we
		// must start running the previous Quicksilver version. Then, the node
		// should auto-upgrade, at which point we can restart the updated
		// Quicksilver validator container.
		config.QuicksilversRepository = previousVersionQuicksilverRepository
		config.QuicksilverTag = previousVersionQuicksilverTag
	}

	return config
}
