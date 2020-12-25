package models

import "errors"

var (
	ErrAPINameIsEmpty           = errors.New("api_name is empty")
	ErrAPIAddressIsEmpty        = errors.New("api_address is empty")
	ErrRegAddressIsEmpty        = errors.New("register_address is empty")
	ErrAPITTLIsEmpty            = errors.New("api_ttl is empty")
	ErrAPIIntervalIsEmpty       = errors.New("api_interval is empty")
	ErrServiceRegAddressIsEmpty = errors.New("serviceaddress is empty")
	ErrBanlancerIsEmtpy         = errors.New("banlancer is empty")
	ErrWeightRatiosIsEmtpy      = errors.New("weighted_ratios is empty")
	ErrBanlancerIntervalIsEmtpy = errors.New("interval is empty")
)
