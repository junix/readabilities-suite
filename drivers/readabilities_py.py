"""Thin native driver for readabilities-py; it does not clean participant output."""

from __future__ import annotations

import json
import sys
from pathlib import Path


def main() -> int:
    from readabilities.web_pages import html_to_markdown, readability_of

    if sys.argv[1:] == ["--doctor"]:
        print(json.dumps({"ready": True, "surface": "readabilities.web_pages.readability_of"}))
        return 0
    if len(sys.argv) < 2:
        print("usage: readabilities_py.py HTML [BASE_URL]", file=sys.stderr)
        return 2
    result = readability_of(Path(sys.argv[1]).read_text(), return_plain_text=False)
    html_content = result.get("content") or ""
    content = html_to_markdown(html_content)
    status = "success" if content.strip() else "no_content"
    metadata = {
        "title": result.get("title") or "",
        "author": result.get("byline") or "",
        "published": result.get("date") or "",
    }
    print(
        json.dumps(
            {
                "schema_version": 1,
                "status": status,
                "content": content,
                "content_format": "markdown",
                "metadata": metadata,
            },
            ensure_ascii=False,
        )
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
