package worker

import (
	"github.com/Supernomad/nexus/nexusd/common"
	"github.com/Supernomad/nexus/nexusd/filter"
	"github.com/Supernomad/nexus/nexusd/iface"
)

type Worker struct {
	filters []filter.Filter
	backend iface.Iface
	log     *common.Logger
	cfg     *common.Config
	done    bool
}

func (worker *Worker) pipeline(queue int) error {
	packet, err := worker.backend.Read(queue)
	if err != nil {
		return err
	}

	for i := 0; i < len(worker.filters); i++ {
		if drop := worker.filters[i].Drop(packet); drop {
			return nil
		}
	}

	return worker.backend.Write(queue, packet)
}

func (worker *Worker) Start(queue int) {
	go func() {
		for !worker.done {
			if err := worker.pipeline(queue); err != nil {
				worker.log.Error.Println("[WORKER]", "Error during work pipeline:", err)
			}
		}
	}()
}

func (worker *Worker) Stop() {
	worker.done = true
}

func New(log *common.Logger, cfg *common.Config, backend iface.Iface, filters []filter.Filter) *Worker {
	return &Worker{
		filters: filters,
		backend: backend,
		log:     log,
		cfg:     cfg,
		done:    false,
	}
}