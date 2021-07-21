from typhoon.components.fetcher.executions.strategies.base_pre_fetch import BasePreFetch


class SimplePreFetch(BasePreFetch):

    def __init__(self, task, config):
        super().__init__(task, config)

    async def http_pre_fetch(self):
        pass

    async def local_pre_fetch(self):
        pass

    async def ftp_pre_fetch(self):
        pass

    async def database_pre_fetch(self):
        pass
