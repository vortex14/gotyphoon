from typhoon.api.components.base import BaseComponent


class Projects(BaseComponent):

    def __init__(self, request, state, paths):
        """type(request) is <class 'aiohttp.web_request.Request'>"""

        super().__init__(request, state, paths)

        self.events = {
            "get_projects": {
                "method": self.get_projects,
                "type": "GET"
            }
        }

    async def get_projects(self):
        return {
            "count_projects": 100
        }