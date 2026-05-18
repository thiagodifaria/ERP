import json
from urllib import error as urlerror
from urllib import parse, request


def redact_url_secret(value: str) -> str:
    redacted_keys = {"access_key", "access_token", "apikey", "apiKey", "key", "token"}
    try:
        parsed = parse.urlsplit(value)
        query = parse.parse_qsl(parsed.query, keep_blank_values=True)
        redacted_query = parse.urlencode([(key, "***" if key in redacted_keys else item_value) for key, item_value in query])
        return parse.urlunsplit((parsed.scheme, parsed.netloc, parsed.path, redacted_query, parsed.fragment))
    except ValueError:
        redacted = value
        for key in redacted_keys:
            redacted = redacted.replace(f"{key}=", f"{key}=***")
        return redacted


def http_json(method: str, url: str, headers: dict[str, str], payload: dict | None = None, timeout: int = 8) -> dict:
    data = None
    request_headers = {"Accept": "application/json", **headers}
    if payload is not None:
        data = json.dumps(payload).encode("utf-8")
        request_headers["Content-Type"] = "application/json"
    req = request.Request(url, data=data, headers=request_headers, method=method)
    try:
        with request.urlopen(req, timeout=timeout) as response:
            body = response.read().decode("utf-8")
            parsed = json.loads(body) if body else {}
            return {"ok": 200 <= response.status < 300, "statusCode": response.status, "body": parsed}
    except urlerror.HTTPError as exc:
        body = exc.read().decode("utf-8")
        try:
            parsed = json.loads(body) if body else {}
        except json.JSONDecodeError:
            parsed = {"message": body[:500]}
        return {"ok": False, "statusCode": exc.code, "body": parsed}
    except (urlerror.URLError, TimeoutError) as exc:
        return {"ok": False, "statusCode": 0, "body": {"message": redact_url_secret(str(exc))}}


def http_form(method: str, url: str, headers: dict[str, str], payload: dict, timeout: int = 8) -> dict:
    data = parse.urlencode(payload).encode("utf-8")
    req = request.Request(url, data=data, headers={"Accept": "application/json", "Content-Type": "application/x-www-form-urlencoded", **headers}, method=method)
    try:
        with request.urlopen(req, timeout=timeout) as response:
            body = response.read().decode("utf-8")
            parsed = json.loads(body) if body else {}
            return {"ok": 200 <= response.status < 300, "statusCode": response.status, "body": parsed}
    except urlerror.HTTPError as exc:
        body = exc.read().decode("utf-8")
        try:
            parsed = json.loads(body) if body else {}
        except json.JSONDecodeError:
            parsed = {"message": body[:500]}
        return {"ok": False, "statusCode": exc.code, "body": parsed}
    except (urlerror.URLError, TimeoutError) as exc:
        return {"ok": False, "statusCode": 0, "body": {"message": redact_url_secret(str(exc))}}
