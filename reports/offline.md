# Readabilities comparison report

Generated: 2026-08-11T22:03:27Z

Profile: `offline`
Result: **23 pass / 25 fail**

| Participant | Pass | Fail | Recall | Noise rejection | Structure | Markdown fidelity | Order | Metadata | Security | Markdown bytes | Words | Noise hints | Duplicate lines | Elapsed ms |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| defuddle | 4 | 12 | 43/43 | 34/38 | 23/33 | 67/70 | 14/14 | 19/19 | 1 | 6537 | 778 | 1 | 0 | 7572 |
| readabilities-py | 3 | 13 | 43/43 | 32/38 | 22/33 | 60/70 | 14/14 | 19/19 | 2 | 6596 | 778 | 1 | 0 | 40543 |
| readabilities-rs | 16 | 0 | 43/43 | 38/38 | 33/33 | 70/70 | 14/14 | 19/19 | 0 | 6788 | 819 | 0 | 0 | 1954 |

## Best-evidence gate

- Eligible on frozen corpus: `true`
- Rust security zero: `true`
- Rust recall at least best participant: `true`
- Rust noise rejection at least best participant: `true`
- Rust structure at least best participant: `true`
- Rust Markdown fidelity at least best participant: `true`
- Rust content order at least best participant: `true`
- Rust metadata at least best participant: `true`
- Recall leaders: `defuddle, readabilities-py, readabilities-rs`
- Noise leaders: `readabilities-rs`
- Structure leaders: `readabilities-rs`
- Markdown leaders: `readabilities-rs`
- Order leaders: `defuddle, readabilities-py, readabilities-rs`
- Metadata leaders: `defuddle, readabilities-py, readabilities-rs`

Quality dimensions are release gates; elapsed time cannot compensate for lost content or retained noise. Eligibility is a corpus claim, not universal superiority.

## Per-case Markdown evidence

