from typhoon.run import ConfigurationComponent
from typhoon.components.processor.processor import Processor
from typhoon.components.processor.api.api import ProcessorApi
from typhoon.components.processor.state.state import StateManager



if __name__ == "__main__":
    component = ConfigurationComponent("processor", Processor, ProcessorApi, StateManager)
    component.run()