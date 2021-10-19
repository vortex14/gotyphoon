from typhoon.components.fetcher.responses.strategies.db.base_handler import BaseHandler


class DbResponseHandler(BaseHandler):

    def __init__(self, response_obj):
        super().__init__(response_obj)

    def handler_default(self):
        return True

    def max_retry_handler(self):
        pass
    def exception_handler(self):
        pass

    def success_handler(self):
        pass