| Case | Participant | Status | Recall | Noise | Structure | Markdown | Order | Bytes | Words | Headings | Links | Code | Noise hints | Duplicates | Site config | SHA-256 |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|
| READ-001 | defuddle | fail | 4/4 | 3/3 | 5/5 | 7/9 | 1/1 | 1301 | 173 | 1 | 2 | 1 | 0 | 0 |  | `6b2c57f0d46f` |
| READ-001 | readabilities-py | fail | 4/4 | 3/3 | 4/5 | 3/9 | 1/1 | 1285 | 172 | 1 | 2 | 1 | 0 | 0 |  | `16b2a4314e5f` |
| READ-001 | readabilities-rs | pass | 4/4 | 3/3 | 5/5 | 9/9 | 1/1 | 1432 | 177 | 2 | 2 | 1 | 0 | 0 |  | `ed0a7fd16eb2` |
| READ-002 | defuddle | fail | 1/1 | 2/2 | 0/1 | 4/4 | 0/0 | 53 | 7 | 0 | 0 | 0 | 0 | 0 |  | `0691eccc362c` |
| READ-002 | readabilities-py | fail | 1/1 | 2/2 | 0/1 | 4/4 | 0/0 | 53 | 7 | 0 | 0 | 0 | 0 | 0 |  | `0691eccc362c` |
| READ-002 | readabilities-rs | pass | 1/1 | 2/2 | 1/1 | 4/4 | 0/0 | 69 | 11 | 1 | 0 | 0 | 0 | 0 |  | `49f31eb8405c` |
| READ-003 | defuddle | fail | 3/3 | 3/3 | 0/1 | 4/4 | 1/1 | 487 | 4 | 0 | 0 | 0 | 0 | 0 |  | `1894f4489e87` |
| READ-003 | readabilities-py | pass | 3/3 | 3/3 | 1/1 | 4/4 | 1/1 | 519 | 6 | 1 | 0 | 0 | 0 | 0 |  | `7fe5e8f3ad47` |
| READ-003 | readabilities-rs | pass | 3/3 | 3/3 | 1/1 | 4/4 | 1/1 | 518 | 6 | 1 | 0 | 0 | 0 | 0 |  | `8ae72837967c` |
| READ-004 | defuddle | pass | 2/2 | 2/2 | 3/3 | 5/5 | 1/1 | 379 | 60 | 2 | 0 | 1 | 0 | 0 |  | `938c70315815` |
| READ-004 | readabilities-py | fail | 2/2 | 2/2 | 3/3 | 4/5 | 1/1 | 375 | 60 | 2 | 0 | 1 | 0 | 0 |  | `a19295848eaf` |
| READ-004 | readabilities-rs | pass | 2/2 | 2/2 | 3/3 | 5/5 | 1/1 | 428 | 64 | 3 | 0 | 1 | 0 | 0 |  | `01e37ffe068a` |
| READ-005 | defuddle | fail | 3/3 | 3/3 | 0/1 | 4/4 | 1/1 | 385 | 47 | 0 | 0 | 0 | 0 | 0 |  | `085121b803ff` |
| READ-005 | readabilities-py | fail | 3/3 | 3/3 | 0/1 | 4/4 | 1/1 | 385 | 47 | 0 | 0 | 0 | 0 | 0 |  | `085121b803ff` |
| READ-005 | readabilities-rs | pass | 3/3 | 3/3 | 1/1 | 4/4 | 1/1 | 419 | 53 | 1 | 0 | 0 | 0 | 0 |  | `9b29173cb9d9` |
| READ-006 | defuddle | fail | 4/4 | 3/3 | 2/3 | 5/5 | 1/1 | 440 | 62 | 0 | 0 | 0 | 0 | 0 |  | `b8c06c7ea645` |
| READ-006 | readabilities-py | fail | 4/4 | 3/3 | 2/3 | 5/5 | 1/1 | 440 | 62 | 0 | 0 | 0 | 0 | 0 |  | `1b60aedebd40` |
| READ-006 | readabilities-rs | pass | 4/4 | 3/3 | 3/3 | 5/5 | 1/1 | 475 | 67 | 1 | 0 | 0 | 0 | 0 |  | `c70600200af0` |
| READ-007 | defuddle | pass | 3/3 | 2/2 | 2/2 | 4/4 | 1/1 | 383 | 49 | 1 | 0 | 1 | 0 | 0 |  | `893c60ac6173` |
| READ-007 | readabilities-py | pass | 3/3 | 2/2 | 2/2 | 4/4 | 1/1 | 383 | 49 | 1 | 0 | 1 | 0 | 0 |  | `893c60ac6173` |
| READ-007 | readabilities-rs | pass | 3/3 | 2/2 | 2/2 | 4/4 | 1/1 | 421 | 56 | 2 | 0 | 1 | 0 | 0 |  | `d24442333f70` |
| READ-008 | defuddle | fail | 3/3 | 2/2 | 2/3 | 5/5 | 1/1 | 281 | 41 | 0 | 0 | 0 | 0 | 0 |  | `aa15813ac42f` |
| READ-008 | readabilities-py | fail | 3/3 | 2/2 | 1/3 | 4/5 | 1/1 | 263 | 33 | 0 | 0 | 0 | 0 | 0 |  | `524206896933` |
| READ-008 | readabilities-rs | pass | 3/3 | 2/2 | 3/3 | 5/5 | 1/1 | 302 | 45 | 1 | 0 | 0 | 0 | 0 |  | `df9357c7021f` |
| READ-009 | defuddle | fail | 3/3 | 2/2 | 1/2 | 4/5 | 1/1 | 315 | 38 | 0 | 1 | 0 | 0 | 0 |  | `c3a2b817302b` |
| READ-009 | readabilities-py | fail | 3/3 | 2/2 | 1/2 | 4/5 | 1/1 | 315 | 38 | 0 | 1 | 0 | 0 | 0 |  | `c3a2b817302b` |
| READ-009 | readabilities-rs | pass | 3/3 | 2/2 | 2/2 | 5/5 | 1/1 | 362 | 43 | 1 | 1 | 0 | 0 | 0 |  | `ab474af29bc1` |
| READ-010 | defuddle | fail | 2/2 | 2/2 | 1/2 | 4/4 | 1/1 | 272 | 23 | 0 | 1 | 0 | 0 | 0 |  | `312ad31aac25` |
| READ-010 | readabilities-py | fail | 2/2 | 2/2 | 1/2 | 4/4 | 1/1 | 268 | 24 | 0 | 2 | 0 | 0 | 0 |  | `9c3f2c2e0408` |
| READ-010 | readabilities-rs | pass | 2/2 | 2/2 | 2/2 | 4/4 | 1/1 | 243 | 25 | 1 | 1 | 0 | 0 | 0 |  | `98de2734ef30` |
| READ-011 | defuddle | pass | 0/0 | 0/0 | 0/0 | 0/0 | 0/0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |  | `` |
| READ-011 | readabilities-py | pass | 0/0 | 0/0 | 0/0 | 0/0 | 0/0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |  | `` |
| READ-011 | readabilities-rs | pass | 0/0 | 0/0 | 0/0 | 0/0 | 0/0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |  | `` |
| READ-012 | defuddle | fail | 3/3 | 2/3 | 1/2 | 4/4 | 1/1 | 502 | 59 | 0 | 1 | 0 | 0 | 0 |  | `376c968773d1` |
| READ-012 | readabilities-py | fail | 3/3 | 1/3 | 1/2 | 4/4 | 1/1 | 557 | 64 | 0 | 1 | 0 | 0 | 0 |  | `786b4196a2fb` |
| READ-012 | readabilities-rs | pass | 3/3 | 3/3 | 2/2 | 4/4 | 1/1 | 454 | 58 | 1 | 0 | 0 | 0 | 0 | medium | `300289324321` |
| READ-013 | defuddle | fail | 3/3 | 2/3 | 2/2 | 4/4 | 1/1 | 456 | 55 | 1 | 0 | 0 | 0 | 0 |  | `3c8d97d14527` |
| READ-013 | readabilities-py | fail | 3/3 | 1/3 | 2/2 | 4/4 | 1/1 | 492 | 58 | 1 | 0 | 0 | 0 | 0 |  | `6795f8c8517f` |
| READ-013 | readabilities-rs | pass | 3/3 | 3/3 | 2/2 | 4/4 | 1/1 | 408 | 51 | 1 | 0 | 0 | 0 | 0 | wikipedia | `67ba17e67ff4` |
| READ-014 | defuddle | pass | 3/3 | 3/3 | 2/2 | 5/5 | 1/1 | 368 | 42 | 2 | 0 | 1 | 0 | 0 |  | `3443119709df` |
| READ-014 | readabilities-py | fail | 3/3 | 3/3 | 2/2 | 4/5 | 1/1 | 346 | 40 | 1 | 0 | 1 | 0 | 0 |  | `19825d262e1e` |
| READ-014 | readabilities-rs | pass | 3/3 | 3/3 | 2/2 | 5/5 | 1/1 | 367 | 42 | 2 | 0 | 1 | 0 | 0 | mdn | `65b1c293a209` |
| READ-015 | defuddle | fail | 3/3 | 3/3 | 1/2 | 4/4 | 1/1 | 455 | 63 | 0 | 0 | 0 | 0 | 0 |  | `c402b2664ada` |
| READ-015 | readabilities-py | fail | 3/3 | 3/3 | 1/2 | 4/4 | 1/1 | 455 | 63 | 0 | 0 | 0 | 0 | 0 |  | `c402b2664ada` |
| READ-015 | readabilities-rs | pass | 3/3 | 3/3 | 2/2 | 4/4 | 1/1 | 491 | 68 | 1 | 0 | 0 | 0 | 0 |  | `c442d35c2fe6` |
| READ-016 | defuddle | fail | 3/3 | 0/2 | 1/2 | 4/4 | 1/1 | 460 | 55 | 0 | 0 | 0 | 1 | 0 |  | `964605d08539` |
| READ-016 | readabilities-py | fail | 3/3 | 0/2 | 1/2 | 4/4 | 1/1 | 460 | 55 | 0 | 0 | 0 | 1 | 0 |  | `d346a311e79b` |
| READ-016 | readabilities-rs | pass | 3/3 | 2/2 | 2/2 | 4/4 | 1/1 | 399 | 53 | 1 | 0 | 0 | 0 | 0 | custom-example | `f0315dd8a51b` |

