from mitmproxy.io import FlowReader
from mitmproxy import http

IDENTIFIER = b"Master"

with open("apple.flow", "rb") as f:
    reader = FlowReader(f)
    for flow in reader.stream():
        if not isinstance(flow, http.HTTPFlow):
            continue  # Skip non-HTTP flows (e.g., DNS)
        if flow.request.content and IDENTIFIER in flow.request.content:
            print(f"Found in request: {flow.request.url}")
            with open("/tmp/found.bin", "wb") as out_file:
                out_file.write(flow.request.content)
            break
