"""Thin native driver for the selected readabilities-py Golden Oracle."""

from __future__ import annotations

import json
import sys
from pathlib import Path


def main() -> int:
    from readabilities.trafilaturas import extract_html

    if sys.argv[1:] == ["--doctor"]:
        print(json.dumps({"ready": True, "surface": "readabilities.trafilaturas.extract_html"}))
        return 0
    if len(sys.argv) < 2:
        print("usage: readabilities_py.py HTML [BASE_URL]", file=sys.stderr)
        return 2
    result = extract_html(
        Path(sys.argv[1]).read_text(encoding="utf-8"),
        base_url=sys.argv[2] if len(sys.argv) > 2 and sys.argv[2] else None,
    )
    content = result["content"]
    status = "success" if content.strip() else "no_content"
    metadata = {
        "title": result["title"],
        "author": result["author"],
        "published": result["published"],
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
