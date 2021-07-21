from typhoon.run import ConfigurationComponent
from typhoon.components.scheduler.scheduler import Scheduler
from typhoon.components.scheduler.api.api import SchedulerApi
from typhoon.components.scheduler.state.state import StateManager


if __name__ == "__main__":
    component = ConfigurationComponent("scheduler", Scheduler, SchedulerApi, StateManager)
    component.run()