#!/usr/bin/env python3
"""同步 optiwin-release 发布页：更新版本号/日期/下载链接，替换 exe 资产（不保留历史版本）。

用法（GitHub Actions sync-site job 调用）:
  python3 scripts/sync_release_page.py \
    --post site/src/blog/optiwin-release/index.mdx \
    --asset-dir site/public/asset \
    --exe dist/OptiWin-v1.7.0.exe \
    --version v1.7.0
"""
import argparse
import os
import re
import shutil
from datetime import UTC, datetime, timedelta


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument('--post', required=True, help='发布页文章路径')
    ap.add_argument('--asset-dir', required=True, help='exe 资产目录')
    ap.add_argument('--exe', required=True, help='新版本 exe 路径')
    ap.add_argument('--version', required=True, help='新版本号（如 v1.7.0）')
    args = ap.parse_args()

    version = args.version.lstrip('v')
    # runner 时间为 UTC，加 8 小时得东八区
    now = datetime.now(UTC) + timedelta(hours=8)
    date_str = now.strftime('%Y-%m-%d')
    date_cn = now.strftime('%Y年%m月%d日')

    content = open(args.post, encoding='utf-8').read()

    # 1. 从标题提取旧版本号
    m = re.search(r'OptiWin v([\d.]+) 发布', content)
    if not m:
        raise SystemExit('错误: 未在文章标题中找到版本号')
    old_ver = m.group(1)

    # 2. 全篇替换旧版本号 → 新版本号（标题/描述/版本表/下载链接）
    content = content.replace(f'v{old_ver}', f'v{version}')

    # 3. frontmatter 日期
    content = re.sub(r'^date: \d{4}-\d{2}-\d{2}$', f'date: {date_str}', content, flags=re.M)

    # 4. 版本信息表发布日期
    content = re.sub(
        r'\| \*\*发布日期\*\* \| [^|]+ \|',
        f'| **发布日期** | {date_cn} |',
        content,
    )

    open(args.post, 'w', encoding='utf-8').write(content)

    # 5. 替换 exe 资产：删除全部旧版本，复制新版本
    for f in os.listdir(args.asset_dir):
        if f.startswith('OptiWin-') and f.endswith('.exe'):
            os.remove(os.path.join(args.asset_dir, f))
    shutil.copy2(args.exe, os.path.join(args.asset_dir, f'OptiWin-v{version}.exe'))

    print(f'已同步发布页: v{old_ver} → v{version} ({date_cn})')


if __name__ == '__main__':
    main()