package main

import (
	"context"
	"fmt"

	"code.emcdtech.com/emcd/sdk/app"

	"code.emcdtech.com/b2b/processing/internal/config"
	"code.emcdtech.com/b2b/processing/internal/metrics"
	"code.emcdtech.com/b2b/processing/internal/repository"
	repositoryMetrics "code.emcdtech.com/b2b/processing/internal/repository/metrics"
	"code.emcdtech.com/b2b/processing/internal/repository/pg"
	"code.emcdtech.com/b2b/processing/internal/worker/expirer"
	invoiceStatusMeter "code.emcdtech.com/b2b/processing/internal/worker/invoice_status_meter"
)

type dependencies struct {
	*app.Deps
	cfg        config.APIConfig
	repository struct {
		invoice repository.Invoice
	}
	worker struct {
		expirer            *expirer.Worker
		invoiceStatusMeter *invoiceStatusMeter.Worker
	}
}

func newDeps(ctx context.Context, cfg config.APIConfig) (*dependencies, error) {
	appDeps, err := app.NewDeps(ctx, cfg.DepsConfig, serviceName, serviceVersion)
	if err != nil {
		return nil, fmt.Errorf("newDeps: %w", err)
	}

	deps := &dependencies{
		Deps: &appDeps,
		cfg:  cfg,
	}

	collector := metrics.New()
	appDeps.Collectors = append(appDeps.Collectors, collector)

	deps.setupRepository(collector)
	deps.setupWorker(collector)

	return deps, nil
}

func (d *dependencies) setupRepository(collector *metrics.AppMetrics) {
	d.repository.invoice = repositoryMetrics.NewInvoice(
		pg.NewInvoice(d.Infrastructure.DB.Pool),
		collector.RepositoryRequestTimeHistogram,
	)
}

func (d *dependencies) setupWorker(collector *metrics.AppMetrics) {
	d.worker.expirer = expirer.NewWorker(d.repository.invoice)
	d.worker.invoiceStatusMeter = invoiceStatusMeter.NewWorker(
		d.repository.invoice,
		collector.InvoiceStatusesGauge,
	)
}

func (d *dependencies) HealthCheckers() []app.HealthChecker {
	return []app.HealthChecker{
		&d.Deps.Infrastructure.DB,
	}
}

func (d *dependencies) Close(ctx context.Context) {
	// nothing to close
}
