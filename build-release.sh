#!/bin/bash
# Execution Ledger release build: wails build → codesign → notarize → staple → zip
# mcp-server-manager/build-release.sh の構成を踏襲（同じApple Developer認証情報を流用）
set -e

# ── 設定（要変更） ──────────────────────────────────────────────────────────
APPLE_ID="chankei613@gmail.com"
TEAM_ID="${TEAM_ID:-}"          # 環境変数 or 直書き（例: "ABC1234567"）
APP_PASSWORD="${APP_PASSWORD:-}" # App-specific password（環境変数推奨）

APP_NAME="execution-ledger"
VERSION=$(grep 'AppVersion' version.go | sed 's/.*"\(.*\)".*/\1/')
APP_PATH="build/bin/${APP_NAME}.app"
ZIP_PATH="${APP_NAME}-${VERSION}.zip"
ENTITLEMENTS="build/darwin/entitlements.plist"

if [ -z "$TEAM_ID" ]; then
  echo "ERROR: TEAM_ID が未設定です。"
  echo "以下のコマンドで Developer ID を確認してください:"
  echo "  security find-identity -v -p codesigning"
  echo ""
  echo "実行例:"
  echo "  TEAM_ID=ABC1234567 APP_PASSWORD=xxxx-xxxx-xxxx-xxxx ./build-release.sh"
  exit 1
fi

IDENTITY="Developer ID Application: keisuke haraguchi (${TEAM_ID})"

echo "==> Building Execution Ledger v${VERSION}..."

xattr -cr build/ 2>/dev/null || true
export PATH="$PATH:$HOME/go/bin"
wails build -platform darwin/universal -o "${APP_NAME}"

echo "==> Code signing..."
codesign \
  --deep \
  --force \
  --verify \
  --verbose \
  --sign "${IDENTITY}" \
  --options runtime \
  --entitlements "${ENTITLEMENTS}" \
  "${APP_PATH}"

codesign --verify --deep --strict "${APP_PATH}"
echo "    Signature OK"

echo "==> Creating zip for notarization..."
ditto -c -k --keepParent "${APP_PATH}" "${ZIP_PATH}"

echo "==> Submitting for notarization (this takes a few minutes)..."
xcrun notarytool submit "${ZIP_PATH}" \
  --apple-id "${APPLE_ID}" \
  --team-id "${TEAM_ID}" \
  --password "${APP_PASSWORD}" \
  --wait

echo "==> Stapling notarization ticket..."
xcrun stapler staple "${APP_PATH}"
xcrun stapler validate "${APP_PATH}"
echo "    Staple OK"

echo "==> Re-zipping stapled app for distribution..."
rm -f "${ZIP_PATH}"
ditto -c -k --keepParent "${APP_PATH}" "${ZIP_PATH}"

echo ""
echo "Done: ${ZIP_PATH}"
echo ""
echo "次のステップ:"
echo "  1. gh release create v${VERSION} ${ZIP_PATH} でGitHub Releaseを作成"
echo "  2. landing/ にzipを同梱してVercelへデプロイ（agent-config-managerと同じ配布パターン）"
