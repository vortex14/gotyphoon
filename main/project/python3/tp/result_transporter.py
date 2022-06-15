from typhoon.run import ConfigurationComponent
from typhoon.components.result_transporter.resulttransporter import ResultTransporter
from typhoon.components.result_transporter.api.api import ResultWorkerApi
from typhoon.components.result_transporter.state.state import StateManager


if __name__ == "__main__":
    component = ConfigurationComponent("result_transporter", ResultTransporter, ResultWorkerApi, StateManager)
    component.run()