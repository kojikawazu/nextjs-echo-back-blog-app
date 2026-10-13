#!/usr/bin/env bash
# scripts/check-secret-files.sh のテスト（#126）。CI（secret-scan.yml）で実スキャンの前に実行し、検出ルールの退行を防ぐ。
#
# スクリプトは git ls-files で対象を引くため、ケースごとに使い捨ての git リポジトリを作り、
# そこへスクリプトをコピーして実行する（追跡済み・未追跡・ignore 済み・履歴の区別は実際の git でしか確かめられない）。
#
# 使い方: scripts/check-secret-files_test.sh
set -euo pipefail

script="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/check-secret-files.sh"

tmp_root=$(mktemp -d)
trap 'rm -rf "$tmp_root"' EXIT

failures=0

# 使い捨てリポジトリを作り、パスを $repo に入れる
new_repo() {
  repo=$(mktemp -d "$tmp_root/repo.XXXXXX")
  git -C "$repo" init -q
  mkdir -p "$repo/scripts"
  cp "$script" "$repo/scripts/"
}

# put <パス> [内容] : ファイルを作る（未追跡のまま）
put() {
  mkdir -p "$(dirname "$repo/$1")"
  printf '%s\n' "${2:-dummy}" > "$repo/$1"
}

# track <パス> [内容] : ファイルを作り、.gitignore を無視して追跡させる（git add -f 相当の混入）
track() {
  put "$@"
  git -C "$repo" add -f -- "$1"
}

commit() {
  git -C "$repo" -c user.name=test -c user.email=test@example.com -c commit.gpgsign=false \
    commit -q -m "$1"
}

# assert_exit <期待する終了コード> <説明> [スクリプトの引数...]
assert_exit() {
  local want=$1 desc=$2 got=0
  shift 2
  out=$("$repo/scripts/check-secret-files.sh" "$@" 2>&1) || got=$?
  if [ "$got" -eq "$want" ]; then
    echo "ok   $desc"
  else
    echo "FAIL $desc (want exit $want, got $got)"
    printf '%s\n' "$out" | sed 's/^/     /'
    failures=$((failures + 1))
  fi
}

# assert_output_has <説明> <行> : 直前の assert_exit の出力に、字下げされた検出行があるか
assert_output_has() {
  if grep -qxF "  $2" <<< "$out"; then
    echo "ok   $1"
  else
    echo "FAIL $1 (missing: $2)"
    printf '%s\n' "$out" | sed 's/^/     /'
    failures=$((failures + 1))
  fi
}

# --- 正常系: 追跡済みの秘匿ファイルをファイル名で検出する ---
for path in \
  config/master.key \
  front/.env \
  .env.local \
  base/.env.production \
  certs/server.pem \
  certs/client.p12 \
  secrets/id_rsa \
  home/.ssh/id_ed25519 \
  gcp/serviceAccountKey.json \
  credentials.json \
  terraform/terraform.tfstate \
  terraform/terraform.tfstate.backup \
  terraform/terraform.tfvars \
  terraform/prod.auto.tfvars \
  設定/.env; do
  new_repo
  track "$path"
  assert_exit 1 "detect tracked name: $path"
done
# 日本語パスが "\350..." とエスケープされず、そのまま出力されること
assert_output_has "non-ASCII path is reported as-is" "設定/.env"

# --- 準正常系: テンプレート・型定義・紛らわしい名前は素通りする ---
for path in \
  backend/.env.example \
  .env.local.example \
  .env.sample \
  .env.template \
  terraform/terraform.tfvars.example \
  src/env.d.ts \
  front/.env.d.ts \
  docs/keyboard.md \
  src/apiKey.ts \
  readme.md \
  terraform/variables.tf \
  terraform/secrets.tf \
  backend/go.sum; do
  new_repo
  track "$path"
  assert_exit 0 "allow tracked name: $path"
done

# --- 準正常系: 追跡対象が空なら OK（スクリプト自身も検出しない）---
new_repo
assert_exit 0 "empty repository"

# --- 正常系: 内容で検出する（任意の名前で保存された鍵）---
# パターンそのものをソースに書くとこのファイル自体が検出されるため、実行時に連結して作る。
begin='-----BEGIN'
new_repo
track keys/server.txt "$begin PRIVATE KEY-----"
assert_exit 1 "detect content: PEM private key"
new_repo
track keys/legacy.txt "$begin RSA PRIVATE KEY-----"
assert_exit 1 "detect content: RSA private key"
new_repo
track keys/deploy "$begin OPENSSH PRIVATE KEY-----"
assert_exit 1 "detect content: OpenSSH private key"
new_repo
track myproj-1a2b3c.json "$(printf '{\n  "type": "service_account",\n  "%s": "x"\n}' 'private_key')"
assert_exit 1 "detect content: service account key JSON"
assert_output_has "content hit is reported" "myproj-1a2b3c.json"

# --- 準正常系: 公開鍵や、キー名への言及だけなら素通りする ---
new_repo
track keys/server.pub "$begin PUBLIC KEY-----"
assert_exit 0 "allow content: public key"
new_repo
track docs/guide.md 'サービスアカウントキーの private_key フィールドを検出する'
assert_exit 0 "allow content: mention of private_key"

# --- 正常系: .gitignore に書き忘れた未追跡の秘匿ファイルを検出する ---
new_repo
put backend/.env
assert_exit 1 "detect untracked name not ignored: backend/.env"
new_repo
put myproj-1a2b3c.json "$begin PRIVATE KEY-----"
assert_exit 1 "detect untracked content not ignored"

# --- 準正常系: .gitignore 済みの未追跡ファイルは素通りする ---
new_repo
put .gitignore $'.env\nterraform/terraform.tfstate\nkeys/'
put .env
put terraform/terraform.tfstate
put keys/server.txt "$begin PRIVATE KEY-----"
assert_exit 0 "allow ignored untracked files"

# --- 異常系: 検出と素通りが混在したら、検出分だけを出力して失敗する ---
new_repo
track backend/.env.example
track terraform/terraform.tfstate
track readme.md
assert_exit 1 "mixed tracked files"
assert_output_has "mixed: detected path is reported" "terraform/terraform.tfstate"
if grep -qE 'backend/\.env\.example|readme\.md' <<< "$out"; then
  echo "FAIL mixed: allowed paths are not reported"
  failures=$((failures + 1))
else
  echo "ok   mixed: allowed paths are not reported"
fi

# --- 正常系: --history は追跡から外した後も履歴上のパスを検出する ---
new_repo
track .env
commit 'add .env'
git -C "$repo" rm -q --cached .env
put .gitignore '.env'
git -C "$repo" add .gitignore
commit 'untrack .env'
assert_exit 0 "removed from index: default mode passes"
assert_exit 1 "removed from index: --history detects" --history
assert_output_has "--history reports the path" ".env"

new_repo
track backend/.env.example
commit 'add template'
assert_exit 0 "--history allows templates" --history

# --- 異常系: 不明な引数・git 管理外では検査せずに失敗する ---
new_repo
assert_exit 2 "unknown argument" --foo
repo=$(mktemp -d "$tmp_root/nogit.XXXXXX")
mkdir -p "$repo/scripts"
cp "$script" "$repo/scripts/"
assert_exit 1 "outside a git repository"

if [ "$failures" -gt 0 ]; then
  echo "$failures 件のテストが失敗しました"
  exit 1
fi
echo "全テスト成功"
