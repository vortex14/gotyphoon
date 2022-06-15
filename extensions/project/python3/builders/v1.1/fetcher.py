from typhoon.run import ConfigurationComponent
from typhoon.components.fetcher.fetcher import Fetcher
from typhoon.components.fetcher.api.api import FetcherApi
from typhoon.components.fetcher.state.state import StateManager


if __name__ == "__main__":
    component = ConfigurationComponent("fetcher", Fetcher, FetcherApi, StateManager)
    component.run()