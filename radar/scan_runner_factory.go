package radar

import (
	"time"

	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/resource"

	"code.cloudfoundry.org/clock"
	"code.cloudfoundry.org/lager"
)

//go:generate counterfeiter . Scanner

type Scanner interface {
	Run(lager.Logger, string) (time.Duration, error)
	Scan(lager.Logger, string) error
	ScanFromVersion(lager.Logger, string, atc.Version) error
}

//go:generate counterfeiter . ScanRunnerFactory

type ScanRunnerFactory interface {
	ScanResourceRunner(lager.Logger, string) IntervalRunner
	ScanResourceTypeRunner(lager.Logger, string) IntervalRunner
}

type scanRunnerFactory struct {
	clock               clock.Clock
	resourceScanner     Scanner
	resourceTypeScanner Scanner
}

func NewScanRunnerFactory(
	resourceFactory resource.ResourceFactory,
	defaultInterval time.Duration,
	dbPipeline db.Pipeline,
	clock clock.Clock,
	externalURL string,
) ScanRunnerFactory {
	resourceScanner := NewResourceScanner(
		clock,
		resourceFactory,
		defaultInterval,
		dbPipeline,
		externalURL,
	)
	resourceTypeScanner := NewResourceTypeScanner(
		resourceFactory,
		defaultInterval,
		dbPipeline,
		externalURL,
	)

	return &scanRunnerFactory{
		clock:               clock,
		resourceScanner:     resourceScanner,
		resourceTypeScanner: resourceTypeScanner,
	}
}

func (sf *scanRunnerFactory) ScanResourceRunner(logger lager.Logger, name string) IntervalRunner {
	return NewIntervalRunner(logger.Session("interval-runner"), sf.clock, name, sf.resourceScanner)
}

func (sf *scanRunnerFactory) ScanResourceTypeRunner(logger lager.Logger, name string) IntervalRunner {
	return NewIntervalRunner(logger.Session("interval-runner"), sf.clock, name, sf.resourceTypeScanner)
}
