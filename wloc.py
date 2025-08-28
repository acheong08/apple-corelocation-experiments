import requests
from mitmproxy.http import HTTPFlow as Flow


class ModifyResponse:
    def response(self, flow: Flow):
        if flow.response is None:
            return
        print(flow.request.url)
        if flow.request.url.endswith("/clls/wloc"):
            # Make request to local server instead
            resp = requests.post(
                "http://127.0.0.1:9090/clls/wloc", data=flow.request.get_content()
            )
            flow.response.set_content(resp.content)
            print("Overwritten")


addons = [ModifyResponse()]
