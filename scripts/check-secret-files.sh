#!/usr/bin/env bash
# 鍵・.env などの秘匿ファイルが Git 管理下に入っていないことを検証する（#126）。
# CI（.github/workflows/secret-scan.yml）とローカルが同じこのスクリプトを実行する。
# 検出パターンを 2 箇所に書き写さないための一元化。
# 手順を CI とローカルに書き分けると、必ずどちらかが腐る。
#
# .gitignore は「未追跡ファイルを追跡させない」側の対策でしかない。一度追跡された
# ファイルには効かず、git add -f や書き漏れも止められない。本スクリプトは
# 「追跡された時点で落とす」検出側を担う。
#
# 対象は「git が見ているファイル」（追跡済み + 未追跡。.gitignore 対象は除く）。
# --others を含めるのは、.gitignore に書き忘れた .env をコミット前に手元で捕まえるため。
# CI では全てコミット済みなので結果は変わらない。
#
# 検査は 2 種類。
#   1. ファイル名: 鍵・証明書・認証情報・環境変数ファイル・Terraform の state / 変数
#   2. 内容: PEM 秘密鍵・GCP サービスアカウントキー（任意のファイル名で保存されるため、名前では捕まらない）
#
# 使い方:
#   scripts/check-secret-files.sh            追跡中 + 未追跡のファイルを検査（CI はこれ）
#   scripts/check-secret-files.sh --history  履歴上の全パスを検査（導入時に一度だけ。ファイル名のみ）
#
# 既定では履歴を走査しない。git ls-files はインデックスを読むだけなので数秒で終わり、
# パスフィルタで分岐させずに常時実行できる。その代わり導入前に混入したものは
# 見つからないため、導入時に --history で一度だけ履歴を照合する。検出パターンを
# 別のコマンドへ書き写さないよう、同じスクリプトのオプションで切り替える。
#
# テスト: scripts/check-secret-files_test.sh
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

# 検出するパス（末尾一致）。鍵・証明書・認証情報・環境変数ファイル。
# Terraform の state には Secret Manager に入れる値（jwt_secret_key 等）が平文で残るため、
# state と変数ファイルも対象にする（本リポジトリは terraform/ を持つ）。
SECRET_PATTERN='(^|/)(\.env(\..+)?|[^/]+\.(key|pem|p12|pfx|jks|keystore|tfstate|tfvars)|[^/]+\.tfstate\.backup|id_rsa|id_ed25519|id_dsa|credentials\.json|serviceAccountKey\.json)$'
# 誤検知になるテンプレートと型定義は除外する。正当な理由で上記に当たるファイルを
# 管理したい場合は、ここへ理由をコメントで添えて足す（ファイル単位で最小に絞る）。
ALLOW_PATTERN='\.(example|sample|template|dist)$|\.env\.d\.ts$'
# 内容で検出する秘密鍵。PEM 秘密鍵（RSA / EC / OPENSSH / ENCRYPTED 等を含む）と、
# サービスアカウントキー JSON の private_key フィールド。
PEM_PATTERN='-----BEGIN ([A-Z]+ )*PRIVATE KEY-----'
SA_KEY_PATTERN='"private_key"[[:space:]]*:'

# git 管理外では「何が追跡されているか」がそもそも定義できない。ここで git ls-files を
# そのまま使うと、失敗がプロセス置換の中に閉じて set -e に捕まらず、対象 0 件の
# 「問題なし」として静かに終わる。黙って検査を飛ばさない。
if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "git 管理下ではないため検査できません（git init 後に実行してください）。" >&2
  exit 1
fi

# 日本語などのパスを "\346..." と引用符付きでエスケープさせない（末尾一致が効かなくなる）
git_paths() { git -c core.quotePath=false "$@"; }
# 検出したパスを字下げして並べる
indent() { while IFS= read -r line; do printf '  %s\n' "$line"; done <<< "$1"; }

# 標準入力のパス一覧から、ファイル名が秘匿ファイルに当たるものを出す
match_names() { grep -E "$SECRET_PATTERN" | grep -vE "$ALLOW_PATTERN" || true; }

# 標準入力のパス一覧から、内容に秘密鍵を含むファイルを出す。
# grep の終了コード 2（読み取りエラー）を「該当なし」と混同しないよう、明示的に失敗させる。
match_contents() {
  local path rc
  while IFS= read -r path; do
    [ -f "$path" ] || continue
    rc=0
    grep -qIE -e "$PEM_PATTERN" -e "$SA_KEY_PATTERN" -- "$path" || rc=$?
    case "$rc" in
      0) printf '%s\n' "$path" ;;
      1) ;;
      *)
        echo "ファイルを読み取れませんでした: $path" >&2
        return 2
        ;;
    esac
  done
}

# パス一覧に対してファイル名・内容の両方を検査し、検出したパスを重複なく出す。
# コマンド置換の中には set -e が引き継がれないため、失敗は || で明示的に伝える。
scan() {
  local paths=$1 by_name by_content
  [ -n "$paths" ] || return 0
  by_name="$(match_names <<< "$paths")"
  by_content="$(match_contents <<< "$paths")" || return 2
  printf '%s\n%s\n' "$by_name" "$by_content" | grep -v '^$' | sort -u || true
}

if [ "${1:-}" = "--history" ]; then
  # 各コミットで追加されたパスを全て引く。git rev-list --objects は blob 単位で重複を
  # 除くため、内容が同じファイル（.env と .env.example が同一など）は片方の名前しか
  # 出ず取りこぼす。--no-renames はリネームを「追加」として数えるため。
  history_hits="$(git_paths log --all --no-renames --diff-filter=A --name-only --pretty=format: | sort -u | match_names)"
  if [ -n "$history_hits" ]; then
    echo "::error::履歴上に秘匿ファイルのパスがあります。"
    indent "$history_hits"
    echo "  追跡から外しても履歴には残る。含まれていた鍵・トークンをローテーションする。"
    exit 1
  fi
  echo "OK: 履歴上に秘匿ファイルのパスはありません"
  exit 0
elif [ "$#" -gt 0 ]; then
  echo "不明な引数: $1（使えるのは --history のみ）" >&2
  exit 2
fi

tracked="$(scan "$(git_paths ls-files --cached)")" || exit 2
untracked="$(scan "$(git_paths ls-files --others --exclude-standard)")" || exit 2

status=0

if [ -n "$tracked" ]; then
  echo "::error::秘匿ファイルが Git 管理下にあります（ファイル名または内容の秘密鍵で検出）。"
  indent "$tracked"
  echo "  git rm --cached で追跡から外し、.gitignore に追加する。"
  echo "  push 済みなら履歴に残っているため、鍵・トークンのローテーションが必要。"
  status=1
fi

if [ -n "$untracked" ]; then
  echo "::error::.gitignore されていない秘匿ファイルがあります（git add すると混入します）。"
  indent "$untracked"
  echo "  .gitignore に追加する。"
  status=1
fi

if [ "$status" -eq 0 ]; then
  echo "OK: 追跡対象・追跡候補に秘匿ファイルはありません"
fi
exit "$status"
