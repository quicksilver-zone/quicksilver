package lockup

import sdkioerrors "cosmossdk.io/errors"

// DONTCOVER

// x/lockup module sentinel errors.
var (
	ErrNotLockOwner                      = sdkioerrors.Register(ModuleName, 1, "msg sender is not the owner of specified lock")
	ErrSyntheticLockupAlreadyExists      = sdkioerrors.Register(ModuleName, 2, "synthetic lockup already exists for same lock and suffix")
	ErrSyntheticDurationLongerThanNative = sdkioerrors.Register(ModuleName, 3, "synthetic lockup duration should be shorter than native lockup duration")
	ErrLockupNotFound                    = sdkioerrors.Register(ModuleName, 4, "lockup not found")
)
