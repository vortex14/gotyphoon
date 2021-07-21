from project.donor.v1 import v1
from typhoon.api.components.base import BaseComponent


class Routes(BaseComponent):

    def __init__(self, config, request, paths):
        """type(request) is <class 'aiohttp.web_request.Request'>"""

        super().__init__(request, config, paths)

        self.components = {
            "v1": v1.V1
        }