## Failures

- `READ-001/defuddle`: case=READ-001 failed_facts=["contains ](https://example.test/details)" "contains ](https://example.test/images/pipeline.png) The relative"] security=[] actual_preview="REQUIRED-CENTRAL-IDEA: A readable-content engine should identify the document's main argument before it thinks about Markdown syntax. This ordering keeps navigation labels, advertisements, and subscription prompts from becoming apparently l"
- `READ-001/readabilities-py`: case=READ-001 failed_facts=["footnote" "image and following prose remain separated" "contains ```rust" "contains ](https://example.test/details)" "contains ](https://example.test/images/pipeline.png) The relative" "contains <sup>1</sup>" "rejects pipeline.png)The relative"] security=[] actual_preview="REQUIRED-CENTRAL-IDEA: A readable-content engine should identify the document's main argument before it thinks about Markdown syntax. This ordering keeps navigation labels, advertisements, and subscription prompts from becoming apparently l"
- `READ-002/defuddle`: case=READ-002 failed_facts=["heading"] security=[] actual_preview="REQUIRED-SHORT: Small pages are still real documents."
- `READ-002/readabilities-py`: case=READ-002 failed_facts=["heading"] security=[] actual_preview="REQUIRED-SHORT: Small pages are still real documents."
- `READ-003/defuddle`: case=READ-003 failed_facts=["heading"] security=[] actual_preview="必须保留：正文提取的第一步是判断页面中什么内容真正属于文章，而不是立刻把所有标签转换成 Markdown。 必须保留：保守的去噪策略同时观察语义容器、文字密度、链接密度和常见的导航模式，并在结果过短时降低清理强度重试。 必须保留：最终安全清理不可关闭，它会删除脚本、事件属性和危险链接，同时保留标题、代码、表格、图片和脚注等有意义的结构。"
- `READ-004/readabilities-py`: case=READ-004 failed_facts=["contains ```json"] security=[] actual_preview="REQUIRED-DOC: The parser accepts HTML and an optional base URL. The local path does not access the network. ## Example ``` {\"source\":\"html\",\"mode\":\"balanced\"} ``` ## Parameters | Name | Meaning | | --- | --- | | source | Input HTML | | mode"
- `READ-005/defuddle`: case=READ-005 failed_facts=["heading"] security=[] actual_preview="REQUIRED-NEWS-LEDE: Marine researchers opened a public observation pier on Tuesday after three years of construction. REQUIRED-NEWS-QUOTE: The director said the facility will help students collect long-term measurements without disrupting c"
- `READ-005/readabilities-py`: case=READ-005 failed_facts=["heading"] security=[] actual_preview="REQUIRED-NEWS-LEDE: Marine researchers opened a public observation pier on Tuesday after three years of construction. REQUIRED-NEWS-QUOTE: The director said the facility will help students collect long-term measurements without disrupting c"
- `READ-006/defuddle`: case=READ-006 failed_facts=["heading"] security=[] actual_preview="REQUIRED-BLOG-ONE: Start by listing the operations that create impact noise, airborne noise, or vibration. - REQUIRED-BLOG-LIST: Separate machines from shared walls. - Seal gaps around doors before adding expensive panels. - Measure again a"
- `READ-006/readabilities-py`: case=READ-006 failed_facts=["heading"] security=[] actual_preview="REQUIRED-BLOG-ONE: Start by listing the operations that create impact noise, airborne noise, or vibration. * REQUIRED-BLOG-LIST: Separate machines from shared walls. * Seal gaps around doors before adding expensive panels. * Measure again a"
- `READ-008/defuddle`: case=READ-008 failed_facts=["heading"] security=[] actual_preview="REQUIRED-MATH: The square of a sum can be expanded without changing its value. $$ a^{2} + 2 a b + b^{2} $$ REQUIRED-MATH-CONTEXT: This representation should survive cleaning even when the page contains unrelated widgets. 1. REQUIRED-MATH-FO"
- `READ-008/readabilities-py`: case=READ-008 failed_facts=["heading" "math" "contains $$\na^{2} + 2 a b + b^{2}\n$$"] security=[] actual_preview="REQUIRED-MATH: The square of a sum can be expanded without changing its value. a2+2ab+b2 REQUIRED-MATH-CONTEXT: This representation should survive cleaning even when the page contains unrelated widgets. 1. REQUIRED-MATH-FOOTNOTE: Variables "
- `READ-009/defuddle`: case=READ-009 failed_facts=["heading" "contains ![A marked field notebook](https://example.test/images/cover.jpg)"] security=[] actual_preview="REQUIRED-IMAGE: Each image needs useful alternative text and an absolute URL in exported content. ![A marked field notebook](/images/cover.jpg) REQUIRED-CAPTION: Measurements recorded at sunrise. REQUIRED-IMAGE-END: Duplicate cover images s"
- `READ-009/readabilities-py`: case=READ-009 failed_facts=["heading" "contains ![A marked field notebook](https://example.test/images/cover.jpg)"] security=[] actual_preview="REQUIRED-IMAGE: Each image needs useful alternative text and an absolute URL in exported content. ![A marked field notebook](/images/cover.jpg) REQUIRED-CAPTION: Measurements recorded at sunrise. REQUIRED-IMAGE-END: Duplicate cover images s"
- `READ-010/defuddle`: case=READ-010 failed_facts=["heading"] security=[iframe_element] actual_preview="REQUIRED-SECURITY: Untrusted HTML must cross a final sanitizer after every extraction backend. REQUIRED-SECURITY-STRUCTURE: Normal paragraph text remains visible. unsafe data link ![safe image](https://example.test/safe.png)"
- `READ-010/readabilities-py`: case=READ-010 failed_facts=["heading"] security=[script_element dangerous_markdown_uri] actual_preview="REQUIRED-SECURITY: Untrusted HTML must cross a final sanitizer after every extraction backend. REQUIRED-SECURITY-STRUCTURE: Normal paragraph text remains visible. [unsafe data link](data:text/html, steal() )![safe image](https://example.tes"
- `READ-012/defuddle`: case=READ-012 failed_facts=["FORBIDDEN-MEDIUM-PROMO" "heading"] security=[] actual_preview="REQUIRED-MEDIUM-OPEN: A site profile should narrow known layout noise without replacing the general readability engine. REQUIRED-MEDIUM-BODY: The portable default still identifies paragraphs, headings, lists, links, and code when no profile"
- `READ-012/readabilities-py`: case=READ-012 failed_facts=["FORBIDDEN-MEDIUM-PROMO" "FORBIDDEN-MEDIUM-RELATED" "heading"] security=[] actual_preview="REQUIRED-MEDIUM-OPEN: A site profile should narrow known layout noise without replacing the general readability engine. REQUIRED-MEDIUM-BODY: The portable default still identifies paragraphs, headings, lists, links, and code when no profile"
- `READ-013/defuddle`: case=READ-013 failed_facts=["FORBIDDEN-WIKI-NAVBOX"] security=[] actual_preview="REQUIRED-WIKI-LEDE: Readable content extraction identifies the primary document while excluding interface chrome. ## History REQUIRED-WIKI-HISTORY: Early systems combined semantic containers with text and link density. - Generic extraction "
- `READ-013/readabilities-py`: case=READ-013 failed_facts=["FORBIDDEN-WIKI-EDIT" "FORBIDDEN-WIKI-NAVBOX"] security=[] actual_preview="REQUIRED-WIKI-LEDE: Readable content extraction identifies the primary document while excluding interface chrome. FORBIDDEN-WIKI-EDIT: edit section ## History REQUIRED-WIKI-HISTORY: Early systems combined semantic containers with text and l"
- `READ-014/readabilities-py`: case=READ-014 failed_facts=["contains ```js"] security=[] actual_preview="REQUIRED-MDN-SUMMARY: The ReadableContent interface represents a cleaned document selected from an HTML snapshot. ## Syntax ``` const article = reader.extract(html); ``` REQUIRED-MDN-SYNTAX: Rendering Markdown happens after the content sele"
- `READ-015/defuddle`: case=READ-015 failed_facts=["heading"] security=[] actual_preview="REQUIRED-GENERIC-OPEN: Most domains will never receive a dedicated site profile, so the default path must remain useful. REQUIRED-GENERIC-BODY: Semantic containers, text density, link density, and conservative retry behavior provide a stron"
- `READ-015/readabilities-py`: case=READ-015 failed_facts=["heading"] security=[] actual_preview="REQUIRED-GENERIC-OPEN: Most domains will never receive a dedicated site profile, so the default path must remain useful. REQUIRED-GENERIC-BODY: Semantic containers, text density, link density, and conservative retry behavior provide a stron"
- `READ-016/defuddle`: case=READ-016 failed_facts=["FORBIDDEN-CUSTOM-SPONSOR" "FORBIDDEN-CUSTOM-TOOLS" "heading"] security=[] actual_preview="REQUIRED-CUSTOM-OPEN: A caller can load a declarative profile without rebuilding the Rust crate. REQUIRED-CUSTOM-BODY: Host matching is boundary aware and selectors are validated before extraction starts. - Choose the content root. - Remove"
- `READ-016/readabilities-py`: case=READ-016 failed_facts=["FORBIDDEN-CUSTOM-SPONSOR" "FORBIDDEN-CUSTOM-TOOLS" "heading"] security=[] actual_preview="REQUIRED-CUSTOM-OPEN: A caller can load a declarative profile without rebuilding the Rust crate. REQUIRED-CUSTOM-BODY: Host matching is boundary aware and selectors are validated before extraction starts. * Choose the content root. * Remove"
